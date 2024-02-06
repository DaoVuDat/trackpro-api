package profilerepo

import (
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	profiledto "github.com/DaoVuDat/trackpro-api/api/resource/profile/dto"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
	"time"
)

type UpdateProfileRepo interface {
	Update(app *ctx.Application, accountId uuid.UUID, update profiledto.ProfileUpdate) (*model.Profile, error)
}

func (store *postgresStore) Update(app *ctx.Application, accountId uuid.UUID, update profiledto.ProfileUpdate) (*model.Profile, error) {
	profileToUpdate := &model.Profile{
		UpdatedAt: time.Now(),
	}

	// Dynamic Query
	var fieldsToUpdate ColumnList
	fieldsToUpdate = append(fieldsToUpdate, Profile.UpdatedAt)

	if update.FirstName.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Profile.FirstName)
		profileToUpdate.FirstName = update.FirstName.String
	}

	if update.LastName.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Profile.LastName)
		profileToUpdate.LastName = update.LastName.String
	}

	if update.ImageUrl.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Profile.ImageURL)
		profileToUpdate.ImageURL = &update.ImageUrl.String
	}

	stmt := Profile.UPDATE(fieldsToUpdate).
		MODEL(profileToUpdate).
		WHERE(Profile.UserID.EQ(UUID(accountId))).
		RETURNING(Profile.AllColumns.Except(Profile.CreatedAt, Profile.UpdatedAt))

	var updatedProfile model.Profile
	err := stmt.Query(store.db, &updatedProfile)
	if err != nil {
		app.Logger.Error().Err(err)
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.FailUpdateError
		}
		return nil, err
	}

	return &updatedProfile, nil
}
