package profiledto

import "github.com/google/uuid"

type ProfileCreate struct {
	UserID    uuid.UUID
	FirstName string
	LastName  string
}
