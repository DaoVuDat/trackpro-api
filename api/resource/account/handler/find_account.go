package accounthandler

import (
	"errors"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	accountservice "github.com/DaoVuDat/trackpro-api/api/resource/account/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func FindAccount(app *ctx.Application) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// parse params
		accountIdString := chi.URLParam(r, "id")
		accountId, err := uuid.Parse(accountIdString)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		// setup repo and service
		accountRepo := accountrepo.NewPostgresStore(app.Db)
		accountService := accountservice.NewFindAccountService(accountRepo)

		account, err := accountService.Find(app, accountId)
		if err != nil {
			if errors.Is(err, common.QueryNoResultErr) {
				app.Render.JSON(w, http.StatusNotFound, common.NotFoundErrorResponse(err))
				return
			}
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"account": account,
		})
	}
}
