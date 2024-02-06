package profilerepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	profiledto "github.com/DaoVuDat/trackpro-api/api/resource/profile/dto"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/go-jet/jet/v2/qrm"
)

type CreateProfileRepo interface {
	Create(app *ctx.Application, create profiledto.ProfileCreate) (*model.Profile, error)
	CreateTX(app *ctx.Application, curCtx context.Context, tx *sql.Tx, create profiledto.ProfileCreate) (*model.Profile, error)
}

func (store *postgresStore) Create(app *ctx.Application, create profiledto.ProfileCreate) (*model.Profile, error) {
	profile := model.Profile{
		UserID:    create.UserID,
		FirstName: create.FirstName,
		LastName:  create.LastName,
	}

	stmt := Profile.
		INSERT(Profile.UserID, Profile.FirstName, Profile.LastName).
		MODEL(profile).
		RETURNING(Profile.AllColumns)

	var dest model.Profile

	err := stmt.Query(store.db, &dest)
	if err != nil {
		app.Logger.Error().Err(err)
		return nil, err
	}

	app.Logger.Debug().Any("Result", dest).Send()

	return &dest, nil
}

func (store *postgresStore) CreateTX(app *ctx.Application, curCtx context.Context, tx *sql.Tx, create profiledto.ProfileCreate) (*model.Profile, error) {
	profile := model.Profile{
		UserID:    create.UserID,
		FirstName: create.FirstName,
		LastName:  create.LastName,
	}

	stmt := Profile.
		INSERT(Profile.UserID, Profile.FirstName, Profile.LastName).
		MODEL(profile).
		RETURNING(Profile.AllColumns)

	var dest model.Profile

	err := stmt.QueryContext(curCtx, tx, &dest)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.FailCreateError
		}
		app.Logger.Error().Err(err)
		return nil, err
	}

	app.Logger.Debug().Any("Result", dest).Send()

	return &dest, nil
}
