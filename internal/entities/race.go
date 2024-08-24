package entities

import (
	"github.com/google/uuid"
)

type Race struct {
	ID   uuid.UUID `db:"id"   json:"id"`
	Name string    `db:"name" json:"name,omitempty"`
}
