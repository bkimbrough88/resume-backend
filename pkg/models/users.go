package models

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/zap"
	"regexp"
)

const (
	ErrorInvalidEmail   = "invalid email"
	ErrorInvalidUserId  = "invalid user_id"
	ErrorNoResultsFound = "no results found"
	UsersTable          = "resume_user"
)

type User struct {
	UserId         string          `json:"user_id"`
	Email          string          `json:"email"`
	Certifications []Certification `json:"certifications,omitempty"`
	Degrees        []Degree        `json:"degrees,omitempty"`
	Experience     []Experience    `json:"experience,omitempty"`
	Github         string          `json:"github,omitempty"`
	GivenName      string          `json:"given_name,omitempty"`
	Location       string          `json:"location,omitempty"`
	Linkedin       string          `json:"linkedin,omitempty"`
	PhoneNumber    string          `json:"phone_number,omitempty"`
	Skills         []Skill         `json:"skills,omitempty"`
	Summary        string          `json:"summary,omitempty"`
	SurName        string          `json:"sur_name,omitempty"`
}

type UserKey struct {
	UserId string `json:"user_id"`
}

func PutUser(user *User, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) error {
	if !isEmail(user.Email) {
		logger.Error("Email is not a valid email", zap.String("email", user.Email))
		return errors.New(ErrorInvalidEmail)
	}

	if len(user.UserId) == 0 {
		logger.Error("UserId is empty")
		return errors.New(ErrorInvalidUserId)
	}

	input, err := getUserPutInput(user)
	if err != nil {
		logger.Error("Failed to construct input for create user", zap.Error(err))
		return err
	}

	_, err = svc.PutItem(input)
	if err != nil {
		logger.Error("Failed to insert new user into database", zap.Error(err))
		return err
	}

	logger.Info("Successfully inserted new user into database")
	return nil
}

func isEmail(email string) bool {
	emailRegex := regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)])")
	return emailRegex.MatchString(email)
}

func getUserPutInput(user *User) (*dynamodb.PutItemInput, error) {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(UsersTable),
	}
	return input, nil
}

func GetUserByKey(key *UserKey, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) (*User, error) {
	input, err := getUserGetItemInput(key)
	if err != nil {
		logger.Error("Failed to get input to query user table for ID", zap.Error(err))
		return nil, err
	}

	result, err := svc.GetItem(input)
	if err != nil {
		logger.Error("Failed to get user with key", zap.Error(err), zap.String("user_id", key.UserId))
		return nil, err
	}

	if len(result.Item) == 0 {
		logger.Error("No results found for with key", zap.String("user_id", key.UserId))
		return nil, errors.New(ErrorNoResultsFound)
	}

	logger.Info("Found user with key", zap.String("user_id", key.UserId))

	user := &User{}
	err = dynamodbattribute.UnmarshalMap(result.Item, user)
	if err != nil {
		logger.Error("Failed to unmarshall dynamo attributes to User object", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func getUserGetItemInput(key *UserKey) (*dynamodb.GetItemInput, error) {
	keyAttr, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		Key:       keyAttr,
		TableName: aws.String(UsersTable),
	}

	return input, nil
}

func DeleteUser(key *UserKey, svc dynamodbiface.DynamoDBAPI, logger *zap.Logger) error {
	input, err := getUserDeleteInput(key)
	if err != nil {
		logger.Error("Failed to get delete input", zap.Error(err))
		return err
	}

	_, err = svc.DeleteItem(input)
	if err != nil {
		logger.Error("Failed to delete user from database", zap.Error(err))
		return err
	}

	return nil
}

func getUserDeleteInput(keyObj *UserKey) (*dynamodb.DeleteItemInput, error) {
	key, err := dynamodbattribute.MarshalMap(keyObj)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: aws.String(UsersTable),
	}
	return input, nil
}
