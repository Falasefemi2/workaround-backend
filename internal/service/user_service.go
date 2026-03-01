package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/email"
	"github.com/falasefemi2/workaround-backend/internal/repository"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
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
