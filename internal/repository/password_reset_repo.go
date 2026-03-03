package repository

import (
	"context"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type PasswordRepo struct {
	queries *db.Queries
}

func NewPasswordRepo(queries *db.Queries) *PasswordRepo {
	return &PasswordRepo{
		queries: queries,
	}
}

func (p *PasswordRepo) CreatePasswordResetToken(
	ctx context.Context,
	arg db.CreatePasswordResetTokenParams,
) (db.PasswordResetToken, error) {
	result, err := p.queries.CreatePasswordResetToken(ctx, arg)
	if err != nil {
		return db.PasswordResetToken{}, err
	}
	return result, nil
}

func (p *PasswordRepo) GetPasswordResetToken(
	ctx context.Context,
	token string,
) (db.PasswordResetToken, error) {
	result, err := p.queries.GetPasswordResetToken(ctx, token)
	if err != nil {
		return db.PasswordResetToken{}, err
	}
	return result, nil
}

func (p *PasswordRepo) MarkTokenUsed(ctx context.Context, token string) error {
	return p.queries.MarkTokenUsed(ctx, token)
}

func (p *PasswordRepo) DeleteExpiredTokens(ctx context.Context) error {
	return p.queries.DeleteExpiredTokens(ctx)
}
