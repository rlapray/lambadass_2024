package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/lambadass-2024/backend/internal/adapters/repositories"
	"github.com/lambadass-2024/backend/internal/entities"
	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/lambadass-2024/backend/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/******************************************************************************
***** Structs
******************************************************************************/

type PetUseCase[T any, U any] struct {
	logger        *zerolog.Logger
	UUIDGenerator utils.UUIDGenerator
	Repository    *repositories.PetRepository[T, U]
}

/******************************************************************************
***** Errors
******************************************************************************/

func (u PetUseCase[T, U]) newError(code, message string, metadata map[string]any, cause error) fault.Fault {
	return fault.NewUseCase(u.logger, "PetUseCase", code, message, metadata, cause)
}

/******************************************************************************
***** Functions
******************************************************************************/

func (u PetUseCase[T, U]) preparePetCreation(id uuid.UUID, name string, raceID uuid.UUID, metadata map[string]any) (entities.Pet, fault.Fault) {
	if id == uuid.Nil {
		u.logger.Debug().Msg("No id provided, generating one...")
		newID, err := u.UUIDGenerator.NewV7()
		if err != nil {
			return entities.Pet{}, u.newError("IDENTIFIER_GENERATION_ERROR", "Cannot generate identifier", metadata, err)
		}
		return entities.Pet{ID: newID, Name: name, Race: entities.Race{ID: raceID}}, nil
	}
	return entities.Pet{ID: id, Name: name, Race: entities.Race{ID: raceID}}, nil
}

func (u PetUseCase[T, U]) Create(id uuid.UUID, name string, raceID uuid.UUID) (entities.Pet, fault.Fault) {
	u.logger.Trace().Msg("Create")
	metadata := map[string]any{
		"id": id,
	}

	p, err := u.preparePetCreation(id, name, raceID, metadata)
	if err != nil {
		return p, err
	}

	p, err = u.Repository.Create(p)
	if err == nil {
		return p, nil
	}
	switch err.Code() {
	case "UNIQUE_VIOLATION":
		return p, u.newError("PET_ID_NOT_UNIQUE", "Pet id not unique", metadata, err)
	default:
		return p, u.newError("PET_CREATION_FAILED", "Pet creation failed", metadata, err)
	}
}

func (u PetUseCase[T, U]) Get(id uuid.UUID) (entities.Pet, fault.Fault) {
	u.logger.Trace().Msg("Get")
	metadata := map[string]any{
		"id": id,
	}

	p, err := u.Repository.Get(id)
	if err == nil {
		return p, nil
	}
	switch err.Code() {
	case "NOT_FOUND":
		return p, u.newError("PET_NOT_FOUND", "Pet not found", metadata, err)
	default:
		return p, u.newError("PET_GET_FAILED", "Cannot get this pet", metadata, err)
	}
}

/******************************************************************************
***** Middlewares
******************************************************************************/

// Setup logger and UUID generator
func (u *PetUseCase[T, U]) OnSetup(_ context.Context, _ *T) fault.Fault {
	ll := log.Logger.With().Str("usecase", "PetUseCases").Logger()
	u.logger = &ll
	u.logger.Trace().Msg("OnSetup")
	u.UUIDGenerator = &utils.GoogleUUIDGenerator{}
	u.logger.Trace().Str("UUIDGenerator", "GoogleUUIDGenerator").Msg("UUIDGenerator is set")
	return nil
}

func (u PetUseCase[T, U]) OnBefore(_ context.Context, _ *T) fault.Fault {
	u.logger.Trace().Msg("OnBefore")
	return nil
}

func (u PetUseCase[T, U]) OnAfter(_ *U, err fault.Fault) fault.Fault {
	u.logger.Trace().Msg("OnAfter")
	return err
}

func (u PetUseCase[T, U]) OnShutdown() {
	u.logger.Trace().Msg("OnShutdown")
}
