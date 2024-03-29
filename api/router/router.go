package router

import (
	accounthandler "github.com/DaoVuDat/trackpro-api/api/resource/account/handler"
	authhandler "github.com/DaoVuDat/trackpro-api/api/resource/auth/handler"
	"github.com/DaoVuDat/trackpro-api/api/resource/healthcheck"
	paymenthandler "github.com/DaoVuDat/trackpro-api/api/resource/payment/handler"
	profilehandler "github.com/DaoVuDat/trackpro-api/api/resource/profile/handler"
	projecthandler "github.com/DaoVuDat/trackpro-api/api/resource/project/handler"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/api/router/middleware"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func SetupRouter(app *ctx.Application) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers

	}))
	router.Use(chimiddleware.StripSlashes)
	router.Use(middleware.LoggingMiddleware(app.Logger))

	// Setup Version 1 Routing
	router.Route("/v1", func(g chi.Router) {
		g.Get("/healthcheck", healthcheck.V1Handler(app))
		g.Group(func(g chi.Router) {
			g.Use(middleware.AnonymousMiddleware(app, true))
			g.Post("/signup", authhandler.SignUp(app))
			g.Post("/login", authhandler.Login(app))
		})

		g.Group(func(g chi.Router) {
			g.Use(middleware.AuthenticatedMiddleware(app, false))
			g.Post("/token/refresh", authhandler.Refresh(app))
		})

		g.Group(func(g chi.Router) {
			g.Use(middleware.AuthenticatedMiddleware(app, true))

			g.Group(func(g chi.Router) {
				g.Use(middleware.IsAdminMiddleware(app))
				g.Post("/payment", paymenthandler.CreatePayment(app))
			})

			g.Route("/account", func(g chi.Router) {
				g.Use(middleware.IsAdminMiddleware(app))
				g.Get("/", accounthandler.ListAccount(app))
				g.Patch("/{id}", accounthandler.UpdateAccount(app))
				g.Get("/{id}", accounthandler.FindAccount(app))
			})

			g.Route("/profile", func(g chi.Router) {
				g.Get("/{id}", profilehandler.FindProfile(app))
				g.Patch("/{id}", profilehandler.UpdateProfile(app))
			})

			g.Route("/project", func(g chi.Router) {
				g.Get("/", projecthandler.ListProject(app))
				g.Get("/{id}", projecthandler.FindProject(app))
				g.Patch("/{id}", projecthandler.UpdateProject(app))
				g.Group(func(g chi.Router) {
					g.Use(middleware.IsAdminMiddleware(app))
					g.Post("/", projecthandler.CreateProject(app))
					g.Delete("/{id}", projecthandler.DeleteProject(app))
				})
			})

		})
	})

	// Not Found
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.Render.JSON(w, http.StatusNotFound, common.NotFoundErrorResponse(nil))
	})
	return router
}
