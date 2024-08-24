package validator

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/*****************************************************************************
***** Structs
******************************************************************************/

type LambdaValidator[T any, U any] struct {
	logger    *zerolog.Logger
	validator *validator.Validate
}

/******************************************************************************
***** Functions
******************************************************************************/

// Validates the given JSON string and populates the provided data structure.
// It returns a fault.Fault if there are any validation errors.
func (t LambdaValidator[T, U]) ValidateJSONIntoStruct(jsonString string, data any) fault.Fault {
	decoder := json.NewDecoder(strings.NewReader(jsonString))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		return fault.NewValidatorFaultFromDecoder(t.logger, err)
	}

	err = t.validator.Struct(data)
	if err != nil {
		return fault.NewValidatorFaultFromStruct(t.logger, err)
	}
	return nil
}

/******************************************************************************
***** Middleware
******************************************************************************/
func (t *LambdaValidator[T, U]) OnSetup(_ context.Context, _ *T) fault.Fault {
	ll := log.Logger.With().Str("commands", "Validator").Logger()
	t.logger = &ll
	t.logger.Info().Msg("Creating validator")
	t.validator = validator.New(validator.WithRequiredStructEnabled())
	t.logger.Trace().Msg("OnSetup")
	return nil
}

func (t *LambdaValidator[T, U]) OnBefore(_ context.Context, _ *T) fault.Fault {
	t.logger.Trace().Msg("OnBefore")
	return nil
}

func (t *LambdaValidator[T, U]) OnAfter(_ *U, err fault.Fault) fault.Fault {
	t.logger.Trace().Msg("OnAfter")
	return err
}

func (t *LambdaValidator[T, U]) OnShutdown() {
	t.logger.Trace().Msg("OnShutdown")
}
