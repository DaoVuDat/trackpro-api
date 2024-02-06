package accountservice

import (
	accountdto "github.com/DaoVuDat/trackpro-api/api/resource/account/dto"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
)

type ListAccountService interface {
	List(app *ctx.Application) ([]accountdto.AccountResponse, error)
}

type listAccountService struct {
	listAccountRepo accountrepo.ListAccountRepo
}

func NewListAccountService(listAccountRepo accountrepo.ListAccountRepo) ListAccountService {
	return &listAccountService{
		listAccountRepo: listAccountRepo,
	}
}

func (service *listAccountService) List(app *ctx.Application) ([]accountdto.AccountResponse, error) {
	accounts, err := service.listAccountRepo.List(app)
	if err != nil {
		return nil, err
	}

	accountsResponse := make([]accountdto.AccountResponse, len(accounts))
	for i, account := range accounts {
		accountsResponse[i].Type = account.Type.String()
		accountsResponse[i].UserId = account.ID.String()
		accountsResponse[i].Status = account.Status.String()
		accountsResponse[i].UserName = account.Username
	}

	return accountsResponse, nil
}
