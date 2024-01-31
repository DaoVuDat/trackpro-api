package main

import (
	"flag"
	"fmt"
	"github.com/ory/graceful"
	"net/http"
	"trackpro/api/router"
	"trackpro/util/config"
	"trackpro/util/cstlogger"
	"trackpro/util/ctx"
)

func main() {
	// Get env mode from flags
	var environment = flag.String("env", "prod", "set environment mode (prod or dev)")
	flag.Parse()

	// Setup custom logger
	logger := cstlogger.NewLogger(*environment)

	// Create application context
	app := &ctx.Application{
		Logger: logger,
	}

	// Get env variables
	app.Config = config.LoadEnvConfigs(app.Logger, ".")

	// inject other dependencies
	r := router.SetupRouter()

	server := graceful.WithDefaults(&http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.ServerAddressPort),
		Handler:      r,
		ReadTimeout:  app.Config.ServerTimeoutRead,
		WriteTimeout: app.Config.ServerTimoutWrite,
		IdleTimeout:  app.Config.ServerTimeoutIdle,
	})

	app.Logger.Info().Msgf("Listening on PORT: %d", app.Config.ServerAddressPort)
	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		app.Logger.Error().Msg("Failed to gracefully shutdown")
	}
	app.Logger.Info().Msg("Server was shutdown gracefully")
}
