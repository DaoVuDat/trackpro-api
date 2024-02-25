package profilehandler

import (
	"encoding/json"
	"errors"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	profiledto "github.com/DaoVuDat/trackpro-api/api/resource/profile/dto"
	profilerepo "github.com/DaoVuDat/trackpro-api/api/resource/profile/repo"
	profileservice "github.com/DaoVuDat/trackpro-api/api/resource/profile/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func UpdateProfile(app *ctx.Application) http.HandlerFunc {
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
		if accessTokenDetail.UserId != accountId.String() {
			app.Render.JSON(w, http.StatusUnauthorized, common.UnauthorizedErrorResponse(nil))
			return
		}

		var profileUpdate profiledto.ProfileUpdate
		err = json.NewDecoder(req.Body).Decode(&profileUpdate)
		if err != nil {
			app.Logger.Error().Msg("Error Decode JSON")
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		// validate request body
		if err = profileUpdate.Validate(); err != nil {
			app.Logger.Error().Err(err)
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		// setup repo and service
		updateProfileRepo := profilerepo.NewPostgresStore(app.Db)
		updateProfileService := profileservice.NewUpdateProfileService(updateProfileRepo)

		profile, err := updateProfileService.Update(app, accountId, profileUpdate)
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

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"profile": profile,
		})
	}
}
