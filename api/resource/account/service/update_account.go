package accountservice

import (
	accountdto "github.com/DaoVuDat/trackpro-api/api/resource/account/dto"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
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
