package authdto

import (
	"github.com/DaoVuDat/trackpro-api/util/regex"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/errgo.v2/errors"
)

type AuthSignUp struct {
	UserName          string `json:"username"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

func (authSignUp AuthSignUp) Validate() error {
	return validation.ValidateStruct(&authSignUp,
		validation.Field(&authSignUp.UserName,
			validation.Required,
			validation.Length(8, 20).Error("username must be large than 7 and less than 21 characters"),
		),
		validation.Field(&authSignUp.FirstName, validation.Required),
		validation.Field(&authSignUp.LastName, validation.Required),
		validation.Field(&authSignUp.Password,
			validation.Required,
			validation.By(func(value interface{}) error {
				v := value.(string)
				if regex.PasswordRegex.MatchString(v) {
					return nil
				}
				return errors.New("invalid password")
			})),
		validation.Field(&authSignUp.ConfirmedPassword,
			validation.Required,
			validation.By(func(value interface{}) error {
				v := value.(string)
				if regex.PasswordRegex.MatchString(v) {
					return nil
				}
				return errors.New("invalid confirmed password")
			})),
	)
}

type AuthLoginTemp struct {
	Token string `json:"token"`
}

type AuthLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (authLogin AuthLogin) Validate() error {
	return validation.ValidateStruct(&authLogin,
		validation.Field(&authLogin.UserName,
			validation.Required,
			validation.Length(8, 20).Error("username must be large than 7 and less than 21 characters"),
		),
		validation.Field(&authLogin.Password,
			validation.Required,
			validation.By(func(value interface{}) error {
				v := value.(string)
				if regex.PasswordRegex.MatchString(v) {
					return nil
				}
				return errors.New("invalid password")
			})),
	)
}
