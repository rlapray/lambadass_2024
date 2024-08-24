package lambda_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/lambadass-2024/backend/internal/commands/validator"
	"github.com/lambadass-2024/backend/internal/fault"
	"github.com/lambadass-2024/backend/internal/frameworks/lambda"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var apiGatewayProxyRequestContext = events.APIGatewayProxyRequestContext{
	RequestID:   "123",
	RequestTime: "time",
}

var metadataDefault = map[string]any{"key": "value"}

type BodyToValidate struct {
	ID   uuid.UUID `json:"id"   validate:"omitempty,uuid4"`
	Name string    `json:"name" validate:"required"`
}

func NewAPIGateway() lambda.APIGatewayClient {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	client := lambda.APIGatewayClient{
		Lambda: lambda.TestNewLambda(&events.APIGatewayProxyRequest{RequestContext: apiGatewayProxyRequestContext}),
	}
	return client
}

/******************************************************************************
***** KO
******************************************************************************/

func Test_APIGateway_KO_withMetadata(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	expectedStatusCode := 404
	expectedCode := "PET_NOT_FOUND"
	expectedMessage := "Cannot find your pet"
	expectedMetadata := metadataDefault

	response, err := apiGateway.KO(expectedStatusCode, expectedCode, expectedMessage, expectedMetadata)
	assert.NotNil(t, response)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
	require.Error(t, err)
	assert.Equal(t, err, fault.NewAPIGateway(&logger, expectedStatusCode, expectedCode, expectedMessage, expectedMetadata, nil))
}

func Test_APIGateway_KO_withoutMetadata(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	expectedStatusCode := 404
	expectedCode := "PET_NOT_FOUND"
	expectedMessage := "Cannot find your pet"
	response, err := apiGateway.KO(expectedStatusCode, expectedCode, expectedMessage, nil)
	assert.NotNil(t, response)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
	require.Error(t, err)
	assert.Equal(t, err, fault.NewAPIGateway(&logger, expectedStatusCode, expectedCode, expectedMessage, nil, nil))
}

func Test_APIGateway_KOFromFault_withMetadataAndError(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	metadata := metadataDefault

	f1 := fault.NewUseCase(&logger, "DummyUseCase", "CODE1", "Code 1", metadata, errors.New("code 1"))

	response, err := apiGateway.KOFromFault(422, f1)
	assert.NotNil(t, response)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
	require.Error(t, err)
	assert.Equal(t, err, fault.NewAPIGateway(&logger, 422, f1.Code(), f1.Message(), f1.Metadata(), f1))
}

func Test_APIGateway_KOFromFault_withMetadataAndWithoutError(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	metadata := metadataDefault

	f1 := fault.NewUseCase(&logger, "DummyUseCase", "CODE1", "Code 1", metadata, nil)

	response, err := apiGateway.KOFromFault(422, f1)
	assert.NotNil(t, response)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
	require.Error(t, err)
	assert.Equal(t, err, fault.NewAPIGateway(&logger, 422, f1.Code(), f1.Message(), f1.Metadata(), f1))
}

func Test_APIGateway_KOFromFault_withoutMetadataWithError(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	f1 := fault.NewUseCase(&logger, "DummyUseCase", "CODE1", "Code 1", nil, errors.New("code 1"))

	response, err := apiGateway.KOFromFault(422, f1)
	assert.NotNil(t, response)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
	require.Error(t, err)
	assert.Equal(t, err, fault.NewAPIGateway(&logger, 422, f1.Code(), f1.Message(), nil, f1))
}

func Test_APIGateway_KOFromFault_withoutMetadataWithoutError(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	f1 := fault.NewUseCase(&logger, "DummyUseCase2", "CODE2", "Code 2", nil, nil)

	response, err := apiGateway.KOFromFault(422, f1)
	assert.NotNil(t, response)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
	require.Error(t, err)
	assert.Equal(t, err, fault.NewAPIGateway(&logger, 422, f1.Code(), f1.Message(), nil, f1))
}

