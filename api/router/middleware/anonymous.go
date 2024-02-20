package middleware

import (
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	authservice "github.com/DaoVuDat/trackpro-api/api/resource/auth/service"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	"net/http"
	"strings"
)

func AnonymousMiddleware(app *ctx.Application, withExpiredTokenValidation bool) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// extract access token
			authHeader := req.Header.Get("Authorization")

			// Extract the token from the Bearer schema
			splitToken := strings.Split(authHeader, " ")
			if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
				next.ServeHTTP(w, req)
				return
			}

			// access token
			accessToken := splitToken[1]

			// validate accessToken
			tokenDetail, err := jwt.ParseToken(app.Logger, accessToken, app.Config.AccessTokenPublicKey, withExpiredTokenValidation)
			if err != nil {
				next.ServeHTTP(w, req)
				return
			}

			// check in redis
			tokenService := authservice.NewTokenService(app.RedisClient)
			_, err = tokenService.GetToken(app, req.Context(), authdto.TokenDetailForRedis{
				UserId:    tokenDetail.UserId,
				TokenType: authdto.Access,
				Token:     *tokenDetail.Token,
			})
			if err != nil {
				next.ServeHTTP(w, req)
				return
			}

			app.Render.JSON(w, http.StatusOK, map[string]interface{}{
				"access_token": accessToken,
				"user_id":      tokenDetail.UserId,
				"role":         tokenDetail.Role,
			})
		})
	}
}
