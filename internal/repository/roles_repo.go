package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type RoleRepo struct {
	queries *db.Queries
}

func NewRoleRepo(queries *db.Queries) *RoleRepo {
	return &RoleRepo{
		queries: queries,
	}
}

func (r *RoleRepo) CreateRole(ctx context.Context, arg db.CreateRoleParams) (db.Role, error) {
	return r.queries.CreateRole(ctx, arg)
}

func (r *RoleRepo) AssignRoleToUser(
	ctx context.Context,
	arg db.AssignRoleToUserParams,
) (db.UserRole, error) {
	return r.queries.AssignRoleToUser(ctx, arg)
}

func (r *RoleRepo) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteRole(ctx, id)
}

func (r *RoleRepo) GetRoleByID(ctx context.Context, id uuid.UUID) (db.Role, error) {
	return r.queries.GetRoleByID(ctx, id)
}

func (r *RoleRepo) GetUserRoles(ctx context.Context, userID pgtype.UUID) ([]db.Role, error) {
	return r.queries.GetUserRoles(ctx, userID)
}

func (r *RoleRepo) ListRoles(ctx context.Context, arg db.ListRolesParams) ([]db.Role, error) {
	return r.queries.ListRoles(ctx, arg)
}

func (r *RoleRepo) RemoveRoleFromUser(ctx context.Context, arg db.RemoveRoleFromUserParams) error {
	return r.queries.RemoveRoleFromUser(ctx, arg)
}

func (r *RoleRepo) UpdateRole(ctx context.Context, arg db.UpdateRoleParams) (db.Role, error) {
	return r.queries.UpdateRole(ctx, arg)
}
