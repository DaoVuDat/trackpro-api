package projectservice

import (
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
)

type UpdateProjectService interface {
	Update(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID, updateProject projectdto.ProjectUpdate) (*projectdto.ProjectResponse, error)
}

type updateProjectService struct {
	updateProjectRepo projectrepo.UpdateProjectRepo
}

func NewUpdateProjectService(updateProjectRepo projectrepo.UpdateProjectRepo) UpdateProjectService {
	return &updateProjectService{
		updateProjectRepo: updateProjectRepo,
	}
}

func (service *updateProjectService) Update(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID, updateProject projectdto.ProjectUpdate) (*projectdto.ProjectResponse, error) {
	updatedProject, err := service.updateProjectRepo.Update(app, projectId, userId, updateProject)
	if err != nil {
		return nil, err
	}

	var projectResponse projectdto.ProjectResponse
	projectResponse.MapFromProjectQuery(*updatedProject)

	return &projectResponse, nil
}
