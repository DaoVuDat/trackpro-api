package profiledto

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
	"regexp"
)

type ProfileCreate struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
}

type ProfileUpdate struct {
	FirstName null.String `json:"first_name"`
	LastName  null.String `json:"last_name"`
	ImageUrl  null.String `json:"image_url"`
}

func (profileUpdate ProfileUpdate) Validate() error {
	return validation.ValidateStruct(&profileUpdate,
		validation.Field(&profileUpdate.FirstName,
			validation.When(
				profileUpdate.FirstName.Valid,
				validation.By(func(value interface{}) error {
					v := value.(null.String)
					if len(v.String) < 1 {
						return errors.New("must be larger than 0 character")
					}

					return nil
				}),
			),
		),
		validation.Field(&profileUpdate.LastName,
			validation.When(
				profileUpdate.LastName.Valid,
				validation.By(func(value interface{}) error {
					v := value.(null.String)
					if len(v.String) < 1 {
						return errors.New("must be larger than 0 character")
					}

					return nil
				}),
			),
		),
		validation.Field(&profileUpdate.ImageUrl,
			validation.When(
				profileUpdate.ImageUrl.Valid,
				validation.By(func(value interface{}) error {
					v := value.(null.String)

					if ok, err := regexp.Match(
						`(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?\/[a-zA-Z0-9]{2,}|((https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?)|(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}(\.[a-zA-Z0-9]{2,})?`,
						[]byte(v.String)); !ok {
						return errors.New("invalid pattern")
					} else if err != nil {
						return err
					}

					return nil
				}),
			),
		),
	)
}

type ProfileResponse struct {
	UserId    string  `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	ImageUrl  *string `json:"image_url,omitempty"`
}
