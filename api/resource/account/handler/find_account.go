package accounthandler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	accountrepo "trackpro/api/resource/account/repo"
	accountservice "trackpro/api/resource/account/service"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
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
