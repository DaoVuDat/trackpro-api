package accountservice

import (
	"trackpro/api/model/project-management/public/model"
	accountdto "trackpro/api/resource/account/dto"
	accountrepo "trackpro/api/resource/account/repo"
	"trackpro/util/ctx"
)

type UpdateAccountService interface {
	Update(app *ctx.Application, accountId string, accountUpdate accountdto.AccountUpdate) (*model.Account, error)
}

type updateAccountService struct {
	updateAccountRepo accountrepo.UpdateAccountRepo
}

func NewUpdateAccountService(updateAccountRepo accountrepo.UpdateAccountRepo) UpdateAccountService {
	return &updateAccountService{
		updateAccountRepo: updateAccountRepo,
	}
}

func (service *updateAccountService) Update(app *ctx.Application, accountId string, accountUpdate accountdto.AccountUpdate) (*model.Account, error) {
	updatedAccount, err := service.updateAccountRepo.Update(app, accountId, accountUpdate)
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}
