package accountrepo

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

type FindAccountRepo interface {
	FindById(app *ctx.Application, accountId uuid.UUID) (*model.Account, error)
	FindByUserName(app *ctx.Application, username string) (*model.Account, error)
}

func (store *postgresStore) FindById(app *ctx.Application, accountId uuid.UUID) (*model.Account, error) {

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

func (store *postgresStore) FindByUserName(app *ctx.Application, username string) (*model.Account, error) {
	// query
	stmt := SELECT(Account.AllColumns.Except(Account.CreatedAt, Account.UpdatedAt)).
		FROM(Account).
		WHERE(Account.Username.EQ(String(username))).LIMIT(1)

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
