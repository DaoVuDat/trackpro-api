package projecthandler

import (
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	projectservice "github.com/DaoVuDat/trackpro-api/api/resource/project/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

func FindProject(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		projectIdString := chi.URLParam(req, "id")
		projectId, err := uuid.Parse(projectIdString)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		// Get profile from access token
		accessTokenDetail := req.Context().Value(authconstant.AccessTokenContextHeader).(*jwt.TokenDetail)
		app.Logger.Debug().Msgf("accessTokenDetail.Role %v", accessTokenDetail.Role)
		app.Logger.Debug().Msgf("model.AccountType_Client.String() %v", model.AccountType_Client.String())
		app.Logger.Debug().Msgf("accessTokenDetail.UserId %v", accessTokenDetail.UserId)

		uid, err := uuid.Parse(accessTokenDetail.UserId)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		payment := false
		paymentString := req.URL.Query().Get("returnPayment")

		if parsedPayment, err := strconv.ParseBool(paymentString); err == nil {
			payment = parsedPayment
		}

		findProjectRepo := projectrepo.NewPostgresStore(app.Db)
		findProjectService := projectservice.NewFindProjectService(findProjectRepo)

		var projectResponse *projectdto.ProjectResponse

		if accessTokenDetail.Role != model.AccountType_Client.String() {
			projectResponse, err = findProjectService.Find(app, projectId, nil, payment)
			if err != nil {
				if errors.Is(err, common.QueryNoResultErr) {
					app.Render.JSON(w, http.StatusNotFound, common.NotFoundErrorResponse(err))
					return
				}
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}
		} else {
			projectResponse, err = findProjectService.Find(app, projectId, &uid, payment)
			if err != nil {
				if errors.Is(err, common.QueryNoResultErr) {
					app.Render.JSON(w, http.StatusNotFound, common.NotFoundErrorResponse(err))
					return
				}
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"project": projectResponse,
		})
	}
}
