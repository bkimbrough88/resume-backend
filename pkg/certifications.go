package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const certificationsTable = "ResumeCertifications"

type Certification struct {
	CertificationName string `json:"certification_name"`
	CertificationBadgeLink string `json:"certification_badge_link,omitempty"`
	Description string `json:"description"`
}

type CertificationKey struct {
	CertificationName string `json:"certification_name"`
}

type Certifications struct {
	Certifications []Certification `json:"certifications"`
}

func GetCertificationsScanInput() *dynamodb.ScanInput {
	return &dynamodb.ScanInput{
		TableName: aws.String(certificationsTable),
	}
}

func ProcessCertificationScanResult(resultItems []map[string]*dynamodb.AttributeValue) Certifications {
	var certs Certifications
	for _, item := range resultItems {
		cert := Certification{}
		err := dynamodbattribute.UnmarshalMap(item, &cert)
		if err != nil {
			// TODO: log something here
		}

		certs.Certifications = append(certs.Certifications, cert)
	}

	return certs
}

func GetCertificationPutInput(cert Certification) (*dynamodb.PutItemInput, error) {
	item, err := dynamodbattribute.MarshalMap(cert)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(certificationsTable),
	}
	return input, nil
}

func GetCertificationDeleteInput(keyObj CertificationKey) (*dynamodb.DeleteItemInput, error) {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key: key,
		TableName: aws.String(certificationsTable),
	}
	return input, nil
}
