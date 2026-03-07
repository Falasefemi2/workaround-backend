package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/repository"
)

var ErrUnauthorizedApprover = errors.New("you don't have permission to act on this approval")

type ApprovalResult struct {
	Approval      db.Approval
	FullyApproved bool
}

type ApprovalService struct {
	repo     *repository.ApprovalRepo
	roleRepo *repository.RoleRepo
}

func NewApprovalService(
	repo *repository.ApprovalRepo,
	roleRepo *repository.RoleRepo,
) *ApprovalService {
	return &ApprovalService{
		repo:     repo,
		roleRepo: roleRepo,
	}
}

func (s *ApprovalService) SetupApprovalChain(
	ctx context.Context,
	arg db.CreateApprovalSetupParams,
) (db.ApprovalSetup, error) {
	return s.repo.CreateApprovalSetup(ctx, arg)
}

func (s *ApprovalService) GetApprovalChain(
	ctx context.Context,
	moduleType pgtype.Text,
) ([]db.ApprovalSetup, error) {
	return s.repo.GetApprovalChain(ctx, moduleType)
}

func (s *ApprovalService) DeleteApprovalSetup(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteApprovalSetup(ctx, id)
}

func (s *ApprovalService) GetPendingApprovals(
	ctx context.Context,
	approverID pgtype.UUID,
) ([]db.Approval, error) {
	return s.repo.GetPendingApprovalByApprover(ctx, approverID)
}

func (s *ApprovalService) GetApprovalsByReference(
	ctx context.Context,
	referenceID uuid.UUID,
) ([]db.Approval, error) {
	return s.repo.GetApprovalByReference(ctx, referenceID)
}

func (s *ApprovalService) SubmitForApproval(
	ctx context.Context,
	moduleType string,
	referenceID uuid.UUID,
	departmentID uuid.UUID,
) (db.Approval, error) {
	// get the approval chain for this module
	chain, err := s.repo.GetApprovalChain(ctx, pgtype.Text{String: moduleType, Valid: true})
	if err != nil {
		return db.Approval{}, err
	}

	if len(chain) == 0 {
		return db.Approval{}, errors.New("no approval chain configured for this module")
	}

	// get level 1 from chain
	var firstSetup db.ApprovalSetup
	for _, setup := range chain {
		if setup.LevelOrder == 1 {
			firstSetup = setup
			break
		}
	}

	// create approval at level 1
	return s.repo.CreateApproval(ctx, db.CreateApprovalParams{
		ModuleType:    pgtype.Text{String: moduleType, Valid: true},
		ReferenceID:   referenceID,
		ApprovalLevel: pgtype.Int4{Int32: 1, Valid: true},
		ApproverID:    firstSetup.RoleID,
	})
}

func (s *ApprovalService) ActOnApproval(
	ctx context.Context,
	referenceID uuid.UUID,
	actorID uuid.UUID,
	status string,
	comment string,
) (ApprovalResult, error) {
	// 1. get pending approval to know what role is required
	pendingApproval, err := s.repo.GetPendingApproval(ctx, referenceID)
	if err != nil {
		return ApprovalResult{}, err
	}

	// 2. verify actor has the right role
	userRoles, err := s.roleRepo.GetUserRoles(ctx, pgtype.UUID{Bytes: actorID, Valid: true})
	if err != nil {
		return ApprovalResult{}, err
	}

	hasRole := false
	for _, role := range userRoles {
		if pendingApproval.ApproverID.Valid && role.ID == uuid.UUID(pendingApproval.ApproverID.Bytes) {
			hasRole = true
			break
		}
	}
	if !hasRole {
		return ApprovalResult{}, ErrUnauthorizedApprover
	}

	// 3. act on the approval
	approval, err := s.repo.ActOnApproval(ctx, db.ActOnApprovalParams{
		ID:      pendingApproval.ID,
		Status:  pgtype.Text{String: status, Valid: true},
		Comment: pgtype.Text{String: comment, Valid: true},
	})
	if err != nil {
		return ApprovalResult{}, err
	}

	// 4. get the full chain for this module
	chain, err := s.repo.GetApprovalChain(ctx, approval.ModuleType)
	if err != nil {
		return ApprovalResult{}, err
	}

	// 5. if rejected → restart chain at level 1
	if status == "rejected" {
		var firstSetup db.ApprovalSetup
		for _, setup := range chain {
			if setup.LevelOrder == 1 {
				firstSetup = setup
				break
			}
		}

		_, err = s.repo.CreateApproval(ctx, db.CreateApprovalParams{
			ModuleType:    approval.ModuleType,
			ReferenceID:   approval.ReferenceID,
			ApprovalLevel: pgtype.Int4{Int32: 1, Valid: true},
			ApproverID:    firstSetup.RoleID,
		})
		if err != nil {
			return ApprovalResult{}, err
		}

		return ApprovalResult{Approval: approval, FullyApproved: false}, nil
	}

	// 6. if approved → check if next level exists
	currentLevel := approval.ApprovalLevel.Int32
	var nextSetup *db.ApprovalSetup
	for _, setup := range chain {
		if setup.LevelOrder == currentLevel+1 {
			nextSetup = &setup
			break
		}
	}

	// 7. if next level exists → create next approval
	if nextSetup != nil {
		_, err = s.repo.CreateApproval(ctx, db.CreateApprovalParams{
			ModuleType:    approval.ModuleType,
			ReferenceID:   approval.ReferenceID,
			ApprovalLevel: pgtype.Int4{Int32: currentLevel + 1, Valid: true},
			ApproverID:    nextSetup.RoleID,
		})
		if err != nil {
			return ApprovalResult{}, err
		}

		return ApprovalResult{Approval: approval, FullyApproved: false}, nil
	}

	// 8. no next level → fully approved
	return ApprovalResult{Approval: approval, FullyApproved: true}, nil
}
