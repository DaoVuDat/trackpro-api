package profiledto

import (
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	"github.com/DaoVuDat/trackpro-api/util/regex"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type ProfileCreate struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
}

type ProfileUpdate struct {
	FirstName   null.String `json:"first_name"`
	LastName    null.String `json:"last_name"`
	ImageUrl    null.String `json:"image_url"`
	About       null.String `json:"about"`
	PhoneNumber null.String `json:"phone_number"`
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

					if regex.EmailRegex.MatchString(v.String) {
						return nil
					}

					return errors.New("invalid pattern")
				}),
			),
		),
		validation.Field(&profileUpdate.About,
			validation.When(
				profileUpdate.About.Valid),
			validation.By(func(value interface{}) error {
				v := value.(null.String)
				if len(v.String) < 1 {
					return errors.New("must be larger than 0 character")
				}

				return nil
			}),
		),
		validation.Field(&profileUpdate.PhoneNumber,
			validation.When(
				profileUpdate.PhoneNumber.Valid),
			validation.By(func(value interface{}) error {
				v := value.(null.String)
				if len(v.String) < 1 {
					return errors.New("must be larger than 0 character")
				}

				return nil
			}),
		),
	)
}

type ProfileResponse struct {
	UserId      string  `json:"user_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	ImageUrl    *string `json:"image_url,omitempty"`
	About       *string `json:"about,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

func (profileResponse *ProfileResponse) MapFromQuery(query model.Profile) {
	profileResponse.UserId = query.UserID.String()
	profileResponse.FirstName = query.FirstName
	profileResponse.LastName = query.LastName
	profileResponse.About = query.About
	profileResponse.ImageUrl = query.ImageURL
	profileResponse.PhoneNumber = query.Phone
}
