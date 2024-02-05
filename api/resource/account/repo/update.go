package accountrepo

import (
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"time"
	"trackpro/api/model/project-management/public/model"
	. "trackpro/api/model/project-management/public/table"
	accountdto "trackpro/api/resource/account/dto"
	"trackpro/util/ctx"
)

type UpdateAccountRepo interface {
	Update(app *ctx.Application, accountId uuid.UUID, update accountdto.AccountUpdate) (*model.Account, error)
}

func (store *postgresStore) Update(app *ctx.Application, accountId uuid.UUID, update accountdto.AccountUpdate) (*model.Account, error) {
	accountToUpdate := model.Account{
		UpdatedAt: time.Now(), // this field could update in DB with trigger / function
	}

	// Create dynamic query
	var fieldsToUpdate ColumnList
	fieldsToUpdate = append(fieldsToUpdate, Account.UpdatedAt)

	if update.Status.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Account.Status)
		var accountStatus model.AccountStatus
		err := accountStatus.Scan(update.Status.String)
		if err != nil {
			return nil, err
		}

		accountToUpdate.Status = accountStatus
	}

	if update.Type.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Account.Type)
		var accountType model.AccountType
		err := accountType.Scan(update.Type.String)
		if err != nil {
			app.Logger.Error().Err(err)
			return nil, err
		}

		accountToUpdate.Type = accountType
	}

	// query
	stmt := Account.UPDATE(fieldsToUpdate).
		MODEL(accountToUpdate).
		WHERE(Account.ID.EQ(UUID(accountId))).
		RETURNING(Account.AllColumns.Except(Account.UpdatedAt, Account.CreatedAt))

	var account model.Account
	err := stmt.Query(store.db, &account)
	if err != nil {
		app.Logger.Error().Err(err)
		return nil, err
	}

	return &account, nil
}
