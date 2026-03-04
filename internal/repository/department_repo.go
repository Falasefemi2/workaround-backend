package repository

import (
	"context"

	"github.com/google/uuid"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type DepartmentRepo struct {
	queries *db.Queries
}

func NewDepartmentRepo(queries *db.Queries) *DepartmentRepo {
	return &DepartmentRepo{
		queries: queries,
	}
}

func (r *DepartmentRepo) CreateDepartment(
	ctx context.Context,
	params db.CreateDepartmentParams,
) (db.Department, error) {
	dept, err := r.queries.CreateDepartment(ctx, params)
	if err != nil {
		return db.Department{}, err
	}
	return dept, nil
}

func (r *DepartmentRepo) DeleteDepartment(ctx context.Context, deptID uuid.UUID) error {
	return r.queries.DeleteDepartment(ctx, deptID)
}

func (r *DepartmentRepo) ListDepartment(
	ctx context.Context,
	params db.ListDepartmentsParams,
) ([]db.Department, error) {
	dept, err := r.queries.ListDepartments(ctx, params)
	if err != nil {
		return nil, err
	}
	return dept, nil
}

func (r *DepartmentRepo) UpdateDepartment(
	ctx context.Context,
	params db.UpdateDepartmentParams,
) (db.Department, error) {
	dept, err := r.queries.UpdateDepartment(ctx, params)
	if err != nil {
		return db.Department{}, err
	}
	return dept, nil
}

func (r *DepartmentRepo) GetDepartmentByID(
	ctx context.Context,
	id uuid.UUID,
) (db.Department, error) {
	dept, err := r.queries.GetDepartmentByID(ctx, id)
	if err != nil {
		return db.Department{}, err
	}
	return dept, nil
}

func (r *DepartmentRepo) AssignHodToDepartment(
	ctx context.Context,
	params db.AssignHodToDepartmentParams,
) (db.Department, error) {
	dept, err := r.queries.AssignHodToDepartment(ctx, params)
	if err != nil {
		return db.Department{}, err
	}
	return dept, nil
}
