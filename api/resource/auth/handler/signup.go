package authhandler

import (
	"encoding/json"
	"errors"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	authservice "github.com/DaoVuDat/trackpro-api/api/resource/auth/service"
	passwordrepo "github.com/DaoVuDat/trackpro-api/api/resource/password/repo"
	profilerepo "github.com/DaoVuDat/trackpro-api/api/resource/profile/repo"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"net/http"
	"strings"
	"time"
)

func SignUp(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var authSignUp authdto.AuthSignUp

		err := json.NewDecoder(req.Body).Decode(&authSignUp)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		if err = authSignUp.Validate(); err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		if strings.Compare(authSignUp.Password, authSignUp.ConfirmedPassword) != 0 {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(common.CPassAndPassNotMatch))
			return
		}

		accountCreateRepo := accountrepo.NewPostgresStore(app.Db)
		profileCreateRepo := profilerepo.NewPostgresStore(app.Db)
		passwordCreateRepo := passwordrepo.NewPostgresStore(app.Db)
		tokenService := authservice.NewTokenService(app.RedisClient)
		signUpService := authservice.NewSignUpService(req.Context(),
			accountCreateRepo,
			profileCreateRepo,
			passwordCreateRepo,
			tokenService,
		)

		accessToken, refreshToken, role, uid, err := signUpService.SignUp(app, authSignUp)
		if err != nil {
			if errors.Is(err, common.FailCreateError) {
				app.Logger.Error().Err(err).Send()
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}

			if errors.Is(err, common.DuplicateValueError) {
				app.Logger.Error().Err(err).Send()
				app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
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
			UserId:      uid,
		})
	}
}
