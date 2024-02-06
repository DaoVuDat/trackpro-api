package healthcheck

import (
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"net/http"
)

func V1Handler(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"Status": "Good",
		})
	}
}
