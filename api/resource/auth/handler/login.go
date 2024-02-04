package authhandler

import (
	"encoding/json"
	"github.com/go-chi/render"
	"net/http"
	authdto "trackpro/api/resource/auth/dto"
	"trackpro/api/resource/auth/service"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
)

func Login(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var authLogin authdto.AuthLoginTemp

		err := json.NewDecoder(req.Body).Decode(&authLogin)
		if err != nil {
			panic(err)
		}

		logInService := authservice.NewLoginService()

		tokenDetail, err := logInService.Login(app, authLogin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, req, common.InternalErrorResponse(err))
			return
		}

		render.JSON(w, req, map[string]interface{}{
			"tokendetail": tokenDetail,
		})
	}
}
