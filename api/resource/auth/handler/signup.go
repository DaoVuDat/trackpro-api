package authhandler

import (
	"encoding/json"
	"errors"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	authservice "github.com/DaoVuDat/trackpro-api/api/resource/auth/service"
	profilerepo "github.com/DaoVuDat/trackpro-api/api/resource/profile/repo"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"net/http"
)

func SignUp(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var authSignUp authdto.AuthSignUp

		err := json.NewDecoder(req.Body).Decode(&authSignUp)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		curCtx := req.Context()
		accountCreateRepo := accountrepo.NewPostgresStore(app.Db)
		profileCreateRepo := profilerepo.NewPostgresStore(app.Db)
		signUpService := authservice.NewSignUpService(curCtx, accountCreateRepo, profileCreateRepo)

		token, err := signUpService.SignUp(app, authSignUp)
		if err != nil {
			if errors.Is(err, common.FailCreateError) {
				app.Logger.Error().Err(err)
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}

			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"Token": token,
		})
	}
}
