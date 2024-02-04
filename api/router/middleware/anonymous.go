package middleware

import (
	"net/http"
	"trackpro/util/ctx"
)

func AnonymousMiddleware(app *ctx.Application) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
		})
	}
}
