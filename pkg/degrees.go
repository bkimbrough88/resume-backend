package pkg

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	degrees = "Degrees"
)

type Degree struct {
	Degree    string `json:"degree"`
	Major     string `json:"major"`
	School    string `json:"school"`
	StartYear int    `json:"start_year"`
	EndYear   int    `json:"end_year,omitempty"`
}

type DegreeKey struct {
	Degree string `json:"degree"`
	Major  string `json:"major"`
	School string `json:"school"`
}

func compareDegrees(updateBuilder expression.UpdateBuilder, currentDegree Degree, updatedDegree Degree, idx int) expression.UpdateBuilder {
	myUpdateBuilder := updateBuilder
	if currentDegree.Degree != updatedDegree.Degree {
		myUpdateBuilder = myUpdateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, degrees, idx, "Degree")), expression.Value(updatedDegree.Degree))
	}

	if currentDegree.Major != updatedDegree.Major {
		myUpdateBuilder = myUpdateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, degrees, idx, "Major")), expression.Value(updatedDegree.Major))
	}

	if currentDegree.School != updatedDegree.School {
		myUpdateBuilder = myUpdateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, degrees, idx, "School")), expression.Value(updatedDegree.School))
	}

	if currentDegree.StartYear != updatedDegree.StartYear {
		myUpdateBuilder = myUpdateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, degrees, idx, "StartYear")), expression.Value(updatedDegree.StartYear))
	}

	if currentDegree.EndYear != updatedDegree.EndYear {
		myUpdateBuilder = myUpdateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, degrees, idx, "EndYear")), expression.Value(updatedDegree.EndYear))
	}

	return myUpdateBuilder
}
