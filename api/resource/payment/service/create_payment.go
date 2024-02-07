package paymentservice

import (
	paymentdto "github.com/DaoVuDat/trackpro-api/api/resource/payment/dto"
	paymentrepo "github.com/DaoVuDat/trackpro-api/api/resource/payment/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
)

type CreatePaymentService interface {
	Create(app *ctx.Application, projectId uuid.UUID, userId uuid.UUID, paymentCreate paymentdto.PaymentCreate) (*paymentdto.PaymentResponse, error)
}

type createPaymentService struct {
	createPaymentRepo paymentrepo.CreatePaymentRepo
}

func NewCreatePaymentService(createPaymentRepo paymentrepo.CreatePaymentRepo) CreatePaymentService {
	return &createPaymentService{
		createPaymentRepo: createPaymentRepo,
	}
}

func (service *createPaymentService) Create(app *ctx.Application, projectId uuid.UUID, userId uuid.UUID, paymentCreate paymentdto.PaymentCreate) (*paymentdto.PaymentResponse, error) {
	payment, err := service.createPaymentRepo.Create(app, projectId, userId, paymentCreate)
	if err != nil {
		return nil, err
	}

	var paymentResponse paymentdto.PaymentResponse

	paymentResponse.MapFromQuery(*payment)

	return &paymentResponse, nil
}
