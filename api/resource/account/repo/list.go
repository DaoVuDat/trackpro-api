package accountrepo

import (
	. "github.com/go-jet/jet/v2/postgres"
	"trackpro/api/model/project-management/public/model"
	. "trackpro/api/model/project-management/public/table"
	"trackpro/util/ctx"
)

type ListAccountRepo interface {
	List(app *ctx.Application) ([]model.Account, error)
}

func (store *postgresStore) List(app *ctx.Application) ([]model.Account, error) {
	// query
	stmt := SELECT(Account.AllColumns.Except(Account.CreatedAt, Account.UpdatedAt)).
		FROM(Account)

	var accounts []model.Account
	err := stmt.Query(store.db, &accounts)
	if err != nil {
		return nil, err
	}

	return accounts, err
}
