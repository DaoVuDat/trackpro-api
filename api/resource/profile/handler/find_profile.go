package profilehandler

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	profilerepo "trackpro/api/resource/profile/repo"
	profileservice "trackpro/api/resource/profile/service"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
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
