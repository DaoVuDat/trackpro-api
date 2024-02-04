package accountdto

import (
	"errors"
	"fmt"
	"github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/guregu/null.v4"
	"trackpro/api/model/project-management/public/model"
)

// This file for unmarshalling in handler

type AccountCreate struct {
	UserName string `json:"user_name"`
}

type AccountUpdate struct {
	Type   null.String `json:"type,omitempty"`
	Status null.String `json:"status,omitempty"`
}

func (accountCreate AccountCreate) Validate() error {
	return validation.ValidateStruct(&accountCreate,
		validation.Field(&accountCreate.UserName, validation.Required, validation.Length(8, 20)))
}

func (accountUpdate AccountUpdate) Validate() error {
	return validation.ValidateStruct(&accountUpdate,
		validation.Field(&accountUpdate.Type,
			validation.When(
				accountUpdate.Type.Valid,
				// custom validation
				validation.By(func(v interface{}) error {
					var accountType model.AccountType
					value := v.(null.String)
					err := accountType.Scan(value.String)
					if err != nil {
						return errors.New(fmt.Sprintf("type must be %s or %s\n", model.AccountType_Admin, model.AccountType_Client))
					}
					return nil
				}),
			),
		),
		validation.Field(&accountUpdate.Status,
			validation.When(
				accountUpdate.Status.Valid,
				// custom validation
				validation.By(func(v interface{}) error {
					var accountStatus model.AccountStatus
					value := v.(null.String)
					err := accountStatus.Scan(value.String)
					if err != nil {
						return errors.New(fmt.Sprintf("status must be %s or %s\n", model.AccountStatus_Pending, model.AccountStatus_Activated))
					}
					return nil
				}),
			),
		),
	)
}
