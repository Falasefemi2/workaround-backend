package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/falasefemi2/workaround-backend/internal/repository"
	"github.com/falasefemi2/workaround-backend/internal/response"
)

type AdminDashboardHandler struct {
	employeeRepo *repository.EmployeeRepo
	offerRepo    *repository.CandidateOfferRepo
}

func NewAdminDashboardHandler(
	employeeRepo *repository.EmployeeRepo,
	offerRepo *repository.CandidateOfferRepo,
) *AdminDashboardHandler {
	return &AdminDashboardHandler{
		employeeRepo: employeeRepo,
		offerRepo:    offerRepo,
	}
}

func (h *AdminDashboardHandler) GetOfferStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.offerRepo.GetOfferStatsByMonth(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not fetch offer stats")
		return
	}
	response.JSON(w, http.StatusOK, stats)
}

func (h *AdminDashboardHandler) GetEmployeeCountByDepartment(
	w http.ResponseWriter,
	r *http.Request,
) {
	counts, err := h.employeeRepo.GetEmployeeCountByDepartment(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not fetch employee counts")
		return
	}
	response.JSON(w, http.StatusOK, counts)
}

func (h *AdminDashboardHandler) RegisterRoutes(
	r chi.Router,
	auth func(http.Handler) http.Handler,
	adminOrHR func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Use(adminOrHR)
		r.Get("/v1/dashboard/offer-stats", h.GetOfferStats)
		r.Get("/v1/dashboard/employee-count", h.GetEmployeeCountByDepartment)
	})
}
