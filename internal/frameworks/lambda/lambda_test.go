package lambda_test

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/lambadass-2024/backend/internal/frameworks/lambda"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var logger zerolog.Logger = zerolog.Logger{}

/******************************************************************************
***** Handler
******************************************************************************/

func Test_Lambda_HandleRequest_Response200(t *testing.T) {
	expectedResponse := events.APIGatewayProxyResponse{StatusCode: 200, Body: "test ok"}

	trezer := lambda.TestNewLambdaWithFunc(func(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, fault.Fault) {
		return events.APIGatewayProxyResponse{StatusCode: 200, Body: "test ok"}, nil
	})

	// Create a sample request
	request := events.APIGatewayProxyRequest{}

	// Call the handleRequest function
	response, err := lambda.TestHandleRequest(&trezer, &request)

	// Assert that the response is of type U and there is no error
	assert.IsType(t, events.APIGatewayProxyResponse{}, response)
	assert.Equal(t, expectedResponse, response)
	assert.NoError(t, err)
}

func Test_Lambda_HandleRequest_ResponseError(t *testing.T) {
	expectedResponse := events.APIGatewayProxyResponse{}
	expectedError := fault.NewAPIGateway(&logger, 500, "ERROR_CODE", "Message", nil, nil)

	trezer := lambda.TestNewLambdaWithFunc(func(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, fault.Fault) {
		return events.APIGatewayProxyResponse{}, expectedError
	})

	// Create a sample request
	request := events.APIGatewayProxyRequest{}

	// Call the handleRequest function
	response, err := lambda.TestHandleRequest(&trezer, &request)

	// Assert that the response is of type U and there is no error
	assert.IsType(t, events.APIGatewayProxyResponse{}, response)
	assert.Equal(t, expectedResponse, response)
	assert.Equal(t, expectedError, err)
}

/******************************************************************************
***** Middleware OK
******************************************************************************/

var (
	agpreq1 = events.APIGatewayProxyRequest{RequestContext: events.APIGatewayProxyRequestContext{RequestID: "42"}}
	agpres1 = events.APIGatewayProxyResponse{StatusCode: 200, Body: "test ok"}
)

type SampleMiddleware1 struct {
	t                *testing.T
	OnSetupCalled    int
	OnBeforeCalled   int
	OnAfterCalled    int
	OnShutdownCalled int
}

func (m *SampleMiddleware1) OnSetup(_ context.Context, firstRequest *events.APIGatewayProxyRequest) fault.Fault {
	assert.NotNil(m.t, firstRequest)
	assert.Equal(m.t, agpreq1, *firstRequest)
	m.OnSetupCalled++
	return nil
}

func (m *SampleMiddleware1) OnBefore(_ context.Context, request *events.APIGatewayProxyRequest) fault.Fault {
	assert.NotNil(m.t, request)
	assert.Equal(m.t, agpreq1, *request)
	m.OnBeforeCalled++
	return nil
}

func (m *SampleMiddleware1) OnAfter(response *events.APIGatewayProxyResponse, err fault.Fault) fault.Fault {
	assert.NotNil(m.t, response)
	assert.Equal(m.t, agpres1, *response)
	m.OnAfterCalled++
	return err
}

func (*SampleMiddleware1) OnShutdown() {}

type SampleMiddlewareB struct {
	t                *testing.T
	OnSetupCalled    int
	OnBeforeCalled   int
	OnAfterCalled    int
	OnShutdownCalled int
}

func (m *SampleMiddlewareB) OnSetup(_ context.Context, firstRequest *events.APIGatewayProxyRequest) fault.Fault {
	assert.NotNil(m.t, firstRequest)
	assert.Equal(m.t, agpreq1, *firstRequest)
	m.OnSetupCalled++
	return nil
}

func (m *SampleMiddlewareB) OnBefore(_ context.Context, request *events.APIGatewayProxyRequest) fault.Fault {
	assert.NotNil(m.t, request)
	assert.Equal(m.t, agpreq1, *request)
	m.OnBeforeCalled++
	return nil
}

