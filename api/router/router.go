package router

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
	accounthandler "trackpro/api/resource/account/handler"
	authhandler "trackpro/api/resource/auth/handler"
	"trackpro/api/resource/healthcheck"
	profilehandler "trackpro/api/resource/profile/handler"
	"trackpro/api/router/common"
	"trackpro/api/router/middleware"
	"trackpro/util/ctx"
)

func SetupRouter(app *ctx.Application) *chi.Mux {
	router := chi.NewRouter()
	router.Use(chimiddleware.StripSlashes)
	router.Use(middleware.LoggingMiddleware(app.Logger))

	// Setup Version 1 Routing
	router.Route("/v1", func(g chi.Router) {
		g.Get("/healthcheck", healthcheck.V1Handler(app))
		g.Post("/signup", authhandler.SignUp(app))
		g.Post("/login", authhandler.Login(app))

		//g.Get("/profile/{id}", profilehandler.FindProfile(app))
		// Protected Routes
		g.Group(func(g chi.Router) {
			//g.Use(jwtauth.Verifier(app.JwtToken))
			//g.Use(jwtauth.Authenticator(app.JwtToken))

			g.Route("/profile", func(g chi.Router) {
				g.Get("/{id}", profilehandler.FindProfile(app))
				g.Patch("/{id}", profilehandler.UpdateProfile(app))
			})
			g.Route("/account", func(g chi.Router) {
				g.Group(func(g chi.Router) {
					g.Use(middleware.IsAdminMiddleware(app.Logger))
					g.Get("/", accounthandler.ListAccount(app))
				})
				g.Get("/{id}", accounthandler.FindAccount(app))
				g.Patch("/{id}", accounthandler.UpdateAccount(app))
			})

		})
	})

	// Not Found
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.Render.JSON(w, http.StatusNotFound, common.NotFoundErrorResponse(nil))
	})
	return router
}
