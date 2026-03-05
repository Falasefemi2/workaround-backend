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
	ErrDeptName     = errors.New("department name already taken")
	ErrCode         = errors.New("department code already taken")
	ErrInvalidHod   = errors.New("user must be an employee to be assigned as HOD")
	ErrDeptNotFound = errors.New("department not found")
)

type DeptService struct {
	repo     *repository.DepartmentRepo
	userRepo *repository.UserRepo
}

func NewDeptService(repo *repository.DepartmentRepo, userRepo *repository.UserRepo) *DeptService {
	return &DeptService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *DeptService) CreateDepartment(
	ctx context.Context,
	params db.CreateDepartmentParams,
) (db.Department, error) {
	dept, err := s.repo.CreateDepartment(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return db.Department{}, ErrCode
		}
		return db.Department{}, err
	}
	return dept, nil
}

func (s *DeptService) DeleteDepartment(ctx context.Context, deptID uuid.UUID) error {
	return s.repo.DeleteDepartment(ctx, deptID)
}

func (s *DeptService) GetDepartmentByID(
	ctx context.Context,
	deptID uuid.UUID,
) (db.Department, error) {
	dept, err := s.repo.GetDepartmentByID(ctx, deptID)
	if err != nil {
		return db.Department{}, err
	}
	return dept, nil
}

func (s *DeptService) ListDepartments(
	ctx context.Context,
	arg db.ListDepartmentsParams,
) ([]db.Department, error) {
	dept, err := s.repo.ListDepartment(ctx, arg)
	if err != nil {
		return nil, err
	}
	return dept, nil
}

func (s *DeptService) UpdateDepartment(
	ctx context.Context,
	params db.UpdateDepartmentParams,
) (db.Department, error) {
	dept, err := s.repo.UpdateDepartment(ctx, params)
	if err != nil {
		return db.Department{}, err
	}
	return dept, nil
}

func (s *DeptService) AssignHodToDepartment(
	ctx context.Context,
	params db.AssignHodToDepartmentParams,
) (db.Department, error) {
	if !params.HodID.Valid {
		return db.Department{}, errors.New("hod_id is required")
	}
	userID := uuid.UUID(params.HodID.Bytes)
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return db.Department{}, ErrUserNotFound
	}
	if user.UserType != "employee" {
		return db.Department{}, ErrInvalidHod
	}
	dept, err := s.repo.AssignHodToDepartment(ctx, params)
	if err != nil {
		return db.Department{}, err
	}

	return dept, nil
}
