package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

func GetSkills() ([]Skill, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(skillsTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	var skills []Skill
	for _, item := range result.Items {
		skill := Skill{}
		err = dynamodbattribute.UnmarshalMap(item, &skill)
		if err != nil {
			// TODO: log something here
		}

		skills = append(skills, skill)
	}

	return skills, nil
}

func PutSkill(skill Skill) error {
	item, err := dynamodbattribute.MarshalMap(skill)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(skillsTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func RemoveSkill(keyObj SkillKey) error {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: key,
		TableName: aws.String(skillsTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}
