package fault

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UseCaseFault struct {
	code       string
	message    string
	middleware string
	metadata   map[string]any
	cause      error
}

func (u UseCaseFault) Code() string {
	return u.code
}

func (UseCaseFault) Layer() Layer {
	return UseCases
}

func (u UseCaseFault) Middleware() string {
	return u.middleware
}

func (u UseCaseFault) Message() string {
	return u.message
}

func (u *UseCaseFault) Metadata() map[string]any {
	return u.metadata
}

func (u UseCaseFault) Cause() error {
	return u.cause
}

func (u UseCaseFault) Error() string {
	return fmt.Sprintf("UseCaseFault [%v] : %v", u.code, u.message)
}

func NewUseCase(logger *zerolog.Logger, middleware, code, message string, metadata map[string]any, cause error) Fault {
	e := logger.Warn().AnErr("cause", cause)
	if metadata != nil {
		if id1, exists := metadata["id"]; exists {
			if id2, ok := id1.(uuid.UUID); ok {
				e = e.Str("object_id", id2.String())
			}
		}
	}
	fault := UseCaseFault{middleware: middleware, code: code, message: message, metadata: metadata, cause: cause}
	e.Err(&fault).Msg(message)
	return &fault
}
