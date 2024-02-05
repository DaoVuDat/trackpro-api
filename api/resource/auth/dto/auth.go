package authdto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type AuthSignUp struct {
	UserName  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type AuthLoginTemp struct {
	Token string `json:"token"`
}

type AuthLogin struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (authSignUp AuthSignUp) Validate() error {
	return validation.ValidateStruct(&authSignUp,
		validation.Field(&authSignUp.UserName,
			validation.Required,
			validation.Length(8, 20).Error("username must be large than 7 and less than 21 characters"),
		),
		validation.Field(&authSignUp.FirstName, validation.Required),
		validation.Field(&authSignUp.LastName, validation.Required),
	)
}
