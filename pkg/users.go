package pkg

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	usersTable = "ResumeUsers"

	listNameFormat            = "%s[%d]"
	listElementNameFormat     = "%s[%d].%s"
	listElementListNameFormat = "%s[%d].%s[%d]"
)

type User struct {
	Id             *uuid.UUID      `json:"id"`
	Certifications []Certification `json:"certifications,omitempty"`
	Degrees        []Degree        `json:"degrees,omitempty"`
	Email          string          `json:"email,omitempty"`
	Experience     []Experience    `json:"experience,omitempty"`
	Github         string          `json:"github,omitempty"`
	GivenName      string          `json:"given_name,omitempty"`
	LastUpdated    *time.Time	    `json:"last_updated,omitempty"`
	Location       string          `json:"location,omitempty"`
	Linkedin       string          `json:"linkedin,omitempty"`
	PhoneNumber    string          `json:"phone_number,omitempty"`
	Skills         []Skill         `json:"skills,omitempty"`
	Summary        string          `json:"summary,omitempty"`
	SurName        string          `json:"sur_name,omitempty"`
}

type UserKey struct {
	Id *uuid.UUID `json:"id"`
}

func CreateUser(user *User, svc *dynamodb.DynamoDB, logger *zap.Logger) error {
	if user.Id == nil {
		newUuid := uuid.New()
		user.Id = &newUuid
	}

	input, err := getUserPutInput(*user)
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

func getUserPutInput(user User) (*dynamodb.PutItemInput, error) {
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

func GetUserById(idStr string, svc *dynamodb.DynamoDB, logger *zap.Logger) (*User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Failed to parse id", zap.Error(err))
		return nil, err
	}

	filter := expression.Name("Id").Equal(expression.Value(id))
	input, err := getUserScanInput(filter)
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
		logger.Error("No results found for ID", zap.String("ID", idStr))
		return nil, fmt.Errorf("no results found")
	} else if *result.Count > 1 {
		logger.Error("Too many results found for ID", zap.String("ID", idStr), zap.Int64("resultsReturned", *result.Count))
		return nil, fmt.Errorf("too many results found")
	} else {
		logger.Info("Found user for ID", zap.String("ID", idStr))

		var user User
		err = dynamodbattribute.UnmarshalMap(result.Items[0], user)
		if err != nil {
			logger.Error("Failed to unmarshall dynamo attributes to User object", zap.Error(err))
			return nil, err
		}

		return &user, nil
	}
}

func getUserScanInput(filter expression.ConditionBuilder) (*dynamodb.ScanInput, error) {
	exp, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.ScanInput{
		TableName:        aws.String(usersTable),
		FilterExpression: exp.Filter(),
	}
	return input, nil
}

