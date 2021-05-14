package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	resume "github.com/bkimbrough88/resume-backend/pkg"
)

type RemoveItemEvent struct {
	Type string `json:"type"`
	Key interface{} `json:"item"`
}

func HandleRemoveRequest(ctx context.Context, request RemoveItemEvent) (bool, error) {
	if strings.EqualFold(request.Type, "certification") {
		if key, ok := request.Key.(resume.CertificationKey); ok {
			if err := resume.RemoveCertification(key); err != nil {
				return false, err
			}
		}
		return true, nil
	} else if strings.EqualFold(request.Type, "degree") {
		if key, ok := request.Key.(resume.DegreeKey); ok {
			if err := resume.RemoveDegree(key); err != nil {
				return false, err
			}
		}
		return true, nil
	} else if strings.EqualFold(request.Type, "experience") {
		if key, ok := request.Key.(resume.ExperienceKey); ok {
			if err := resume.RemoveExperience(key); err != nil {
				return false, err
			}
		}
		return true, nil
	} else if strings.EqualFold(request.Type, "skill") {
		if key, ok := request.Key.(resume.SkillKey); ok {
			if err := resume.RemoveSkill(key); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	return false, fmt.Errorf("type '%s' was not an expected type", request.Type)
}

func main() {
	lambda.Start(HandleRemoveRequest)
}
