package accountservice

import (
	"github.com/google/uuid"
	accountdto "trackpro/api/resource/account/dto"
	accountrepo "trackpro/api/resource/account/repo"
	"trackpro/util/ctx"
)

type FindAccountService interface {
	Find(app *ctx.Application, accountId uuid.UUID) (*accountdto.AccountResponse, error)
}

type findAccountService struct {
	findAccountRepo accountrepo.FindAccountRepo
}

func NewFindAccountService(findAccountRepo accountrepo.FindAccountRepo) FindAccountService {
	return &findAccountService{
		findAccountRepo: findAccountRepo,
	}
}

func (service *findAccountService) Find(app *ctx.Application, accountId uuid.UUID) (*accountdto.AccountResponse, error) {
	account, err := service.findAccountRepo.Find(app, accountId)
	if err != nil {
		return nil, err
	}

	accountResponse := &accountdto.AccountResponse{
		UserId:   account.ID.String(),
		UserName: account.Username,
		Type:     account.Type.String(),
		Status:   account.Status.String(),
	}

	return accountResponse, nil
}
