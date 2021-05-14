package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go.uber.org/zap"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	resume "github.com/bkimbrough88/resume-backend/pkg"
)

type GetItemsEvent struct {
	Type string `json:"type"`
}

func HandleGetRequest(request GetItemsEvent) (string, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	logger, loggerErr := zap.NewProduction()
	if loggerErr != nil {
		log.Printf("failed to initialize the logger. Error: %s", loggerErr.Error())
		return "", loggerErr
	}

	if strings.EqualFold(request.Type, "certification") {
		input := resume.GetCertificationsScanInput()
		result, err := svc.Scan(input)
		if err != nil {
			logger.Error("Failed to scan Certifications table", zap.Error(err))
			return "", err
		}

		certs := resume.ProcessCertificationScanResult(result.Items)
		certsJson, err := json.Marshal(certs)
		if err != nil {
			logger.Error("Failed to marshal certifications object into JSON", zap.Error(err))
			return "", err
		}

		return string(certsJson), nil
	} else if strings.EqualFold(request.Type, "degree") {
		input := resume.GetDegreesScanInput()
		result, err := svc.Scan(input)
		if err != nil {
			logger.Error("Failed to scan Degrees table", zap.Error(err))
			return "", err
		}

		degrees := resume.ProcessDegreeScanResult(result.Items)
		degreesJson, err := json.Marshal(degrees)
		if err != nil {
			logger.Error("Failed to marshal degrees object into JSON", zap.Error(err))
			return "", err
		}

		return string(degreesJson), nil
	} else if strings.EqualFold(request.Type, "experience") {
		input := resume.GetExperiencesScanInput()
		result, err := svc.Scan(input)
		if err != nil {
			logger.Error("Failed to scan Experience table", zap.Error(err))
			return "", err
		}

		experiences := resume.ProcessExperienceScanResult(result.Items)
		experiencesJson, err := json.Marshal(experiences)
		if err != nil {
			logger.Error("Failed to marshal experiences object into JSON", zap.Error(err))
			return "", err
		}

		return string(experiencesJson), nil
	} else if strings.EqualFold(request.Type, "generalInfo") {
		input := resume.GetGeneralInfoScanInput()
		result, err := svc.Scan(input)
		if err != nil {
			logger.Error("Failed to scan GeneralInfo table", zap.Error(err))
			return "", err
		}

		generalInfo := resume.ProcessGeneralInfoScanResult(result.Items)
		generalInfoJson, err := json.Marshal(generalInfo)
		if err != nil {
			logger.Error("Failed to marshal generalInfo object into JSON", zap.Error(err))
			return "", err
		}

		return string(generalInfoJson), nil
	} else if strings.EqualFold(request.Type, "skill") {
		input := resume.GetSkillsScanInput()
		result, err := svc.Scan(input)
		if err != nil {
			logger.Error("Failed to scan Skills table", zap.Error(err))
			return "", err
		}

		skills := resume.ProcessSkillScanResult(result.Items)
		skillsJson, err := json.Marshal(skills)
		if err != nil {
			logger.Error("Failed to marshal skills object into JSON", zap.Error(err))
			return "", err
		}

		return string(skillsJson), nil
	}

	logger.Error("type was not an expected type", zap.String("type", request.Type))
	return "", fmt.Errorf("type '%s' was not an expected type", request.Type)
}

func main() {
	lambda.Start(HandleGetRequest)
}
