package accountrepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	accountdto "github.com/DaoVuDat/trackpro-api/api/resource/account/dto"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-jet/jet/v2/qrm"
	"strings"
)

type CreateAccountRepo interface {
	Create(application *ctx.Application, create accountdto.AccountCreate) (*model.Account, error)
	CreateTX(application *ctx.Application, curCtx context.Context, tx *sql.Tx, create accountdto.AccountCreate) (*model.Account, error)
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

func (store *postgresStore) CreateTX(app *ctx.Application, curCtx context.Context, tx *sql.Tx, create accountdto.AccountCreate) (*model.Account, error) {
	accountToCreate := model.Account{
		Username: create.UserName,
	}

	stmt := Account.INSERT(Account.Username).
		MODEL(accountToCreate).
		RETURNING(Account.AllColumns.
			Except(Account.CreatedAt, Account.UpdatedAt))

	var dest model.Account

	err := stmt.QueryContext(curCtx, tx, &dest)
	if err != nil {
		app.Logger.Error().Err(err).Send()

		if strings.Contains(err.Error(), "(SQLSTATE 23505)") {
			return nil, common.DuplicateValueError
		}

		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.FailCreateError
		}
		return nil, err
	}

	app.Logger.Debug().Any("acc created", dest).Send()

	return &dest, nil
}
