package pkg

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"go.uber.org/zap"
)

const (
	usersTable = "ResumeUsers"

	listNameFormat            = "%s[%d]"
	listElementNameFormat     = "%s[%d].%s"
	listElementListNameFormat = "%s[%d].%s[%d]"
)

type User struct {
	Email          string          `json:"email"`
	Certifications []Certification `json:"certifications,omitempty"`
	Degrees        []Degree        `json:"degrees,omitempty"`
	Experience     []Experience    `json:"experience,omitempty"`
	Github         string          `json:"github,omitempty"`
	GivenName      string          `json:"given_name,omitempty"`
	LastUpdated    *time.Time      `json:"last_updated,omitempty"`
	Location       string          `json:"location,omitempty"`
	Linkedin       string          `json:"linkedin,omitempty"`
	PhoneNumber    string          `json:"phone_number,omitempty"`
	Skills         []Skill         `json:"skills,omitempty"`
	Summary        string          `json:"summary,omitempty"`
	SurName        string          `json:"sur_name,omitempty"`
}

type UserKey struct {
	Email string `json:"email"`
}

func CreateUser(user *User, svc *dynamodb.DynamoDB, logger *zap.Logger) error {
	if isEmail(user.Email) {
		logger.Error("Email is not a valid email", zap.String("email", user.Email))
		return fmt.Errorf("invalid email")
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

	logger.Info("Successfully inserted new usr into database")
	return nil
}

func isEmail(email string) bool {
	emailRegex := regexp.MustCompile("(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9]))\\.){3}(?:(2(5[0-5]|[0-4][0-9])|1[0-9][0-9]|[1-9]?[0-9])|[a-z0-9-]*[a-z0-9]:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)])")
	return emailRegex.MatchString(email)
}

func getUserPutInput(user *User) (*dynamodb.PutItemInput, error) {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(usersTable),
	}
	return input, nil
}

func GetUserByKey(key *UserKey, svc *dynamodb.DynamoDB, logger *zap.Logger) (*User, error) {
	filter := expression.Name("Email").Equal(expression.Value(key.Email))
	input, err := getUserScanInput(&filter)
	if err != nil {
		logger.Error("Failed to get input to query user table for ID", zap.Error(err))
		return nil, err
	}

	result, err := svc.Scan(input)
	if err != nil {
		logger.Error("Failed to scan user table for ID", zap.Error(err))
		return nil, err
	}

	if *result.Count == 0 {
		logger.Error("No results found for email", zap.String("email", key.Email))
		return nil, fmt.Errorf("no results found")
	} else if *result.Count > 1 {
		logger.Error("Too many results found for ID", zap.String("email", key.Email), zap.Int64("resultsReturned", *result.Count))
		return nil, fmt.Errorf("too many results found")
	} else {
		logger.Info("Found user for ID", zap.String("email", key.Email))

		var user User
		err = dynamodbattribute.UnmarshalMap(result.Items[0], user)
		if err != nil {
			logger.Error("Failed to unmarshall dynamo attributes to User object", zap.Error(err))
			return nil, err
		}

		return &user, nil
	}
}

func getUserScanInput(filter *expression.ConditionBuilder) (*dynamodb.ScanInput, error) {
	exp, err := expression.NewBuilder().WithFilter(*filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		FilterExpression:          exp.Filter(),
		TableName:                 aws.String(usersTable),
	}
	return input, nil
}

func UpdateUser(key *UserKey, updatedUser *User, svc *dynamodb.DynamoDB, logger *zap.Logger) error {
	if key.Email != updatedUser.Email {
		logger.Info("Updating the user's email, must re-create the user")

		if isEmail(updatedUser.Email) {
			logger.Error("Email is not a valid email", zap.String("email", updatedUser.Email))
			return fmt.Errorf("invalid email")
		}

		if err := DeleteUser(key, svc, logger); err != nil {
			logger.Error("Failed to delete the user from the database")
			return err
		}

		if err := CreateUser(updatedUser, svc, logger); err != nil {
			logger.Error("Failed to re-create the user with the updated user")
			return err
		}

		return nil
	}

	currentUser, err := GetUserByKey(key, svc, logger)
	if err != nil {
		logger.Error("Failed to get the current user from the database", zap.Error(err))
		return err
	}

	updateBuilder, err := getUserUpdateBuilder(currentUser, updatedUser)
	if err != nil {
		logger.Error("Failed to construct update", zap.Error(err))
		return err
	}

	expr, err := expression.NewBuilder().WithUpdate(*updateBuilder).Build()
	if err != nil {
		logger.Error("Failed to build update expression", zap.Error(err))
		return err
	}

	dynamoKey, err := dynamodbattribute.MarshalMap(key)
	if err != nil {
		logger.Error("Failed to convert key into dynamo attribute map", zap.Error(err))
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key:                       dynamoKey,
		TableName:                 aws.String(usersTable),
		UpdateExpression:          expr.Update(),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		logger.Error("Failed to update user in database", zap.Error(err))
		return err
	}

	return nil
}

