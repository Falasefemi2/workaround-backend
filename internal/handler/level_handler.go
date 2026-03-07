package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/response"
	"github.com/falasefemi2/workaround-backend/internal/service"
	appvalidator "github.com/falasefemi2/workaround-backend/internal/validator"
)

type LevelHandler struct {
	service *service.LevelService
}

func NewLevelHandler(service *service.LevelService) *LevelHandler {
	return &LevelHandler{
		service: service,
	}
}

type CreateLevelRequest struct {
	Name                    string  `json:"name"                      validate:"required"`
	Code                    string  `json:"code"                      validate:"required"`
	AnnualLeaveDays         int32   `json:"annual_leave_days"         validate:"required"`
	MinimumLeaveDays        int32   `json:"minimum_leave_days"        validate:"required"`
	TotalAnnualLeaveDays    int32   `json:"total_annual_leave_days"   validate:"required"`
	LeaveExpirationInterval int32   `json:"leave_expiration_interval" validate:"required"`
	AnnualGross             float64 `json:"annual_gross"              validate:"required"`
	BasicSalary             float64 `json:"basic_salary"              validate:"required"`
	TransportAllowance      float64 `json:"transport_allowance"       validate:"required"`
	DomesticAllowance       float64 `json:"domestic_allowance"        validate:"required"`
	UtilityAllowance        float64 `json:"utility_allowance"         validate:"required"`
	LunchSubsidy            float64 `json:"lunch_subsidy"             validate:"required"`
	SupportTotal            float64 `json:"support_total"             validate:"required"`
}

type UpdateLevelRequest struct {
	Name                    string  `json:"name"`
	Code                    string  `json:"code"`
	AnnualLeaveDays         int32   `json:"annual_leave_days"`
	MinimumLeaveDays        int32   `json:"minimum_leave_days"`
	TotalAnnualLeaveDays    int32   `json:"total_annual_leave_days"`
	LeaveExpirationInterval int32   `json:"leave_expiration_interval"`
	AnnualGross             float64 `json:"annual_gross"`
	BasicSalary             float64 `json:"basic_salary"`
	TransportAllowance      float64 `json:"transport_allowance"`
	DomesticAllowance       float64 `json:"domestic_allowance"`
	UtilityAllowance        float64 `json:"utility_allowance"`
	LunchSubsidy            float64 `json:"lunch_subsidy"`
	SupportTotal            float64 `json:"support_total"`
}

func (h *LevelHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)
		r.Post("/v1/levels", h.CreateLevel)
		r.Get("/v1/levels", h.ListLevels)
		r.Get("/v1/levels/search", h.SearchLevels)
		r.Get("/v1/levels/{id}", h.GetLevelByID)
		r.Put("/v1/levels/{id}", h.UpdateLevel)
		r.Delete("/v1/levels/{id}", h.DeleteLevel)
	})
}

// CreateLevel godoc
// @Summary Create a new level
// @Description Creates a new level in the system
// @Tags Levels
// @Accept json
// @Produce json
// @Param request body CreateLevelRequest true "Level payload"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/levels [post]
func (h *LevelHandler) CreateLevel(w http.ResponseWriter, r *http.Request) {
	var req CreateLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	annualGross, ok := scanNumericOrError(w, req.AnnualGross, "invalid annual_gross")
	if !ok {
		return
	}
	basicSalary, ok := scanNumericOrError(w, req.BasicSalary, "invalid basic_salary")
	if !ok {
		return
	}
	transportAllowance, ok := scanNumericOrError(w, req.TransportAllowance, "invalid transport_allowance")
	if !ok {
		return
	}
	domesticAllowance, ok := scanNumericOrError(w, req.DomesticAllowance, "invalid domestic_allowance")
	if !ok {
		return
	}
	utilityAllowance, ok := scanNumericOrError(w, req.UtilityAllowance, "invalid utility_allowance")
	if !ok {
		return
	}
	lunchSubsidy, ok := scanNumericOrError(w, req.LunchSubsidy, "invalid lunch_subsidy")
	if !ok {
		return
	}
	supportTotal, ok := scanNumericOrError(w, req.SupportTotal, "invalid support_total")
	if !ok {
		return
	}

	level, err := h.service.CreateLevel(r.Context(), db.CreateLevelParams{
		Name:                    req.Name,
		Code:                    req.Code,
		AnnualLeaveDays:         req.AnnualLeaveDays,
		MinimumLeaveDays:        req.MinimumLeaveDays,
		TotalAnnualLeaveDays:    req.TotalAnnualLeaveDays,
		LeaveExpirationInterval: req.LeaveExpirationInterval,
		AnnualGross:             annualGross,
		BasicSalary:             basicSalary,
		TransportAllowance:      transportAllowance,
		DomesticAllowance:       domesticAllowance,
		UtilityAllowance:        utilityAllowance,
		LunchSubsidy:            lunchSubsidy,
		SupportTotal:            supportTotal,
	})
	if err != nil {
		if errors.Is(err, service.ErrLevelCodeTaken) {
			response.Error(w, http.StatusConflict, "level code already in use")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create level")
		return
	}

	response.JSON(w, http.StatusCreated, level)
}

// GetLevelByID godoc
// @Summary Get level by id
// @Description Retrieves a level by its id
// @Tags Levels
// @Accept json
// @Produce json
// @Param id path string true "Level ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/levels/{id} [get]
func (h *LevelHandler) GetLevelByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid level id")
		return
	}

	level, err := h.service.GetLevelByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "level not found")
		return
	}

	response.JSON(w, http.StatusOK, level)
}

