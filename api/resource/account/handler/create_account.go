package accounthandler

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
	"trackpro/api/resource/account/dto"
	"trackpro/util/ctx"
)

func CreateAccount(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var createAccountData dto.AccountCreate
		app.Logger.Debug().Msgf("%v", "dsa")

		err := json.NewDecoder(req.Body).Decode(&createAccountData)
		if err != nil {
			panic(err)
		}
		app.Logger.Debug().Msgf("%v", createAccountData)
		render.JSON(w, req, map[string]interface{}{
			"Account": "Good",
		})
	}
}
