package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const generalInfoTable = "ResumeGeneralInfo"

type GeneralInfo struct {
	Email string `json:"email"`
	Github string `json:"github,omitempty"`
	GivenName string `json:"given_name"`
	Location string `json:"location"`
	Linkedin string `json:"linkedin,omitempty"`
	PhoneNumber string `json:"phone_number"`
	Summary string `json:"summary"`
	SurName string `json:"sur_name"`
}

func GetGeneralInfo() (*GeneralInfo, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(generalInfoTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	if *result.Count >= 1 {
		generalInfo := GeneralInfo{}
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &generalInfo)
		if err != nil {
			// TODO: log something here
		}
		return &generalInfo, nil
	} else {
		return nil, fmt.Errorf("no items exist in the %s table", generalInfoTable)
	}
}

func PutGeneralInfo(generalInfo GeneralInfo) error {
	item, err := dynamodbattribute.MarshalMap(generalInfo)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(generalInfoTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