// ListLevels godoc
// @Summary List levels
// @Description Returns a paginated list of levels
// @Tags Levels
// @Accept json
// @Produce json
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/levels [get]
func (h *LevelHandler) ListLevels(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)

	levels, err := h.service.ListLevels(r.Context(), db.ListLevelsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not list levels")
		return
	}

	response.JSON(w, http.StatusOK, levels)
}

// SearchLevels godoc
// @Summary Search levels
// @Description Searches levels by query string
// @Tags Levels
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/levels/search [get]
func (h *LevelHandler) SearchLevels(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.Error(w, http.StatusBadRequest, "search query is required")
		return
	}

	limit, offset := parsePagination(r)

	levels, err := h.service.SearchLevels(r.Context(), db.SearchLevelsParams{
		Column1: pgtype.Text{String: query, Valid: query != ""},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not search levels")
		return
	}

	response.JSON(w, http.StatusOK, levels)
}

// UpdateLevel godoc
// @Summary Update level
// @Description Updates an existing level
// @Tags Levels
// @Accept json
// @Produce json
// @Param id path string true "Level ID"
// @Param request body UpdateLevelRequest true "Level payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/levels/{id} [put]
func (h *LevelHandler) UpdateLevel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid level id")
		return
	}

	var req UpdateLevelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	annualGross, ok := scanNumericOrError(w, req.AnnualGross, "invalid annual_gross")
	if !ok {
		return
	}
	basicSalary, ok := scanNumericOrError(w, req.BasicSalary, "invalid basic_salary")
	if !ok {
		return
	}
	transportAllowance, ok := scanNumericOrError(w, req.TransportAllowance, "invalid transport_allowance")
	if !ok {
		return
	}
	domesticAllowance, ok := scanNumericOrError(w, req.DomesticAllowance, "invalid domestic_allowance")
	if !ok {
		return
	}
	utilityAllowance, ok := scanNumericOrError(w, req.UtilityAllowance, "invalid utility_allowance")
	if !ok {
		return
	}
	lunchSubsidy, ok := scanNumericOrError(w, req.LunchSubsidy, "invalid lunch_subsidy")
	if !ok {
		return
	}
	supportTotal, ok := scanNumericOrError(w, req.SupportTotal, "invalid support_total")
	if !ok {
		return
	}

	level, err := h.service.UpdateLevel(r.Context(), db.UpdateLevelParams{
		ID:                      id,
		Name:                    req.Name,
		Code:                    req.Code,
		AnnualLeaveDays:         req.AnnualLeaveDays,
		MinimumLeaveDays:        req.MinimumLeaveDays,
		TotalAnnualLeaveDays:    req.TotalAnnualLeaveDays,
		LeaveExpirationInterval: req.LeaveExpirationInterval,
		AnnualGross:             annualGross,
		BasicSalary:             basicSalary,
		TransportAllowance:      transportAllowance,
		DomesticAllowance:       domesticAllowance,
		UtilityAllowance:        utilityAllowance,
		LunchSubsidy:            lunchSubsidy,
		SupportTotal:            supportTotal,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not update level")
		return
	}

	response.JSON(w, http.StatusOK, level)
}

// DeleteLevel godoc
// @Summary Delete level
// @Description Deletes a level by id
// @Tags Levels
// @Accept json
// @Produce json
// @Param id path string true "Level ID"
// @Success 204 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/levels/{id} [delete]
func (h *LevelHandler) DeleteLevel(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid level id")
		return
	}

	if err := h.service.DeleteLevel(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "could not delete level")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
