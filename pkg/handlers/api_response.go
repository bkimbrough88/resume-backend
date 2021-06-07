package handlers

import (
	"encoding/json"
	"github.com/bkimbrough88/resume-backend/pkg/models"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

func apiResponse(status int, body interface{}, logger *zap.Logger) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, err := json.Marshal(body)
	if err != nil {
		logger.Error("Failed to marshal body", zap.Error(err), zap.Any("body", body))
	}

	resp.Body = string(stringBody)
	return &resp, nil
}

func getErrorStatusCode(err error) int {
	if err.Error() == models.ErrorInvalidEmail || err.Error() == models.ErrorInvalidUserId {
		return http.StatusBadRequest
	}

	if err.Error() == models.ErrorNoResultsFound {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
