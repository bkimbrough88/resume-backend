package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	resume "github.com/bkimbrough88/resume-backend/pkg"
)

type PutItemEvent struct {
	Type string `json:"type"`
	Item interface{} `json:"item"`
}

func HandlePutRequest(ctx context.Context, request PutItemEvent) (bool, error) {
	if strings.EqualFold(request.Type, "certification") {
		if cert, ok := request.Item.(resume.Certification); ok {
			if err := resume.PutCertification(cert); err != nil {
				return false, err
			}
		}
		return true, nil
	} else if strings.EqualFold(request.Type, "degree") {
		if degree, ok := request.Item.(resume.Degree); ok {
			if err := resume.PutDegree(degree); err != nil {
				return false, err
			}
		}
		return true, nil
	} else if strings.EqualFold(request.Type, "experience") {
		if experience, ok := request.Item.(resume.Experience); ok {
			if err := resume.PutExperience(experience); err != nil {
				return false, err
			}
		}
		return true, nil
	} else if strings.EqualFold(request.Type, "generalInfo") {
		if generalInfo, ok := request.Item.(resume.GeneralInfo); ok {
			if err := resume.PutGeneralInfo(generalInfo); err != nil {
				return false, err
			}
		}
		return true, nil
	} else if strings.EqualFold(request.Type, "skill") {
		if skill, ok := request.Item.(resume.Skill); ok {
			if err := resume.PutSkill(skill); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	return false, fmt.Errorf("type '%s' was not an expected type", request.Type)
}

func main() {
	lambda.Start(HandlePutRequest)
}
