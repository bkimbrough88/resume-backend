package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const skillsTable = "ResumeSkills"

type Skill struct {
	Name string `json:"name"`
	YearsOfExperience *int `json:"years_of_experience,omitempty"`
}

type SkillKey struct {
	Name string `json:"name"`
}

type Skills struct {
	Skills []Skill `json:"skills"`
}

func GetSkillsScanInput() *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		TableName: aws.String(skillsTable),
	}
}

func ProcessSkillScanResult(resultItems []map[string]*dynamodb.AttributeValue) Skills {
	var skills Skills
	for _, item := range resultItems {
		skill := Skill{}
		err := dynamodbattribute.UnmarshalMap(item, &skill)
		if err != nil {
			// TODO: log something here
		}

		skills.Skills = append(skills.Skills, skill)
	}

	return skills
}

func GetSkillPutInput(skill Skill) (*dynamodb.PutItemInput, error) {
	item, err := dynamodbattribute.MarshalMap(skill)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(skillsTable),
	}
	return input, nil
}

func GetSkillDeleteInput(keyObj SkillKey) (*dynamodb.DeleteItemInput, error) {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key: key,
		TableName: aws.String(skillsTable),
	}
	return input, nil
}
