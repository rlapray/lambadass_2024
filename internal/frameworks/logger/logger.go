package logger

import (
	"context"
	"os"

	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

/******************************************************************************
***** Structs
******************************************************************************/

type Client[T any, U any] struct {
	Logger zerolog.Logger
}

/******************************************************************************
***** Middleware
******************************************************************************/

func (m *Client[T, U]) preSetup() {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.999Z07:00"
	if os.Getenv("ENVIRONMENT") == "LOCAL" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	m.Logger = log.Logger.With().Str("framework", "LOGGER").Logger()
	m.Logger.Trace().Msg("OnSetup")
}

func (m *Client[T, U]) OnSetup(_ context.Context, _ *T) fault.Fault {
	m.preSetup()
	m.Logger = log.Logger.With().
		Str("type", "unknown").
		Logger()
	m.Logger.Debug().Msg("Setup logger ok (any)")
	return nil
}

func (m *Client[T, U]) OnBefore(_ context.Context, _ *T) fault.Fault {
	m.Logger.Trace().Msg("OnBefore")
	return nil
}

func (m *Client[T, U]) OnAfter(_ *U, err fault.Fault) fault.Fault {
	m.Logger.Trace().Msg("OnAfter")
	return err
}

func (m *Client[T, U]) OnShutdown() {
	m.Logger.Trace().
		Msg("OnShutdown")
}

/******************************************************************************
***** Functions
******************************************************************************/

// With creates a child logger with the field added to its context.
func (*Client[T, U]) With() zerolog.Context {
	return log.Logger.With()
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func (*Client[T, U]) Err(err error) *zerolog.Event {
	return log.Logger.Err(err)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func (*Client[T, U]) Trace() *zerolog.Event {
	return log.Logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (*Client[T, U]) Debug() *zerolog.Event {
	return log.Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (*Client[T, U]) Info() *zerolog.Event {
	return log.Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (*Client[T, U]) Warn() *zerolog.Event {
	return log.Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (*Client[T, U]) Error() *zerolog.Event {
	return log.Logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func (*Client[T, U]) Fatal() *zerolog.Event {
	return log.Logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func (*Client[T, U]) Panic() *zerolog.Event {
	return log.Logger.Panic()
}
