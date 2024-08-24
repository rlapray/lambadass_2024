package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/lambadass-2024/backend/internal/adapters/repositories"
	validatorcommand "github.com/lambadass-2024/backend/internal/commands/validator"
	"github.com/lambadass-2024/backend/internal/fault"
	lambdaframework "github.com/lambadass-2024/backend/internal/frameworks/lambda"
	loggerframework "github.com/lambadass-2024/backend/internal/frameworks/logger"
	sqlframework "github.com/lambadass-2024/backend/internal/frameworks/sql"
	"github.com/lambadass-2024/backend/internal/usecases"
)

var (
	Logger        = loggerframework.APIGatewayClient{}
	Lambda        = lambdaframework.APIGatewayClient{}
	SQL           = sqlframework.GenericClient[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]{}
	PetRepository = repositories.PetRepository[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]{SQL: &SQL}
	PetUseCase    = usecases.PetUseCase[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]{Repository: &PetRepository}
	Validator     = validatorcommand.LambdaValidator[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]{}
)

type Body struct {
	ID uuid.UUID `json:"id" validate:"required,uuid4"`
}

func HandleRequest(
	_ context.Context,
	request events.APIGatewayProxyRequest, //nolint: gocritic // provided by aws
) (events.APIGatewayProxyResponse, fault.Fault) {
	var id uuid.UUID
	var err error

	if len(request.QueryStringParameters) > 1 {
		return Lambda.KOFromValidatorFault(
			fault.NewValidatorFault(&Logger.Client.Logger, "BAD_REQUEST", "Provide only ID", nil, nil))
	}

	if idstring, exists := request.QueryStringParameters["id"]; exists {
		id, err = uuid.Parse(idstring)
		if err != nil {
			return Lambda.KOFromValidatorFault(
				fault.NewValidatorFault(&Logger.Client.Logger, "BAD_REQUEST", "Can't parse ID", nil, err))
		}

		pet, err2 := PetUseCase.Get(id)
		if err2 == nil {
			return Lambda.OK(pet)
		}
		switch err2.Code() {
		case "PET_NOT_FOUND":
			return Lambda.KOFromFault(404, err2)
		default:
			return Lambda.KOFromFault(500, err2)
		}
	}
	return Lambda.KOFromValidatorFault(
		fault.NewValidatorFault(&Logger.Client.Logger, "BAD_REQUEST", "Provide ID", nil, nil))
}
