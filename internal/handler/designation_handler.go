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

type DesignationHandler struct {
	service *service.DesignationService
}

func NewDesignationHandler(service *service.DesignationService) *DesignationHandler {
	return &DesignationHandler{
		service: service,
	}
}

type CreateDesignationRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

type UpdateDesignationRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (h *DesignationHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)

		r.Post("/v1/designations", h.CreateDesignation)
		r.Get("/v1/designations", h.ListDesignations)
		r.Get("/v1/designations/{id}", h.GetDesignationByID)
		r.Put("/v1/designations/{id}", h.UpdateDesignation)
		r.Delete("/v1/designations/{id}", h.DeleteDesignation)
	})
}

func (h *DesignationHandler) CreateDesignation(w http.ResponseWriter, r *http.Request) {
	var req CreateDesignationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	designation, err := h.service.CreateDesignation(
		r.Context(),
		db.CreateDesignationParams{
			Name: req.Name,
			Code: pgtype.Text{String: req.Code, Valid: true},
		},
	)
	if err != nil {

		if errors.Is(err, service.ErrDesignationName) {
			response.Error(w, http.StatusConflict, "designation already exists")
			return
		}
		if errors.Is(err, service.ErrDesignationCode) {
			response.Error(w, http.StatusConflict, "designation code already exists")
			return
		}

		response.Error(w, http.StatusInternalServerError, "failed to create designation")
		return
	}

	response.JSON(w, http.StatusCreated, designation)
}

func (h *DesignationHandler) DeleteDesignation(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	designationID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid designation id")
		return
	}

	err = h.service.DeleteDesignation(r.Context(), designationID)
	if err != nil {

		if errors.Is(err, service.ErrDesignationNotFound) {
			response.Error(w, http.StatusNotFound, "designation not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "could not delete designation")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DesignationHandler) GetDesignationByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	designationID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid designation id")
		return
	}

	designation, err := h.service.GetDesignationByID(r.Context(), designationID)
	if err != nil {

		if errors.Is(err, service.ErrDesignationNotFound) {
			response.Error(w, http.StatusNotFound, "designation not found")
			return
		}

		response.Error(w, http.StatusInternalServerError, "could not fetch designation")
		return
	}

	response.JSON(w, http.StatusOK, designation)
}

func (h *DesignationHandler) ListDesignations(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(10)
	offset := int32(0)

	if limitStr != "" {

		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid limit")
			return
		}

		limit = int32(parsedLimit)
	}

	if offsetStr != "" {

		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid offset")
			return
		}

		offset = int32(parsedOffset)
	}

	designations, err := h.service.ListDesignations(
		r.Context(),
		db.ListDesignationsParams{
			Limit:  limit,
			Offset: offset,
		},
	)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not list designations")
		return
	}

	response.JSON(w, http.StatusOK, designations)
}

func (h *DesignationHandler) UpdateDesignation(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	designationID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid designation id")
		return
	}

	var req UpdateDesignationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	designation, err := h.service.UpdateDesignation(
		r.Context(),
		db.UpdateDesignationParams{
			ID:   designationID,
			Name: req.Name,
			Code: pgtype.Text{String: req.Code, Valid: true},
		},
	)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not update designation")
		return
	}

	response.JSON(w, http.StatusOK, designation)
}
