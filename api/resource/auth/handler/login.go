package authhandler

import (
	"encoding/json"
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
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		logInService := authservice.NewLoginService()

		tokenDetail, err := logInService.Login(app, authLogin)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"tokendetail": tokenDetail,
		})
	}
}
