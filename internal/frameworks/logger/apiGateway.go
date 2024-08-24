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

type APIGatewayClient struct {
	Client[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]
}

/******************************************************************************
***** Middleware
******************************************************************************/

func (m *APIGatewayClient) OnSetup(_ context.Context, firstRequest *events.APIGatewayProxyRequest) fault.Fault {
	m.preSetup()

	log.Logger = log.Logger.With().
		Str("type", "APIGatewayProxy").
		Str("method", firstRequest.HTTPMethod).
		Str("path", firstRequest.Path).
		Str("request", firstRequest.RequestContext.RequestID).
		Logger()

	m.Logger = m.Logger.With().
		Str("type", "APIGatewayProxy").
		Str("method", firstRequest.HTTPMethod).
		Str("path", firstRequest.Path).
		Str("request", firstRequest.RequestContext.RequestID).
		Logger()
	m.Logger.Debug().Msg("Setup logger ok (MiddlewareAPIGateway)")
	return nil
}

// Unsafe if the request contains private informations
func (m *APIGatewayClient) OnBefore(_ context.Context, request *events.APIGatewayProxyRequest) fault.Fault {
	m.Logger.Trace().Msg("OnBefore")
	m.Logger.Trace().Interface("request", request).Msg("Request log")
	return nil
}
