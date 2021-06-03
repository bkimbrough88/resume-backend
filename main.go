package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	resume "github.com/bkimbrough88/resume-backend/pkg"
	"go.uber.org/zap"
)

var svc resume.DynamoService
var logger *zap.Logger

type Event struct {
	Action string          `json:"action"`
	Key    *resume.UserKey `json:"id,omitempty"`
	User   *resume.User    `json:"user,omitempty"`
}

type EventResponse struct {
	StatusCode int          `json:"status_code"`
	Error      string       `json:"error,omitempty"`
	User       *resume.User `json:"user,omitempty"`
}

func init() {
	loggerProduction, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initiate logger. Error: %s", err.Error())
	}
	logger = loggerProduction

	sess := session.Must(session.NewSession())
	dynamo := dynamodb.New(sess)
	svc = resume.NewDynamoService(dynamo)
}

func HandleRequest(request Event) (string, error) {
	if strings.EqualFold(request.Action, "create") {
		if request.User == nil {
			logger.Error("No user data was provided", zap.Any("request", request))
			return "", getErrorResponse(http.StatusBadRequest, "no user data was provided")
		}

		err := resume.CreateUser(request.User, svc, logger)
		if err != nil {
			logger.Error("Failed to create user", zap.Error(err), zap.Any("request", request))
			return "", constructErrorResponse(err)
		}

		return getSuccessResponse(http.StatusCreated, nil)
	} else if strings.EqualFold(request.Action, "get") {
		if request.Key == nil {
			logger.Error("No key was provided", zap.Any("request", request))
			return "", getErrorResponse(http.StatusBadRequest, "no key was provided")
		}

		user, err := resume.GetUserByKey(request.Key, svc, logger)
		if err != nil {
			logger.Error("Failed to get user", zap.Error(err), zap.Any("request", request))
			return "", constructErrorResponse(err)
		}

		return getSuccessResponse(http.StatusOK, user)
	} else if strings.EqualFold(request.Action, "update") {
		if request.Key == nil {
			logger.Error("No key was provided", zap.Any("request", request))
			return "", getErrorResponse(http.StatusBadRequest, "no key was provided")
		}

		if request.User == nil {
			logger.Error("No user to update was provided", zap.Any("request", request))
			return "", getErrorResponse(http.StatusBadRequest, "no user to update was provided")
		}

		err := resume.UpdateUser(request.Key, request.User, svc, logger)
		if err != nil {
			logger.Error("Failed to update user", zap.Error(err), zap.Any("request", request))
			return "", constructErrorResponse(err)
		}

		return getSuccessResponse(http.StatusAccepted, nil)
	} else if strings.EqualFold(request.Action, "delete") {
		if request.Key == nil {
			logger.Error("No key was provided", zap.Any("request", request))
			return "", getErrorResponse(http.StatusBadRequest, "no key was provided")
		}

		err := resume.DeleteUser(request.Key, svc, logger)
		if err != nil {
			logger.Error("Failed to delete user", zap.Error(err), zap.Any("request", request))
			return "", constructErrorResponse(err)
		}

		return getSuccessResponse(http.StatusAccepted, nil)
	}

	logger.Error("Action not recognized", zap.Any("request", request))
	return "", getErrorResponse(http.StatusBadRequest, "action not recognized")
}

func constructErrorResponse(err error) error {
	if err.Error() == "invalid email" {
		return getErrorResponse(http.StatusBadRequest, err.Error())
	}

	if err.Error() == "no results found" || err.Error() == "too many results found" {
		return getErrorResponse(http.StatusNotFound, err.Error())
	}

	return getErrorResponse(http.StatusInternalServerError, err.Error())
}

func getErrorResponse(statusCode int, reason string) error {
	response := EventResponse{
		StatusCode: statusCode,
		Error:      reason,
	}

	if body, err := json.Marshal(response); err != nil {
		logger.Error("Failed to marshal response JSON", zap.Error(err))
		return fmt.Errorf("{ \"status_code\": %d, \"error\": \"%s\" }", statusCode, reason)
	} else {
		return fmt.Errorf("%s", body)
	}
}

func getSuccessResponse(statusCode int, user *resume.User) (string, error) {
	response := EventResponse{
		StatusCode: statusCode,
		User:       user,
	}

	if body, err := json.Marshal(response); err != nil {
		logger.Error("Failed to marshal response JSON", zap.Error(err))
		if user == nil {
			return fmt.Sprintf("{ \"status_code\": %d }", statusCode), nil
		}
		return "", constructErrorResponse(err)
	} else {
		return string(body), nil
	}
}

func main() {
	lambda.Start(HandleRequest)
}
