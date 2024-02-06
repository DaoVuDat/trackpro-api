package projecthandler

import (
	"encoding/json"
	"errors"
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	projectservice "github.com/DaoVuDat/trackpro-api/api/resource/project/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"net/http"
)

func CreateProject(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var projectCreate projectdto.ProjectCreate

		if err := json.NewDecoder(req.Body).Decode(&projectCreate); err != nil {
			app.Logger.Error().Msg("Error Decode JSON")
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		if err := projectCreate.Validate(); err != nil {
			app.Logger.Error().Err(err)
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		createProjectRepo := projectrepo.NewPostgresStore(app.Db)
		createProjectService := projectservice.NewCreateProjectService(createProjectRepo)

		project, err := createProjectService.Create(app, projectCreate)
		if err != nil {
			if errors.Is(err, common.FailCreateError) {
				app.Logger.Error().Err(err)
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}
			app.Logger.Error().Err(err)
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"project": project,
		})
	}
}
