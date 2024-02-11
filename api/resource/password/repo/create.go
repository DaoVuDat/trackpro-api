package passwordrepo

import (
	"context"
	"database/sql"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	passworddto "github.com/DaoVuDat/trackpro-api/api/resource/password/dto"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
)

type CreatePasswordRepo interface {
	Create(app *ctx.Application, passwordCreate passworddto.PasswordCreate) (bool, error)
	CreateTX(app *ctx.Application, curCtx context.Context, tx *sql.Tx, passwordCreate passworddto.PasswordCreate) (bool, error)
}

func (store *postgresStore) Create(app *ctx.Application, passwordCreate passworddto.PasswordCreate) (bool, error) {
	passwordToCreate := model.Password{
		UserID:       passwordCreate.UserId,
		HashPassword: passwordCreate.HashedPassword,
	}

	stmt := Password.
		INSERT(Password.AllColumns).
		MODEL(passwordToCreate)

	_, err := stmt.Exec(store.db)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (store *postgresStore) CreateTX(app *ctx.Application, curCtx context.Context, tx *sql.Tx, passwordCreate passworddto.PasswordCreate) (bool, error) {
	passwordToCreate := model.Password{
		UserID:       passwordCreate.UserId,
		HashPassword: passwordCreate.HashedPassword,
	}

	stmt := Password.
		INSERT(Password.AllColumns).
		MODEL(passwordToCreate)

	_, err := stmt.ExecContext(curCtx, tx)
	if err != nil {
		return false, err
	}

	return true, nil
}
