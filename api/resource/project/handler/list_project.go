package projecthandler

import (
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	projectservice "github.com/DaoVuDat/trackpro-api/api/resource/project/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

func ListProject(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userIdString := req.URL.Query().Get("uid")

		userId, err := uuid.Parse(userIdString)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		payment := false
		paymentString := req.URL.Query().Get("returnPayment")

		if parsedPayment, err := strconv.ParseBool(paymentString); err == nil {
			payment = parsedPayment
		}

		listProjectRepo := projectrepo.NewPostgresStore(app.Db)
		listProjectService := projectservice.NewListProjectService(listProjectRepo)

		projects, err := listProjectService.List(app, userId, payment)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"projects": projects,
		})

	}
}
