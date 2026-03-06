package repository

import (
	"context"

	"github.com/google/uuid"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type UnitRepo struct {
	queries *db.Queries
}

func NewUnitRepo(queries *db.Queries) *UnitRepo {
	return &UnitRepo{
		queries: queries,
	}
}

func (r *UnitRepo) CreateUnit(
	ctx context.Context,
	params db.CreateUnitParams,
) (db.Unit, error) {
	return r.queries.CreateUnit(ctx, params)
}

func (r *UnitRepo) AssignUnitLead(
	ctx context.Context,
	params db.AssignUnitLeadParams,
) (db.Unit, error) {
	return r.queries.AssignUnitLead(ctx, params)
}

func (r *UnitRepo) DeleteUnit(ctx context.Context, unitID uuid.UUID) error {
	return r.queries.DeleteUnit(ctx, unitID)
}

func (r *UnitRepo) GetUnitByID(
	ctx context.Context,
	unitID uuid.UUID,
) (db.Unit, error) {
	return r.queries.GetUnitByID(ctx, unitID)
}

func (r *UnitRepo) ListUnits(
	ctx context.Context,
	params db.ListUnitsParams,
) ([]db.Unit, error) {
	return r.queries.ListUnits(ctx, params)
}

func (r *UnitRepo) UpdateUnit(
	ctx context.Context,
	params db.UpdateUnitParams,
) (db.Unit, error) {
	return r.queries.UpdateUnit(ctx, params)
}
