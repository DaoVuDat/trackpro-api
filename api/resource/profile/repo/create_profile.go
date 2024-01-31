package profilerepo

import (
	"trackpro/api/model/project-management/public/model"
	. "trackpro/api/model/project-management/public/table"
	profiledto "trackpro/api/resource/profile/dto"
	"trackpro/util/ctx"
)

type CreateProfileRepo interface {
	Create(app *ctx.Application, create profiledto.ProfileCreate) (*model.Profile, error)
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
