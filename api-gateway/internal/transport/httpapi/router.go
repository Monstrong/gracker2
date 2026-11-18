package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/monstrong/gracker2/api-gateway/internal/health"
	"github.com/monstrong/gracker2/api-gateway/internal/middleware"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/handlers"
	"github.com/monstrong/gracker2/pkg/jwtparser"
	"github.com/monstrong/gracker2/pkg/logger"
)

func NewRouter(
	authHandler *handlers.AuthHandler,
	torrentHandler *handlers.TorrentHandler,
	tokenManager jwtparser.TokenManager,
    chacker *health.Checker,
	l logger.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(middleware.RequestLogger(l))
	r.Use(middleware.Recoverer(l))

	r.Get("/health", Health)
    r.Get("/ready", Ready(chacker))

    // Public auth routes
    r.Route("/auth", func(r chi.Router) {
        r.Post("/register", authHandler.Register)
        r.Post("/login", authHandler.Login)
        r.Post("/refresh", authHandler.Refresh)
        r.Post("/logout", authHandler.Logout)
    })

    // Public user routes
    r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
        authHandler.GetUser(w, r, chi.URLParam(r, "id"))
    })

    // Authenticated routes
    r.Group(func(r chi.Router) {
        r.Use(middleware.AuthMiddleware(tokenManager))

        r.Post("/auth/logout-all", authHandler.LogoutEverywhere)

        r.Get("/users/me", authHandler.GetMe)
        r.Patch("/users/me", authHandler.UpdateUser)
        r.Patch("/users/me/password", authHandler.ChangePassword)
        r.Delete("/users", authHandler.DeleteUser)
        // Admin routes
        r.Get("/admin/users/{id}", func(w http.ResponseWriter, r *http.Request) {
            authHandler.AdminGetUser(w, r, chi.URLParam(r, "id"))
        })
        r.Patch("/admin/users/role", authHandler.RoleUpdate)
        r.Post("/admin/users/ban", authHandler.BanUser)
        r.Delete("/admin/users/ban", authHandler.UnbanUser)
    })

	return r
}
