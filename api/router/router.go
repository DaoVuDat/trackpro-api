package router

import (
	"github.com/uptrace/bunrouter"
	accounthandler "trackpro/api/resource/account/handler"
	"trackpro/api/resource/healthcheck"
	"trackpro/api/router/middleware"
	"trackpro/util/ctx"
)

func SetupRouter(app *ctx.Application) *bunrouter.Router {
	router := bunrouter.New(
		// Install Middleware for All Routing
		bunrouter.Use(middleware.LoggingMiddleware(app.Logger)),
	)

	accHandler := accounthandler.New(app)

	// Setup Version 1 Routing
	router.WithGroup("/v1", func(g *bunrouter.Group) {
		g.GET("/healthcheck", healthcheck.CheckV1)
		g.GET("/account", accHandler.CreateAccount)
	})

	return router
}
