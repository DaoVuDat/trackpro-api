package middleware

import (
	"errors"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	authservice "github.com/DaoVuDat/trackpro-api/api/resource/auth/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	jwt2 "github.com/lestrrat-go/jwx/v2/jwt"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

func AuthenticatedMiddleware(app *ctx.Application, withExpiredTokenValidation bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// extract access token
			authHeader := req.Header.Get("Authorization")

			// Extract the token from the Bearer schema
			splitToken := strings.Split(authHeader, " ")
			if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
				app.Render.JSON(w, http.StatusUnauthorized,
					common.UnauthorizedErrorResponse(errors.New("invalid authentication credentials")))
				return
			}

			// access token
			accessToken := splitToken[1]

			// validate accessToken
			tokenDetail, err := jwt.ParseToken(app.Logger, accessToken, app.Config.AccessTokenPublicKey, withExpiredTokenValidation)
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

			// check access token in redis
			tokenService := authservice.NewTokenService(app.RedisClient)
			_, err = tokenService.GetToken(app, req.Context(), authdto.TokenDetailForRedis{
				UserId:    tokenDetail.UserId,
				TokenType: authdto.Access,
				Token:     *tokenDetail.Token,
			})
			if err != nil {
				if !errors.Is(err, common.TokenExpired) {
					app.Render.JSON(w, http.StatusInternalServerError,
						common.InternalErrorResponse(err))
					return
				}
			}

			// add token detail to context
			contextWithPayload := context.WithValue(req.Context(), authconstant.AccessTokenContextHeader, tokenDetail)
			req = req.WithContext(contextWithPayload)

			next.ServeHTTP(w, req)
		})
	}
}