func UpdateUser(idStr string, updatedUser *User, svc *dynamodb.DynamoDB, logger *zap.Logger) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Failed to parse id", zap.Error(err))
		return err
	}

	currentUser, err := GetUserById(idStr, svc, logger)
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

	userKey := UserKey{Id: &id}
	key, err := dynamodbattribute.MarshalMap(userKey)
	if err != nil {
		logger.Error("Failed to convert key into dynamo attribute map", zap.Error(err))
		return err
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Key:                       key,
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
		updateBuilder = updateBuilder.Set(expression.Name("Email"), expression.Value(updatedUser.Email))
	}

	if currentUser.Github != updatedUser.Github {
		updateBuilder = updateBuilder.Set(expression.Name("Github"), expression.Value(updatedUser.Github))
	}

	if currentUser.GivenName != updatedUser.GivenName {
		updateBuilder = updateBuilder.Set(expression.Name("GivenName"), expression.Value(updatedUser.GivenName))
	}

	if currentUser.Location != updatedUser.Location {
		updateBuilder = updateBuilder.Set(expression.Name("Location"), expression.Value(updatedUser.Location))
	}

	if currentUser.Linkedin != updatedUser.Linkedin {
		updateBuilder = updateBuilder.Set(expression.Name("Linkedin"), expression.Value(updatedUser.Linkedin))
	}

	if currentUser.PhoneNumber != updatedUser.PhoneNumber {
		updateBuilder = updateBuilder.Set(expression.Name("PhoneNumber"), expression.Value(updatedUser.PhoneNumber))
	}

	if currentUser.Summary != updatedUser.Summary {
		updateBuilder = updateBuilder.Set(expression.Name("Summary"), expression.Value(updatedUser.Summary))
	}

	if currentUser.SurName != updatedUser.SurName {
		updateBuilder = updateBuilder.Set(expression.Name("SurName"), expression.Value(updatedUser.SurName))
	}

	currentCertsCount := len(currentUser.Certifications)
	updatedCertsCount := len(updatedUser.Certifications)
	for idx, currentCert := range currentUser.Certifications {
		if idx < updatedCertsCount-1 {
			updateBuilder = compareCertifications(updateBuilder, currentCert, updatedUser.Certifications[idx], idx)
		} else {
			updateBuilder = updateBuilder.Remove(expression.Name(fmt.Sprintf(listNameFormat, certifications, idx)))
		}
	}
	for idx := currentCertsCount; idx < updatedCertsCount; idx++ {
		newCert, err := dynamodbattribute.MarshalMap(updatedUser.Certifications[idx])
		if err != nil {
			return nil, err
		}

		updateBuilder = updateBuilder.Add(expression.Name(fmt.Sprintf(listNameFormat, certifications, idx)), expression.Value(newCert))
	}

	currentDegreesCount := len(currentUser.Degrees)
	updatedDegreesCount := len(updatedUser.Degrees)
	for idx, currentDegree := range currentUser.Degrees {
		if idx < updatedDegreesCount-1 {
			updateBuilder = compareDegrees(updateBuilder, currentDegree, updatedUser.Degrees[idx], idx)
		} else {
			updateBuilder = updateBuilder.Remove(expression.Name(fmt.Sprintf(listNameFormat, degrees, idx)))
		}
	}
	for idx := currentDegreesCount; idx < updatedDegreesCount; idx++ {
		newDegree, err := dynamodbattribute.MarshalMap(updatedUser.Degrees[idx])
		if err != nil {
			return nil, err
		}

		updateBuilder = updateBuilder.Add(expression.Name(fmt.Sprintf(listNameFormat, degrees, idx)), expression.Value(newDegree))
	}

	currentExperienceCount := len(currentUser.Experience)
	updatedExperienceCount := len(updatedUser.Experience)
	for idx, currentExperience := range currentUser.Experience {
		if idx < updatedExperienceCount-1 {
			updateBuilder = compareExperience(updateBuilder, currentExperience, updatedUser.Experience[idx], idx)
		} else {
			updateBuilder = updateBuilder.Remove(expression.Name(fmt.Sprintf(listNameFormat, experience, idx)))
		}
	}
	for idx := currentExperienceCount; idx < updatedExperienceCount; idx++ {
		newExperience, err := dynamodbattribute.MarshalMap(updatedUser.Experience[idx])
		if err != nil {
			return nil, err
		}

		updateBuilder = updateBuilder.Add(expression.Name(fmt.Sprintf(listNameFormat, experience, idx)), expression.Value(newExperience))
	}

	currentSkillsCount := len(currentUser.Skills)
	updatedSkillsCount := len(updatedUser.Skills)
	for idx, currentSkill := range currentUser.Skills {
		if idx < updatedSkillsCount-1 {
			updateBuilder = compareSkills(updateBuilder, currentSkill, updatedUser.Skills[idx], idx)
		} else {
			updateBuilder = updateBuilder.Remove(expression.Name(fmt.Sprintf(listNameFormat, certifications, idx)))
		}
	}
	for idx := currentSkillsCount; idx < updatedSkillsCount; idx++ {
		newSkills, err := dynamodbattribute.MarshalMap(updatedUser.Skills[idx])
		if err != nil {
			return nil, err
		}

		updateBuilder = updateBuilder.Add(expression.Name(fmt.Sprintf(listNameFormat, degrees, idx)), expression.Value(newSkills))
	}

	return &updateBuilder, nil
}

func DeleteUser(idStr string, svc *dynamodb.DynamoDB, logger *zap.Logger) error {
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Error("Failed to parse id", zap.Error(err))
		return err
	}

	key := UserKey{Id: &id}
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

func getUserDeleteInput(keyObj UserKey) (*dynamodb.DeleteItemInput, error) {
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
