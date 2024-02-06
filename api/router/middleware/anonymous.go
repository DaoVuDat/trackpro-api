package middleware

import (
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"net/http"
)

func AnonymousMiddleware(app *ctx.Application) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
		})
	}
}
