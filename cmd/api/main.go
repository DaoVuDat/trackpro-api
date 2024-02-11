package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/DaoVuDat/trackpro-api/api/router"
	"github.com/DaoVuDat/trackpro-api/util/config"
	"github.com/DaoVuDat/trackpro-api/util/cstlogger"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/graceful"
	"github.com/redis/go-redis/v9"
	"github.com/unrolled/render"
	"net/http"
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
		Config: config.LoadEnvConfigs(logger, "."), // Get env variables
	}

	// Setup DB
	cfg, err := pgxpool.ParseConfig(app.Config.DBDsn)
	if err != nil {
		app.Logger.Error().Err(err)
		panic(1)
	}
	cfg.MaxConns = app.Config.DBMaxConnection
	cfg.MinConns = app.Config.DBMinConnection
	cfg.MaxConnLifetime = app.Config.DBMaxConnectionLifeTime
	cfg.MaxConnIdleTime = app.Config.DBMaxConnectionIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		app.Logger.Error().Err(err)
		panic(1)
	}
	defer pool.Close()
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	app.Db = db

	// Render
	renderer := render.New()
	app.Render = renderer

	// Setup Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     app.Config.RedisAddressAndPort,
		Password: app.Config.RedisPassword,
		DB:       app.Config.RedisDB,
	})
	app.RedisClient = rdb

	// Setup Route
	r := router.SetupRouter(app)

	// Start Server
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
