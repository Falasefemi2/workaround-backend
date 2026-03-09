package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/email"
	"github.com/falasefemi2/workaround-backend/internal/repository"
)

type CandidateOfferService struct {
	repo         *repository.CandidateOfferRepo
	employeeRepo *repository.EmployeeRepo
	userRepo     *repository.UserRepo
	smtpCfg      email.SMTPConfig
}

func NewCandidateOfferService(
	repo *repository.CandidateOfferRepo,
	employeeRepo *repository.EmployeeRepo,
	userRepo *repository.UserRepo,
	smtpCfg email.SMTPConfig,
) *CandidateOfferService {
	return &CandidateOfferService{
		repo:         repo,
		employeeRepo: employeeRepo,
		userRepo:     userRepo,
		smtpCfg:      smtpCfg,
	}
}

func (s *CandidateOfferService) CreateCandidate(
	ctx context.Context,
	arg db.CreateCandidateParams,
) (db.Candidate, error) {
	existing, err := s.repo.GetCandidateByEmail(ctx, arg.Email)
	if err == nil && existing.ID != uuid.Nil {
		return db.Candidate{}, ErrEmailTaken
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return db.Candidate{}, err
	}

	candidate, err := s.repo.CreateCandidate(ctx, arg)
	if err != nil {
		return db.Candidate{}, err
	}

	go email.SendEmail(s.smtpCfg, email.EmailParams{
		To:      candidate.Email,
		Subject: "Application Received",
		Body:    "Your application has been received.",
	})

	return candidate, nil
}

func (s *CandidateOfferService) CreateOffer(
	ctx context.Context,
	arg db.CreateOfferParams,
) (db.Offer, error) {
	return s.repo.CreateOffer(ctx, arg)
}

func (s *CandidateOfferService) AcceptOffer(ctx context.Context, offerID uuid.UUID) error {
	offer, err := s.repo.GetOfferByID(ctx, offerID)
	if err != nil {
		return err
	}
	candidateID := uuid.UUID(offer.CandidateID.Bytes)
	candidate, err := s.repo.GetCandidateByID(ctx, candidateID)
	if err != nil {
		return err
	}
	plainPassword, err := GeneratePassword(12)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user, err := s.userRepo.CreateUser(ctx, db.CreateUserParams{
		Email:        candidate.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    candidate.FirstName,
		LastName:     candidate.LastName,
		UserType:     "employee",
		Status:       pgtype.Text{String: "active", Valid: true},
	})
	if err != nil {
		return err
	}
	employeeNumber := fmt.Sprintf("EMP%s", uuid.New().String()[:6])
	_, err = s.employeeRepo.CreateEmployee(ctx, db.CreateEmployeeParams{
		UserID:           pgtype.UUID{Bytes: user.ID, Valid: true},
		EmployeeNumber:   employeeNumber,
		DepartmentID:     offer.DepartmentID,
		DesignationID:    offer.DesignationID,
		LevelID:          offer.LevelID,
		EmploymentStatus: pgtype.Text{String: "active", Valid: true},
		DateOfEmployment: pgtype.Date{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return err
	}
	_, err = s.repo.UpdateOfferStatus(ctx, db.UpdateOfferStatusParams{
		ID:     offerID,
		Status: pgtype.Text{String: "accepted", Valid: true},
	})
	if err != nil {
		return err
	}
	_, err = s.repo.UpdateCandidateStatus(ctx, db.UpdateCandidateStatusParams{
		ID:     candidate.ID,
		Status: pgtype.Text{String: "accepted", Valid: true},
	})
	if err != nil {
		return err
	}
	go email.SendEmail(s.smtpCfg, email.EmailParams{
		To:      candidate.Email,
		Subject: "Welcome to Workaround",
		Body: fmt.Sprintf(
			"Congratulations! Your offer has been accepted. Your temporary password is: <b>%s</b>. Please login and change it.",
			plainPassword,
		),
	})

	return nil
}

func (s *CandidateOfferService) RejectOffer(
	ctx context.Context,
	candidateID uuid.UUID,
) (db.Candidate, error) {
	return s.repo.UpdateCandidateStatus(ctx, db.UpdateCandidateStatusParams{
		ID:     candidateID,
		Status: pgtype.Text{String: "rejected", Valid: true},
	})
}

func (s *CandidateOfferService) GetCandidateByID(
	ctx context.Context,
	id uuid.UUID,
) (db.Candidate, error) {
	return s.repo.GetCandidateByID(ctx, id)
}

func (s *CandidateOfferService) ListCandidates(
	ctx context.Context,
	arg db.ListCandidatesParams,
) ([]db.Candidate, error) {
	return s.repo.ListCandidates(ctx, arg)
}

func (s *CandidateOfferService) ListOffers(
	ctx context.Context,
	arg db.ListOffersParams,
) ([]db.Offer, error) {
	return s.repo.ListOffers(ctx, arg)
}

func (s *CandidateOfferService) GetOffersByCandidate(
	ctx context.Context,
	candidateID uuid.UUID,
) ([]db.Offer, error) {
	return s.repo.GetOfferByCandidate(ctx, pgtype.UUID{
		Bytes: candidateID,
		Valid: true,
	})
}

func (s *CandidateOfferService) DeleteCandidate(
	ctx context.Context,
	id uuid.UUID,
) error {
	_, err := s.repo.GetCandidateByID(ctx, id)
	if err != nil {
		return err
	}
	offers, err := s.repo.GetOfferByCandidate(ctx, pgtype.UUID{
		Bytes: id,
		Valid: true,
	})
	if err != nil {
		return err
	}
	if len(offers) > 0 {
		return errors.New("cannot delete candidate with existing offers")
	}

	return s.repo.DeleteCandidate(ctx, id)
}

func (s *CandidateOfferService) UpdateCandidate(
	ctx context.Context,
	arg db.UpdateCandidateParams,
) (db.Candidate, error) {
	return s.repo.UpdateCandidate(ctx, arg)
}
