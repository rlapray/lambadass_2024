package fault

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type RepositoryFault struct {
	code       string
	message    string
	middleware string
	metadata   map[string]any
	cause      error
}

func (e RepositoryFault) Code() string {
	return e.code
}

func (RepositoryFault) Layer() Layer {
	return Adapters
}

func (e RepositoryFault) Middleware() string {
	return e.middleware
}

func (e RepositoryFault) Message() string {
	return e.message
}

func (e RepositoryFault) Metadata() map[string]any {
	return e.metadata
}

func (e RepositoryFault) Cause() error {
	return e.cause
}

func (e RepositoryFault) Error() string {
	return fmt.Sprintf("RepositoryFault [%v] : %v", e.code, e.message)
}

func NewRepository(logger *zerolog.Logger, middleware, code, message string, metadata map[string]any, cause error) Fault {
	e := logger.Warn().AnErr("cause", cause)
	if metadata != nil {
		if id1, exists := metadata["id"]; exists {
			if id2, ok := id1.(uuid.UUID); ok {
				e = e.Str("object_id", id2.String())
			}
		}
	}
	fault := RepositoryFault{middleware: middleware, code: code, message: message, metadata: metadata, cause: cause}
	e.Err(&fault).Msg(message)
	return &fault
}
