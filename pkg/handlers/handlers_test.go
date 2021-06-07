package handlers

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	mocks "github.com/bkimbrough88/resume-backend/pkg"
	"github.com/bkimbrough88/resume-backend/pkg/models"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	user   *models.User
	svc    dynamodbiface.DynamoDBAPI
)

func setupHandler(t *testing.T) {
	logger, _ = zap.NewDevelopment()
	user = &models.User{
		UserId: "user1",
		Email:  "user1@domain.com",
	}

	svc = mocks.DynamoServiceMock{}

	mocks.DeleteItemMock = func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
		return &dynamodb.DeleteItemOutput{}, nil
	}

	attr, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		t.Fatalf("Failed to marshal user int Dynamo attribute map: %s", err.Error())
	}
	mocks.GetItemMock = func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{Item: attr}, nil
	}

	mocks.PutItemMock = func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{}, nil
	}

	mocks.UpdateItemMock = func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
		return &dynamodb.UpdateItemOutput{}, nil
	}
}

func TestGetUser(t *testing.T) {
	setupHandler(t)

	event := events.APIGatewayProxyRequest{
		Resource:   "/user/{id}",
		Path:       "/v1/user/user1",
		HTTPMethod: "GET",
		PathParameters: map[string]string{
			"id": "user1",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourceID:   "GET /user/{id}",
			Stage:        "v1",
			ResourcePath: "/user/{id}",
			HTTPMethod:   "GET",
		},
		Body:            "",
		IsBase64Encoded: false,
	}
	if res, err := GetUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else if http.StatusOK != res.StatusCode {
		t.Errorf("Expected status code to be %d, but was %d", http.StatusOK, res.StatusCode)
	}

	mocks.GetItemMock = func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{},
		}, nil
	}
	if res, err := GetUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else if http.StatusNotFound != res.StatusCode {
		t.Errorf("Expected status code to be %d, but was %d", http.StatusNotFound, res.StatusCode)
	}

	event = events.APIGatewayProxyRequest{
		Resource:   "/user",
		Path:       "/v1/user",
		HTTPMethod: "GET",
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourceID:   "GET /user",
			Stage:        "v1",
			ResourcePath: "/user",
			HTTPMethod:   "GET",
		},
		Body:            "",
		IsBase64Encoded: false,
	}
	if res, err := GetUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else {
		if http.StatusBadRequest != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusBadRequest, res.StatusCode)
		}

		errorBody := &ErrorBody{}
		if jsonErr := json.Unmarshal([]byte(res.Body), errorBody); jsonErr != nil {
			t.Errorf("Failed to covert body to error body objec: %s", jsonErr.Error())
		} else if ErrorUserIdNotProvided != *errorBody.ErrorMsg {
			t.Errorf("Expected error to be '%s', but was '%s'", ErrorUserIdNotProvided, *errorBody.ErrorMsg)
		}
	}
}

