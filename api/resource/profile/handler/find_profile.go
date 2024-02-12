package profilehandler

import (
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	profilerepo "github.com/DaoVuDat/trackpro-api/api/resource/profile/repo"
	profileservice "github.com/DaoVuDat/trackpro-api/api/resource/profile/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func FindProfile(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// parse params
		accountIdString := chi.URLParam(req, "id")
		accountId, err := uuid.Parse(accountIdString)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		// Get profile from access token
		accessTokenDetail := req.Context().Value(authconstant.AccessTokenContextHeader).(*jwt.TokenDetail)
		app.Logger.Debug().Msgf("accessTokenDetail.Role %v", accessTokenDetail.Role)
		app.Logger.Debug().Msgf("model.AccountType_Client.String() %v", model.AccountType_Client.String())
		app.Logger.Debug().Msgf("accessTokenDetail.UserId %v", accessTokenDetail.UserId)
		app.Logger.Debug().Msgf("accountId.String() %v", accountId.String())
		if accessTokenDetail.Role == model.AccountType_Client.String() && accessTokenDetail.UserId != accountId.String() {
			app.Render.JSON(w, http.StatusUnauthorized, common.UnauthorizedErrorResponse(nil))
			return
		}

		// setup repo and service
		findProfileRepo := profilerepo.NewPostgresStore(app.Db)
		findProfileService := profileservice.NewFindProfileService(findProfileRepo)

		profile, err := findProfileService.Find(app, accountId)
		if err != nil {
			if errors.Is(err, common.QueryNoResultErr) {
				app.Render.JSON(w, http.StatusNotFound, common.NotFoundErrorResponse(err))
				return
			}
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"profile": profile,
		})
	}
}
