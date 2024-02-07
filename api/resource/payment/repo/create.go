package paymentrepo

import (
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	paymentdto "github.com/DaoVuDat/trackpro-api/api/resource/payment/dto"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type CreatePaymentRepo interface {
	Create(app *ctx.Application, projectId uuid.UUID, uid uuid.UUID, createPayment paymentdto.PaymentCreate) (*model.PaymentHistory, error)
}

func (store *postgresStore) Create(app *ctx.Application, projectId uuid.UUID, uid uuid.UUID, createPayment paymentdto.PaymentCreate) (*model.PaymentHistory, error) {
	paymentToCreate := model.PaymentHistory{
		ProjectID: projectId,
		Amount:    int32(createPayment.Amount),
		UserID:    uid,
	}

	stmt := PaymentHistory.
		INSERT(PaymentHistory.ProjectID, PaymentHistory.UserID, PaymentHistory.Amount).
		MODEL(paymentToCreate).
		RETURNING(PaymentHistory.Amount, PaymentHistory.ID, PaymentHistory.CreatedAt, PaymentHistory.ProjectID)

	var payment model.PaymentHistory

	err := stmt.Query(store.db, &payment)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.FailCreateError
		}
		return nil, err
	}

	return &payment, nil
}
