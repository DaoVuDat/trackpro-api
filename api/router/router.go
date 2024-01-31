package router

import (
	"github.com/uptrace/bunrouter"
	"trackpro/api/resource/healthcheck"
)

func SetupRouter() *bunrouter.Router {
	router := bunrouter.New()

	// Setup Version 1 Routing
	router.WithGroup("/v1", func(g *bunrouter.Group) {
		g.GET("/healthcheck", healthcheck.CheckV1)
	})

	return router
}
