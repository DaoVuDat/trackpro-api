package authhandler

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
	accountrepo "trackpro/api/resource/account/repo"
	authdto "trackpro/api/resource/auth/dto"
	signupservice "trackpro/api/resource/auth/service"
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

		accountCreateRepo := accountrepo.NewPostgresStore(app.Db)
		profileCreateRepo := profilerepo.NewPostgresStore(app.Db)
		signUpService := signupservice.NewSignUpService(accountCreateRepo, profileCreateRepo)

		err = signUpService.SignUp(app, authSignUp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, req, common.InternalErrorResponse(err))
			return
		}

		render.JSON(w, req, map[string]interface{}{
			"Account": "Good",
		})
	}
}
