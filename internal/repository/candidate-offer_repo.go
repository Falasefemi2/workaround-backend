package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
)

type CandidateOfferRepo struct {
	queries *db.Queries
}

func NewCandidateOfferRepo(queries *db.Queries) *CandidateOfferRepo {
	return &CandidateOfferRepo{
		queries: queries,
	}
}

func (r *CandidateOfferRepo) CreateCandidate(
	ctx context.Context,
	arg db.CreateCandidateParams,
) (db.Candidate, error) {
	return r.queries.CreateCandidate(ctx, arg)
}

func (r *CandidateOfferRepo) CreateOffer(
	ctx context.Context,
	arg db.CreateOfferParams,
) (db.Offer, error) {
	return r.queries.CreateOffer(ctx, arg)
}

func (r *CandidateOfferRepo) DeleteCandidate(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteCandidate(ctx, id)
}

func (r *CandidateOfferRepo) GetCandidateByEmail(
	ctx context.Context,
	email string,
) (db.Candidate, error) {
	return r.queries.GetCandidateByEmail(ctx, email)
}

func (r *CandidateOfferRepo) GetCandidateByID(
	ctx context.Context,
	id uuid.UUID,
) (db.Candidate, error) {
	return r.queries.GetCandidateByID(ctx, id)
}

func (r *CandidateOfferRepo) GetOfferByID(ctx context.Context, id uuid.UUID) (db.Offer, error) {
	return r.queries.GetOfferByID(ctx, id)
}

func (r *CandidateOfferRepo) GetOfferByCandidate(
	ctx context.Context,
	candidateID pgtype.UUID,
) ([]db.Offer, error) {
	return r.queries.GetOffersByCandidate(ctx, candidateID)
}

func (r *CandidateOfferRepo) ListCandidates(
	ctx context.Context,
	arg db.ListCandidatesParams,
) ([]db.Candidate, error) {
	return r.queries.ListCandidates(ctx, arg)
}

func (r *CandidateOfferRepo) ListOffers(
	ctx context.Context,
	arg db.ListOffersParams,
) ([]db.Offer, error) {
	return r.queries.ListOffers(ctx, arg)
}

func (r *CandidateOfferRepo) UpdateCandidate(
	ctx context.Context,
	arg db.UpdateCandidateParams,
) (db.Candidate, error) {
	return r.queries.UpdateCandidate(ctx, arg)
}

func (r *CandidateOfferRepo) UpdateCandidateStatus(
	ctx context.Context,
	arg db.UpdateCandidateStatusParams,
) (db.Candidate, error) {
	return r.queries.UpdateCandidateStatus(ctx, arg)
}

func (r *CandidateOfferRepo) UpdateOfferStatus(
	ctx context.Context,
	arg db.UpdateOfferStatusParams,
) (db.Offer, error) {
	return r.queries.UpdateOfferStatus(ctx, arg)
}
