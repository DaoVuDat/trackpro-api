package ctx

import (
	"database/sql"
	"github.com/DaoVuDat/trackpro-api/util/config"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/unrolled/render"
)

// Application Hold dependencies for our HTTP handlers, helpers, and middleware.
type Application struct {
	Logger      *zerolog.Logger
	Config      *config.EnvConfigs
	Db          *sql.DB
	Render      *render.Render
	RedisClient *redis.Client
}
