package accounthandler

import (
	"github.com/go-chi/render"
	"net/http"
	"trackpro/util/ctx"
)

func UpdateAccount(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		render.JSON(w, req, map[string]interface{}{
			"Update Account": "Good",
		})
	}
}
