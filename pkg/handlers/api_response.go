package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

func apiResponse(status int, body interface{}, logger *zap.Logger) (*events.APIGatewayV2HTTPResponse, error) {
	resp := events.APIGatewayV2HTTPResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, err := json.Marshal(body)
	if err != nil {
		logger.Error("Failed to marshal body", zap.Error(err), zap.Any("body", body))
	}

	resp.Body = string(stringBody)
	return &resp, nil
}

func getErrorStatusCode(err error) int {
	if err.Error() == "invalid email" || err.Error() == "invalid user_id" {
		return http.StatusBadRequest
	}

	if err.Error() == "no results found" {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