func TestPutUser(t *testing.T) {
	setupHandler(t)

	userStr, jsonErr := json.Marshal(user)
	if jsonErr != nil {
		t.Fatalf("Failed to convert user to JSON")
	}

	event := events.APIGatewayProxyRequest{
		Resource:   "/user",
		Path:       "/v1/user",
		HTTPMethod: "POST",
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourceID:   "POST /user",
			Stage:        "v1",
			ResourcePath: "/user",
			HTTPMethod:   "POST",
		},
		Body:            string(userStr),
		IsBase64Encoded: false,
	}
	if res, err := PutUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else if http.StatusAccepted != res.StatusCode {
		t.Errorf("Expected status code to be %d, but was %d", http.StatusAccepted, res.StatusCode)
	}

	user.Email = "not an email"
	userStr, _ = json.Marshal(user)
	event.Body = string(userStr)
	if res, err := PutUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else {
		if http.StatusBadRequest != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusBadRequest, res.StatusCode)
		}

		errorBody := &ErrorBody{}
		if jsonErr := json.Unmarshal([]byte(res.Body), errorBody); jsonErr != nil {
			t.Errorf("Failed to covert body to error body object: %s", jsonErr.Error())
		} else if models.ErrorInvalidEmail != *errorBody.ErrorMsg {
			t.Errorf("Expected error to be '%s', but was '%s'", models.ErrorInvalidEmail, *errorBody.ErrorMsg)
		}
	}

	event.Body = ""
	if res, err := PutUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else {
		if http.StatusBadRequest != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusBadRequest, res.StatusCode)
		}

		errorBody := &ErrorBody{}
		if jsonErr := json.Unmarshal([]byte(res.Body), errorBody); jsonErr != nil {
			t.Errorf("Failed to covert body to error body objec: %s", jsonErr.Error())
		} else if ErrorUserNotProvided != *errorBody.ErrorMsg {
			t.Errorf("Expected error to be '%s', but was '%s'", ErrorUserNotProvided, *errorBody.ErrorMsg)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	setupHandler(t)

	event := events.APIGatewayProxyRequest{
		Resource:   "/user/{id}",
		Path:       "/v1/user/user1",
		HTTPMethod: "DELETE",
		PathParameters: map[string]string{
			"id": "user1",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourceID:   "DELETE /user/{id}",
			Stage:        "v1",
			ResourcePath: "/user/{id}",
			HTTPMethod:   "DELETE",
		},
		Body:            "",
		IsBase64Encoded: false,
	}
	if res, err := DeleteUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else if http.StatusAccepted != res.StatusCode {
		t.Errorf("Expected status code to be %d, but was %d", http.StatusAccepted, res.StatusCode)
	}

	expectedError := "some error"
	mocks.DeleteItemMock = func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
		return nil, errors.New(expectedError)
	}
	if res, err := DeleteUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else {
		if http.StatusInternalServerError != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusInternalServerError, res.StatusCode)
		}

		errorBody := &ErrorBody{}
		if jsonErr := json.Unmarshal([]byte(res.Body), errorBody); jsonErr != nil {
			t.Errorf("Failed to covert body to error body objec: %s", jsonErr.Error())
		} else if expectedError != *errorBody.ErrorMsg {
			t.Errorf("Expected error to be '%s', but was '%s'", expectedError, *errorBody.ErrorMsg)
		}
	}

	event = events.APIGatewayProxyRequest{
		Resource:   "/user",
		Path:       "/v1/user",
		HTTPMethod: "DELETE",
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourceID:   "DELETE /user",
			Stage:        "v1",
			ResourcePath: "/user",
			HTTPMethod:   "DELETE",
		},
		Body:            "",
		IsBase64Encoded: false,
	}
	if res, err := DeleteUser(event, svc, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else {
		if http.StatusBadRequest != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusBadRequest, res.StatusCode)
		}

		errorBody := &ErrorBody{}
		if jsonErr := json.Unmarshal([]byte(res.Body), errorBody); jsonErr != nil {
			t.Errorf("Failed to covert body to error body objec: %s", jsonErr.Error())
		} else if ErrorUserIdNotProvided != *errorBody.ErrorMsg {
			t.Errorf("Expected error to be '%s', but was '%s'", ErrorUserIdNotProvided, *errorBody.ErrorMsg)
		}
	}
}

func TestMethodNotAllowed(t *testing.T) {
	setupHandler(t)

	event := events.APIGatewayProxyRequest{
		Resource:   "/user/{id}",
		Path:       "/v1/user/user1",
		HTTPMethod: "PATCH",
		PathParameters: map[string]string{
			"id": "user1",
		},
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourceID:   "PATCH /user/{id}",
			Stage:        "v1",
			ResourcePath: "/user/{id}",
			HTTPMethod:   "PATCH",
		},
		Body:            "",
		IsBase64Encoded: false,
	}
	if res, err := UnhandledMethod(event, logger); err != nil {
		t.Errorf("Failed to get a response for GetUser: %s", err.Error())
	} else if res == nil {
		t.Errorf("Expected to have a response, but it was nil")
	} else {
		if http.StatusMethodNotAllowed != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusMethodNotAllowed, res.StatusCode)
		}

		errorBody := &ErrorBody{}
		if jsonErr := json.Unmarshal([]byte(res.Body), errorBody); jsonErr != nil {
			t.Errorf("Failed to covert body to error body objec: %s", jsonErr.Error())
		} else if ErrorMethodNotAllowed != *errorBody.ErrorMsg {
			t.Errorf("Expected error to be '%s', but was '%s'", ErrorMethodNotAllowed, *errorBody.ErrorMsg)
		}
	}
}
