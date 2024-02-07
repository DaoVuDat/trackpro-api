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
	Find(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID, returnPayment bool) (*projectdto.ProjectQuery, error)
}

func (store *postgresStore) Find(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID, returnPayment bool) (*projectdto.ProjectQuery, error) {

	// Setup dynamic SELECT statement
	var selectColumns ProjectionList
	selectColumns = ProjectionList{
		Project.AllColumns.Except(Project.CreatedAt, Project.UpdatedAt),
		Account.Username,
	}

	// Setup dynamic FROM statement
	var fromStatement ReadableTable
	fromStatement = Project.INNER_JOIN(Account, Account.ID.EQ(Project.UserID))

	// Setup dynamic WHERE statement
	whereStatement := Bool(true)
	whereStatement = whereStatement.AND(Project.ID.EQ(UUID(projectId)))

	if userId != nil {
		whereStatement = whereStatement.AND(Project.UserID.EQ(UUID(*userId)))
	}

	if returnPayment {
		selectColumns = ProjectionList{
			Project.AllColumns.Except(Project.CreatedAt, Project.UpdatedAt),
			Account.Username,
			PaymentHistory.ID, PaymentHistory.Amount, PaymentHistory.CreatedAt,
		}
		fromStatement = fromStatement.LEFT_JOIN(PaymentHistory, Project.ID.EQ(PaymentHistory.ProjectID))
	}

	stmt := SELECT(selectColumns).
		FROM(fromStatement).
		WHERE(whereStatement)

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
