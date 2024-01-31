package accountservice

import (
	accountdto "trackpro/api/resource/account/dto"
	accountrepo "trackpro/api/resource/account/repo"

	"trackpro/util/ctx"
)

type CreateAccountService interface {
	CreateAccount(application *ctx.Application, create accountdto.AccountCreate) error
}

type createAccountService struct {
	createAccountRepo accountrepo.CreateAccountRepo
}

func NewCreateAccountService(createAccountRepo accountrepo.CreateAccountRepo) CreateAccountService {
	return &createAccountService{
		createAccountRepo: createAccountRepo,
	}
}

func (service *createAccountService) CreateAccount(application *ctx.Application, create accountdto.AccountCreate) error {
	application.Logger.Debug().Msg("Create Account")

	_, err := service.createAccountRepo.Create(application, create)
	if err != nil {
		return err
	}

	return nil
}
