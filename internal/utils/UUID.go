package utils

import (
	"github.com/google/uuid"
)

/******************************************************************************
***** Structs
******************************************************************************/

type UUIDGenerator interface {
	NewV7() (uuid.UUID, error)
}

type GoogleUUIDGenerator struct{}

func (*GoogleUUIDGenerator) NewV7() (uuid.UUID, error) {
	return uuid.NewV7()
}

/******************************************************************************
***** Mocks
******************************************************************************/

type MockUUIDGenerator struct{}

func (*MockUUIDGenerator) NewV7() (uuid.UUID, error) {
	return uuid.MustParse("11111111-1111-1111-1111-111111111111"), nil
}
