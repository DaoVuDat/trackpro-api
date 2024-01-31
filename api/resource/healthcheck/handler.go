package healthcheck

import (
	"github.com/uptrace/bunrouter"
	"net/http"
)

func CheckV1(w http.ResponseWriter, req bunrouter.Request) error {
	return bunrouter.JSON(w, bunrouter.H{
		"Status": "Good",
	})
}
