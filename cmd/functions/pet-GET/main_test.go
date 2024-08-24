package main_test

import (
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	. "github.com/lambadass-2024/backend/cmd/functions/pet-GET/handler"
	"github.com/lambadass-2024/backend/internal/adapters/repositories"
	"github.com/lambadass-2024/backend/internal/entities"
	lambdaframework "github.com/lambadass-2024/backend/internal/frameworks/lambda"
	sqlframework "github.com/lambadass-2024/backend/internal/frameworks/sql"
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
		Use(&Validator)
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

	id := "752cd6644267493eb8311d4587abf5b3"

	petIn := entities.Pet{ID: uuid.MustParse(id)}
	petOutReal := []entities.Pet{}
	petOutReal = append(petOutReal, entities.Pet{ID: uuid.MustParse(id), Name: "bang",
		Race: entities.Race{ID: uuid.MustParse("752cd6644267493eb8311d4587abf000"), Name: "hbzf"},
	})

	key := sqlframework.SelectMapKey{Q: repositories.PetSQLGet, DA: petIn}
	sqlMock.MockSelectMap(key, nil, petOutReal)

	request := events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"id": id}}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"id\":\"752cd664-4267-493e-b831-1d4587abf5b3\",\"name\":\"bang\",\"race\":{\"id\":\"752cd664-4267-493e-b831-1d4587abf000\",\"name\":\"hbzf\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKONoResult(t *testing.T) {
	lambda := Before()

	id := "752cd6644267493eb8311d4587abf5b3"

	petIn := entities.Pet{ID: uuid.MustParse(id)}
	petOutReal := []entities.Pet{}

	key := sqlframework.SelectMapKey{Q: repositories.PetSQLGet, DA: petIn}
	sqlMock.MockSelectMap(key, nil, petOutReal)

	request := events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"id": id}}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":404,\"code\":\"PET_NOT_FOUND\",\"message\":\"Pet not found\",\"metadata\":{\"id\":\"752cd664-4267-493e-b831-1d4587abf5b3\",\"requestId\":\"\",\"requestTime\":\"\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKONoID(t *testing.T) {
	lambda := Before()

	id := "752cd6644267493eb8311d4587abf5b3"

	petIn := entities.Pet{ID: uuid.MustParse(id)}
	petOutReal := []entities.Pet{}

	key := sqlframework.SelectMapKey{Q: repositories.PetSQLGet, DA: petIn}
	sqlMock.MockSelectMap(key, nil, petOutReal)

	request := events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{}}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":400,\"code\":\"BAD_REQUEST\",\"message\":\"Provide ID\",\"metadata\":{\"requestId\":\"\",\"requestTime\":\"\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKOWithIDAndOtherFields(t *testing.T) {
	lambda := Before()

	id := "752cd6644267493eb8311d4587abf5b3"

	petIn := entities.Pet{ID: uuid.MustParse(id)}
	petOutReal := []entities.Pet{}
	petOutReal = append(petOutReal, entities.Pet{ID: uuid.MustParse(id), Name: "bang"})

	key := sqlframework.SelectMapKey{Q: repositories.PetSQLGet, DA: petIn}
	sqlMock.MockSelectMap(key, nil, petOutReal)

	request := events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"id": id, "string1": "dad"}}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":400,\"code\":\"BAD_REQUEST\",\"message\":\"Provide only ID\",\"metadata\":{\"requestId\":\"\",\"requestTime\":\"\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKOWithInternalProblem(t *testing.T) {
	lambda := Before()

	id := "752cd6644267493eb8311d4587abf5b3"

	petIn := entities.Pet{ID: uuid.MustParse(id)}
	petOutReal := []entities.Pet{}
	petOutReal = append(petOutReal, entities.Pet{ID: uuid.MustParse(id), Name: "bang"})
	petOutReal = append(petOutReal, entities.Pet{ID: uuid.MustParse(id), Name: "bang"})

	key := sqlframework.SelectMapKey{Q: repositories.PetSQLGet, DA: petIn}
	sqlMock.MockSelectMap(key, nil, petOutReal)

	request := events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"id": id}}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":500,\"code\":\"PET_GET_FAILED\",\"message\":\"Cannot get this pet\",\"metadata\":{\"id\":\"752cd664-4267-493e-b831-1d4587abf5b3\",\"requestId\":\"\",\"requestTime\":\"\"}}", response.Body)

	assert.NoError(t, f)
}

func TestPetPostKOWithIDMalformed(t *testing.T) {
	lambda := Before()

	id := "752cd6644267493eb8311d4587abf5b3"

	petIn := entities.Pet{ID: uuid.MustParse(id)}
	petOutReal := []entities.Pet{}
	petOutReal = append(petOutReal, entities.Pet{ID: uuid.MustParse(id), Name: "bang"})

	key := sqlframework.SelectMapKey{Q: repositories.PetSQLGet, DA: petIn}
	sqlMock.MockSelectMap(key, nil, petOutReal)

	request := events.APIGatewayProxyRequest{QueryStringParameters: map[string]string{"id": "752cd6644267493eb8311d4587abf5b3aaaaaaaaaaa"}}

	response, f := lambda.TestHandleRequest(HandleRequest, &request)
	assert.NotNil(t, response)
	assert.Equal(t, "{\"statusCode\":400,\"code\":\"BAD_REQUEST\",\"message\":\"Can't parse ID\",\"metadata\":{\"requestId\":\"\",\"requestTime\":\"\"}}", response.Body)

	assert.NoError(t, f)
}
