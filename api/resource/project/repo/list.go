package projectrepo

import (
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type ListProjectRepo interface {
	List(app *ctx.Application, UID uuid.UUID, withUid bool, returnPayment bool) ([]projectdto.ProjectQuery, error)
	ListByUID(app *ctx.Application, userId uuid.UUID, returnPayment bool) ([]projectdto.ProjectQuery, error)
}

func (store *postgresStore) List(app *ctx.Application, UID uuid.UUID, withUid bool, returnPayment bool) ([]projectdto.ProjectQuery, error) {
	var selectColumns ProjectionList
	selectColumns = ProjectionList{
		Project.AllColumns.Except(Project.CreatedAt, Project.UpdatedAt),
		Account.Username,
	}

	var fromStatement ReadableTable
	fromStatement = Project.INNER_JOIN(Account, Account.ID.EQ(Project.UserID))

	whereStatement := Bool(true)
	if withUid {
		whereStatement = whereStatement.AND(Project.UserID.EQ(UUID(UID)))
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

	var projects []projectdto.ProjectQuery

	err := stmt.Query(store.db, &projects)
	if err != nil {
		return nil, err
	}

	return projects, err
}

func (store *postgresStore) ListByUID(app *ctx.Application, userId uuid.UUID, returnPayment bool) ([]projectdto.ProjectQuery, error) {
	var selectColumns ProjectionList
	selectColumns = ProjectionList{
		Project.AllColumns.Except(Project.CreatedAt, Project.UpdatedAt),
		Account.Username,
	}

	var fromStatement ReadableTable
	fromStatement = Project.INNER_JOIN(Account, Account.ID.EQ(Project.UserID))

	whereStatement := Bool(true)
	whereStatement = whereStatement.AND(Project.UserID.EQ(UUID(userId)))

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

	var projects []projectdto.ProjectQuery

	err := stmt.Query(store.db, &projects)
	if err != nil {
		return nil, err
	}

	return projects, err
}
