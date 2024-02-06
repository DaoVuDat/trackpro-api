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

type UpdateProjectRepo interface {
	Update(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID, updateProject projectdto.ProjectUpdate) (*projectdto.ProjectQuery, error)
}

func (store *postgresStore) Update(app *ctx.Application, projectId uuid.UUID, userId *uuid.UUID, updateProject projectdto.ProjectUpdate) (*projectdto.ProjectQuery, error) {
	var projectToUpdate model.Project

	var fieldsToUpdate ColumnList

	if updateProject.Name.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Project.Name)
		projectToUpdate.Name = &updateProject.Name.String
	}

	if updateProject.Description.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Project.Description)
		projectToUpdate.Description = &updateProject.Description.String
	}

	if updateProject.Price.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Project.Price)
		price := int32(updateProject.Price.Int64)
		projectToUpdate.Price = &price
	}

	if updateProject.StartTime.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Project.StartTime)
		projectToUpdate.StartTime = &updateProject.StartTime.Time
	}

	if updateProject.EndTime.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Project.EndTime)
		projectToUpdate.EndTime = &updateProject.EndTime.Time
	}

	if updateProject.Status.Valid {
		fieldsToUpdate = append(fieldsToUpdate, Project.Status)
		var projectStatus model.ProjectStatus
		err := projectStatus.Scan(updateProject.Status.String)
		if err != nil {
			app.Logger.Error().Err(err)
			return nil, err
		}

		projectToUpdate.Status = projectStatus
	}

	condition := Bool(true)
	condition = condition.AND(Project.ID.EQ(UUID(projectId)))
	if userId != nil {
		condition = condition.AND(Project.UserID.EQ(UUID(*userId)))
	}

	cte := CTE("after_update")
	cteUserIdColumn := Project.UserID.From(cte)

	stmt := WITH(cte.AS(Project.UPDATE(fieldsToUpdate).
		MODEL(projectToUpdate).
		WHERE(condition).
		RETURNING(Project.AllColumns.Except(Project.CreatedAt, Project.UpdatedAt))),
	)(
		SELECT(cte.AllColumns(), Account.Username).
			FROM(cte.INNER_JOIN(Account, Account.ID.EQ(cteUserIdColumn))),
	)

	var project projectdto.ProjectQuery

	err := stmt.Query(store.db, &project)
	if err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, common.FailUpdateError
		}
		return nil, err
	}

	return &project, nil
}
