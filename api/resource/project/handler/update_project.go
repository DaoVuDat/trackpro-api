package projecthandler

import (
	"encoding/json"
	"errors"
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	projectservice "github.com/DaoVuDat/trackpro-api/api/resource/project/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
)

func UpdateProject(app *ctx.Application) http.HandlerFunc {
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

		var projectUpdate projectdto.ProjectUpdate
		if err = json.NewDecoder(req.Body).Decode(&projectUpdate); err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(errors.New("invalid 'price'")))
				return
			}
			if strings.Contains(err.Error(), "parsing time") {
				app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(errors.New("time: use ISO8601 format")))
				return
			}

			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		if err = projectUpdate.Validate(); err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		updateProjectRepo := projectrepo.NewPostgresStore(app.Db)
		updateProjectService := projectservice.NewUpdateProjectService(updateProjectRepo)

		project, err := updateProjectService.Update(app, projectId, userId, projectUpdate)
		if err != nil {
			if errors.Is(err, common.FailUpdateError) {
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"project": project,
		})
	}
}
