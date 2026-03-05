package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/repository"
)

var (
	ErrDesignationName     = errors.New("designation name already taken")
	ErrDesignationCode     = errors.New("designation code already taken")
	ErrDesignationNotFound = errors.New("designation not found")
)

type DesignationService struct {
	repo *repository.DesignationRepo
}

func NewDesignationService(repo *repository.DesignationRepo) *DesignationService {
	return &DesignationService{
		repo: repo,
	}
}

func (s *DesignationService) CreateDesignation(
	ctx context.Context,
	params db.CreateDesignationParams,
) (db.Designation, error) {
	des, err := s.repo.CreateDesignation(ctx, params)
	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "designations_name_key":
				return db.Designation{}, ErrDesignationName
			case "designations_code_key":
				return db.Designation{}, ErrDesignationCode
			default:
				return db.Designation{}, ErrDesignationName
			}
		}

		return db.Designation{}, err
	}

	return des, nil
}

func (s *DesignationService) DeleteDesignation(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repo.DeleteDesignation(ctx, id)
}

func (s *DesignationService) GetDesignationByID(
	ctx context.Context,
	id uuid.UUID,
) (db.Designation, error) {
	des, err := s.repo.GetDesignationByID(ctx, id)
	if err != nil {
		return db.Designation{}, err
	}

	return des, nil
}

func (s *DesignationService) ListDesignations(
	ctx context.Context,
	arg db.ListDesignationsParams,
) ([]db.Designation, error) {
	des, err := s.repo.ListDesignations(ctx, arg)
	if err != nil {
		return nil, err
	}

	return des, nil
}

func (s *DesignationService) UpdateDesignation(
	ctx context.Context,
	params db.UpdateDesignationParams,
) (db.Designation, error) {
	des, err := s.repo.UpdateDesignation(ctx, params)
	if err != nil {
		return db.Designation{}, err
	}

	return des, nil
}
