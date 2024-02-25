package profileservice

import (
	profiledto "github.com/DaoVuDat/trackpro-api/api/resource/profile/dto"
	profilerepo "github.com/DaoVuDat/trackpro-api/api/resource/profile/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
)

type FindProfileService interface {
	Find(app *ctx.Application, accountId uuid.UUID) (*profiledto.ProfileResponse, error)
}

type findProfileService struct {
	findProfileRepo profilerepo.FindProfileRepo
}

func NewFindProfileService(findProfileRepo profilerepo.FindProfileRepo) FindProfileService {
	return &findProfileService{
		findProfileRepo: findProfileRepo,
	}
}

func (service *findProfileService) Find(app *ctx.Application, accountId uuid.UUID) (*profiledto.ProfileResponse, error) {
	profile, err := service.findProfileRepo.Find(app, accountId)
	if err != nil {
		return nil, err
	}

	var profileResponse profiledto.ProfileResponse
	profileResponse.MapFromQuery(*profile)

	return &profileResponse, nil
}
