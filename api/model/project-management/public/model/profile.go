//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
	ImageURL  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
