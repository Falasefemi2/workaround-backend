package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type ApprovalRepo struct {
	queries *db.Queries
}

func NewApprovalRepo(queries *db.Queries) *ApprovalRepo {
	return &ApprovalRepo{
		queries: queries,
	}
}

func (r *ApprovalRepo) ActOnApproval(
	ctx context.Context,
	arg db.ActOnApprovalParams,
) (db.Approval, error) {
	return r.queries.ActOnApproval(ctx, arg)
}

func (r *ApprovalRepo) CreateApproval(
	ctx context.Context,
	arg db.CreateApprovalParams,
) (db.Approval, error) {
	return r.queries.CreateApproval(ctx, arg)
}

func (r *ApprovalRepo) CreateApprovalSetup(
	ctx context.Context,
	arg db.CreateApprovalSetupParams,
) (db.ApprovalSetup, error) {
	return r.queries.CreateApprovalSetup(ctx, arg)
}

func (r *ApprovalRepo) DeleteApprovalSetup(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteApprovalSetup(ctx, id)
}

func (r *ApprovalRepo) GetApprovalChain(
	ctx context.Context,
	moduleType pgtype.Text,
) ([]db.ApprovalSetup, error) {
	return r.queries.GetApprovalChain(ctx, moduleType)
}

func (r *ApprovalRepo) GetApprovalByReference(
	ctx context.Context,
	referenceID uuid.UUID,
) ([]db.Approval, error) {
	return r.queries.GetApprovalsByReference(ctx, referenceID)
}

func (r *ApprovalRepo) GetPendingApproval(
	ctx context.Context,
	referenceID uuid.UUID,
) (db.Approval, error) {
	return r.queries.GetPendingApproval(ctx, referenceID)
}

func (r *ApprovalRepo) GetPendingApprovalByApprover(
	ctx context.Context,
	approvalID pgtype.UUID,
) ([]db.Approval, error) {
	return r.queries.GetPendingApprovalsByApprover(ctx, approvalID)
}
