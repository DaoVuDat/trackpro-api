package router

import (
	"github.com/go-chi/chi/v5"
	accounthandler "trackpro/api/resource/account/handler"
	authhandler "trackpro/api/resource/auth/handler"
	"trackpro/api/resource/healthcheck"
	"trackpro/api/router/middleware"
	"trackpro/util/ctx"
)

func SetupRouter(app *ctx.Application) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.LoggingMiddleware(app.Logger))

	// Setup Version 1 Routing
	router.Route("/v1", func(g chi.Router) {
		g.Get("/healthcheck", healthcheck.CheckV1)
		g.Post("/signup", authhandler.SignUp(app))
		g.Route("/account", func(g chi.Router) {
			g.Get("/", accounthandler.ListAccount(app))
			g.Post("/", accounthandler.CreateAccount(app))
			g.Patch("/account/:id", accounthandler.UpdateAccount(app))
		})
	})

	return router
}
