package main_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	. "github.com/lambadass-2024/backend/cmd/functions/pet-POST/handler"
	"github.com/lambadass-2024/backend/internal/adapters/repositories"
	"github.com/lambadass-2024/backend/internal/entities"
	"github.com/lambadass-2024/backend/internal/fault"
	lambdaframework "github.com/lambadass-2024/backend/internal/frameworks/lambda"
	sqlframework "github.com/lambadass-2024/backend/internal/frameworks/sql"
	"github.com/lambadass-2024/backend/internal/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

/******************************************************************************
***** Tests preparation
******************************************************************************/

var sqlMock = sqlframework.MockClient[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]{}

func Before() *lambdaframework.Lambda[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse] {
	//zerolog.SetGlobalLevel(zerolog.Disabled)
	PetRepository.SQL = &sqlMock
	PetUseCase.Repository = &PetRepository

	return Lambda.
		Use(&Logger).
		Use(&Lambda).
		Use(&sqlMock).
		Use(&PetRepository).
		Use(&PetUseCase).
		Use(&Validator).
		Use(&Mockerie[events.APIGatewayProxyRequest, events.APIGatewayProxyResponse]{})
}

type Mockerie[T any, U any] struct{}

func (u *Mockerie[T, U]) OnSetup(_ context.Context, _ *T) fault.Fault {
	PetUseCase.UUIDGenerator = &utils.MockUUIDGenerator{}
	return nil
}
func (u Mockerie[T, U]) OnBefore(_ context.Context, _ *T) fault.Fault {
	return nil
}

func (u Mockerie[T, U]) OnAfter(_ *U, err fault.Fault) fault.Fault {
	return err
}

