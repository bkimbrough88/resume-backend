package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/bkimbrough88/resume-backend/pkg/models"
)

const (
	ErrorMethodNotAllowed  = "method not allowed"
	ErrorUserIdNotProvided = "userId not provided"
	ErrorUserNotProvided   = "user not provided in body"
)

type SuccessBody struct {
	User *models.User `json:"user,omitempty"`
}

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) (*events.APIGatewayProxyResponse, error) {
	userId := req.PathParameters["id"]
	if len(userId) > 0 {
		key := &models.UserKey{UserId: userId}
		user, err := models.GetUserByKey(key, svc, logger)
		if err != nil {
			return apiResponse(getErrorStatusCode(err), ErrorBody{aws.String(err.Error())}, logger)
		}

		return apiResponse(http.StatusOK, SuccessBody{User: user}, logger)
	} else {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(ErrorUserIdNotProvided)}, logger)
	}
}

func PutUser(req events.APIGatewayProxyRequest, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) (*events.APIGatewayProxyResponse, error) {
	if len(req.Body) > 0 {
		user := &models.User{}
		if err := json.Unmarshal([]byte(req.Body), user); err != nil {
			logger.Error("Failed to unmarshal body into User object", zap.Error(err), zap.String("body", req.Body))
			return apiResponse(getErrorStatusCode(err), ErrorBody{aws.String(err.Error())}, logger)
		}

		if err := models.PutUser(user, svc, logger); err != nil {
			return apiResponse(getErrorStatusCode(err), ErrorBody{aws.String(err.Error())}, logger)
		}

		return apiResponse(http.StatusAccepted, SuccessBody{}, logger)
	} else {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(ErrorUserNotProvided)}, logger)
	}
}

func DeleteUser(req events.APIGatewayProxyRequest, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) (*events.APIGatewayProxyResponse, error) {
	userId := req.PathParameters["id"]
	if len(userId) > 0 {
		key := &models.UserKey{UserId: userId}
		if err := models.DeleteUser(key, svc, logger); err != nil {
			return apiResponse(getErrorStatusCode(err), ErrorBody{aws.String(err.Error())}, logger)
		}

		return apiResponse(http.StatusAccepted, SuccessBody{}, logger)
	} else {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(ErrorUserIdNotProvided)}, logger)
	}
}

func UnhandledMethod(req events.APIGatewayProxyRequest, logger *zap.Logger) (*events.APIGatewayProxyResponse, error) {
	logger.Warn("Method not allowed", zap.String("method", req.HTTPMethod))
	return apiResponse(http.StatusMethodNotAllowed, ErrorBody{ErrorMsg: aws.String(ErrorMethodNotAllowed)}, logger)
}
