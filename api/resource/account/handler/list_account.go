package accounthandler

import (
	"net/http"
	accountrepo "trackpro/api/resource/account/repo"
	accountservice "trackpro/api/resource/account/service"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
)

func ListAccount(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		// setup repo and service
		listAccountRepo := accountrepo.NewPostgresStore(app.Db)
		listAccountService := accountservice.NewListAccountService(listAccountRepo)

		accounts, err := listAccountService.List(app)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"accounts": accounts,
		})
	}
}
