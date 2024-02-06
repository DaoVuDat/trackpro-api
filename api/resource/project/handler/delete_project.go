package projecthandler

import (
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	projectservice "github.com/DaoVuDat/trackpro-api/api/resource/project/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func DeleteProject(app *ctx.Application) http.HandlerFunc {
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

		deleteProjectRepo := projectrepo.NewPostgresStore(app.Db)
		deleteProjectService := projectservice.NewDeleteProjectService(deleteProjectRepo)

		if err = deleteProjectService.Delete(app, projectId, userId); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"success": "ok",
		})
	}
}
