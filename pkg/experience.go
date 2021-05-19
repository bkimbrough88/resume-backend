package pkg

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	experience = "Experience"
)

type Experience struct {
	Company          string   `json:"company"`
	JobTitle         string   `json:"job_title"`
	StartMonth       string   `json:"start_month"`
	StartYear        int      `json:"start_year"`
	EndMonth         string   `json:"end_month,omitempty"`
	EndYear          int      `json:"end_year,omitempty"`
	Responsibilities []string `json:"responsibilities,omitempty"`
}

type ExperienceKey struct {
	Company  string `json:"company"`
	JobTitle string `json:"job_title"`
}

func compareExperience(updateBuilder expression.UpdateBuilder, currentExperience Experience, updatedExperience Experience, idx int) {
	if currentExperience.Company != updatedExperience.Company {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, experience, idx, "Company")), expression.Value(updatedExperience.Company))
	}

	if currentExperience.JobTitle != updatedExperience.JobTitle {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, experience, idx, "JobTitle")), expression.Value(updatedExperience.JobTitle))
	}

	if currentExperience.StartMonth != updatedExperience.StartMonth {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, experience, idx, "StartMonth")), expression.Value(updatedExperience.StartMonth))
	}

	if currentExperience.StartYear != updatedExperience.StartYear {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, experience, idx, "StartYear")), expression.Value(updatedExperience.StartYear))
	}

	if currentExperience.EndMonth != updatedExperience.EndMonth {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, experience, idx, "EndMonth")), expression.Value(updatedExperience.EndMonth))
	}

	if currentExperience.EndYear != updatedExperience.EndYear {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, experience, idx, "EndYear")), expression.Value(updatedExperience.EndYear))
	}

	currentResponsibilitiesCount := len(currentExperience.Responsibilities)
	updatedResponsibilitiesCount := len(updatedExperience.Responsibilities)
	for responsibilityIdx, currentResponsibility := range currentExperience.Responsibilities {
		if responsibilityIdx < updatedResponsibilitiesCount {
			if currentResponsibility != updatedExperience.Responsibilities[responsibilityIdx] {
				updateBuilder.Set(expression.Name(fmt.Sprintf(listElementListNameFormat, experience, idx, "Responsibilities", responsibilityIdx)), expression.Value(updatedExperience.Responsibilities[responsibilityIdx]))
			}
		} else {
			updateBuilder.Remove(expression.Name(fmt.Sprintf(listElementListNameFormat, experience, idx, "Responsibilities", responsibilityIdx)))
		}
	}
	for responsibilityIdx := currentResponsibilitiesCount; responsibilityIdx < updatedResponsibilitiesCount; responsibilityIdx++ {
		updateBuilder.Add(expression.Name(fmt.Sprintf(listElementListNameFormat, experience, idx, "Responsibilities", responsibilityIdx)), expression.Value(updatedExperience.Responsibilities[responsibilityIdx]))
	}
}
