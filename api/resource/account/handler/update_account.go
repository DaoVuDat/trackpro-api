package accounthandler

import (
	"encoding/json"
	"errors"
	accountdto "github.com/DaoVuDat/trackpro-api/api/resource/account/dto"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	accountservice "github.com/DaoVuDat/trackpro-api/api/resource/account/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func UpdateAccount(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var updateAccountData accountdto.AccountUpdate

		// Get {id} param
		accountIdString := chi.URLParam(req, "id")
		accountId, err := uuid.Parse(accountIdString)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		err = json.NewDecoder(req.Body).Decode(&updateAccountData)
		if err != nil {
			app.Logger.Error().Msg("Error Decode JSON")
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		// validate input
		err = updateAccountData.Validate()
		if err != nil {
			app.Logger.Error().Err(err)
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		// create repo and service
		updateAccountRepo := accountrepo.NewPostgresStore(app.Db)
		updateAccountService := accountservice.NewUpdateAccountService(updateAccountRepo)

		updatedAccount, err := updateAccountService.Update(app, accountId, updateAccountData)
		if err != nil {
			if errors.Is(err, common.FailUpdateError) {
				app.Logger.Error().Err(err)
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}
			app.Logger.Error().Err(err)
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, []map[string]interface{}{
			{"accounts": updatedAccount},
		})
	}
}
