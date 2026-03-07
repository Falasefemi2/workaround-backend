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

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler(service *service.RoleService) *RoleHandler {
	return &RoleHandler{
		service: service,
	}
}

type CreateRoleRequest struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateeRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AssignRoleRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
	RoleID uuid.UUID `json:"role_id" validate:"required"`
}

type RemoveRoleRequest struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
}

func (h *RoleHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)
		r.Post("/v1/roles", h.CreateRole)
		r.Get("/v1/roles", h.ListRoles)
		r.Get("/v1/roles/{id}", h.GetRoleByID)
		r.Put("/v1/roles/{id}", h.UpdateRole)
		r.Delete("/v1/roles/{id}", h.DeleteRole)
		r.Post("/v1/roles/assign", h.AssignRoleToUser)
		r.Delete("/v1/roles/remove", h.RemoveRoleFromUser)
		r.Get("/v1/users/{id}/roles", h.GetUserRoles)
	})
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := appvalidator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}
	params := db.CreateRoleParams{
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	}
	role, err := h.service.CreateRole(r.Context(), params)
	if err != nil {
		if errors.Is(err, service.ErrRoleNameTaken) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create role")
		return
	}
	response.JSON(w, http.StatusCreated, role)
}

func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	offset, _ := strconv.ParseInt(offsetStr, 10, 32)
	params := db.ListRolesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	roles, err := h.service.ListRoles(r.Context(), params)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list roles")
		return
	}
	response.JSON(w, http.StatusOK, roles)
}

func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	role, err := h.service.GetRoleByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusNotFound, "role not found")
		return
	}
	response.JSON(w, http.StatusOK, role)
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req UpdateeRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	params := db.UpdateRoleParams{
		ID:          id,
		Name:        req.Name,
		Description: pgtype.Text{String: req.Description, Valid: req.Description != ""},
	}
	role, err := h.service.UpdateRole(r.Context(), params)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update role")
		return
	}
	response.JSON(w, http.StatusOK, role)
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	err = h.service.DeleteRole(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to delete role")
		return
	}
	response.JSON(w, http.StatusOK, "role deleted")
}

func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req AssignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := appvalidator.Validate(req); err != nil {
		response.ValidationError(w, err)
		return
	}
	var userID, roleID pgtype.UUID
	userID.Scan(req.UserID)
	roleID.Scan(req.RoleID)
	params := db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	}
	userRole, err := h.service.AssignRoleToUser(r.Context(), params)
	if err != nil {
		if errors.Is(err, service.ErrRoleGiven) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to assign role")
		return
	}
	response.JSON(w, http.StatusCreated, userRole)
}

func (h *RoleHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	var req RemoveRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	var userID, roleID pgtype.UUID
	userID.Scan(req.UserID)
	roleID.Scan(req.RoleID)
	params := db.RemoveRoleFromUserParams{
		UserID: userID,
		RoleID: roleID,
	}
	err := h.service.RemoveRoleFromUser(r.Context(), params)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to remove role")
		return
	}
	response.JSON(w, http.StatusOK, "role removed")
}

func (h *RoleHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}
	userID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
	roles, err := h.service.GetUserRoles(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get roles")
		return
	}
	response.JSON(w, http.StatusOK, roles)
}
