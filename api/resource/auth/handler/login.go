package authhandler

import (
	"encoding/json"
	"errors"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	authservice "github.com/DaoVuDat/trackpro-api/api/resource/auth/service"
	passwordrepo "github.com/DaoVuDat/trackpro-api/api/resource/password/repo"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"net/http"
	"time"
)

func Login(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var authLogin authdto.AuthLogin

		err := json.NewDecoder(req.Body).Decode(&authLogin)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		if err = authLogin.Validate(); err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		tokenService := authservice.NewTokenService(app.RedisClient)
		findAccountRepo := accountrepo.NewPostgresStore(app.Db)
		findPasswordRepo := passwordrepo.NewPostgresStore(app.Db)
		logInService := authservice.NewLoginService(
			findAccountRepo,
			findPasswordRepo,
			tokenService,
		)

		accessToken, refreshToken, role, err := logInService.Login(app, req.Context(), authLogin)
		if err != nil {
			if errors.Is(err, common.QueryNoResultErr) {
				app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(errors.New("invalid credentials")))
				return
			}
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		// refresh token cookie
		cookie := http.Cookie{
			Expires:  time.Now().Add(app.Config.RefreshTokenExpiredIn),
			Secure:   false,
			HttpOnly: true,
			Path:     "/v1/token/refresh",
			SameSite: http.SameSiteLaxMode,
			Name:     authconstant.RefreshTokenCookieHeader,
			Value:    refreshToken,
		}

		http.SetCookie(w, &cookie)

		app.Render.JSON(w, http.StatusOK, authdto.AuthResponse{
			AccessToken: accessToken,
			Role:        role,
		})
	}
}
