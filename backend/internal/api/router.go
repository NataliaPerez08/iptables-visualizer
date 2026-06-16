package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/anomalyco/iptables-visualizer/internal/api/handlers"
	"github.com/anomalyco/iptables-visualizer/internal/api/middleware"
	"github.com/anomalyco/iptables-visualizer/internal/config"
	"github.com/anomalyco/iptables-visualizer/internal/engine"
	"github.com/anomalyco/iptables-visualizer/internal/repository"
)

func NewRouter(db *sql.DB, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.RequestLoggingMiddleware)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	userRepo := repository.NewUserRepository(db)
	policyRepo := repository.NewPolicyRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	compiler := engine.NewCompiler()
	validator := engine.NewValidator()

	authHandler := handlers.NewAuthHandler(userRepo, auditRepo, cfg.JWT)
	policyHandler := handlers.NewPolicyHandler(policyRepo, auditRepo, compiler, validator)
	auditHandler := handlers.NewAuditHandler(auditRepo)

	authMw := middleware.AuthMiddleware(cfg.JWT.Secret)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/auth/login", authHandler.Login)

		r.Group(func(r chi.Router) {
			r.Use(authMw)

			r.Get("/auth/me", authHandler.Me)

			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.RoleMiddleware("admin"))
				r.Post("/", authHandler.CreateUser)
				r.Get("/", authHandler.ListUsers)
			})

			r.Route("/policies", func(r chi.Router) {
				r.Get("/", policyHandler.List)
				r.Post("/", policyHandler.Create)
				r.Get("/{id}", policyHandler.Get)
				r.Post("/{id}/validate", policyHandler.Validate)
				r.Post("/{id}/dry-run", policyHandler.DryRun)

				r.Group(func(r chi.Router) {
					r.Use(middleware.RoleMiddleware("admin", "editor"))
					r.Put("/{id}", policyHandler.Update)
					r.Post("/{id}/apply", policyHandler.Apply)
				})

				r.With(middleware.RoleMiddleware("admin")).Delete("/{id}", policyHandler.Delete)
			})

			r.Route("/audit", func(r chi.Router) {
				r.Use(middleware.RoleMiddleware("admin"))
				r.Get("/", auditHandler.Query)
			})
		})
	})

	return r
}
