package router

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
	accounthandler "trackpro/api/resource/account/handler"
	authhandler "trackpro/api/resource/auth/handler"
	"trackpro/api/resource/healthcheck"
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
		g.Get("/healthcheck", healthcheck.V1Handler)
		g.Post("/signup", authhandler.SignUp(app))

		// Protected Routes
		g.Group(func(g chi.Router) {
			//g.Use(jwtauth.Verifier(app.JwtToken))
			//g.Use(jwtauth.Authenticator(app.JwtToken))
			g.Route("/account", func(g chi.Router) {
				g.Get("/", accounthandler.ListAccount(app))
				g.Post("/", accounthandler.CreateAccount(app))
				g.Patch("/:id", accounthandler.UpdateAccount(app))
			})
		})
	})

	// Not Found
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, common.NotFoundErrorResponse(nil))
	})
	return router
}