func (m *SampleMiddlewareB) OnAfter(response *events.APIGatewayProxyResponse, err fault.Fault) fault.Fault {
	assert.NotNil(m.t, response)
	assert.Equal(m.t, agpres1, *response)
	m.OnAfterCalled++
	return err
}

func (*SampleMiddlewareB) OnShutdown() {}

func Test_Lambda_Use(t *testing.T) { //nolint:dupl //it's tests, who cares about duplication
	// Create a Trezer instance

	lmbd := lambda.TestNewLambdaWithFunc(func(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, fault.Fault) {
		return agpres1, nil
	})
	// Create a sample middleware1
	middleware1 := SampleMiddleware1{t: t}
	middleware1b := SampleMiddlewareB{t: t}

	// Add the middleware to the Trezer instance
	lmbd.Use(&middleware1)
	lmbd.Use(&middleware1b)

	// Assert that the middleware is added to the middlewares slice
	assert.Contains(t, lambda.TestGetMiddlewares(&lmbd), &middleware1)

	// Call the handleRequest function
	_, _ = lambda.TestHandleRequest(&lmbd, &agpreq1)
	assert.Equal(t, 1, middleware1.OnSetupCalled)
	assert.Equal(t, 1, middleware1.OnBeforeCalled)
	assert.Equal(t, 1, middleware1.OnAfterCalled)
	assert.Equal(t, 1, middleware1b.OnSetupCalled)
	assert.Equal(t, 1, middleware1b.OnBeforeCalled)
	assert.Equal(t, 1, middleware1b.OnAfterCalled)

	_, _ = lambda.TestHandleRequest(&lmbd, &agpreq1)
	assert.Equal(t, 1, middleware1.OnSetupCalled) // OnSetup is called only once
	assert.Equal(t, 2, middleware1.OnBeforeCalled)
	assert.Equal(t, 2, middleware1.OnAfterCalled)
	assert.Equal(t, 1, middleware1b.OnSetupCalled) // OnSetup is called only once
	assert.Equal(t, 2, middleware1b.OnBeforeCalled)
	assert.Equal(t, 2, middleware1b.OnAfterCalled)
}

/******************************************************************************
***** Middleware onSetup fail
******************************************************************************/

type SampleMiddleware2 struct {
	t                *testing.T
	OnSetupCalled    int
	OnBeforeCalled   int
	OnAfterCalled    int
	OnShutdownCalled int
}

func (m *SampleMiddleware2) OnSetup(_ context.Context, firstRequest *events.APIGatewayProxyRequest) fault.Fault {
	assert.NotNil(m.t, firstRequest)
	assert.Equal(m.t, agpreq1, *firstRequest)
	m.OnSetupCalled++
	return fault.NewUseCase(&logger, "SampleMiddleware2", "CODE1", "Code 1", nil, nil)
}

func (m *SampleMiddleware2) OnBefore(_ context.Context, _ *events.APIGatewayProxyRequest) fault.Fault {
	assert.Fail(m.t, "OnBefore should not be called")
	return nil
}

func (m *SampleMiddleware2) OnAfter(_ *events.APIGatewayProxyResponse, err fault.Fault) fault.Fault {
	m.OnAfterCalled++
	return err
}

func (*SampleMiddleware2) OnShutdown() {}

func Test_Lambda_OnSetupFail(t *testing.T) { //nolint:dupl //it's tests, who cares about duplication
	// Create a Trezer instance

	lmbd := lambda.TestNewLambdaWithFunc(func(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, fault.Fault) {
		return agpres1, nil
	})
	// Create a sample middleware
	middleware := SampleMiddleware2{t: t}
	middleware1b := SampleMiddlewareB{t: t}

	// Add the middleware to the Trezer instance
	lmbd.Use(&middleware)
	lmbd.Use(&middleware1b)

	// Assert that the middleware is added to the middlewares slice
	assert.Contains(t, lambda.TestGetMiddlewares(&lmbd), &middleware)

	// Call the handleRequest function
	_, _ = lambda.TestHandleRequest(&lmbd, &agpreq1)
	assert.Equal(t, 1, middleware.OnSetupCalled)
	assert.Equal(t, 0, middleware.OnBeforeCalled) // onSetup fail, no onBefore
	assert.Equal(t, 0, middleware.OnAfterCalled)  // no onAfter too

	assert.Equal(t, 0, middleware1b.OnSetupCalled) // the next middleware will never be called
	assert.Equal(t, 0, middleware1b.OnBeforeCalled)
	assert.Equal(t, 0, middleware1b.OnAfterCalled)

	_, _ = lambda.TestHandleRequest(&lmbd, &agpreq1)
	assert.Equal(t, 2, middleware.OnSetupCalled) // onSetup fail, no onBefore, forever
	assert.Equal(t, 0, middleware.OnBeforeCalled)
	assert.Equal(t, 0, middleware.OnAfterCalled)

	assert.Equal(t, 0, middleware1b.OnSetupCalled)
	assert.Equal(t, 0, middleware1b.OnBeforeCalled)
	assert.Equal(t, 0, middleware1b.OnAfterCalled)
}

