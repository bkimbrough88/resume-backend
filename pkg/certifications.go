package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

func GetCertifications() (*Certifications, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(certificationsTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	result, err := svc.Scan(input)
	if err != nil {
		return nil, err
	}

	var certs Certifications
	for _, item := range result.Items {
		cert := Certification{}
		err = dynamodbattribute.UnmarshalMap(item, &cert)
		if err != nil {
			// TODO: log something here
		}

		certs.Certifications = append(certs.Certifications, cert)
	}

	return &certs, nil
}

func PutCertification(cert Certification) error {
	item, err := dynamodbattribute.MarshalMap(cert)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: item,
		TableName: aws.String(certificationsTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func RemoveCertification(keyObj CertificationKey) error {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: key,
		TableName: aws.String(certificationsTable),
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	_, err = svc.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}
