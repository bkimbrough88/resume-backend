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

var svc *dynamodb.DynamoDB
var logger *zap.Logger

type Event struct {
	Action string       `json:"action"`
	Id     *string      `json:"id,omitempty"`
	User   *resume.User `json:"user,omitempty"`
}

func init() {
	loggerProduction, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initiate logger. Error: %s", err.Error())
	}
	logger = loggerProduction

	sess := session.Must(session.NewSession())
	svc = dynamodb.New(sess)
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
		if request.Id == nil {
			logger.Error("No id was provided", zap.Any("request", request))
			return "", fmt.Errorf("no id was provided")
		}

		user, err := resume.GetUserById(*request.Id, svc, logger)
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
		if request.Id == nil {
			logger.Error("No id was provided", zap.Any("request", request))
			return "", fmt.Errorf("no id was provided")
		}

		if request.User == nil {
			logger.Error("No user to update was provided", zap.Any("request", request))
			return "", fmt.Errorf("no user to update was provided")
		}

		err := resume.UpdateUser(*request.Id, request.User, svc, logger)
		if err != nil {
			logger.Error("Failed to update user", zap.Error(err), zap.Any("request", request))
			return "", err
		}

		return "", nil
	} else if strings.EqualFold(request.Action, "delete") {
		if request.Id == nil {
			logger.Error("No id was provided", zap.Any("request", request))
			return "", fmt.Errorf("no id was provided")
		}

		err := resume.DeleteUser(*request.Id, svc, logger)
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
