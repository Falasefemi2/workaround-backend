package repository

import (
	"context"

	"github.com/google/uuid"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type DesignationRepo struct {
	queries *db.Queries
}

func NewDesignationRepo(queries *db.Queries) *DesignationRepo {
	return &DesignationRepo{
		queries: queries,
	}
}

func (r *DesignationRepo) CreateDesignation(
	ctx context.Context,
	params db.CreateDesignationParams,
) (db.Designation, error) {
	des, err := r.queries.CreateDesignation(ctx, params)
	if err != nil {
		return db.Designation{}, err
	}
	return des, nil
}

func (r *DesignationRepo) DeleteDesignation(
	ctx context.Context,
	id uuid.UUID,
) error {
	return r.queries.DeleteDesignation(ctx, id)
}

func (r *DesignationRepo) ListDesignations(
	ctx context.Context,
	params db.ListDesignationsParams,
) ([]db.Designation, error) {
	des, err := r.queries.ListDesignations(ctx, params)
	if err != nil {
		return nil, err
	}
	return des, nil
}

func (r *DesignationRepo) UpdateDesignation(
	ctx context.Context,
	params db.UpdateDesignationParams,
) (db.Designation, error) {
	des, err := r.queries.UpdateDesignation(ctx, params)
	if err != nil {
		return db.Designation{}, err
	}
	return des, nil
}

func (r *DesignationRepo) GetDesignationByID(
	ctx context.Context,
	id uuid.UUID,
) (db.Designation, error) {
	des, err := r.queries.GetDesignationByID(ctx, id)
	if err != nil {
		return db.Designation{}, err
	}
	return des, nil
}
