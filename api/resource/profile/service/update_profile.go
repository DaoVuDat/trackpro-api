package profileservice

import (
	profiledto "github.com/DaoVuDat/trackpro-api/api/resource/profile/dto"
	profilerepo "github.com/DaoVuDat/trackpro-api/api/resource/profile/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
)

type UpdateProfileService interface {
	Update(app *ctx.Application, accountId uuid.UUID, update profiledto.ProfileUpdate) (*profiledto.ProfileResponse, error)
}

type updateProfileService struct {
	updateProfileRepo profilerepo.UpdateProfileRepo
}

func NewUpdateProfileService(updateProfileRepo profilerepo.UpdateProfileRepo) UpdateProfileService {
	return &updateProfileService{
		updateProfileRepo: updateProfileRepo,
	}
}

func (service *updateProfileService) Update(app *ctx.Application, accountId uuid.UUID, update profiledto.ProfileUpdate) (*profiledto.ProfileResponse, error) {
	updatedProfile, err := service.updateProfileRepo.Update(app, accountId, update)
	if err != nil {
		return nil, err
	}

	profileResponse := profiledto.ProfileResponse{
		UserId:    updatedProfile.UserID.String(),
		FirstName: updatedProfile.FirstName,
		LastName:  updatedProfile.LastName,
		ImageUrl:  updatedProfile.ImageURL,
	}

	return &profileResponse, nil
}
