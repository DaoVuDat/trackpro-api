package ctx

import (
	"github.com/rs/zerolog"
	"trackpro/util/config"
)

// Application Hold dependencies for our HTTP handlers, helpers, and middleware.
type Application struct {
	Logger *zerolog.Logger
	Config *config.EnvConfigs
}
