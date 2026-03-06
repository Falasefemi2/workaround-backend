package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/repository"
)

var ErrLevelCodeTaken = errors.New("code taken already")

type LevelService struct {
	repo *repository.LevelRepo
}

func NewLevelService(repo *repository.LevelRepo) *LevelService {
	return &LevelService{
		repo: repo,
	}
}

func (s *LevelService) CreateLevel(
	ctx context.Context,
	params db.CreateLevelParams,
) (db.Level, error) {
	level, err := s.repo.CreateLevel(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return db.Level{}, ErrLevelCodeTaken
		}
		return db.Level{}, err
	}
	return level, nil
}

func (s *LevelService) DeleteLevel(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteLevel(ctx, id)
}

func (s *LevelService) GetLevelByID(ctx context.Context, id uuid.UUID) (db.Level, error) {
	level, err := s.repo.GetLevelByID(ctx, id)
	if err != nil {
		return db.Level{}, err
	}
	return level, nil
}

func (s *LevelService) ListLevels(
	ctx context.Context,
	arg db.ListLevelsParams,
) ([]db.Level, error) {
	levels, err := s.repo.ListLevel(ctx, arg)
	if err != nil {
		return nil, err
	}
	return levels, nil
}

func (s *LevelService) SearchLevels(
	ctx context.Context,
	arg db.SearchLevelsParams,
) ([]db.Level, error) {
	levels, err := s.repo.SearchLevel(ctx, arg)
	if err != nil {
		return nil, err
	}
	return levels, nil
}

func (s *LevelService) UpdateLevel(
	ctx context.Context,
	arg db.UpdateLevelParams,
) (db.Level, error) {
	level, err := s.repo.UpdateLevel(ctx, arg)
	if err != nil {
		return db.Level{}, err
	}
	return level, nil
}
