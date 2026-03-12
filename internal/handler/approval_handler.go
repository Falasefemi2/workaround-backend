package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/middleware"
	"github.com/falasefemi2/workaround-backend/internal/response"
	"github.com/falasefemi2/workaround-backend/internal/service"
	appvalidator "github.com/falasefemi2/workaround-backend/internal/validator"
)

type ApprovalHandler struct {
	service *service.ApprovalService
}

func NewApprovalHandler(service *service.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{
		service: service,
	}
}

type CreateApprovalSetupRequest struct {
	ModuleType   string `json:"module_type"   validate:"required"`
	DepartmentID string `json:"department_id"`
	LevelOrder   int32  `json:"level_order"   validate:"required"`
	RoleID       string `json:"role_id"       validate:"required"`
}

type SubmitForApprovalRequest struct {
	ModuleType   string    `json:"module_type"   validate:"required"`
	ReferenceID  uuid.UUID `json:"reference_id"  validate:"required"`
	DepartmentID uuid.UUID `json:"department_id" validate:"required"`
}

type ActOnApprovalRequest struct {
	Status  string `json:"status"  validate:"required,oneof=approved rejected"`
	Comment string `json:"comment"`
}

func (h *ApprovalHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)
		r.Post("/v1/approvals/setup", h.CreateApprovalSetup)
		r.Get("/v1/approvals/setup", h.GetApprovalChain)
		r.Delete("/v1/approvals/setup/{id}", h.DeleteApprovalSetup)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Post("/v1/approvals/submit", h.SubmitForApproval)
		r.Get("/v1/approvals/pending", h.GetPendingApprovals)
		r.Get("/v1/approvals/reference/{reference_id}", h.GetApprovalsByReference)
		r.Post("/v1/approvals/reference/{reference_id}/act", h.ActOnApproval)
	})
}

// CreateApprovalSetup godoc
// @Summary Create approval setup
// @Description Creates an approval chain level for a module and department
// @Tags Approvals
// @Accept json
// @Produce json
// @Param request body CreateApprovalSetupRequest true "Approval setup payload"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/approvals/setup [post]
func (h *ApprovalHandler) CreateApprovalSetup(w http.ResponseWriter, r *http.Request) {
	var req CreateApprovalSetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := appvalidator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	departmentID, ok := scanUUIDOrError(w, req.DepartmentID, "invalid department id")
	if !ok {
		return
	}

	roleID, ok := scanUUIDOrError(w, req.RoleID, "invalid role id")
	if !ok {
		return
	}

	setup, err := h.service.SetupApprovalChain(r.Context(), db.CreateApprovalSetupParams{
		ModuleType:   pgtype.Text{String: req.ModuleType, Valid: req.ModuleType != ""},
		DepartmentID: departmentID,
		LevelOrder:   req.LevelOrder,
		RoleID:       roleID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create approval setup")
		return
	}

	response.JSON(w, http.StatusCreated, setup)
}

// GetApprovalChain godoc
// @Summary Get approval chain
// @Description Retrieves approval chain by module type
// @Tags Approvals
// @Accept json
// @Produce json
// @Param module_type query string true "Module type"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/approvals/setup [get]
func (h *ApprovalHandler) GetApprovalChain(w http.ResponseWriter, r *http.Request) {
	moduleType := r.URL.Query().Get("module_type")
	if moduleType == "" {
		response.Error(w, http.StatusBadRequest, "module_type is required")
		return
	}

	chain, err := h.service.GetApprovalChain(
		r.Context(),
		pgtype.Text{String: moduleType, Valid: true},
	)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get approval chain")
		return
	}

	response.JSON(w, http.StatusOK, chain)
}

// DeleteApprovalSetup godoc
// @Summary Delete approval setup
// @Description Deletes an approval setup level by id
// @Tags Approvals
// @Accept json
// @Produce json
// @Param id path string true "Approval setup ID"
// @Success 204 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/approvals/setup/{id} [delete]
func (h *ApprovalHandler) DeleteApprovalSetup(w http.ResponseWriter, r *http.Request) {
	setupID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid approval setup id")
		return
	}

	if err := h.service.DeleteApprovalSetup(r.Context(), setupID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete approval setup")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SubmitForApproval godoc
// @Summary Submit for approval
// @Description Submits a module record for approval
// @Tags Approvals
// @Accept json
// @Produce json
// @Param request body SubmitForApprovalRequest true "Submit for approval payload"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/approvals/submit [post]
func (h *ApprovalHandler) SubmitForApproval(w http.ResponseWriter, r *http.Request) {
	var req SubmitForApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := appvalidator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	approval, err := h.service.SubmitForApproval(
		r.Context(),
		req.ModuleType,
		req.ReferenceID,
		req.DepartmentID,
	)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to submit for approval")
		return
	}

	response.JSON(w, http.StatusCreated, approval)
}

// GetPendingApprovals godoc
// @Summary Get pending approvals
// @Description Returns pending approvals for the authenticated approver
// @Tags Approvals
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/approvals/pending [get]
func (h *ApprovalHandler) GetPendingApprovals(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userIDStr == "" {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id in token")
		return
	}

	pending, err := h.service.GetPendingApprovals(
		r.Context(),
		pgtype.UUID{Bytes: userID, Valid: true},
	)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get pending approvals")
		return
	}

	response.JSON(w, http.StatusOK, pending)
}

// GetApprovalsByReference godoc
// @Summary Get approvals by reference
// @Description Retrieves approvals for a module reference
// @Tags Approvals
// @Accept json
// @Produce json
// @Param reference_id path string true "Reference ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/approvals/reference/{reference_id} [get]
func (h *ApprovalHandler) GetApprovalsByReference(w http.ResponseWriter, r *http.Request) {
	referenceID, err := uuid.Parse(chi.URLParam(r, "reference_id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid reference id")
		return
	}

	approvals, err := h.service.GetApprovalsByReference(r.Context(), referenceID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get approvals")
		return
	}

	response.JSON(w, http.StatusOK, approvals)
}

// ActOnApproval godoc
// @Summary Act on approval
// @Description Approves or rejects an approval request
// @Tags Approvals
// @Accept json
// @Produce json
// @Param reference_id path string true "Reference ID"
// @Param request body ActOnApprovalRequest true "Approval action payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/approvals/reference/{reference_id}/act [post]
func (h *ApprovalHandler) ActOnApproval(w http.ResponseWriter, r *http.Request) {
	referenceID, err := uuid.Parse(chi.URLParam(r, "reference_id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid reference id")
		return
	}

	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userIDStr == "" {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	actorID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id in token")
		return
	}

	var req ActOnApprovalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := appvalidator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}

	result, err := h.service.ActOnApproval(
		r.Context(),
		referenceID,
		actorID,
		req.Status,
		req.Comment,
	)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorizedApprover) {
			response.Error(w, http.StatusForbidden, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, "failed to act on approval")
		return
	}

	response.JSON(w, http.StatusOK, result)
}
