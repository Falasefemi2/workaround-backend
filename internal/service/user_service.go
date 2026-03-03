package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/email"
	"github.com/falasefemi2/workaround-backend/internal/repository"
)

var (
	ErrEmailTaken   = errors.New("email already in use")
	ErrInvalidCreds = errors.New("invalid email or password")
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidToken = errors.New("invalid or expired token")
)

const defaultCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789" +
	"!@#$%^&*()-_=+"

type UserService struct {
	repo         *repository.UserRepo
	jwtSecret    string
	smtpCfg      email.SMTPConfig
	passwordRepo *repository.PasswordRepo
}

func NewUserService(
	repo *repository.UserRepo,
	smtpCfg email.SMTPConfig,
	jwtSecret string,
	passwordRepo *repository.PasswordRepo,
) *UserService {
	return &UserService{
		repo:         repo,
		jwtSecret:    jwtSecret,
		smtpCfg:      smtpCfg,
		passwordRepo: passwordRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	_, err := s.repo.GetUserByEmail(ctx, params.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return db.User{}, err
	}
	if err == nil {
		return db.User{}, ErrEmailTaken
	}

	plainPassword, err := GeneratePassword(12)
	if err != nil {
		return db.User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, err
	}

	params.PasswordHash = string(hashedPassword)

	user, err := s.repo.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, err
	}

	err = email.SendEmail(s.smtpCfg, email.EmailParams{
		To:      user.Email,
		Subject: "Welcome to Workaround",
		Body: fmt.Sprintf(
			"Hello %s, your temporary password is: <b>%s</b>. Please login and change it.",
			user.Email,
			plainPassword,
		),
	})
	if err != nil {
		return db.User{}, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteUser(ctx, userID)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (db.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	user, err := s.repo.ListUser(ctx, arg)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	user, err := s.repo.UpdateUser(ctx, params)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, userEmail, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, userEmail)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidCreds
	}
	if err != nil {
		return "", err // real error
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCreds
	}
	return s.generateToken(user.ID.String(), user.UserType)
}

func (s *UserService) ForgotPassword(ctx context.Context, userEmail string) error {
	user, err := s.repo.GetUserByEmail(ctx, userEmail)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)

	_, err = s.passwordRepo.CreatePasswordResetToken(ctx, db.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(time.Hour), Valid: true},
	})
	if err != nil {
		return err
	}

	resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)

	return email.SendEmail(s.smtpCfg, email.EmailParams{
		To:      user.Email,
		Subject: "Reset your password",
		Body: fmt.Sprintf(
			"Click this link to reset your password: <a href='%s'>%s</a>. Link expires in 1 hour.",
			resetLink,
			resetLink,
		),
	})
}

func (s *UserService) ResetPassword(ctx context.Context, token, newPassword string) error {
	resetToken, err := s.passwordRepo.GetPasswordResetToken(ctx, token)
	if err != nil {
		return ErrInvalidToken
	}
	if time.Now().After(resetToken.ExpiresAt.Time) {
		return ErrInvalidToken
	}
	// if resetToken.Used {
	// 	return ErrInvalidToken
	// }
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.repo.UpdatePassword(ctx, db.UpdatePasswordParams{
		ID:           resetToken.UserID,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		return err
	}
	return s.passwordRepo.MarkTokenUsed(ctx, token)
}

func (s *UserService) generateToken(userID, userType string) (string, error) {
	claims := jwt.MapClaims{
		"sub":       userID,
		"user_type": userType,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func GeneratePassword(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("password length must be greater than zero")
	}

	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(defaultCharset)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		password[i] = defaultCharset[num.Int64()]
	}

	return string(password), nil
}
