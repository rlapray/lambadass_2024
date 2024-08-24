package fault

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type SQLFault struct {
	code     string
	message  string
	metadata map[string]any
	cause    error
}

func (e SQLFault) Code() string {
	return e.code
}

func (SQLFault) Layer() Layer {
	return Adapters
}

func (SQLFault) Middleware() string {
	return "Sql"
}

func (e SQLFault) Message() string {
	return e.message
}

func (e SQLFault) Metadata() map[string]any {
	return e.metadata
}

func (e SQLFault) Cause() error {
	return e.cause
}

func (e SQLFault) Error() string {
	return fmt.Sprintf("SqlFault [%v] : %v", e.code, e.message)
}

func NewSQL(logger *zerolog.Logger, code, message string, metadata map[string]any, cause error) Fault {
	e := logger.Warn().AnErr("cause", cause)
	var msg string
	if metadata != nil {
		if dur, exists := metadata["duration"]; exists {
			if d, ok := dur.(int64); ok {
				e = e.Int64("duration", d)
				durationObj := time.Duration(d) * time.Millisecond
				msg = fmt.Sprintf("%v (%v)", message, durationObj.String())
			}
		}
	}
	fault := SQLFault{code: code, message: message, metadata: metadata, cause: cause}
	e.Err(&fault).Msg(msg)
	return &fault
}
