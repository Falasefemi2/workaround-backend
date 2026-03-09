package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type EmployeeRepo struct {
	queries *db.Queries
}

func NewEmployeeRepo(queries *db.Queries) *EmployeeRepo {
	return &EmployeeRepo{
		queries: queries,
	}
}

func (r *EmployeeRepo) CreateEmployee(
	ctx context.Context,
	arg db.CreateEmployeeParams,
) (db.Employee, error) {
	return r.queries.CreateEmployee(ctx, arg)
}

func (r *EmployeeRepo) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteEmployee(ctx, id)
}

func (r *EmployeeRepo) GetEmployeeByID(ctx context.Context, id uuid.UUID) (db.Employee, error) {
	return r.queries.GetEmployeeByID(ctx, id)
}

func (r *EmployeeRepo) GetEmployeeByNumber(
	ctx context.Context,
	employeeNumber string,
) (db.Employee, error) {
	return r.queries.GetEmployeeByNumber(ctx, employeeNumber)
}

func (r *EmployeeRepo) GetEmployeeByUserID(
	ctx context.Context,
	userID pgtype.UUID,
) (db.Employee, error) {
	return r.queries.GetEmployeeByUserID(ctx, userID)
}

func (r *EmployeeRepo) ListEmployees(
	ctx context.Context,
	arg db.ListEmployeesParams,
) ([]db.Employee, error) {
	return r.queries.ListEmployees(ctx, arg)
}

func (r *EmployeeRepo) UpdateEmployee(
	ctx context.Context,
	arg db.UpdateEmployeeParams,
) (db.Employee, error) {
	return r.queries.UpdateEmployee(ctx, arg)
}
