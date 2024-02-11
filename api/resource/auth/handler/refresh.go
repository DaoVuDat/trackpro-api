package authhandler

import (
	"errors"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	authservice "github.com/DaoVuDat/trackpro-api/api/resource/auth/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	jwt2 "github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"time"
)

func Refresh(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// extract cookie
		cookie, err := req.Cookie(authconstant.RefreshTokenCookieHeader)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				app.Render.JSON(w, http.StatusUnauthorized, common.InternalErrorResponse(err))
				return
			}
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}
		refreshToken := cookie.Value

		tokenService := authservice.NewTokenService(app.RedisClient)
		// check refresh token
		refreshTokenDetail, err := jwt.ParseToken(app.Logger, refreshToken, app.Config.RefreshTokenPublicKey, true)
		if err != nil {
			return
		}
		if err != nil {
			if errors.Is(err, jwt2.ErrTokenExpired()) {
				app.Render.JSON(w, http.StatusUnauthorized,
					common.UnauthorizedErrorResponse(err))
				return
			}
			app.Render.JSON(w, http.StatusInternalServerError,
				common.InternalErrorResponse(err))
			return
		}

		// get access token from
		accessTokenDetail := req.Context().Value(authconstant.AccessTokenContextHeader).(*jwt.TokenDetail)

		newAccessToken, newRefreshToken, role, err := tokenService.RefreshAccessToken(app, req.Context(), authdto.PrivateClaimsForToken{
			UserId: refreshTokenDetail.UserId,
			Role:   refreshTokenDetail.Role,
		}, authdto.TokenDetailForRedis{
			UserId:    accessTokenDetail.UserId,
			TokenType: authdto.Access,
			Token:     *accessTokenDetail.Token,
		}, authdto.TokenDetailForRedis{
			UserId:    refreshTokenDetail.UserId,
			TokenType: authdto.Refresh,
			Token:     *refreshTokenDetail.Token,
		})
		if err != nil {
			if errors.Is(err, common.OldAccessTokenAndRefreshTokenNotMatch) {
				// we maybe revoke all token of this user
			}

			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		// refresh token cookie
		newCookie := http.Cookie{
			Expires:  time.Now().Add(app.Config.RefreshTokenExpiredIn),
			Secure:   false,
			HttpOnly: true,
			Path:     "/v1/token/refresh",
			SameSite: http.SameSiteLaxMode,
			Name:     authconstant.RefreshTokenCookieHeader,
			Value:    newRefreshToken,
		}

		http.SetCookie(w, &newCookie)

		app.Render.JSON(w, http.StatusOK, authdto.AuthResponse{
			AccessToken: newAccessToken,
			Role:        role,
		})
	}
}
