package projectrepo

import (
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type DeleteProjectRepo interface {
	Delete(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID) error
}

func (store *postgresStore) Delete(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID) error {
	condition := Bool(true)
	condition = condition.AND(Project.ID.EQ(UUID(projectId)))

	if userId != nil {
		condition = condition.AND(Project.UserID.EQ(UUID(*userId)))
	}

	stmt := Project.
		DELETE().
		WHERE(condition).RETURNING(Project.ID)

	query, _ := stmt.Sql()
	app.Logger.Debug().Msg(query)

	var project model.Project

	err := stmt.Query(store.db, &project)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return common.FailDeleteError
		}

		return err
	}

	if _, err = uuid.Parse(project.ID.String()); err != nil {
		return common.FailDeleteError
	}

	return nil
}
