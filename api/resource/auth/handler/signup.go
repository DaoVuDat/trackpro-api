package authhandler

import (
	"encoding/json"
	"github.com/go-chi/render"
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
			panic(err)
		}

		curCtx := req.Context()
		accountCreateRepo := accountrepo.NewPostgresStore(app.Db)
		profileCreateRepo := profilerepo.NewPostgresStore(app.Db)
		signUpService := authservice.NewSignUpService(curCtx, accountCreateRepo, profileCreateRepo)

		token, err := signUpService.SignUp(app, authSignUp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, req, common.InternalErrorResponse(err))
			return
		}

		render.JSON(w, req, map[string]interface{}{
			"Token": token,
		})
	}
}