/******************************************************************************
***** Middleware onBefore fail
******************************************************************************/

type SampleMiddleware3 struct {
	t                *testing.T
	OnSetupCalled    int
	OnBeforeCalled   int
	OnAfterCalled    int
	OnShutdownCalled int
}

func (m *SampleMiddleware3) OnSetup(_ context.Context, firstRequest *events.APIGatewayProxyRequest) fault.Fault {
	assert.NotNil(m.t, firstRequest)
	assert.Equal(m.t, agpreq1, *firstRequest)
	m.OnSetupCalled++
	return nil
}

func (m *SampleMiddleware3) OnBefore(_ context.Context, request *events.APIGatewayProxyRequest) fault.Fault {
	assert.NotNil(m.t, request)
	assert.Equal(m.t, agpreq1, *request)
	m.OnBeforeCalled++
	return fault.NewUseCase(&logger, "SampleMiddleware3", "CODE2", "Code 2", nil, nil)
}

func (m *SampleMiddleware3) OnAfter(_ *events.APIGatewayProxyResponse, err fault.Fault) fault.Fault {
	m.OnAfterCalled++
	return err
}

func (*SampleMiddleware3) OnShutdown() {}

func Test_Lambda_OnBeforeFail(t *testing.T) { //nolint:dupl //it's tests, who cares about duplication
	// Create a Trezer instance

	lmbd := lambda.TestNewLambdaWithFunc(func(_ context.Context, _ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, fault.Fault) {
		return agpres1, nil
	})
	// Create a sample middleware
	middleware := SampleMiddleware3{t: t}
	middlewareB := SampleMiddlewareB{t: t}

	// Add the middleware to the Trezer instance
	lmbd.Use(&middleware)
	lmbd.Use(&middlewareB)

	// Assert that the middleware is added to the middlewares slice
	assert.Contains(t, lambda.TestGetMiddlewares(&lmbd), &middleware)

	// Call the handleRequest function
	_, _ = lambda.TestHandleRequest(&lmbd, &agpreq1)
	assert.Equal(t, 1, middleware.OnSetupCalled)
	assert.Equal(t, 1, middleware.OnBeforeCalled)
	assert.Equal(t, 0, middleware.OnAfterCalled) // onbefore failed, onafter will not be called

	assert.Equal(t, 1, middlewareB.OnSetupCalled)
	assert.Equal(t, 0, middlewareB.OnBeforeCalled) // previous before failed, next one is not called
	assert.Equal(t, 0, middlewareB.OnAfterCalled)  // as the previous middleware failed, this one will not be called
	// to release ressources, use onShutdown, not onAfter

	_, _ = lambda.TestHandleRequest(&lmbd, &agpreq1)
	assert.Equal(t, 1, middleware.OnSetupCalled)
	assert.Equal(t, 2, middleware.OnBeforeCalled)
	assert.Equal(t, 0, middleware.OnAfterCalled) // // onbefore failed, onafter will not be called, forever

	assert.Equal(t, 1, middlewareB.OnSetupCalled)
	assert.Equal(t, 0, middlewareB.OnBeforeCalled) // previous before failed, next one is not called, forever
	assert.Equal(t, 0, middlewareB.OnAfterCalled)  // as the previous middleware failed, this one will not be called, forever
	// to release ressources, use onShutdown, not onAfter
}
