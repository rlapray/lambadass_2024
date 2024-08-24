package logger

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/rs/zerolog/log"
)

/******************************************************************************
***** Structs
******************************************************************************/

type MiddlewareSQS struct {
	Client[events.SQSEvent, events.SQSEventResponse]
}

/******************************************************************************
***** Middleware
******************************************************************************/

func (m *MiddlewareSQS) OnSetup(_ context.Context, firstRequest events.SQSEvent) fault.Fault {
	m.preSetup()

	log.Logger = log.Logger.With().
		Str("type", "SQSEvent").
		Int("count", len(firstRequest.Records)).
		Logger()
	m.Logger = m.Logger.With().
		Str("type", "SQSEvent").
		Int("count", len(firstRequest.Records)).
		Logger()

	log.Logger.Debug().Msg("Setup logger ok (MiddlewareSQS)")
	return nil
}
