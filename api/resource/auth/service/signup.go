package signupservice

import (
	accountrepo "trackpro/api/resource/account/repo"
	authdto "trackpro/api/resource/auth/dto"
	profilerepo "trackpro/api/resource/profile/repo"
	"trackpro/util/ctx"
)

type SignUpService interface {
	SignUp(app *ctx.Application, authSignUp authdto.AuthSignUp) error
}

type signupService struct {
	accountCreateRepo accountrepo.CreateAccountRepo
	profileCreateRepo profilerepo.CreateProfileRepo
}

func NewSignUpService(
	accountCreateRepo accountrepo.CreateAccountRepo,
	profileCreateRepo profilerepo.CreateProfileRepo,
) SignUpService {
	return &signupService{
		accountCreateRepo: accountCreateRepo,
		profileCreateRepo: profileCreateRepo,
	}
}

func (service *signupService) SignUp(app *ctx.Application, authSignUp authdto.AuthSignUp) error {
	// open transaction for creating account then profile

	panic(1)
}
