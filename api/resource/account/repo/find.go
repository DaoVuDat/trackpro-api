package accountrepo

import (
	"errors"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"trackpro/api/model/project-management/public/model"
	. "trackpro/api/model/project-management/public/table"
	"trackpro/api/router/common"
	"trackpro/util/ctx"
)

type FindAccountRepo interface {
	Find(app *ctx.Application, accountId uuid.UUID) (*model.Account, error)
}

func (store *postgresStore) Find(app *ctx.Application, accountId uuid.UUID) (*model.Account, error) {

	// query
	stmt := SELECT(Account.AllColumns.Except(Account.CreatedAt, Account.UpdatedAt)).
		FROM(Account).
		WHERE(Account.ID.EQ(UUID(accountId))).LIMIT(1)

	var account model.Account
	err := stmt.Query(store.db, &account)
	if err != nil {
		app.Logger.Error().Err(err)
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.QueryNoResultErr
		}

		return nil, err
	}

	return &account, nil
}
