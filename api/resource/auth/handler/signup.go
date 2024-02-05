package authhandler

import (
	"encoding/json"
	"net/http"
	accountrepo "trackpro/api/resource/account/repo"
	authdto "trackpro/api/resource/auth/dto"
	"trackpro/api/resource/auth/service"
	profilerepo "trackpro/api/resource/profile/repo"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
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
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"Token": token,
		})
	}
}
