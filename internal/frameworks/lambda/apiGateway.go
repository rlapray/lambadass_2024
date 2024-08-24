package lambda

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/rs/zerolog"
)

/******************************************************************************
***** Structs
******************************************************************************/

type APIGatewayClient struct {
	Lambda[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]
}

type HTTPResponseKOBody struct {
	StatusCode int            `json:"statusCode"`
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	Metadata   map[string]any `json:"metadata"`
}

/******************************************************************************
***** Functions
******************************************************************************/

// Generate response body string for use for APIGatewayProxyResponse
func (t APIGatewayClient) newErrorResponseBody(statusCode int, code, message string, additionalMetadata map[string]any) (string, error) {
	if additionalMetadata != nil {
		additionalMetadata["requestId"] = t.request.RequestContext.RequestID
		additionalMetadata["requestTime"] = t.request.RequestContext.RequestTime
	} else {
		additionalMetadata = map[string]any{
			"requestId":   t.request.RequestContext.RequestID,
			"requestTime": t.request.RequestContext.RequestTime,
		}
	}

	resObject := HTTPResponseKOBody{StatusCode: statusCode, Code: code, Message: message, Metadata: additionalMetadata}
	resJSON, err := json.Marshal(resObject)
	if err != nil {
		return "", fault.NewAPIGateway(t.logger, 500, code, message,
			map[string]any{"marshall": map[string]any{"message": err.Error()}}, nil)
	}
	return string(resJSON), nil
}

// KO generate a (APIGatewayProxyResponse,fault.Fault) tuple for your lambda meaning there was an error.
//
// Example :
//
//	func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, fault.Fault) {
//		pet, err := petRepository.find(request.PathParameters["id"])
//		if err != nil {
//			return trezer.KO(404, "PET_NOT_FOUND", "Cannot find your pet", nil)
//		}
//		return trezer.OK(pet)
//	}
func (t APIGatewayClient) KO(statusCode int, code, message string, metadata map[string]any) (events.APIGatewayProxyResponse, fault.Fault) {
	return events.APIGatewayProxyResponse{}, fault.NewAPIGateway(t.logger, statusCode, code, message, metadata, nil)
}

// KOFromFault generate a (APIGatewayProxyResponse,fault.Fault) tuple for your lambda from any fault
func (t APIGatewayClient) KOFromFault(statusCode int, flt fault.Fault) (events.APIGatewayProxyResponse, fault.Fault) {
	return events.APIGatewayProxyResponse{}, fault.NewAPIGatewayFromFault(t.logger, statusCode, flt)
}

// KOFromValidatorFault generate a (APIGatewayProxyResponse,fault.Fault) tuple for your lambda from a validator fault
func (t APIGatewayClient) KOFromValidatorFault(flt fault.Fault) (events.APIGatewayProxyResponse, fault.Fault) {
	return events.APIGatewayProxyResponse{}, fault.NewAPIGatewayFromValidatorFault(t.logger, flt)
}

// OK generate a APIGatewayProxyResponse for your lambda, marshaling your response object into JSON.
//
// Example :
//
//	func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//		pet, err := petRepository.find(request.PathParameters["id"])
//		if err != nil {
//			return trezer.KO(404, "PET_NOT_FOUND", "Cannot find your pet", nil)
//		}
//		return trezer.OK(pet)
//	}
func (t APIGatewayClient) OK(obj any) (events.APIGatewayProxyResponse, fault.Fault) {
	if obj == nil {
		t.logger.Trace().Msg("Response object is nil (204)")
		return events.APIGatewayProxyResponse{StatusCode: 204}, nil
	}
	resJSON, err := json.Marshal(obj)
	if err != nil {
		return events.APIGatewayProxyResponse{},
			fault.NewAPIGateway(t.logger, 500, "ERROR_MARSHALL_JSON", "Error while marshaling an object to JSON", map[string]any{
				"marshall": map[string]any{"message": err.Error()},
			}, err)
	}
	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(resJSON)}, nil
}

/******************************************************************************
***** Middleware
******************************************************************************/

// OnSetup is called *before* the first API Gateway request is processed.
func (t APIGatewayClient) OnSetup(_ context.Context, _ *events.APIGatewayProxyRequest) fault.Fault {
	t.logger.Trace().Msg("OnSetup")
	return nil
}

// OnBefore is called before *each* API Gateway request is processed.
func (t APIGatewayClient) OnBefore(_ context.Context, _ *events.APIGatewayProxyRequest) fault.Fault {
	t.logger.Trace().Msg("OnBefore")
	return nil
}

// OnAfter is called after *each* API Gateway response is generated.
func (t APIGatewayClient) OnAfter(response *events.APIGatewayProxyResponse, err fault.Fault) fault.Fault {
	t.logger.Trace().Msg("OnAfter")
	if response == nil { // Impossible without changing Lambda::handleRequest
		return fault.NewAPIGateway(t.logger, 500, "API_GATEWAY_NIL_RESPONSE", "APIGatewayClient::OnAfter received a nil response", nil, err)
	}

	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}
	response.Headers["requestId"] = t.request.RequestContext.RequestID
	response.Headers["requestTime"] = t.request.RequestContext.RequestTime

	if err != nil {
		apigf, ok := err.(*fault.APIGatewayProxyFault)
		if ok {
			t.logger.Trace().Msg("Error type is an ApiGatewayFault")
			res, _ := t.newErrorResponseBody(apigf.StatusCode, apigf.Code(), apigf.Message(), apigf.Metadata())
			response.Body = res
			response.StatusCode = apigf.StatusCode
			return nil
		}
		t.logger.Warn().Msg("Error type is a Fault but should be an ApiGatewayFault with a status code, so choosing 500 by default")
		res, _ := t.newErrorResponseBody(500, err.Code(), err.Message(), err.Metadata())
		response.Body = res
		response.StatusCode = 500
		return nil
	}
	return nil
}

// OnShutdown is called when the lambda is killed by AWS
func (t APIGatewayClient) OnShutdown() {
	t.logger.Trace().Msg("OnShutdown")
}

/******************************************************************************
***** Test
******************************************************************************/

// This function should only be used in tests
func TestNewLambda(req *events.APIGatewayProxyRequest) Lambda[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse] {
	l := zerolog.Logger{}
	return Lambda[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]{
		request: req,
		logger:  &l,
	}
}

// This function should only be used in tests
func TestNewLambdaWithFunc(
	f HandlerFunc[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse],
) Lambda[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse] {
	return Lambda[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]{
		handler: f,
	}
}

// This function should only be used in tests
func TestHandleRequest(
	trezer *Lambda[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse],
	request *events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	return trezer.handleRequest(context.Background(), *request)
}

// This function should only be used in tests
func TestGetMiddlewares(
	trezer *Lambda[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse],
) []MiddlewareInterface[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse] {
	return trezer.middlewares
}
