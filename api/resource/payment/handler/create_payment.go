package paymenthandler

import (
	"encoding/json"
	paymentdto "github.com/DaoVuDat/trackpro-api/api/resource/payment/dto"
	paymentrepo "github.com/DaoVuDat/trackpro-api/api/resource/payment/repo"
	paymentservice "github.com/DaoVuDat/trackpro-api/api/resource/payment/service"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
	"net/http"
)

func CreatePayment(app *ctx.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		projectIdString := req.URL.Query().Get("pid")
		userIdString := req.URL.Query().Get("uid")

		// parsed projectId and userId in uuid
		projectId, err := uuid.Parse(projectIdString)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		userId, err := uuid.Parse(userIdString)
		if err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		var createPaymentRequest paymentdto.PaymentCreate
		if err = json.NewDecoder(req.Body).Decode(&createPaymentRequest); err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		if err = createPaymentRequest.Validate(); err != nil {
			app.Render.JSON(w, http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}

		createPaymentRepo := paymentrepo.NewPostgresStore(app.Db)
		createPaymentService := paymentservice.NewCreatePaymentService(createPaymentRepo)

		payment, err := createPaymentService.Create(app, projectId, userId, createPaymentRequest)
		if err != nil {
			app.Render.JSON(w, http.StatusInternalServerError, common.InternalErrorResponse(err))
			return
		}

		app.Render.JSON(w, http.StatusOK, map[string]interface{}{
			"payment": payment,
		})
	}
}
