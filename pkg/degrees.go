package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const degreesTable = "ResumeDegrees"

type Degree struct {
	Degree string `json:"degree"`
	EndYear int `json:"end_year,omitempty"`
	Major string `json:"major"`
	School string `json:"school"`
	StartYear int `json:"start_year"`
}

type DegreeKey struct {
	Degree string `json:"degree"`
	Major string `json:"major"`
	School string `json:"school"`
}

type Degrees struct {
	Degrees []Degree `json:"degrees"`
}

func GetDegreesScanInput() *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		TableName: aws.String(degreesTable),
	}
}

func ProcessDegreeScanResult(resultItems []map[string]*dynamodb.AttributeValue) Degrees {
	var degrees Degrees
	for _, item := range resultItems {
		degree := Degree{}
		err := dynamodbattribute.UnmarshalMap(item, &degree)
		if err != nil {
			// TODO: log something here
		}

		degrees.Degrees = append(degrees.Degrees, degree)
	}

	return degrees
}

func GetDegreePutInput(degree Degree) (*dynamodb.PutItemInput, error) {
	item, err := dynamodbattribute.MarshalMap(degree)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(degreesTable),
	}
	return input, nil
}

func GetDegreeDeleteInput(keyObj DegreeKey) (*dynamodb.DeleteItemInput, error) {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key: key,
		TableName: aws.String(degreesTable),
	}
	return input, nil
}

func PutDegree(degree Degree) error {
	item, err := dynamodbattribute.MarshalMap(degree)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(degreesTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func RemoveDegree(keyObj DegreeKey) error {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: key,
		TableName: aws.String(degreesTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}
