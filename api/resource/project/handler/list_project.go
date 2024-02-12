package projecthandler

import (
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	authconstant "github.com/DaoVuDat/trackpro-api/api/resource/auth/constant"
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	projectservice "github.com/DaoVuDat/trackpro-api/api/resource/project/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

func ListProject(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		onlyUidString := req.URL.Query().Get("onlyUid")

		onlyUid := false
		if onlyUiParsed, err := strconv.ParseBool(onlyUidString); err == nil {
			onlyUid = onlyUiParsed
		}

		// Get profile from access token
		accessTokenDetail := req.Context().Value(authconstant.AccessTokenContextHeader).(*jwt.TokenDetail)
		app.Logger.Debug().Msgf("accessTokenDetail.Role %v", accessTokenDetail.Role)
		app.Logger.Debug().Msgf("model.AccountType_Client.String() %v", model.AccountType_Client.String())
		app.Logger.Debug().Msgf("accessTokenDetail.UserId %v", accessTokenDetail.UserId)

		uid, err := uuid.Parse(accessTokenDetail.UserId)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		payment := false
		paymentString := req.URL.Query().Get("returnPayment")

		if parsedPayment, err := strconv.ParseBool(paymentString); err == nil {
			payment = parsedPayment
		}

		listProjectRepo := projectrepo.NewPostgresStore(app.Db)
		listProjectService := projectservice.NewListProjectService(listProjectRepo)

		var projects []projectdto.ProjectResponse

		if accessTokenDetail.Role != model.AccountType_Client.String() {
			projects, err = listProjectService.List(app, uid, onlyUid, payment)
			if err != nil {
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}
		} else {
			projects, err = listProjectService.ListByUID(app, uid, payment)
			if err != nil {
				app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
				return
			}
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"projects": projects,
		})

	}
}
