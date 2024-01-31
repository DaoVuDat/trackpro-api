package healthcheck

import (
	"github.com/go-chi/render"

	"net/http"
)

func CheckV1(w http.ResponseWriter, req *http.Request) {
	render.JSON(w, req, map[string]interface{}{
		"Status": "Good",
	})
}
