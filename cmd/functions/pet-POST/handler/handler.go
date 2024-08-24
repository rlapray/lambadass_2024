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
	ID     uuid.UUID `json:"id"     validate:"omitempty,uuid"`
	RaceID uuid.UUID `json:"raceId" validate:"uuid,required"`
	Name   string    `json:"name"   validate:"required"`
}

func HandleRequest(
	_ context.Context,
	request events.APIGatewayProxyRequest, //nolint: gocritic // provided by aws
) (events.APIGatewayProxyResponse, fault.Fault) {
	var data Body
	err := Validator.ValidateJSONIntoStruct(request.Body, &data)
	if err != nil {
		return Lambda.KOFromValidatorFault(err)
	}

	pet, err2 := PetUseCase.Create(data.ID, data.Name, data.RaceID)

	if err2 == nil {
		return Lambda.OK(pet)
	}

	switch err2.Code() {
	case "PET_ID_NOT_UNIQUE":
		return Lambda.KOFromFault(422, err2)
	default:
		return Lambda.KOFromFault(500, err2)
	}
}
