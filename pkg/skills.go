package pkg

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	skills = "Skills"
)

type Skill struct {
	Name              string `json:"name"`
	YearsOfExperience int    `json:"years_of_experience,omitempty"`
}

type SkillKey struct {
	Name string `json:"name"`
}

func compareSkills(updateBuilder expression.UpdateBuilder, currentSkill Skill, updatedSkill Skill, idx int) {
	if currentSkill.Name != updatedSkill.Name {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, skills, idx, "Name")), expression.Value(updatedSkill.Name))
	}

	if currentSkill.YearsOfExperience != updatedSkill.YearsOfExperience {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, skills, idx, "YearsOfExperience")), expression.Value(updatedSkill.YearsOfExperience))
	}
}
