package lambda

// https://docs.aws.amazon.com/lambda/latest/dg/golang-exceptions.html

import (
	"context"
	"reflect"
	"time"

	lambdaaws "github.com/aws/aws-lambda-go/lambda"
	"github.com/lambadass-2024/backend/internal/fault"
	baselogger "github.com/lambadass-2024/backend/internal/frameworks/logger"
	"github.com/rs/zerolog"
)

/******************************************************************************
***** Structs
******************************************************************************/

// HandlerFunc is a function type that defines the signature of the Lambda function handler.
type HandlerFunc[T any, U any] func(context.Context, T) (U, fault.Fault)

// Lambda is a struct that represents the Lambda function framework.
// It allows you to add middlewares and will handle requests.
type Lambda[T any, U any] struct {
	handler     HandlerFunc[T, U]
	startedOnce bool
	middlewares []MiddlewareInterface[T, U]
	logger      *zerolog.Logger
	request     *T
	startTime   int64
}

/******************************************************************************
***** Middleware
******************************************************************************/

// MiddlewareInterface is an interface that defines the methods to be implemented by middleware functions.
// Middleware functions can be used to perform :
//
// - setup (first call)
//
// - pre-processing (before each request)
//
// - post-processing (after each request)
//
// - and cleanup tasks (when the lambda is killed with SIGTERM by AWS).
type MiddlewareInterface[T any, U any] interface {
	// Called on first request
	// Call middlewares in the order they were added
	OnSetup(ctx context.Context, firstRequest *T) fault.Fault

	// Called before each request
	// Call middlewares in the order they were added
	OnBefore(ctx context.Context, request *T) fault.Fault

	// Called after each request
	// Call middlewares in the *reverse* order they were added
	OnAfter(response *U, flt fault.Fault) fault.Fault

	// Called on SIGTERM (lambda killed by AWS)
	// Call middlewares in the *reverse* order they were added
	OnShutdown()
}

func (t *Lambda[T, U]) onSetupHandler(
	ctx context.Context, request T, workingMiddlewares []MiddlewareInterface[T, U],
) ([]MiddlewareInterface[T, U], fault.Fault) {
	if !t.startedOnce {
		t.logger = &zerolog.Logger{}
		t.logger.Debug().Msg("OnSetup")
		for i, mw := range workingMiddlewares {
			err := mw.OnSetup(ctx, &request)
			if i == 0 { // Special case : first middleware should be the logger
				l := baselogger.APIGatewayClient{}
				ll := l.With().Str("framework", "LAMBDA").Logger()
				t.logger = &ll
			}
			if err != nil {
				t.logger.Error().Err(err).Msgf("OnSetup encountered an error (%v/%v middlewares added)", i+1, len(workingMiddlewares))
				return workingMiddlewares[:i], err
			}
		}
		t.startedOnce = true
		t.logger.Info().Msgf("%v middleware(s) added", len(workingMiddlewares))
	}
	return workingMiddlewares, nil
}

func (t *Lambda[T, U]) onBeforeHandler(
	ctx context.Context, request T, workingMiddlewares []MiddlewareInterface[T, U],
) ([]MiddlewareInterface[T, U], fault.Fault) {
	t.logger.Debug().Msg("onBeforeHandler")
	for i, mw := range workingMiddlewares {
		err := mw.OnBefore(ctx, &request)
		if err != nil {
			t.logger.Error().Err(err).Msgf("OnBefore encountered an error (%v/%v middlewares triggered)", i+1, len(workingMiddlewares))
			return workingMiddlewares[:i], err
		}
	}
	return workingMiddlewares, nil
}

func (t *Lambda[T, U]) onAfterHandler(res *U, err fault.Fault, workingMiddlewares []MiddlewareInterface[T, U]) (U, fault.Fault) {
	t.logger.Debug().Msg("onAfterHandler")
	for i := len(workingMiddlewares) - 1; i >= 0; i-- {
		err = workingMiddlewares[i].OnAfter(res, err)
	}
	return *res, err
}

/******************************************************************************
***** Functions
******************************************************************************/

// Use adds a middleware to the Trezer framework.
func (t *Lambda[T, U]) Use(mw MiddlewareInterface[T, U]) *Lambda[T, U] {
	t.middlewares = append(t.middlewares, mw)
	return t
}

// Start starts the lambda with the specified handler function.
//
// Also configure the execution of OnShutdown for all middleware when the lambda
// will receive SIGTERM, in the *reverse* order middlewares were added.
func (t *Lambda[T, U]) Start(handler HandlerFunc[T, U]) {
	t.handler = handler
	t.startTime = time.Now().UnixMilli()

	lambdaaws.StartWithOptions(t.handleRequest, lambdaaws.WithEnableSIGTERM(func() {
		t.logger.Debug().Msg("Received SIGTERM, starting shutdown hooks... (1/2)")
		for i := len(t.middlewares) - 1; i >= 0; i-- {
			t.middlewares[i].OnShutdown()
		}
		t.logger.Debug().Msg("Received SIGTERM, all shutdown hooks triggered (2/2)")
		uptime := time.Now().UnixMilli() - t.startTime
		uptimeDuration := time.Duration(uptime) * time.Millisecond
		uptimeString := uptimeDuration.String()
		t.logger.Info().Int64("uptime", uptime).Msgf("Uptime: %s", uptimeString)
	}))
}

// handleRequest is the internal function that handles the Lambda function request.
//
// It executes the middleware functions :
//
// - OnSetup will be executed once for each middleware, in the order they were added with Use
//
// - OnBefore will be executed for each request for each middleware, in the order they were added with Use
//
// - OnAfter will be executed for each request for each middleware, in the *reverse* order they were added with Use
//
// - OnShutdown is not executed here, see Start
func (t *Lambda[T, U]) handleRequest(ctx context.Context, request T) (U, error) {
	t.request = &request
	workingMiddlewares, err := t.onSetupHandler(ctx, request, t.middlewares)
	if err != nil {
		empty := reflect.New(reflect.TypeFor[U]()).Interface().(*U) //nolint:revive //if this doesn't work, reflection is broken and this is unrecoverable
		return t.onAfterHandler(empty, err, workingMiddlewares)
	}

	workingMiddlewares, err = t.onBeforeHandler(ctx, request, workingMiddlewares)
	if err != nil {
		empty := reflect.New(reflect.TypeFor[U]()).Interface().(*U) //nolint:revive //if this doesn't work, reflection is broken and this is unrecoverable
		return t.onAfterHandler(empty, err, workingMiddlewares)
	}

	t.logger.Trace().Msg("Entering handler...")
	res, err := t.handler(ctx, request)
	t.logger.Trace().Err(err).Interface("response", res).Msg("Exited handler")

	finalResponse, finalError := t.onAfterHandler(&res, err, workingMiddlewares)
	t.logger.Trace().AnErr("finalError", finalError).Interface("finalResponse", res).Msg("Returning from handleRequest")
	return finalResponse, finalError
}

// This function should only be used in tests
func (t *Lambda[T, U]) TestHandleRequest(handler HandlerFunc[T, U], request *T) (U, error) {
	t.handler = handler
	return t.handleRequest(context.Background(), *request)
}