func Test_APIGateway_KOFromValidatorFault(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}
	v := validator.LambdaValidator[int, int]{}
	_ = v.OnSetup(context.Background(), nil)
	jsonString := "{\"id\": \"752cd6644267493eb8311d4587abf5b3\", \"name\":42}"
	var bodyToValidate BodyToValidate
	f2 := v.ValidateJSONIntoStruct(jsonString, bodyToValidate)

	response, err := apiGateway.KOFromValidatorFault(f2)
	assert.NotNil(t, response)
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
	require.Error(t, err)
	assert.Equal(t, err, fault.NewAPIGateway(&logger, 500, "UNEXPECTED_INPUT_VALIDATION_ERROR", "Validation raised an unexpected error", nil, f2))
}

/******************************************************************************
***** OK
******************************************************************************/
func Test_APIGateway_OK_200(t *testing.T) {
	apiGateway := NewAPIGateway()

	// Test case 1: Successful response 200
	expectedObj := metadataDefault

	response, err := apiGateway.OK(expectedObj)

	require.NoError(t, err)
	assert.Equal(t, 200, response.StatusCode)

	var responseBody map[string]any
	err2 := json.Unmarshal([]byte(response.Body), &responseBody)
	require.NoError(t, err2)
	assert.Equal(t, expectedObj, responseBody)
}

func Test_APIGateway_OK_204(t *testing.T) {
	apiGateway := NewAPIGateway()

	// Test case 2: Successful response 204
	response, err := apiGateway.OK(nil)
	require.NoError(t, err)
	assert.Equal(t, 204, response.StatusCode)
}

func Test_APIGateway_OK_BadObject(t *testing.T) {
	apiGateway := NewAPIGateway()

	// Test case 2: Successful response 204
	response, err := apiGateway.OK(make(chan struct{}))
	assert.Equal(t, events.APIGatewayProxyResponse{}, response)
	require.Error(t, err)

	agpe, ok := err.(*fault.APIGatewayProxyFault)
	assert.True(t, ok)
	assert.Equal(t, "ERROR_MARSHALL_JSON", agpe.Code())
}

/******************************************************************************
***** OnAfter
******************************************************************************/
func Test_APIGateway_OnAfter_ResponseNil(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	err := fault.NewAPIGateway(&logger, 500, "CODE", "Message", nil, nil)

	err = apiGateway.OnAfter(nil, err)
	require.Error(t, err)
}

func Test_APIGateway_OnAfter_Success(t *testing.T) {
	apiGateway := NewAPIGateway()

	response := &events.APIGatewayProxyResponse{StatusCode: 200, Body: "hello test"}
	expectedResponse := events.APIGatewayProxyResponse{StatusCode: 200, Body: "hello test"}

	err2 := apiGateway.OnAfter(response, nil)
	require.NoError(t, err2)

	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, expectedResponse.Body, response.Body)
}

func Test_APIGateway_OnAfter_SuccessRaisedError(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	response := &events.APIGatewayProxyResponse{}
	f1 := fault.NewUseCase(&logger, "DummyUseCase1", "ERROR1", "Error 1", nil, nil)

	expectedBody := "{\"statusCode\":500,\"code\":\"ERROR1\",\"message\":\"Error 1\",\"metadata\":{\"requestId\":\"123\",\"requestTime\":\"time\"}}"

	err2 := apiGateway.OnAfter(response, f1)
	require.NoError(t, err2)

	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, expectedBody, response.Body)
}

func Test_APIGateway_OnAfter_SuccessRaisedAPIGatewayProxyFaultError(t *testing.T) {
	apiGateway := NewAPIGateway()
	logger := zerolog.Logger{}

	response := &events.APIGatewayProxyResponse{}
	f1 := fault.NewAPIGateway(&logger, 422, "ERROR1", "Error 1", nil, nil)

	expectedBody := "{\"statusCode\":422,\"code\":\"ERROR1\",\"message\":\"Error 1\",\"metadata\":{\"requestId\":\"123\",\"requestTime\":\"time\"}}"

	err2 := apiGateway.OnAfter(response, f1)
	require.NoError(t, err2)

	assert.Equal(t, 422, response.StatusCode)
	assert.Equal(t, expectedBody, response.Body)
}
