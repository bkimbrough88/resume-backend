package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	resume "github.com/bkimbrough88/resume-backend/pkg"
	"go.uber.org/zap"
)

type RemoveItemEvent struct {
	Type string `json:"type"`
	Key interface{} `json:"item"`
}

func HandleRemoveRequest(request RemoveItemEvent) (bool, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	logger, loggerErr := zap.NewProduction()
	if loggerErr != nil {
		log.Printf("failed to initialize the logger. Error: %s", loggerErr.Error())
		return false, loggerErr
	}

	if strings.EqualFold(request.Type, "certification") {
		if key, ok := request.Key.(resume.CertificationKey); ok {
			input, err := resume.GetCertificationDeleteInput(key)
			if err != nil {
				logger.Error("Failed to get the certification delete input", zap.Error(err))
				return false, err
			}

			if _, err := svc.DeleteItem(input); err != nil {
				logger.Error("Failed to delete certification from database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("key is not a CertificationKey type", zap.Any("key", request.Key))
			return false, fmt.Errorf("key is not a CertificationKey type")
		}

		return true, nil
	} else if strings.EqualFold(request.Type, "degree") {
		if key, ok := request.Key.(resume.DegreeKey); ok {
			input, err := resume.GetDegreeDeleteInput(key)
			if err != nil {
				logger.Error("Failed to get the degree delete input", zap.Error(err))
				return false, err
			}

			if _, err := svc.DeleteItem(input); err != nil {
				logger.Error("Failed to delete degree from database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("key is not a DegreeKey type", zap.Any("key", request.Key))
			return false, fmt.Errorf("key is not a DegreeKey type")
		}

		return true, nil
	} else if strings.EqualFold(request.Type, "experience") {
		if key, ok := request.Key.(resume.ExperienceKey); ok {
			input, err := resume.GetExperienceDeleteInput(key)
			if err != nil {
				logger.Error("Failed to get the experience delete input", zap.Error(err))
				return false, err
			}

			if _, err := svc.DeleteItem(input); err != nil {
				logger.Error("Failed to delete experience from database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("key is not a ExperienceKey type", zap.Any("key", request.Key))
			return false, fmt.Errorf("key is not a ExperienceKey type")
		}

		return true, nil
	} else if strings.EqualFold(request.Type, "skill") {
		if key, ok := request.Key.(resume.SkillKey); ok {
			input, err := resume.GetSkillDeleteInput(key)
			if err != nil {
				logger.Error("Failed to get the skill delete input", zap.Error(err))
				return false, err
			}

			if _, err := svc.DeleteItem(input); err != nil {
				logger.Error("Failed to delete skill from database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("key is not a SkillKey type", zap.Any("key", request.Key))
			return false, fmt.Errorf("key is not a SkillKey type")
		}

		return true, nil
	}

	logger.Error("type was not an expected type", zap.String("type", request.Type))
	return false, fmt.Errorf("type '%s' was not an expected type", request.Type)
}

func main() {
	lambda.Start(HandleRemoveRequest)
}
