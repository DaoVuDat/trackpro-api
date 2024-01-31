package accounthandler

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
	accountdto "trackpro/api/resource/account/dto"

	accountrepo "trackpro/api/resource/account/repo"
	accountservice "trackpro/api/resource/account/service"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
)

func CreateAccount(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var createAccountData accountdto.AccountCreate

		err := json.NewDecoder(req.Body).Decode(&createAccountData)
		if err != nil {
			panic(err)
		}
		app.Logger.Debug().Any("username", createAccountData.UserName)
		createRepo := accountrepo.NewPostgresStore(app.Db)
		createService := accountservice.NewCreateAccountService(createRepo)

		err = createService.CreateAccount(app, createAccountData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, req, common.InternalErrorResponse(err))
			return
		}

		render.JSON(w, req, map[string]interface{}{
			"Account": "Good",
		})
	}
}
