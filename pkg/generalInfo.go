package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
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

func GetGeneralInfoScanInput() *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		TableName: aws.String(generalInfoTable),
	}
}

func ProcessGeneralInfoScanResult(resultItems []map[string]*dynamodb.AttributeValue) GeneralInfo {
	generalInfo := GeneralInfo{}
	if len(resultItems) >= 1 {
		err := dynamodbattribute.UnmarshalMap(resultItems[0], &generalInfo)
		if err != nil {
			// TODO: log something here
		}
	}

	return generalInfo
}

func GetGeneralInfoPutInput(generalInfo GeneralInfo) (*dynamodb.PutItemInput, error) {
	item, err := dynamodbattribute.MarshalMap(generalInfo)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(generalInfoTable),
	}
	return input, nil
}
