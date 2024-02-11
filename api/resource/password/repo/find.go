package passwordrepo

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

type FindPasswordRepo interface {
	Find(app *ctx.Application, userId uuid.UUID) (string, error)
}

func (store *postgresStore) Find(app *ctx.Application, userId uuid.UUID) (string, error) {

	stmt := SELECT(Password.HashPassword).
		FROM(Password).
		WHERE(Password.UserID.EQ(UUID(userId)))

	var password model.Password

	err := stmt.Query(store.db, &password)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return "", common.QueryNoResultErr
		}
		return "", err
	}

	return password.HashPassword, nil
}
