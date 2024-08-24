package entities

import (
	"github.com/google/uuid"
)

type Pet struct {
	ID   uuid.UUID `db:"id"   json:"id"`
	Name string    `db:"name" json:"name,omitempty"`
	Race Race      `db:"race" json:"race,omitempty"`
}
