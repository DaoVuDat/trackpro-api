package projectservice

import (
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
)

type ListProjectService interface {
	List(app *ctx.Application, userId uuid.UUID, onlyUID bool, returnPayment bool) ([]projectdto.ProjectResponse, error)
	ListByUID(app *ctx.Application, userId uuid.UUID, returnPayment bool) ([]projectdto.ProjectResponse, error)
}

type listProjectService struct {
	listProjectRepo projectrepo.ListProjectRepo
}

func NewListProjectService(listProjectRepo projectrepo.ListProjectRepo) ListProjectService {
	return &listProjectService{
		listProjectRepo: listProjectRepo,
	}
}

func (service *listProjectService) List(app *ctx.Application, userId uuid.UUID, onlyUID bool, returnPayment bool) ([]projectdto.ProjectResponse, error) {
	projects, err := service.listProjectRepo.List(app, userId, onlyUID, returnPayment)
	if err != nil {
		return nil, err
	}

	projectsResponse := make([]projectdto.ProjectResponse, len(projects))
	for i, project := range projects {
		var projectResponse projectdto.ProjectResponse
		projectResponse.MapFromProjectQuery(project)
		projectsResponse[i] = projectResponse
	}

	return projectsResponse, nil
}

func (service *listProjectService) ListByUID(app *ctx.Application, userId uuid.UUID, returnPayment bool) ([]projectdto.ProjectResponse, error) {
	projects, err := service.listProjectRepo.ListByUID(app, userId, returnPayment)
	if err != nil {
		return nil, err
	}

	projectsResponse := make([]projectdto.ProjectResponse, len(projects))
	for i, project := range projects {
		var projectResponse projectdto.ProjectResponse
		projectResponse.MapFromProjectQuery(project)
		projectsResponse[i] = projectResponse
	}

	return projectsResponse, nil
}
