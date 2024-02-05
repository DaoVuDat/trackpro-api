package projectdto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gopkg.in/guregu/null.v4"
	"time"
)

type ProjectCreate struct {
	UserId      string      `json:"user_id"`
	Name        null.String `json:"name"`
	Description null.String `json:"description"`
	Price       null.Int    `json:"price"`
	StartTime   null.Time   `json:"start_time"`
	EndTime     null.Time   `json:"end_time"`
}

func (projectCreate ProjectCreate) Validate() error {
	return validation.ValidateStruct(&projectCreate)
}

type ProjectUpdate struct {
}

func (projectUpdate ProjectUpdate) Validate() error {
	return validation.Validate(&projectUpdate)
}

type ProjectResponse struct {
	Id          string    `json:"id"`
	UserId      string    `json:"user_id"`
	UserName    string    `json:"username"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}
