package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/response"
	"github.com/falasefemi2/workaround-backend/internal/service"
	"github.com/falasefemi2/workaround-backend/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type CreateUserRequest struct {
	Email     string `json:"email"      validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Phone     string `json:"phone"`
	UserType  string `json:"user_type"  validate:"required,oneof=admin hr employee candidate"`
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/v1/auth/users", h.CreateUser)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := validator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	user, err := h.service.CreateUser(r.Context(), db.CreateUserParams{
		Email:     req.Email,
		FirstName: pgtype.Text{String: req.FirstName, Valid: req.FirstName != ""},
		LastName:  pgtype.Text{String: req.LastName, Valid: req.LastName != ""},
		Phone:     pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
		UserType:  req.UserType,
	})
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			response.Error(w, http.StatusConflict, "email already in use")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	response.JSON(w, http.StatusCreated, user)
}
