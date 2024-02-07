package paymentdto

import (
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
)

type PaymentCreate struct {
	Amount int `json:"amount"`
}

func (paymentCreate PaymentCreate) Validate() error {
	return validation.ValidateStruct(&paymentCreate,
		validation.Field(&paymentCreate.Amount, validation.Required),
	)
}

type PaymentResponse struct {
	Id        int64     `json:"id"`
	ProjectId string    `json:"project_id,omitempty"`
	Amount    int       `json:"amount,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (paymentResponse *PaymentResponse) MapFromQuery(query model.PaymentHistory) {
	paymentResponse.Amount = int(query.Amount)
	paymentResponse.CreatedAt = query.CreatedAt
	paymentResponse.Id = query.ID
	if query.ProjectID == uuid.Nil {
		paymentResponse.ProjectId = ""
	} else {
		paymentResponse.ProjectId = query.ProjectID.String()
	}

}
