package accounthandler

import (
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	accountservice "github.com/DaoVuDat/trackpro-api/api/resource/account/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"net/http"
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
