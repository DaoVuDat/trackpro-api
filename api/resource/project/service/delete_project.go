package projectservice

import (
	projectrepo "github.com/DaoVuDat/trackpro-api/api/resource/project/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/google/uuid"
)

type DeleteProjectService interface {
	Delete(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID) error
}

type deleteProjectService struct {
	deleteProjectRepo projectrepo.DeleteProjectRepo
}

func NewDeleteProjectService(deleteProjectRepo projectrepo.DeleteProjectRepo) DeleteProjectService {
	return &deleteProjectService{
		deleteProjectRepo: deleteProjectRepo,
	}
}

func (service *deleteProjectService) Delete(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID) error {
	err := service.deleteProjectRepo.Delete(app, projectId, userId)

	return err
}
