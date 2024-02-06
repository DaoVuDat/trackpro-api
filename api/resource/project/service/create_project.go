package projectservice

import (
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
)

type CreateProjectService interface {
	Create(app *ctx.Application, createProject projectdto.ProjectCreate) (*projectdto.ProjectResponse, error)
}

type createProjectService struct {
	createProjectRepo projectrepo.CreateProjectRepo
}

func NewCreateProjectService(createProjectRepo projectrepo.CreateProjectRepo) CreateProjectService {
	return &createProjectService{
		createProjectRepo: createProjectRepo,
	}
}

func (service *createProjectService) Create(app *ctx.Application, createProject projectdto.ProjectCreate) (*projectdto.ProjectResponse, error) {
	project, err := service.createProjectRepo.Create(app, createProject)
	if err != nil {
		return nil, err
	}

	var projectResponse projectdto.ProjectResponse
	projectResponse.MapFromProjectQuery(*project)

	return &projectResponse, nil
}