func (u Mockerie[T, U]) OnShutdown() {
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

/******************************************************************************
***** Tests
******************************************************************************/

func TestPetPostOKWithID(t *testing.T) {
	lambda := Before()

	pet := entities.Pet{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3"), Name: "a", Race: entities.Race{ID: uuid.MustParse("752cd6644267493eb8311d4587abf000")}}

	key := sqlframework.ExecMapKey{Q: repositories.PetSQLCreate, D: pet}
	value := sqlframework.ExecOneRowAffectedMapValue{F: nil}
	sqlMock.MockExecOneRowAffectedMap(key, value)

	request := events.APIGatewayProxyRequest{Body: `{"id": "752cd6644267493eb8311d4587abf5b3", "name":"a", "raceId": "752cd6644267493eb8311d4587abf000"}`}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"id\":\"752cd664-4267-493e-b831-1d4587abf5b3\",\"name\":\"a\",\"race\":{\"id\":\"752cd664-4267-493e-b831-1d4587abf000\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostOKWithoutID(t *testing.T) {
	lambda := Before()

	pet := entities.Pet{ID: uuid.MustParse("11111111111111111111111111111111"), Name: "a", Race: entities.Race{ID: uuid.MustParse("752cd6644267493eb8311d4587abf000")}}

	key := sqlframework.ExecMapKey{Q: repositories.PetSQLCreate, D: pet}
	value := sqlframework.ExecOneRowAffectedMapValue{F: nil}
	sqlMock.MockExecOneRowAffectedMap(key, value)

	request := events.APIGatewayProxyRequest{Body: `{"name":"a", "raceId": "752cd6644267493eb8311d4587abf000"}`}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"id\":\"11111111-1111-1111-1111-111111111111\",\"name\":\"a\",\"race\":{\"id\":\"752cd664-4267-493e-b831-1d4587abf000\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKONotUnique(t *testing.T) {
	lambda := Before()
	logger := zerolog.Logger{}

	pet := entities.Pet{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3"), Name: "a", Race: entities.Race{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3")}}

	key := sqlframework.ExecMapKey{Q: repositories.PetSQLCreate, D: pet}
	value := sqlframework.ExecOneRowAffectedMapValue{F: fault.NewSQL(&logger, "UNIQUE_VIOLATION", "", nil, nil)}
	sqlMock.MockExecOneRowAffectedMap(key, value)

	request := events.APIGatewayProxyRequest{Body: `{"id": "752cd6644267493eb8311d4587abf5b3", "name":"a", "raceId": "752cd6644267493eb8311d4587abf5b3"}`}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":422,\"code\":\"PET_ID_NOT_UNIQUE\",\"message\":\"Pet id not unique\",\"metadata\":{\"id\":\"752cd664-4267-493e-b831-1d4587abf5b3\",\"requestId\":\"\",\"requestTime\":\"\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKOUnexpectedError(t *testing.T) {
	lambda := Before()
	logger := zerolog.Logger{}

	pet := entities.Pet{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3"), Name: "a", Race: entities.Race{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3")}}

	key := sqlframework.ExecMapKey{Q: repositories.PetSQLCreate, D: pet}
	value := sqlframework.ExecOneRowAffectedMapValue{F: fault.NewSQL(&logger, "DUMMY_ERROR_UNEXPECTED", "", nil, nil)}
	sqlMock.MockExecOneRowAffectedMap(key, value)

	request := events.APIGatewayProxyRequest{Body: `{"id": "752cd6644267493eb8311d4587abf5b3", "name":"a", "raceId": "752cd6644267493eb8311d4587abf5b3"}`}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":500,\"code\":\"PET_CREATION_FAILED\",\"message\":\"Pet creation failed\",\"metadata\":{\"id\":\"752cd664-4267-493e-b831-1d4587abf5b3\",\"requestId\":\"\",\"requestTime\":\"\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKOValidationBadUUID(t *testing.T) {
	lambda := Before()

	pet := entities.Pet{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3"), Name: "a", Race: entities.Race{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3")}}

	key := sqlframework.ExecMapKey{Q: repositories.PetSQLCreate, D: pet}
	value := sqlframework.ExecOneRowAffectedMapValue{F: nil}
	sqlMock.MockExecOneRowAffectedMap(key, value)

	request := events.APIGatewayProxyRequest{Body: `{"id": "752cd6644267493eb8311d4587abf5b3_____", "name":"a", "raceId": "752cd6644267493eb8311d4587abf5b3"}`}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":400,\"code\":\"BAD_REQUEST\",\"message\":\"Cannot unmarshall the provided JSON\",\"metadata\":{\"requestId\":\"\",\"requestTime\":\"\",\"unmarshall\":{\"message\":\"invalid UUID length: 37\"}}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKOValidationNoName(t *testing.T) {
	lambda := Before()

	pet := entities.Pet{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3"), Name: "a", Race: entities.Race{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3")}}

	key := sqlframework.ExecMapKey{Q: repositories.PetSQLCreate, D: pet}
	value := sqlframework.ExecOneRowAffectedMapValue{F: nil}
	sqlMock.MockExecOneRowAffectedMap(key, value)

	request := events.APIGatewayProxyRequest{Body: `{"id": "752cd6644267493eb8311d4587abf5b3", "raceId": "752cd6644267493eb8311d4587abf5b3"}`}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":400,\"code\":\"BAD_REQUEST\",\"message\":\"Validation failed\",\"metadata\":{\"requestId\":\"\",\"requestTime\":\"\",\"validation\":[{\"message\":\"Key: 'Body.Name' Error:Field validation for 'Name' failed on the 'required' tag\",\"field\":\"Name\",\"namespace\":\"Body.Name\",\"tag\":\"required\",\"value\":\"\"}]}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKOValidationNoRace(t *testing.T) {
	lambda := Before()

	pet := entities.Pet{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3"), Name: "a", Race: entities.Race{ID: uuid.MustParse("752cd6644267493eb8311d4587abf5b3")}}

	key := sqlframework.ExecMapKey{Q: repositories.PetSQLCreate, D: pet}
	value := sqlframework.ExecOneRowAffectedMapValue{F: nil}
	sqlMock.MockExecOneRowAffectedMap(key, value)

	request := events.APIGatewayProxyRequest{Body: `{"id": "752cd6644267493eb8311d4587abf5b3", "name": "ziiugf"}`}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":400,\"code\":\"BAD_REQUEST\",\"message\":\"Validation failed\",\"metadata\":{\"requestId\":\"\",\"requestTime\":\"\",\"validation\":[{\"message\":\"Key: 'Body.RaceID' Error:Field validation for 'RaceID' failed on the 'required' tag\",\"field\":\"RaceID\",\"namespace\":\"Body.RaceID\",\"tag\":\"required\",\"value\":\"00000000-0000-0000-0000-000000000000\"}]}}", response.Body)

	assert.NoError(t, f)
}
