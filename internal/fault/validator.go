package fault

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type ValidatorFault struct {
	code     string
	message  string
	metadata map[string]any
	cause    error
}

type ValidationError struct {
	Message         string `json:"message"`
	Field           string `json:"field"`
	StructNamespace string `json:"namespace"`
	Tag             string `json:"tag"`
	Value           any    `json:"value"`
}

var regexpWrongType = regexp.MustCompile(`json: cannot unmarshal (.+?) into (.+?) of type (.+)`)

func (u ValidatorFault) Code() string {
	return u.code
}

func (ValidatorFault) Layer() Layer {
	return Commands
}

func (ValidatorFault) Middleware() string {
	return "Validator"
}

func (u ValidatorFault) Message() string {
	return u.message
}

func (u *ValidatorFault) Metadata() map[string]any {
	return u.metadata
}

func (u ValidatorFault) Cause() error {
	return u.cause
}

func (u ValidatorFault) Error() string {
	return fmt.Sprintf("ValidatorFault [%v] : %v", u.code, u.message)
}

func NewValidatorFault(logger *zerolog.Logger, code, message string, metadata map[string]any, cause error) Fault {
	e := logger.Warn().AnErr("cause", cause)
	fault := ValidatorFault{code: code, message: message, metadata: metadata, cause: cause}
	e.Err(&fault).Msg(message)
	return &fault
}

func NewValidatorFaultFromStruct(logger *zerolog.Logger, cause error) Fault {
	e := logger.Warn().AnErr("cause", cause)

	errs, ok := cause.(validator.ValidationErrors)
	if ok {
		ves := make([]ValidationError, len((errs)))
		for i, err := range errs {
			ves[i] = ValidationError{
				Message:         err.Error(),
				Field:           err.Field(),
				StructNamespace: err.StructNamespace(),
				Tag:             err.Tag(),
				Value:           err.Value(),
			}
		}
		fault := ValidatorFault{
			code:    "BAD_REQUEST",
			message: "Validation failed",
			metadata: map[string]any{
				"validation": ves,
			}, cause: cause,
		}
		e.Err(&fault).Msg("Request body validation error")
		return &fault
	}

	fault := ValidatorFault{code: "UNEXPECTED_INPUT_VALIDATION_ERROR", message: "Validation raised an unexpected error", metadata: nil, cause: cause}
	e.Err(&fault).Msg("Request body validation error")
	return &fault
}

func NewValidatorFaultFromDecoder(logger *zerolog.Logger, cause error) Fault {
	e := logger.Warn().AnErr("cause", cause)

	var code string
	var message string
	switch {
	case strings.Contains(cause.Error(), "json: unknown field"):
		code = "UNKNOWN_FIELD"
		message = "Cannot unmarshall the provided JSON : unknown field"
	case strings.Contains(cause.Error(), "looking for beginning of value") ||
		strings.Contains(cause.Error(), "unexpected end of JSON input") ||
		strings.Contains(cause.Error(), "invalid character") ||
		strings.Contains(cause.Error(), "unexpected EOF"):
		code = "MALFORMED_JSON"
		message = "Cannot unmarshall the provided JSON because its malformed"
	case strings.Contains(cause.Error(), "EOF"):
		code = "EMPTY_JSON"
		message = "Cannot unmarshall the provided JSON because it's empty"
	case regexpWrongType.FindStringSubmatch(cause.Error()) != nil:
		code = "WRONG_TYPE"
		message = "Cannot unmarshall the provided JSON because a wrong type is used"
	case strings.Contains(cause.Error(), "json: Unmarshal(nil"):
		code = "INTERNAL_MARSHALING_ERROR"
		message = "Cannot unmarshall the provided JSON because of an internal error"
		e.Str("error_reason", "You provided a nil pointer to the Decode function")
	case strings.Contains(cause.Error(), "json: Unmarshal(non-pointer"):
		code = "INTERNAL_MARSHALING_ERROR"
		message = "Cannot unmarshall the provided JSON because of an internal error"
		e.Str("error_reason", "You provided a real object to the Decode function, expected a pointer")
	default:
		code = "BAD_REQUEST"
		message = "Cannot unmarshall the provided JSON"
	}

	fault := ValidatorFault{code: code, message: message, metadata: map[string]any{
		"unmarshall": map[string]any{"message": cause.Error()},
	}, cause: cause}
	e.Err(&fault).Msg(message)
	return &fault
}
