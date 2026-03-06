package repository

import (
	"context"

	"github.com/google/uuid"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type LevelRepo struct {
	queries *db.Queries
}

func NewLevelRepo(queries *db.Queries) *LevelRepo {
	return &LevelRepo{
		queries: queries,
	}
}

func (r *LevelRepo) CreateLevel(
	ctx context.Context,
	params db.CreateLevelParams,
) (db.Level, error) {
	return r.queries.CreateLevel(ctx, params)
}

func (r *LevelRepo) DeleteLevel(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteLevel(ctx, id)
}

func (r *LevelRepo) GetLevelByCode(ctx context.Context, code string) (db.Level, error) {
	return r.queries.GetLevelByCode(ctx, code)
}

func (r *LevelRepo) GetLevelByID(ctx context.Context, id uuid.UUID) (db.Level, error) {
	return r.queries.GetLevelByID(ctx, id)
}

func (r *LevelRepo) ListLevel(ctx context.Context, arg db.ListLevelsParams) ([]db.Level, error) {
	return r.queries.ListLevels(ctx, arg)
}

func (r *LevelRepo) SearchLevel(
	ctx context.Context,
	arg db.SearchLevelsParams,
) ([]db.Level, error) {
	return r.queries.SearchLevels(ctx, arg)
}

func (r *LevelRepo) UpdateLevel(ctx context.Context, arg db.UpdateLevelParams) (db.Level, error) {
	return r.queries.UpdateLevel(ctx, arg)
}
