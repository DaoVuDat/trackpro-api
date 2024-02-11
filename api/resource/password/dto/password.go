package passworddto

import "github.com/google/uuid"

type PasswordCreate struct {
	UserId         uuid.UUID
	HashedPassword string
}
