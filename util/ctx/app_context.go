package ctx

import (
	"database/sql"
	"github.com/rs/zerolog"
	"github.com/unrolled/render"
	"trackpro/util/config"
)

// Application Hold dependencies for our HTTP handlers, helpers, and middleware.
type Application struct {
	Logger *zerolog.Logger
	Config *config.EnvConfigs
	Db     *sql.DB
	Render *render.Render

	// Temp
	DataCache map[string][]byte
}
