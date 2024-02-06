package profilerepo

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

type FindProfileRepo interface {
	Find(app *ctx.Application, accountId uuid.UUID) (*model.Profile, error)
}

func (store *postgresStore) Find(app *ctx.Application, accountId uuid.UUID) (*model.Profile, error) {
	// query
	stmt := SELECT(Profile.AllColumns.Except(Profile.CreatedAt, Profile.UpdatedAt)).
		FROM(Profile).
		WHERE(Profile.UserID.EQ(UUID(accountId)))

	var profile model.Profile
	err := stmt.Query(store.db, &profile)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.QueryNoResultErr
		}
	}

	return &profile, nil
}
