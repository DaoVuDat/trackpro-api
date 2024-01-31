package middleware

import (
	"github.com/rs/zerolog"
	"net/http"
)

func LoggingMiddleware(logger *zerolog.Logger) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger.Info().
				Str("IP", req.RemoteAddr).
				Str("Host", req.Host).
				Str("Protocol", req.Proto).
				Str("URI", req.RequestURI).
				Str("Method", req.Method).
				Str("UserAgent", req.UserAgent()).
				Send()
			next.ServeHTTP(w, req)
		})
	}
}
