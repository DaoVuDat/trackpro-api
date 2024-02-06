package projecthandler

import (
	"errors"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	projectservice "github.com/DaoVuDat/trackpro-api/api/resource/project/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func FindProject(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		projectIdString := chi.URLParam(req, "id")
		projectId, err := uuid.Parse(projectIdString)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		userIdString := req.URL.Query().Get("uid")
		var userId *uuid.UUID
		if len(userIdString) > 0 {
			parsedUuid, err := uuid.Parse(userIdString)
			if err != nil {
				app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
				return
			}

			userId = &parsedUuid
		}

		findProjectRepo := projectrepo.NewPostgresStore(app.Db)
		findProjectService := projectservice.NewFindProjectService(findProjectRepo)

		projectResponse, err := findProjectService.Find(app, projectId, userId)
		if err != nil {
			if errors.Is(err, common.QueryNoResultErr) {
				app.Render.JSON(w, http.StatusNotFound, common.NotFoundErrorResponse(err))
				return
			}
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"project": projectResponse,
		})
	}
}
