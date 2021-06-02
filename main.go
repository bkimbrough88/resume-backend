package main

import (
	"encoding/json"
	"fmt"
	"log"
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
			return "", fmt.Errorf("no user data was provided")
		}

		err := resume.CreateUser(request.User, svc, logger)
		if err != nil {
			logger.Error("Failed to create user", zap.Error(err), zap.Any("request", request))
			return "", err
		}

		return "", nil
	} else if strings.EqualFold(request.Action, "get") {
		if request.Key == nil {
			logger.Error("No key was provided", zap.Any("request", request))
			return "", fmt.Errorf("no key was provided")
		}

		user, err := resume.GetUserByKey(request.Key, svc, logger)
		if err != nil {
			logger.Error("Failed to get user", zap.Error(err), zap.Any("request", request))
			return "", err
		}

		userJson, err := json.Marshal(user)
		if err != nil {
			logger.Error("Failed to marshal user to JSON", zap.Error(err))
			return "", err
		}

		return string(userJson), nil
	} else if strings.EqualFold(request.Action, "update") {
		if request.Key == nil {
			logger.Error("No id was provided", zap.Any("request", request))
			return "", fmt.Errorf("no id was provided")
		}

		if request.User == nil {
			logger.Error("No user to update was provided", zap.Any("request", request))
			return "", fmt.Errorf("no user to update was provided")
		}

		err := resume.UpdateUser(request.Key, request.User, svc, logger)
		if err != nil {
			logger.Error("Failed to update user", zap.Error(err), zap.Any("request", request))
			return "", err
		}

		return "", nil
	} else if strings.EqualFold(request.Action, "delete") {
		if request.Key == nil {
			logger.Error("No id was provided", zap.Any("request", request))
			return "", fmt.Errorf("no id was provided")
		}

		err := resume.DeleteUser(request.Key, svc, logger)
		if err != nil {
			logger.Error("Failed to delete user", zap.Error(err), zap.Any("request", request))
			return "", err
		}

		return "", nil
	}

	logger.Error("Action not recognized", zap.Any("request", request))
	return "", fmt.Errorf("action not recognized")
}

func main() {
	lambda.Start(HandleRequest)
}
