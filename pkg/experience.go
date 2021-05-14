package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const experienceTable = "ResumeExperience"

type Experience struct {
	Company string `json:"company"`
	EndMonth string `json:"end_month,omitempty"`
	EndYear int `json:"end_year,omitempty"`
	JobTitle string `json:"job_title"`
	Responsibilities []string `json:"responsibilities"`
	StartMonth string `json:"start_month"`
	StartYear int `json:"start_year"`
}

type ExperienceKey struct {
	Company string `json:"company"`
	JobTitle string `json:"job_title"`
	StartYear int `json:"start_year"`
}

type Experiences struct {
	Experiences []Experience `json:"experiences"`
}

func GetExperiencesScanInput() *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		TableName: aws.String(experienceTable),
	}
}

func ProcessExperienceScanResult(resultItems []map[string]*dynamodb.AttributeValue) Experiences {
	var experiences Experiences
	for _, item := range resultItems {
		experience := Experience{}
		err := dynamodbattribute.UnmarshalMap(item, &experience)
		if err != nil {
			// TODO: log something here
		}

		experiences.Experiences = append(experiences.Experiences, experience)
	}

	return experiences
}

func GetExperiencePutInput(experience Experience) (*dynamodb.PutItemInput, error) {
	item, err := dynamodbattribute.MarshalMap(experience)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(experienceTable),
	}
	return input, nil
}

func GetExperienceDeleteInput(keyObj ExperienceKey) (*dynamodb.DeleteItemInput, error) {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key: key,
		TableName: aws.String(experienceTable),
	}
	return input, nil
}