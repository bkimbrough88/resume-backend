package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/bkimbrough88/resume-backend/pkg/handlers"
	"go.uber.org/zap"
	"log"
	"os"
)

var (
	svc    dynamodbiface.DynamoDBAPI
	logger *zap.Logger
)

func main() {
	loggerProduction, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initiate logger. Error: %s", err.Error())
	}
	logger = loggerProduction

	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		logger.Error("Failed to establish new AWS session", zap.Error(err))
		return
	}
	svc = dynamodb.New(awsSession)
	lambda.Start(handler)
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	logger.Info("Received request", zap.Any("request", req))
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetUser(req, svc, logger)
	case "POST":
		return handlers.PutUser(req, svc, logger)
	case "DELETE":
		return handlers.DeleteUser(req, svc, logger)
	default:
		return handlers.UnhandledMethod(req, logger)
	}
}
