package accountrepo

import (
	"trackpro/api/model/project-management/public/model"
	. "trackpro/api/model/project-management/public/table"
	accountdto "trackpro/api/resource/account/dto"

	"trackpro/util/ctx"
)

type CreateAccountRepo interface {
	Create(application *ctx.Application, create accountdto.AccountCreate) (*model.Account, error)
}

func (store *postgresStore) Create(app *ctx.Application, create accountdto.AccountCreate) (*model.Account, error) {
	account := model.Account{
		Username: create.UserName,
	}

	stmt := Account.INSERT(Account.Username).MODEL(account).RETURNING(Account.AllColumns)

	var dest model.Account

	err := stmt.Query(store.db, &dest)
	if err != nil {
		app.Logger.Error().Err(err)
		return nil, err
	}

	app.Logger.Debug().Any("Result", dest).Send()

	return &account, nil
}