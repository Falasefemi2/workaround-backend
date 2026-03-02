package repository

import (
	"context"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/google/uuid"
)

type UserRepo struct {
	queries *db.Queries
}

func NewUserRepo(queries *db.Queries) *UserRepo {
	return &UserRepo{
		queries: queries,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return r.queries.DeleteUser(ctx, userID)
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID uuid.UUID) (db.User, error) {
	user, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (r *UserRepo) ListUser(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	user, err := r.queries.ListUsers(ctx, arg)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	user, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
