package projectrepo

import (
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type ListProjectRepo interface {
	List(app *ctx.Application, userId uuid.UUID) ([]projectdto.ProjectQuery, error)
}

func (store *postgresStore) List(app *ctx.Application, userId uuid.UUID) ([]projectdto.ProjectQuery, error) {
	stmt := SELECT(Project.AllColumns.Except(Project.CreatedAt, Project.UpdatedAt), Account.Username).
		FROM(Project.INNER_JOIN(Account, Account.ID.EQ(Project.UserID))).
		WHERE(Project.UserID.EQ(UUID(userId)))

	var projects []projectdto.ProjectQuery

	err := stmt.Query(store.db, &projects)
	if err != nil {
		return nil, err
	}

	return projects, err
}
