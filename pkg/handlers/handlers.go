package handlers

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/bkimbrough88/resume-backend/pkg/models"
)

type SuccessBody struct {
	User *models.User `json:"user,omitempty"`
}

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayV2HTTPRequest, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) (*events.APIGatewayV2HTTPResponse, error) {
	userId := req.PathParameters["id"]
	if len(userId) > 0 {
		key := &models.UserKey{UserId: userId}
		user, err := models.GetUserByKey(key, svc, logger)
		if err != nil {
			return apiResponse(getErrorStatusCode(err), ErrorBody{aws.String(err.Error())}, logger)
		}

		return apiResponse(http.StatusOK, SuccessBody{User: user}, logger)
	} else {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String("userId not provided")}, logger)
	}
}

func PutUser(req events.APIGatewayV2HTTPRequest, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) (*events.APIGatewayV2HTTPResponse, error) {
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
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String("user not provided in body")}, logger)
	}
}

func DeleteUser(req events.APIGatewayV2HTTPRequest, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) (*events.APIGatewayV2HTTPResponse, error) {
	userId := req.PathParameters["id"]
	if len(userId) > 0 {
		key := &models.UserKey{UserId: userId}
		if err := models.DeleteUser(key, svc, logger); err != nil {
			return apiResponse(getErrorStatusCode(err), ErrorBody{aws.String(err.Error())}, logger)
		}

		return apiResponse(http.StatusAccepted, SuccessBody{}, logger)
	} else {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String("userId not provided")}, logger)
	}
}

func UnhandledMethod(req events.APIGatewayV2HTTPRequest, logger *zap.Logger) (*events.APIGatewayV2HTTPResponse, error) {
	logger.Warn("Method not allowed", zap.String("method", req.RequestContext.HTTP.Method))
	return apiResponse(http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"), logger)
}