func getUserUpdateBuilder(currentUser *User, updatedUser *User) (*expression.UpdateBuilder, error) {
	updateBuilder := expression.Set(expression.Name("LastUpdated"), expression.Value(time.Now().UTC()))

	if currentUser.Email != updatedUser.Email {
		updateBuilder.Set(expression.Name("Email"), expression.Value(updatedUser.Email))
	}

	if currentUser.Github != updatedUser.Github {
		updateBuilder.Set(expression.Name("Github"), expression.Value(updatedUser.Github))
	}

	if currentUser.GivenName != updatedUser.GivenName {
		updateBuilder.Set(expression.Name("GivenName"), expression.Value(updatedUser.GivenName))
	}

	if currentUser.Location != updatedUser.Location {
		updateBuilder.Set(expression.Name("Location"), expression.Value(updatedUser.Location))
	}

	if currentUser.Linkedin != updatedUser.Linkedin {
		updateBuilder.Set(expression.Name("Linkedin"), expression.Value(updatedUser.Linkedin))
	}

	if currentUser.PhoneNumber != updatedUser.PhoneNumber {
		updateBuilder.Set(expression.Name("PhoneNumber"), expression.Value(updatedUser.PhoneNumber))
	}

	if currentUser.Summary != updatedUser.Summary {
		updateBuilder.Set(expression.Name("Summary"), expression.Value(updatedUser.Summary))
	}

	if currentUser.SurName != updatedUser.SurName {
		updateBuilder.Set(expression.Name("SurName"), expression.Value(updatedUser.SurName))
	}

	currentCertsCount := len(currentUser.Certifications)
	updatedCertsCount := len(updatedUser.Certifications)
	for idx, currentCert := range currentUser.Certifications {
		if idx < updatedCertsCount {
			compareCertifications(updateBuilder, currentCert, updatedUser.Certifications[idx], idx)
		} else {
			updateBuilder.Remove(expression.Name(fmt.Sprintf(listNameFormat, certifications, idx)))
		}
	}
	for idx := currentCertsCount; idx < updatedCertsCount; idx++ {
		updateBuilder.Add(expression.Name(fmt.Sprintf(listNameFormat, certifications, idx)), expression.Value(updatedUser.Certifications[idx]))
	}

	currentDegreesCount := len(currentUser.Degrees)
	updatedDegreesCount := len(updatedUser.Degrees)
	for idx, currentDegree := range currentUser.Degrees {
		if idx < updatedDegreesCount {
			compareDegrees(updateBuilder, currentDegree, updatedUser.Degrees[idx], idx)
		} else {
			updateBuilder.Remove(expression.Name(fmt.Sprintf(listNameFormat, degrees, idx)))
		}
	}
	for idx := currentDegreesCount; idx < updatedDegreesCount; idx++ {
		updateBuilder.Add(expression.Name(fmt.Sprintf(listNameFormat, degrees, idx)), expression.Value(updatedUser.Degrees[idx]))
	}

	currentExperienceCount := len(currentUser.Experience)
	updatedExperienceCount := len(updatedUser.Experience)
	for idx, currentExperience := range currentUser.Experience {
		if idx < updatedExperienceCount {
			compareExperience(updateBuilder, currentExperience, updatedUser.Experience[idx], idx)
		} else {
			updateBuilder.Remove(expression.Name(fmt.Sprintf(listNameFormat, experience, idx)))
		}
	}
	for idx := currentExperienceCount; idx < updatedExperienceCount; idx++ {
		updateBuilder.Add(expression.Name(fmt.Sprintf(listNameFormat, experience, idx)), expression.Value(updatedUser.Experience[idx]))
	}

	currentSkillsCount := len(currentUser.Skills)
	updatedSkillsCount := len(updatedUser.Skills)
	for idx, currentSkill := range currentUser.Skills {
		if idx < updatedSkillsCount {
			compareSkills(updateBuilder, currentSkill, updatedUser.Skills[idx], idx)
		} else {
			updateBuilder.Remove(expression.Name(fmt.Sprintf(listNameFormat, certifications, idx)))
		}
	}
	for idx := currentSkillsCount; idx < updatedSkillsCount; idx++ {
		updateBuilder.Add(expression.Name(fmt.Sprintf(listNameFormat, degrees, idx)), expression.Value(updatedUser.Skills[idx]))
	}

	return &updateBuilder, nil
}

func DeleteUser(key *UserKey, svc *dynamodb.DynamoDB, logger *zap.Logger) error {
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
		TableName: aws.String(usersTable),
	}
	return input, nil
}
