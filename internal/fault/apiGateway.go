package fault

import (
	"fmt"

	"github.com/rs/zerolog"
)

type APIGatewayProxyFault struct {
	StatusCode int
	code       string
	message    string
	metadata   map[string]any
	cause      error
}

func (u APIGatewayProxyFault) Code() string {
	return u.code
}

func (APIGatewayProxyFault) Layer() Layer {
	return Commands
}

func (APIGatewayProxyFault) Middleware() string {
	return "APIGatewayProxy"
}

func (u APIGatewayProxyFault) Message() string {
	return u.message
}

func (u *APIGatewayProxyFault) Metadata() map[string]any {
	return u.metadata
}

func (u APIGatewayProxyFault) Cause() error {
	return u.cause
}

func (u APIGatewayProxyFault) Error() string {
	return fmt.Sprintf("APIGatewayProxyFault [%v] : %v", u.code, u.message)
}

func NewAPIGateway(logger *zerolog.Logger, statusCode int, code, message string, metadata map[string]any, cause error) Fault {
	e := logger.Warn().AnErr("cause", cause)
	fault := APIGatewayProxyFault{StatusCode: statusCode, code: code, message: message, metadata: metadata, cause: cause}
	e.Err(&fault).Msg(message)
	return &fault
}

func NewAPIGatewayFromFault(logger *zerolog.Logger, statusCode int, cause Fault) Fault {
	e := logger.Warn().AnErr("cause", cause)
	fault := APIGatewayProxyFault{StatusCode: statusCode, code: cause.Code(), message: cause.Message(), metadata: cause.Metadata(), cause: cause}
	e.Err(&fault).Msg(cause.Message())
	return &fault
}

func NewAPIGatewayFromValidatorFault(logger *zerolog.Logger, cause Fault) Fault {
	e := logger.Warn().AnErr("cause", cause)

	var statusCode int
	switch cause.Code() {
	case "UNKNOWN_FIELD":
		statusCode = 400
	case "MALFORMED_JSON":
		statusCode = 400
	case "EMPTY_JSON":
		statusCode = 400
	case "WRONG_TYPE":
		statusCode = 400
	case "BAD_REQUEST":
		statusCode = 400
	case "INTERNAL_MARSHALING_ERROR":
		statusCode = 500
	default:
		statusCode = 500
	}

	fault := APIGatewayProxyFault{StatusCode: statusCode, code: cause.Code(), message: cause.Message(), metadata: cause.Metadata(), cause: cause}
	e.Err(&fault).Msg(cause.Message())
	return &fault
}
