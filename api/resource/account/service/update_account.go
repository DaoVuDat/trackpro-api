package accountservice

import (
	"github.com/google/uuid"
	accountdto "trackpro/api/resource/account/dto"
	accountrepo "trackpro/api/resource/account/repo"
	"trackpro/util/ctx"
)

type UpdateAccountService interface {
	Update(app *ctx.Application, accountId uuid.UUID, accountUpdate accountdto.AccountUpdate) (*accountdto.AccountResponse, error)
}

type updateAccountService struct {
	updateAccountRepo accountrepo.UpdateAccountRepo
}

func NewUpdateAccountService(updateAccountRepo accountrepo.UpdateAccountRepo) UpdateAccountService {
	return &updateAccountService{
		updateAccountRepo: updateAccountRepo,
	}
}

func (service *updateAccountService) Update(app *ctx.Application, accountId uuid.UUID, accountUpdate accountdto.AccountUpdate) (*accountdto.AccountResponse, error) {
	updatedAccount, err := service.updateAccountRepo.Update(app, accountId, accountUpdate)
	if err != nil {
		return nil, err
	}

	accountResponse := &accountdto.AccountResponse{
		UserId:   updatedAccount.ID.String(),
		UserName: updatedAccount.Username,
		Type:     updatedAccount.Type.String(),
		Status:   updatedAccount.Status.String(),
	}

	return accountResponse, nil
}
