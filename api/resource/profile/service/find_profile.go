package profileservice

import (
	"github.com/google/uuid"
	profiledto "trackpro/api/resource/profile/dto"
	profilerepo "trackpro/api/resource/profile/repo"
	"trackpro/util/ctx"
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

	profileResponse := &profiledto.ProfileResponse{
		UserId:    profile.UserID.String(),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		ImageUrl:  profile.ImageURL,
	}

	return profileResponse, nil
}
