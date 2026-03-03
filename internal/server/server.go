package server

import (
	"net/http"

	db "github.com/falasefemi2/workaround-backend/db/generated"
	"github.com/falasefemi2/workaround-backend/internal/config"
	"github.com/falasefemi2/workaround-backend/internal/email"
	"github.com/falasefemi2/workaround-backend/internal/handler"
	appmw "github.com/falasefemi2/workaround-backend/internal/middleware"
	"github.com/falasefemi2/workaround-backend/internal/repository"
	"github.com/falasefemi2/workaround-backend/internal/service"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.Server.AllowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	queries := db.New(pool)

	userRepo := repository.NewUserRepo(queries)
	passwordRepo := repository.NewPasswordRepo(queries)
	userService := service.NewUserService(userRepo, email.SMTPConfig{
		Host:     cfg.Email.Host,
		Port:     cfg.Email.Port,
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
	}, cfg.Primary.JWTSecret, passwordRepo)

	userHandler := handler.NewUserHandler(userService)
	authMiddleware := appmw.RequireAuth(cfg.Primary.JWTSecret)
	adminOrHR := appmw.RequireRoles("admin", "hr")
	rateLimiter := appmw.NewRateLimiter(5, 10)
	userHandler.RegisterRoutes(
		r,
		authMiddleware,
		adminOrHR,
		rateLimiter.Limit,
	)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})

	return r
}
