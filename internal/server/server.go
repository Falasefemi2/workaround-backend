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
	deptRepo := repository.NewDepartmentRepo(queries)
	designationRepo := repository.NewDesignationRepo(queries)
	unitRepo := repository.NewUnitRepo(queries)
	levelRepo := repository.NewLevelRepo(queries)

	userService := service.NewUserService(
		userRepo,
		email.SMTPConfig{
			Host:     cfg.Email.Host,
			Port:     cfg.Email.Port,
			Username: cfg.Email.Username,
			Password: cfg.Email.Password,
		},
		cfg.Primary.JWTSecret,
		passwordRepo,
	)

	deptService := service.NewDeptService(deptRepo, userRepo)
	designationService := service.NewDesignationService(designationRepo)
	unitService := service.NewUnitService(unitRepo, userRepo)
	levelService := service.NewLevelService(levelRepo)

	userHandler := handler.NewUserHandler(userService)
	deptHandler := handler.NewDeptHandler(deptService)
	designationHandler := handler.NewDesignationHandler(designationService)
	unitHandler := handler.NewUnitHandler(unitService)
	levelHandler := handler.NewLevelHandler(levelService)

	authMiddleware := appmw.RequireAuth(cfg.Primary.JWTSecret)
	adminOrHR := appmw.RequireRoles("admin", "hr")
	rateLimiter := appmw.NewRateLimiter(5, 10)

	userHandler.RegisterRoutes(
		r,
		authMiddleware,
		adminOrHR,
		rateLimiter.Limit,
	)

	deptHandler.RegisterRoutes(
		r,
		authMiddleware,
		adminOrHR,
	)

	designationHandler.RegisterRoutes(
		r,
		authMiddleware,
		adminOrHR,
	)

	unitHandler.RegisterRoutes(
		r,
		authMiddleware,
		adminOrHR,
	)
	levelHandler.RegisterRoutes(r, authMiddleware, adminOrHR)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			http.Error(w, "failed to write response", http.StatusInternalServerError)
			return
		}
	})

	return r
}
