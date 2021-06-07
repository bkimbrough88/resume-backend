package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/bkimbrough88/resume-backend/pkg/models"
	"go.uber.org/zap"
)

const (
	contentType                = "Content-Type"
	contentTypeApplicationJson = "application/json"
	emptyResponse              = "{}"
)

func setupApiResponse() {
	logger, _ = zap.NewDevelopment()

	user = &models.User{
		UserId: "user",
		Email:  "user@domain.com",
	}
}

func TestApiResponse(t *testing.T) {
	setupApiResponse()

	successNoUser := SuccessBody{}
	if res, err := apiResponse(http.StatusOK, successNoUser, logger); err != nil {
		t.Errorf("Failed to get API response: %s", err.Error())
	} else {
		if http.StatusOK != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusOK, res.StatusCode)
		}

		if contentTypeApplicationJson != res.Headers[contentType] {
			t.Errorf("Expected %s header to be '%s', but was '%s'", contentType, contentTypeApplicationJson, res.Headers[contentType])
		}

		if emptyResponse != res.Body {
			t.Errorf("Expected body to be '%s', but was '%s'", emptyResponse, res.Body)
		}
	}

	successWithUser := SuccessBody{User: user}
	if res, err := apiResponse(http.StatusOK, successWithUser, logger); err != nil {
		t.Errorf("Failed to get API response: %s", err.Error())
	} else {
		if http.StatusOK != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusOK, res.StatusCode)
		}

		if contentTypeApplicationJson != res.Headers[contentType] {
			t.Errorf("Expected %s header to be '%s', but was '%s'", contentType, contentTypeApplicationJson, res.Headers[contentType])
		}

		resUser := &SuccessBody{}
		if unmarshalErr := json.Unmarshal([]byte(res.Body), resUser); unmarshalErr != nil {
			t.Errorf("Failed to unmarshal response body: %s", unmarshalErr.Error())
		} else if resUser.User != nil {
			if user.UserId != resUser.User.UserId {
				t.Errorf("Expected user_id to be '%s', but was '%s'", user.UserId, resUser.User.UserId)
			}

			if user.Email != resUser.User.Email {
				t.Errorf("Expected email to be '%s', but was '%s'", user.Email, resUser.User.Email)
			}
		} else {
			t.Errorf("Expected user to be present, but was nil")
		}
	}

	errorBody := ErrorBody{ErrorMsg: aws.String("bad request")}
	if res, err := apiResponse(http.StatusBadRequest, errorBody, logger); err != nil {
		t.Errorf("Failed to get API response: %s", err.Error())
	} else {
		if http.StatusBadRequest != res.StatusCode {
			t.Errorf("Expected status code to be %d, but was %d", http.StatusBadRequest, res.StatusCode)
		}

		if contentTypeApplicationJson != res.Headers[contentType] {
			t.Errorf("Expected %s header to be '%s', but was '%s'", contentType, contentTypeApplicationJson, res.Headers[contentType])
		}

		resErrorBody := &ErrorBody{}
		if unmarshalErr := json.Unmarshal([]byte(res.Body), resErrorBody); unmarshalErr != nil {
			t.Errorf("Failed to unmarshal response body: %s", unmarshalErr.Error())
		} else if resErrorBody.ErrorMsg != nil {
			if *errorBody.ErrorMsg != *resErrorBody.ErrorMsg {
				t.Errorf("Expected error message to be '%s', but was '%s'", *errorBody.ErrorMsg, *resErrorBody.ErrorMsg)
			}
		} else {
			t.Errorf("Expected error message to be present, but was nil")
		}
	}
}

func TestGetErrorStatusCode(t *testing.T) {
	setupApiResponse()

	if code := getErrorStatusCode(errors.New(models.ErrorInvalidEmail)); http.StatusBadRequest != code {
		t.Errorf("Expected status code for error '%s' to be %d, but was %d", models.ErrorInvalidEmail, http.StatusBadRequest, code)
	}

	if code := getErrorStatusCode(errors.New(models.ErrorInvalidUserId)); http.StatusBadRequest != code {
		t.Errorf("Expected status code for error '%s' to be %d, but was %d", models.ErrorInvalidUserId, http.StatusBadRequest, code)
	}

	if code := getErrorStatusCode(errors.New(models.ErrorNoResultsFound)); http.StatusNotFound != code {
		t.Errorf("Expected status code for error '%s' to be %d, but was %d", models.ErrorNoResultsFound, http.StatusNotFound, code)
	}

	if code := getErrorStatusCode(errors.New("some other error")); http.StatusInternalServerError != code {
		t.Errorf("Expected status code for error 'some other error' to be %d, but was %d", http.StatusInternalServerError, code)
	}
}
