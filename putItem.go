package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.uber.org/zap"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	resume "github.com/bkimbrough88/resume-backend/pkg"
)

type PutItemEvent struct {
	Type string `json:"type"`
	Item interface{} `json:"item"`
}

func HandlePutRequest(request PutItemEvent) (bool, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	logger, loggerErr := zap.NewProduction()
	if loggerErr != nil {
		log.Printf("failed to initialize the logger. Error: %s", loggerErr.Error())
		return false, loggerErr
	}

	if strings.EqualFold(request.Type, "certification") {
		if cert, ok := request.Item.(resume.Certification); ok {
			input, err := resume.GetCertificationPutInput(cert)
			if err != nil {
				logger.Error("Failed to get Certification put input", zap.Error(err))
				return false, err
			}

			if _, err = svc.PutItem(input); err != nil {
				logger.Error("Failed to get insert Certification into database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("item is not a Certification type", zap.Any("item", request.Item))
			return false, fmt.Errorf("item is not a Certification type")
		}

		return true, nil
	} else if strings.EqualFold(request.Type, "degree") {
		if degree, ok := request.Item.(resume.Degree); ok {
			input, err := resume.GetDegreePutInput(degree)
			if err != nil {
				logger.Error("Failed to get Degree put input", zap.Error(err))
				return false, err
			}

			if _, err = svc.PutItem(input); err != nil {
				logger.Error("Failed to get insert Degree into database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("item is not a Degree type", zap.Any("item", request.Item))
			return false, fmt.Errorf("item is not a Degree type")
		}

		return true, nil
	} else if strings.EqualFold(request.Type, "experience") {
		if experience, ok := request.Item.(resume.Experience); ok {
			input, err := resume.GetExperiencePutInput(experience)
			if err != nil {
				logger.Error("Failed to get Experience put input", zap.Error(err))
				return false, err
			}

			if _, err = svc.PutItem(input); err != nil {
				logger.Error("Failed to get insert Experience into database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("item is not an Experience type", zap.Any("item", request.Item))
			return false, fmt.Errorf("item is not an Experience type")
		}

		return true, nil
	} else if strings.EqualFold(request.Type, "generalInfo") {
		if generalInfo, ok := request.Item.(resume.GeneralInfo); ok {
			input, err := resume.GetGeneralInfoPutInput(generalInfo)
			if err != nil {
				logger.Error("Failed to get GeneralInfo put input", zap.Error(err))
				return false, err
			}

			if _, err = svc.PutItem(input); err != nil {
				logger.Error("Failed to get insert GeneralInfo into database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("item is not a GeneralInfo type", zap.Any("item", request.Item))
			return false, fmt.Errorf("item is not a GeneralInfo type")
		}

		return true, nil
	} else if strings.EqualFold(request.Type, "skill") {
		if skill, ok := request.Item.(resume.Skill); ok {
			input, err := resume.GetSkillPutInput(skill)
			if err != nil {
				logger.Error("Failed to get Skill put input", zap.Error(err))
				return false, err
			}

			if _, err = svc.PutItem(input); err != nil {
				logger.Error("Failed to get insert Certification into database", zap.Error(err))
				return false, err
			}
		} else {
			logger.Error("item is not a Skill type", zap.Any("item", request.Item))
			return false, fmt.Errorf("item is not a Skill type")
		}

		return true, nil
	}

	logger.Error("type was not an expected type", zap.String("type", request.Type))
	return false, fmt.Errorf("type '%s' was not an expected type", request.Type)
}

func main() {
	lambda.Start(HandlePutRequest)
}
