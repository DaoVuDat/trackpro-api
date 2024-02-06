package projectrepo

import (
	"errors"
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type FindProjectRepo interface {
	Find(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID) (*projectdto.ProjectQuery, error)
}

func (store *postgresStore) Find(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID) (*projectdto.ProjectQuery, error) {
	condition := Bool(true)
	condition = condition.AND(Project.ID.EQ(UUID(projectId)))

	if userId != nil {
		condition = condition.AND(Project.UserID.EQ(UUID(*userId)))
	}

	stmt := SELECT(Project.AllColumns.Except(Project.CreatedAt, Project.UpdatedAt), Account.Username).
		FROM(Project.INNER_JOIN(Account, Account.ID.EQ(Project.UserID))).
		WHERE(condition)

	var project projectdto.ProjectQuery

	err := stmt.Query(store.db, &project)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.QueryNoResultErr
		}
		return nil, err
	}

	return &project, nil
}
