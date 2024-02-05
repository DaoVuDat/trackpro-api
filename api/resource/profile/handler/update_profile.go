package profilehandler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	profiledto "trackpro/api/resource/profile/dto"
	profilerepo "trackpro/api/resource/profile/repo"
	profileservice "trackpro/api/resource/profile/service"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
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
			if errors.Is(err, common.QueryNoResultErr) {
				app.Logger.Error().Err(err)
				app.Render.JSON(w, http.StatusNotFound, common.NotFoundErrorResponse(err))
				return
			}
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"profile": profile,
		})
	}
}
