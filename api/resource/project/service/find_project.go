package projectservice

import (
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
)

type FindProjectService interface {
	Find(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID) (*projectdto.ProjectResponse, error)
}

type findProjectService struct {
	findProjectRepo projectrepo.FindProjectRepo
}

func NewFindProjectService(findProojectRepo projectrepo.FindProjectRepo) FindProjectService {
	return &findProjectService{
		findProjectRepo: findProojectRepo,
	}
}

func (service *findProjectService) Find(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID) (*projectdto.ProjectResponse, error) {
	project, err := service.findProjectRepo.Find(app, projectId, userId)
	if err != nil {
		return nil, err
	}

	var projectResponse projectdto.ProjectResponse
	projectResponse.MapFromProjectQuery(*project)

	return &projectResponse, nil
}
