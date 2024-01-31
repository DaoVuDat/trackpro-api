package accounthandler

import (
	"github.com/uptrace/bunrouter"
	"net/http"
)

func (handler *Handler) CreateAccount(w http.ResponseWriter, r bunrouter.Request) error {
	w.Header().Set("set-cookie", "dadsa")
	return bunrouter.JSON(w, bunrouter.H{
		"Account": "Good",
	})
}
