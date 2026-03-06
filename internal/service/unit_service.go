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
	ErrUnitNameTaken = errors.New("unit name already taken")
	ErrUnitNotFound  = errors.New("unit not found")
	ErrInvalidLead   = errors.New("user must be an employee to be assigned as unit lead")
)

type UnitService struct {
	repo     *repository.UnitRepo
	userRepo *repository.UserRepo
}

func NewUnitService(
	repo *repository.UnitRepo,
	userRepo *repository.UserRepo,
) *UnitService {
	return &UnitService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *UnitService) CreateUnit(
	ctx context.Context,
	params db.CreateUnitParams,
) (db.Unit, error) {
	unit, err := s.repo.CreateUnit(ctx, params)
	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return db.Unit{}, ErrUnitNameTaken
		}

		return db.Unit{}, err
	}

	return unit, nil
}

func (s *UnitService) DeleteUnit(ctx context.Context, unitID uuid.UUID) error {
	return s.repo.DeleteUnit(ctx, unitID)
}

func (s *UnitService) GetUnitByID(
	ctx context.Context,
	unitID uuid.UUID,
) (db.Unit, error) {
	unit, err := s.repo.GetUnitByID(ctx, unitID)
	if err != nil {
		return db.Unit{}, err
	}

	return unit, nil
}

func (s *UnitService) ListUnits(
	ctx context.Context,
	params db.ListUnitsParams,
) ([]db.Unit, error) {
	return s.repo.ListUnits(ctx, params)
}

func (s *UnitService) UpdateUnit(
	ctx context.Context,
	params db.UpdateUnitParams,
) (db.Unit, error) {
	return s.repo.UpdateUnit(ctx, params)
}

func (s *UnitService) AssignUnitLead(
	ctx context.Context,
	params db.AssignUnitLeadParams,
) (db.Unit, error) {
	if !params.UnitLeadID.Valid {
		return db.Unit{}, errors.New("unit_lead_id is required")
	}

	userID := uuid.UUID(params.UnitLeadID.Bytes)

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return db.Unit{}, ErrUserNotFound
	}

	if user.UserType != "employee" {
		return db.Unit{}, ErrInvalidLead
	}

	unit, err := s.repo.AssignUnitLead(ctx, params)
	if err != nil {
		return db.Unit{}, err
	}

	return unit, nil
}
