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
	"github.com/falasefemi2/workaround-backend/internal/middleware"
	"github.com/falasefemi2/workaround-backend/internal/response"
	"github.com/falasefemi2/workaround-backend/internal/service"
	appvalidator "github.com/falasefemi2/workaround-backend/internal/validator"
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

type authRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type forgetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token"        validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type AuthRequest = authRequest
type ForgetPasswordRequest = forgetPasswordRequest
type ResetPasswordRequest = resetPasswordRequest

func (h *UserHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
	rateLimit func(http.Handler) http.Handler,
) {
	r.Post("/v1/auth/forgot-password", h.ForgotPassword)
	r.Post("/v1/auth/reset-password", h.ResetPassword)
	r.Post("/v1/auth/logout", h.Logout)

	r.Group(func(r chi.Router) {
		r.Use(rateLimit)
		r.Post("/v1/auth/login", h.Login)
	})

	r.Post("/v1/auth/logout", h.Logout)

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)

		r.Post("/v1/auth/users", h.CreateUser)

		r.Get("/v1/users", h.ListUsers)
		r.Get("/v1/users/email/{email}", h.GetUserByEmail)
		r.Get("/v1/users/{id}", h.GetUserByID)

		r.Put("/v1/users/{id}", h.UpdateUser)
		r.Delete("/v1/users/{id}", h.DeleteUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Get("/v1/users/me", h.Me)
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User payload"
// @Success 201 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/auth/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := appvalidator.Validate(req); errs != nil {
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

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes a user by id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		response.Error(w, http.StatusBadRequest, "user id is required")
		return
	}

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.service.DeleteUser(r.Context(), userID); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserByEmail godoc
// @Summary Get user by email
// @Description Retrieves a user by email address
// @Tags Users
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/users/email/{email} [get]
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		response.Error(w, http.StatusBadRequest, "email is required")
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not fetch user")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// GetUserByID godoc
// @Summary Get user by id
// @Description Retrieves a user by id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		response.Error(w, http.StatusBadRequest, "user id is required")
		return
	}

	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "could not fetch user")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// ListUsers godoc
// @Summary List users
// @Description Returns a paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param limit query int false "Page size"
// @Param offset query int false "Page offset"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.service.ListUsers(r.Context(), db.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not list users")
		return
	}

	response.JSON(w, http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Update user
// @Description Updates the authenticated user profile
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body db.UpdateUserParams true "User payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userIDStr == "" {
		response.Error(w, http.StatusBadRequest, "invalid user ID in context")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	var req db.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.ID = userID

	user, err := h.service.UpdateUser(r.Context(), req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not update user")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// Login godoc
// @Summary Login user
// @Description Authenticates user and sets auth cookie
// @Tags Users
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Login payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/auth/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCreds) {
			response.Error(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		response.Error(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	setAuthCookie(w, token)
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged in successfully"})
}

// Logout godoc
// @Summary Logout user
// @Description Clears auth cookie and logs out user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/auth/logout [post]
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1, // delete the cookie immediately
	})
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// setAuthCookie sets the JWT as an httpOnly cookie
func setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true, // not accessible via JavaScript
		Path:     "/",
		MaxAge:   86400, // 24 hours in seconds
		SameSite: http.SameSiteLaxMode,
		// Secure: true  // uncomment in production (requires HTTPS)
	})
}

// Me godoc
// @Summary Get authenticated user
// @Description Returns authenticated user id from token context
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /v1/users/me [get]
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"user_id": userID,
	})
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Sends password reset instructions if email exists
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ForgetPasswordRequest true "Forgot password payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/auth/forgot-password [post]
func (h *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	err := h.service.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		map[string]string{"message": "if your email exists you will receive a reset link"},
	)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Resets password using reset token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset password payload"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/auth/reset-password [post]
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	err := h.service.ResetPassword(r.Context(), req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			response.Error(w, http.StatusBadRequest, "invalid or expired token")
			return
		}
		response.Error(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	response.JSON(
		w,
		http.StatusOK,
		map[string]string{"message": "password reset successfully"},
	)
}
