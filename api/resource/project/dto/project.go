package projectdto

import (
	"errors"
	"fmt"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	paymentdto "github.com/DaoVuDat/trackpro-api/api/resource/payment/dto"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"time"
)

type ProjectCreate struct {
	UserId      string      `json:"user_id"`
	Name        null.String `json:"name"`
	Description null.String `json:"description"`
	Price       null.Int    `json:"price"`
	StartTime   null.Time   `json:"start_time"`
	EndTime     null.Time   `json:"end_time"`
}

func (projectCreate ProjectCreate) Validate() error {
	return validation.ValidateStruct(&projectCreate,
		validation.Field(&projectCreate.UserId,
			validation.Required,
			validation.By(func(v interface{}) error {
				value := v.(string)
				_, err := uuid.Parse(value)
				if err != nil {
					return errors.New("must be valid uuid")
				}
				return nil
			}),
		),
		validation.Field(&projectCreate.Price,
			validation.When(projectCreate.Price.Valid,
				validation.Min(100),
			),
		),
	)
}

type ProjectUpdate struct {
	Name        null.String `json:"name"`
	Description null.String `json:"description"`
	Price       null.Int    `json:"price"`
	StartTime   null.Time   `json:"start_time,omitempty"`
	EndTime     null.Time   `json:"end_time,omitempty"`
	Status      null.String `json:"status"`
}

func (projectUpdate ProjectUpdate) Validate() error {
	return validation.ValidateStruct(&projectUpdate,
		validation.Field(&projectUpdate.Price,
			validation.When(projectUpdate.Price.Valid,
				validation.Min(100),
			),
		),
		validation.Field(&projectUpdate.Name,
			validation.When(projectUpdate.Name.Valid,
				validation.By(func(v interface{}) error {
					value := v.(null.String)
					if len(value.String) < 1 {
						return errors.New("must larger than 0 character")
					}
					return nil
				}),
			),
		),
		validation.Field(&projectUpdate.Description,
			validation.When(projectUpdate.Price.Valid,
				validation.By(func(v interface{}) error {
					value := v.(null.String)
					if len(value.String) < 1 {
						return errors.New("must larger than 0 character")
					}
					return nil
				}),
			),
		),
		validation.Field(&projectUpdate.Status,
			validation.When(projectUpdate.Status.Valid,
				validation.By(func(v interface{}) error {
					value := v.(null.String)

					var projectStatus model.ProjectStatus
					err := projectStatus.Scan(value.String)
					if err != nil {
						return errors.New(fmt.Sprintf("must be %s, %s, or %s\n",
							model.ProjectStatus_Registering,
							model.ProjectStatus_Progressing,
							model.ProjectStatus_Finished))
					}
					return nil
				}),
			),
		),
	)
}

type ProjectQuery struct {
	model.Project
	model.Account
	Payments []model.PaymentHistory
}

type ProjectResponse struct {
	Id          string                       `json:"id"`
	UserId      string                       `json:"user_id"`
	UserName    string                       `json:"username"`
	ProjectName *string                      `json:"project_name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	Price       *int                         `json:"price,omitempty"`
	Status      string                       `json:"status"`
	StartTime   *time.Time                   `json:"start_time,omitempty"`
	EndTime     *time.Time                   `json:"end_time,omitempty"`
	Payment     []paymentdto.PaymentResponse `json:"payment,omitempty"`
}

func (project *ProjectResponse) MapFromProjectQuery(query ProjectQuery) {
	project.Id = query.Project.ID.String()
	project.ProjectName = query.Name
	project.Description = query.Description
	project.UserId = query.UserID.String()
	project.Status = query.Project.Status.String()
	project.UserName = query.Username
	price := int(*query.Price)
	project.Price = &price
	project.StartTime = query.StartTime
	project.EndTime = query.EndTime

	if len(query.Payments) > 0 {
		payments := make([]paymentdto.PaymentResponse, len(query.Payments))
		for i, paymentQuery := range query.Payments {
			var payment paymentdto.PaymentResponse
			payment.MapFromQuery(paymentQuery)
			payments[i] = payment
		}
		project.Payment = payments
	} else {
		project.Payment = make([]paymentdto.PaymentResponse, 0)
	}
}
