package accounthandler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	accountdto "trackpro/api/resource/account/dto"
	accountrepo "trackpro/api/resource/account/repo"
	accountservice "trackpro/api/resource/account/service"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
)

func UpdateAccount(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var updateAccountData accountdto.AccountUpdate

		// Get {id} param
		accountId := chi.URLParam(req, "id")
		app.Logger.Debug().Msg(accountId)

		err := json.NewDecoder(req.Body).Decode(&updateAccountData)
		if err != nil {
			app.Logger.Error().Msg("Error Decode JSON")
			panic(1)
		}

		// validate input
		err = updateAccountData.Validate()
		if err != nil {
			app.Logger.Error().Err(err)
			render.JSON(w, req, common.InternalErrorResponse(err))
			return
		}

		// create repo and service
		updateAccountRepo := accountrepo.NewPostgresStore(app.Db)
		updateAccountService := accountservice.NewUpdateAccountService(updateAccountRepo)

		updatedAccount, err := updateAccountService.Update(app, accountId, updateAccountData)
		if err != nil {
			app.Logger.Error().Err(err)
			render.JSON(w, req, common.InternalErrorResponse(err))
			return
		}

		render.JSON(w, req, []map[string]interface{}{
			{"accounts": updatedAccount},
		})
	}
}
