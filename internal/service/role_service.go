package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/repository"
)

var (
	ErrRoleNameTaken = errors.New("role name already taken")
	ErrRoleGiven     = errors.New("role already assigned to user")
)

type RoleService struct {
	repo     *repository.RoleRepo
	userRepo *repository.UserRepo
}

func NewRoleService(repo *repository.RoleRepo, userRepo *repository.UserRepo) *RoleService {
	return &RoleService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *RoleService) CreateRole(ctx context.Context, params db.CreateRoleParams) (db.Role, error) {
	role, err := s.repo.CreateRole(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return db.Role{}, ErrRoleNameTaken
		}
		return db.Role{}, err
	}
	return role, nil
}

func (s *RoleService) AssignRoleToUser(
	ctx context.Context,
	params db.AssignRoleToUserParams,
) (db.UserRole, error) {
	if !params.UserID.Valid {
		return db.UserRole{}, errors.New("user_id is required")
	}

	userID := uuid.UUID(params.UserID.Bytes)
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return db.UserRole{}, ErrUserNotFound
	}

	userRole, err := s.repo.AssignRoleToUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return db.UserRole{}, ErrRoleGiven
		}
		return db.UserRole{}, err
	}

	return userRole, nil
}

func (s *RoleService) RemoveRoleFromUser(
	ctx context.Context,
	arg db.RemoveRoleFromUserParams,
) error {
	return s.repo.RemoveRoleFromUser(ctx, arg)
}

func (s *RoleService) GetUserRoles(ctx context.Context, userID pgtype.UUID) ([]db.Role, error) {
	roles, err := s.repo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) GetRoleByID(ctx context.Context, id uuid.UUID) (db.Role, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return db.Role{}, err
	}
	return role, nil
}

func (s *RoleService) ListRoles(ctx context.Context, arg db.ListRolesParams) ([]db.Role, error) {
	roles, err := s.repo.ListRoles(ctx, arg)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, arg db.UpdateRoleParams) (db.Role, error) {
	role, err := s.repo.UpdateRole(ctx, arg)
	if err != nil {
		return db.Role{}, err
	}
	return role, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteRole(ctx, id)
}
