package middleware

import (
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	"net/http"
	"strings"
)

func IsAdminMiddleware(app *ctx.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// get access token from
			accessTokenDetail := req.Context().Value(authconstant.AccessTokenContextHeader).(*jwt.TokenDetail)

			app.Logger.Debug().Any("Access Token", accessTokenDetail.Token).Send()
			app.Logger.Debug().Any("User Id", accessTokenDetail.UserId).Send()
			app.Logger.Debug().Any("Role", accessTokenDetail.Role).Send()

			if strings.Compare(accessTokenDetail.Role, model.AccountType_Admin.String()) != 0 {
				app.Render.JSON(w, http.StatusUnauthorized,
					common.UnauthorizedErrorResponse(errors.New("invalid authentication credentials")))
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
