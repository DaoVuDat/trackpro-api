package middleware

import (
	"github.com/rs/zerolog"
	"github.com/uptrace/bunrouter"
	"net/http"
)

func LoggingMiddleware(logger *zerolog.Logger) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			logger.Info().
				Str("IP", req.RemoteAddr).
				Str("Host", req.Host).
				Str("URI", req.RequestURI).
				Str("Method", req.Method).
				Str("Protocol", req.Proto).
				Str("UserAgent", req.UserAgent()).
				Send()
			return next(w, req)
		}
	}
}
