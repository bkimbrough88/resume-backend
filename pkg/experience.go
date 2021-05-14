package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

func GetExperiences() ([]Experience, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(experienceTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	var experiences []Experience
	for _, item := range result.Items {
		experience := Experience{}
		err = dynamodbattribute.UnmarshalMap(item, &experience)
		if err != nil {
			// TODO: log something here
		}

		experiences = append(experiences, experience)
	}

	return experiences, nil
}

func PutExperience(experience Experience) error {
	item, err := dynamodbattribute.MarshalMap(experience)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(experienceTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func RemoveExperience(keyObj ExperienceKey) error {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: key,
		TableName: aws.String(experienceTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}

