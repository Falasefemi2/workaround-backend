package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/response"
	"github.com/falasefemi2/workaround-backend/internal/service"
	appvalidator "github.com/falasefemi2/workaround-backend/internal/validator"
)

type UnitHandler struct {
	service *service.UnitService
}

func NewUnitHandler(service *service.UnitService) *UnitHandler {
	return &UnitHandler{
		service: service,
	}
}

type CreateUnitRequest struct {
	DepartmentID string `json:"department_id" validate:"required"`
	Name         string `json:"name"          validate:"required"`
	UnitLeadID   string `json:"unit_lead_id"`
}

type UpdateUnitRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *UnitHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)

		r.Post("/v1/units", h.CreateUnit)
		r.Get("/v1/units", h.ListUnits)
		r.Get("/v1/units/{id}", h.GetUnitByID)
		r.Put("/v1/units/{id}", h.UpdateUnit)
		r.Delete("/v1/units/{id}", h.DeleteUnit)
		r.Put("/v1/units/{id}/lead", h.AssignUnitLead)
	})
}

func (h *UnitHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	var req CreateUnitRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	var deptUUID pgtype.UUID
	if err := deptUUID.Scan(req.DepartmentID); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid department_id format")
		return
	}
	var leadID pgtype.UUID

	if req.UnitLeadID != "" {
		if err := leadID.Scan(req.UnitLeadID); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid unit_lead_id")
			return
		}
	}

	unit, err := h.service.CreateUnit(r.Context(), db.CreateUnitParams{
		DepartmentID: deptUUID,
		Name:         req.Name,
		UnitLeadID:   leadID,
	})
	if err != nil {
		if errors.Is(err, service.ErrUnitNameTaken) {
			response.Error(w, http.StatusConflict, "unit name already exists")
			return
		}

		response.Error(w, http.StatusInternalServerError, "failed to create unit")
		return
	}

	response.JSON(w, http.StatusCreated, unit)
}

func (h *UnitHandler) DeleteUnit(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	unitID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid unit id")
		return
	}

	if err := h.service.DeleteUnit(r.Context(), unitID); err != nil {
		response.Error(w, http.StatusInternalServerError, "could not delete unit")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UnitHandler) GetUnitByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	unitID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid unit id")
		return
	}

	unit, err := h.service.GetUnitByID(r.Context(), unitID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not fetch unit")
		return
	}

	response.JSON(w, http.StatusOK, unit)
}

func (h *UnitHandler) ListUnits(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(10)
	offset := int32(0)

	if limitStr != "" {
		val, err := strconv.Atoi(limitStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = int32(val)
	}

	if offsetStr != "" {
		val, err := strconv.Atoi(offsetStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = int32(val)
	}

	units, err := h.service.ListUnits(r.Context(), db.ListUnitsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not list units")
		return
	}

	response.JSON(w, http.StatusOK, units)
}

func (h *UnitHandler) UpdateUnit(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	unitID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid unit id")
		return
	}

	var req UpdateUnitRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	unit, err := h.service.UpdateUnit(r.Context(), db.UpdateUnitParams{
		ID:   unitID,
		Name: req.Name,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not update unit")
		return
	}

	response.JSON(w, http.StatusOK, unit)
}

func (h *UnitHandler) AssignUnitLead(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	unitID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid unit id")
		return
	}

	type assignRequest struct {
		UnitLeadID string `json:"unit_lead_id" validate:"required"`
	}

	var req assignRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	var leadID pgtype.UUID

	if err := leadID.Scan(req.UnitLeadID); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid unit_lead_id")
		return
	}

	unit, err := h.service.AssignUnitLead(r.Context(), db.AssignUnitLeadParams{
		ID:         unitID,
		UnitLeadID: leadID,
	})
	if err != nil {

		if errors.Is(err, service.ErrInvalidLead) {
			response.Error(w, http.StatusBadRequest, "user must be an employee")
			return
		}

		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "could not assign unit lead")
		return
	}

	response.JSON(w, http.StatusOK, unit)
}
