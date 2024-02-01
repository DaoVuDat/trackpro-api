package healthcheck

import (
	"github.com/go-chi/render"

	"net/http"
)

func V1Handler(w http.ResponseWriter, req *http.Request) {
	render.JSON(w, req, map[string]interface{}{
		"Status": "Good",
	})
}
