package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/lambadass-2024/backend/internal/entities"
	"github.com/lambadass-2024/backend/internal/fault"
	sqlframework "github.com/lambadass-2024/backend/internal/frameworks/sql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	PetSQLCreate = "INSERT INTO pet(id, name, race_id) VALUES(:id, :name, :race.id)"
	PetSQLGet    = `SELECT p.id, p.name, r.id as "race.id", r.name as "race.name" FROM pet p inner join race r on p.race_id = r.id WHERE p.id = :id`
)

/******************************************************************************
***** Structs
******************************************************************************/

type PetRepository[T any, U any] struct {
	logger *zerolog.Logger
	SQL    sqlframework.Client[T, U]
}

/******************************************************************************
***** Errors
******************************************************************************/

func (r PetRepository[T, U]) newError(code, message string, metadata map[string]any, cause error) fault.Fault {
	return fault.NewRepository(r.logger, "PetRepository", code, message, metadata, cause)
}

/******************************************************************************
***** Functions
******************************************************************************/

func (r PetRepository[T, U]) Create(pet entities.Pet) (entities.Pet, fault.Fault) {
	metadata := map[string]any{
		"id": pet.ID,
	}
	err := r.SQL.ExecOneRowAffected(PetSQLCreate, pet)
	if err != nil {
		switch err.Code() {
		case "UNIQUE_VIOLATION":
			return entities.Pet{}, r.newError("UNIQUE_VIOLATION", "Pet id not unique", metadata, err)
		default:
			return entities.Pet{}, r.newError("INSERT_ERROR", "Error while inserting Pet", metadata, err)
		}
	}
	r.logger.Debug().Msgf("Pet %v created", pet.ID)
	return pet, nil
}

func (r PetRepository[T, U]) Get(id uuid.UUID) (entities.Pet, fault.Fault) {
	metadata := map[string]any{
		"id": id,
	}
	petIn := entities.Pet{ID: id}
	petsOut := []entities.Pet{}

	err := r.SQL.Select(PetSQLGet, petIn, &petsOut)
	r.logger.Warn().Interface("hop", petsOut).Msg("Debug warn")

	if err != nil {
		return entities.Pet{}, r.newError("SELECT_ERROR", "Error while selecting Pet", metadata, err)
	}
	if len(petsOut) == 0 {
		return entities.Pet{}, r.newError("NOT_FOUND", "Pet not found", metadata, err)
	}
	if len(petsOut) > 1 {
		return entities.Pet{}, r.newError("TOO_MANY_PETS", "Multiple pets found", metadata, err)
	}
	r.logger.Debug().Msgf("Pet %v found", petsOut[0].ID)
	return petsOut[0], nil
}

/******************************************************************************
***** Middlewares
******************************************************************************/

// Setup the logger
func (r *PetRepository[T, U]) OnSetup(_ context.Context, _ *T) fault.Fault {
	ll := log.Logger.With().Str("repository", "PetRepository").Logger()
	r.logger = &ll
	r.logger.Trace().Msg("OnSetup")
	return nil
}

func (r PetRepository[T, U]) OnBefore(_ context.Context, _ *T) fault.Fault {
	r.logger.Trace().Msg("OnBefore")
	return nil
}

// OnAfter is called after *each* API Gateway response is generated.
func (r PetRepository[T, U]) OnAfter(_ *U, err fault.Fault) fault.Fault {
	r.logger.Trace().Msg("OnAfter")
	return err
}

func (r PetRepository[T, U]) OnShutdown() {
	r.logger.Trace().Msg("OnShutdown")
}
