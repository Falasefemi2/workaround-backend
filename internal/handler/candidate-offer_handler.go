package handler

import (
	"encoding/json"
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

type CandidateOfferHandler struct {
	service *service.CandidateOfferService
}

func NewCandidateOfferHandler(service *service.CandidateOfferService) *CandidateOfferHandler {
	return &CandidateOfferHandler{
		service: service,
	}
}

func (h *CandidateOfferHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)

		r.Post("/v1/candidates", h.CreateCandidate)
		r.Get("/v1/candidates", h.ListCandidates)
		r.Get("/v1/candidates/{id}", h.GetCandidateByID)
		r.Put("/v1/candidates/{id}", h.UpdateCandidate)
		r.Delete("/v1/candidates/{id}", h.DeleteCandidate)

		r.Post("/v1/offers", h.CreateOffer)
		r.Get("/v1/offers", h.ListOffers)
		r.Get("/v1/candidates/{id}/offers", h.GetOffersByCandidate)

		r.Put("/v1/offers/{id}/accept", h.AcceptOffer)
		r.Put("/v1/offers/{id}/reject", h.RejectOffer)
	})
}

type CreateCandidateRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Email     string `json:"email"      validate:"required,email"`
	Phone     string `json:"phone"`
}

type UpdateCandidateRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Phone     string `json:"phone"`
}

type CreateOfferRequest struct {
	CandidateID       string `json:"candidate_id"        validate:"required"`
	DepartmentID      string `json:"department_id"       validate:"required"`
	DesignationID     string `json:"designation_id"      validate:"required"`
	LevelID           string `json:"level_id"            validate:"required"`
	ProposedStartDate string `json:"proposed_start_date" validate:"required"`
	NewStartDate      string `json:"new_start_date"`
	OfferLetterUrl    string `json:"offer_letter_url"`
	CreatedBy         string `json:"created_by"          validate:"required"`
}

func (h *CandidateOfferHandler) CreateCandidate(w http.ResponseWriter, r *http.Request) {
	var req CreateCandidateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	candidate, err := h.service.CreateCandidate(r.Context(), db.CreateCandidateParams{
		FirstName: pgtype.Text{String: req.FirstName, Valid: true},
		LastName:  pgtype.Text{String: req.LastName, Valid: true},
		Email:     req.Email,
		Phone:     pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
	})
	if err != nil {

		if err == service.ErrEmailTaken {
			response.Error(w, http.StatusConflict, "email already exists")
			return
		}

		response.Error(w, http.StatusInternalServerError, "failed to create candidate")
		return
	}

	response.JSON(w, http.StatusCreated, candidate)
}

func (h *CandidateOfferHandler) ListCandidates(w http.ResponseWriter, r *http.Request) {
	limit := int32(10)
	offset := int32(0)

	if val := r.URL.Query().Get("limit"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = int32(i)
	}

	if val := r.URL.Query().Get("offset"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = int32(i)
	}

	candidates, err := h.service.ListCandidates(r.Context(), db.ListCandidatesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not list candidates")
		return
	}

	response.JSON(w, http.StatusOK, candidates)
}

func (h *CandidateOfferHandler) GetCandidateByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid candidate id")
		return
	}

	candidate, err := h.service.GetCandidateByID(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "candidate not found")
		return
	}

	response.JSON(w, http.StatusOK, candidate)
}

func (h *CandidateOfferHandler) UpdateCandidate(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid candidate id")
		return
	}

	var req UpdateCandidateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	candidate, err := h.service.UpdateCandidate(r.Context(), db.UpdateCandidateParams{
		ID:        id,
		FirstName: pgtype.Text{String: req.FirstName, Valid: true},
		LastName:  pgtype.Text{String: req.LastName, Valid: true},
		Phone:     pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not update candidate")
		return
	}

	response.JSON(w, http.StatusOK, candidate)
}

func (h *CandidateOfferHandler) DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid candidate id")
		return
	}

	err = h.service.DeleteCandidate(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not delete candidate")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CandidateOfferHandler) CreateOffer(w http.ResponseWriter, r *http.Request) {
	var req CreateOfferRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := appvalidator.Validate(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}

	candidateID, ok := scanUUIDOrError(w, req.CandidateID, "invalid candidate_id")
	if !ok {
		return
	}

	deptID, ok := scanUUIDOrError(w, req.DepartmentID, "invalid department_id")
	if !ok {
		return
	}

	desID, ok := scanUUIDOrError(w, req.DesignationID, "invalid designation_id")
	if !ok {
		return
	}

	levelID, ok := scanUUIDOrError(w, req.LevelID, "invalid level_id")
	if !ok {
		return
	}

	creatorID, ok := scanUUIDOrError(w, req.CreatedBy, "invalid created_by")
	if !ok {
		return
	}

	offer, err := h.service.CreateOffer(r.Context(), db.CreateOfferParams{
		CandidateID:    candidateID,
		DepartmentID:   deptID,
		DesignationID:  desID,
		LevelID:        levelID,
		OfferLetterUrl: pgtype.Text{String: req.OfferLetterUrl, Valid: req.OfferLetterUrl != ""},
		CreatedBy:      creatorID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create offer")
		return
	}

	response.JSON(w, http.StatusCreated, offer)
}

func (h *CandidateOfferHandler) AcceptOffer(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid offer id")
		return
	}

	if err := h.service.AcceptOffer(r.Context(), id); err != nil {
		response.Error(w, http.StatusInternalServerError, "could not accept offer")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "offer accepted",
	})
}

func (h *CandidateOfferHandler) RejectOffer(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid offer id")
		return
	}

	_, err = h.service.RejectOffer(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not reject offer")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "offer rejected",
	})
}

func (h *CandidateOfferHandler) GetOffersByCandidate(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid candidate id")
		return
	}

	offers, err := h.service.GetOffersByCandidate(r.Context(), id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not fetch offers")
		return
	}

	response.JSON(w, http.StatusOK, offers)
}

func (h *CandidateOfferHandler) ListOffers(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)

	offers, err := h.service.ListOffers(r.Context(), db.ListOffersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not list offers")
		return
	}

	response.JSON(w, http.StatusOK, offers)
}
