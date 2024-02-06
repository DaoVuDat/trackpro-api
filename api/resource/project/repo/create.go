package projectrepo

import (
	"errors"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	. "github.com/DaoVuDat/trackpro-api/api/model/project-management/public/table"
	projectdto "github.com/DaoVuDat/trackpro-api/api/resource/project/dto"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

type CreateProjectRepo interface {
	Create(app *ctx.Application, createProject projectdto.ProjectCreate) (*projectdto.ProjectQuery, error)
}

func (store *postgresStore) Create(app *ctx.Application, createProject projectdto.ProjectCreate) (*projectdto.ProjectQuery, error) {
	userId, err := uuid.Parse(createProject.UserId)
	if err != nil {
		return nil, common.UUIDBadRequest
	}

	projectToCreate := model.Project{
		UserID: userId,
	}
	var fieldsToAdd ColumnList
	fieldsToAdd = append(fieldsToAdd, Project.UserID)

	if createProject.Name.Valid {
		fieldsToAdd = append(fieldsToAdd, Project.Name)
		projectToCreate.Name = &createProject.Name.String
	}

	if createProject.Description.Valid {
		fieldsToAdd = append(fieldsToAdd, Project.Description)
		projectToCreate.Description = &createProject.Description.String
	}

	if createProject.Price.Valid {
		fieldsToAdd = append(fieldsToAdd, Project.Price)
		price := int32(createProject.Price.Int64)
		projectToCreate.Price = &price
	}

	if createProject.StartTime.Valid {
		fieldsToAdd = append(fieldsToAdd, Project.StartTime)
		projectToCreate.StartTime = &createProject.StartTime.Time
	}

	if createProject.EndTime.Valid {
		fieldsToAdd = append(fieldsToAdd, Project.EndTime)
		projectToCreate.EndTime = &createProject.EndTime.Time
	}

	// Create CTE and export Columns for that CTE
	cte := CTE("after_insert")
	cteUserIdColumn := Project.UserID.From(cte)

	stmt := WITH(
		cte.AS(
			Project.INSERT(fieldsToAdd).
				MODEL(projectToCreate).
				RETURNING(Project.AllColumns.Except(Project.CreatedAt, Project.UpdatedAt)),
		),
	)(
		SELECT(cte.AllColumns(), Account.Username).
			FROM(cte.INNER_JOIN(Account, cteUserIdColumn.EQ(Account.ID))),
	)

	var project projectdto.ProjectQuery

	err = stmt.Query(store.db, &project)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.FailCreateError
		}
		return nil, err
	}

	return &project, nil
}
