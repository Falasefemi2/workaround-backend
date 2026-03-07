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

type DeptHandler struct {
	service *service.DeptService
}

func NewDeptHandler(service *service.DeptService) *DeptHandler {
	return &DeptHandler{
		service: service,
	}
}

type CreateDeptRequest struct {
	Name  string `json:"name"   validate:"required"`
	Code  string `json:"code"   validate:"required"`
	HodID string `json:"hod_id"`
}

type UpdateDepartmentRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type AssignHodRequest struct {
	HodID string `json:"hod_id" validate:"required"`
}

func (h *DeptHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)

		r.Post("/v1/departments", h.CreateDepartment)
		r.Get("/v1/departments", h.ListDepartments)
		r.Get("/v1/departments/{id}", h.GetDepartmentByID)
		r.Put("/v1/departments/{id}", h.UpdateDepartment)
		r.Delete("/v1/departments/{id}", h.DeleteDepartment)
		r.Put("/v1/departments/{id}/hod", h.AssignHodToDepartment)
	})
}

// CreateDepartment godoc
// @Summary Create a new department
// @Description Creates a new department in the system
// @Tags Departments
// @Accept json
// @Produce json
// @Param request body CreateDeptRequest true "Department payload"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/departments [post]
func (h *DeptHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req CreateDeptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	hodID := pgtype.UUID{}
	if req.HodID != "" {
		var ok bool
		hodID, ok = scanUUIDOrError(w, req.HodID, "invalid hod_id format")
		if !ok {
			return
		}
	}

	dept, err := h.service.CreateDepartment(r.Context(), db.CreateDepartmentParams{
		Name:  req.Name,
		Code:  req.Code,
		HodID: hodID,
	})
	if err != nil {
		if errors.Is(err, service.ErrDeptName) {
			response.Error(w, http.StatusConflict, "name already in use")
			return
		}
		if errors.Is(err, service.ErrCode) {
			response.Error(w, http.StatusConflict, "code already in use")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create department")
		return
	}
	response.JSON(w, http.StatusCreated, dept)
}

// DeleteDepartment godoc
// @Summary Delete department
// @Description Deletes a department by id
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 204 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/departments/{id} [delete]
func (h *DeptHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		response.Error(w, http.StatusBadRequest, "department id is required")
		return
	}

	deptID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid department id")
		return
	}

	if err := h.service.DeleteDepartment(r.Context(), deptID); err != nil {
		if errors.Is(err, service.ErrDeptNotFound) {
			response.Error(w, http.StatusNotFound, "department not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not delete department")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetDepartmentByID godoc
// @Summary Get department by id
// @Description Retrieves a department by id
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/departments/{id} [get]
func (h *DeptHandler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		response.Error(w, http.StatusBadRequest, "department id is required")
		return
	}

	deptID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid department id")
		return
	}

	dept, err := h.service.GetDepartmentByID(r.Context(), deptID)
	if err != nil {
		if errors.Is(err, service.ErrDeptNotFound) {
			response.Error(w, http.StatusNotFound, "department not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not fetch department")
		return
	}

	response.JSON(w, http.StatusOK, dept)
}

// ListDepartments godoc
// @Summary List departments
// @Description Returns a paginated list of departments
// @Tags Departments
// @Accept json
// @Produce json
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/departments [get]
func (h *DeptHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(10) // default
	offset := int32(0) // default

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

	depts, err := h.service.ListDepartments(r.Context(), db.ListDepartmentsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not list departments")
		return
	}

	response.JSON(w, http.StatusOK, depts)
}

// UpdateDepartment godoc
// @Summary Update department
// @Description Updates an existing department
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Param request body UpdateDepartmentRequest true "Department payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/departments/{id} [put]
func (h *DeptHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		response.Error(w, http.StatusBadRequest, "department id is required")
		return
	}

	deptID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid department id")
		return
	}

	var req UpdateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	dept, err := h.service.UpdateDepartment(r.Context(), db.UpdateDepartmentParams{
		ID:   deptID,
		Name: req.Name,
		Code: req.Code,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not update department")
		return
	}

	response.JSON(w, http.StatusOK, dept)
}

// AssignHodToDepartment godoc
// @Summary Assign HOD to department
// @Description Assigns a head of department to a department
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Param request body AssignHodRequest true "HOD payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/departments/{id}/hod [put]
func (h *DeptHandler) AssignHodToDepartment(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		response.Error(w, http.StatusBadRequest, "department id is required")
		return
	}

	deptID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid department id")
		return
	}

	var req AssignHodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	hodID, ok := scanUUIDOrError(w, req.HodID, "invalid hod_id format")
	if !ok {
		return
	}

	dept, err := h.service.AssignHodToDepartment(r.Context(), db.AssignHodToDepartmentParams{
		ID:    deptID,
		HodID: hodID,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidHod) {
			response.Error(
				w,
				http.StatusBadRequest,
				"user must be an employee to be assigned as HOD",
			)
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not assign HOD")
		return
	}

	response.JSON(w, http.StatusOK, dept)
}
