package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter creates a chi router with auth endpoints and CORS.
func NewRouter(jwtSecret, corsOrigin string) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware(corsOrigin))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/request-otp", HandleRequestOTP)
		r.Post("/auth/verify-otp", HandleVerifyOTP)
		r.Post("/auth/qr/init", HandleQRInit)
		r.Get("/auth/qr/status", HandleQRPoll)
		r.Post("/auth/refresh", HandleRefresh)
		r.Post("/auth/logout", HandleLogout)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(JWTAuth(jwtSecret))

			r.Get("/transactions", HandleListTransactions)
			r.Post("/transactions", HandleCreateTransaction)
			r.Put("/transactions/{id}", HandleUpdateTransaction)
			r.Delete("/transactions/{id}", HandleDeleteTransaction)

			r.Get("/wallets", HandleListWallets)
			r.Get("/contacts", HandleListContacts)

			r.Get("/budgets", HandleListBudgets)
			r.Post("/budgets", HandleSetBudget)
			r.Delete("/budgets/{categoryID}", HandleDeleteBudget)
			r.Get("/budgets/alerts", HandleBudgetAlerts)

			r.Get("/summary/charts", HandleChartData)
			r.Get("/categories", HandleListCategories)
			r.Get("/profile", HandleGetProfile)
		})
	})

	return r
}

func corsMiddleware(origin string) func(next http.Handler) http.Handler {
	if origin == "" {
		origin = "*"
	}
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{origin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}
