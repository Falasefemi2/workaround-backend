package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/email"
	"github.com/falasefemi2/workaround-backend/internal/repository"
)

var (
	ErrEmailTaken   = errors.New("email already in use")
	ErrInvalidCreds = errors.New("invalid email or password")
)

const defaultCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789" +
	"!@#$%^&*()-_=+"

type UserService struct {
	repo      *repository.UserRepo
	jwtSecret string
	smtpCfg   email.SMTPConfig
}

func NewUserService(
	repo *repository.UserRepo,
	smtpCfg email.SMTPConfig,
	jwtSecret string,
) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
		smtpCfg:   smtpCfg,
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

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
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
