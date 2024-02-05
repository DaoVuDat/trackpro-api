package profilerepo

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
