package accountservice

import (
	accountdto "github.com/DaoVuDat/trackpro-api/api/resource/account/dto"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
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
