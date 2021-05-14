package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	resume "github.com/bkimbrough88/resume-backend/pkg"
)

type GetItemsEvent struct {
	Type string `json:"type"`
}

func HandleGetRequest(ctx context.Context, request GetItemsEvent) (string, error) {
	if strings.EqualFold(request.Type, "certification") {
		if certs, err := resume.GetCertifications(); err != nil {
			return "", err
		} else {
			if result, err := json.Marshal(certs); err != nil {
				return "", err
			} else {
				return string(result), nil
			}
		}
	} else if strings.EqualFold(request.Type, "degree") {
		if degrees, err := resume.GetDegrees(); err != nil {
			return "", err
		} else {
			if result, err := json.Marshal(degrees); err != nil {
				return "", err
			} else {
				return string(result), nil
			}
		}
	} else if strings.EqualFold(request.Type, "experience") {
		if experiences, err := resume.GetExperiences(); err != nil {
			return "", err
		} else {
			if result, err := json.Marshal(experiences); err != nil {
				return "", err
			} else {
				return string(result), nil
			}
		}
	} else if strings.EqualFold(request.Type, "generalInfo") {
		if generalInfo, err := resume.GetGeneralInfo(); err != nil {
			return "", err
		} else {
			if result, err := json.Marshal(generalInfo); err != nil {
				return "", err
			} else {
				return string(result), nil
			}
		}
	} else if strings.EqualFold(request.Type, "skill") {
		if skills, err := resume.GetSkills(); err != nil {
			return "", err
		} else {
			if result, err := json.Marshal(skills); err != nil {
				return "", err
			} else {
				return string(result), nil
			}
		}
	}

	return "", fmt.Errorf("type '%s' was not an expected type", request.Type)
}

func main() {
	lambda.Start(HandleGetRequest)
}
