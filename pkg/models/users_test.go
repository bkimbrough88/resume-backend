package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	mocks "github.com/bkimbrough88/resume-backend/pkg"
	"go.uber.org/zap"
)

var (
	user   *User
	logger *zap.Logger
)

func setup(t *testing.T) {
	logger, _ = zap.NewDevelopment()
	user = &User{
		Certifications: []Certification{
			{
				Name:         "Some Cert",
				BadgeLink:    "https://example.com",
				DateAchieved: "10-28-2019",
				DateExpires:  "10-28-2022",
			},
		},
		Degrees: []Degree{
			{
				Degree:    "BS",
				Major:     "CS",
				School:    "University",
				StartYear: 2017,
				EndYear:   2021,
			},
		},
		Email: "user@domain.com",
		Experience: []Experience{
			{
				Company:    "Co",
				JobTitle:   "SRE",
				StartMonth: "May",
				StartYear:  2020,
				EndMonth:   "June",
				EndYear:    2020,
				Responsibilities: []string{
					"foo",
					"bar",
				},
			},
		},
		Github:      "https://github.com/user",
		GivenName:   "John",
		Location:    "Place, State",
		Linkedin:    "https://www.linkedin.com/in/user",
		PhoneNumber: "999-999-9999",
		Skills: []Skill{
			{
				Name:              "Go",
				YearsOfExperience: 2,
			},
		},
		Summary: "My awesome summary",
		SurName: "Doe",
		UserId:  "user1",
	}

	mocks.DeleteItemMock = func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
		return &dynamodb.DeleteItemOutput{}, nil
	}

	attr, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		t.Fatalf("Failed to marshal user int Dynamo attribute map: %s", err.Error())
	}
	mocks.GetItemMock = func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{Item: attr}, nil
	}

	mocks.PutItemMock = func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{}, nil
	}

	mocks.UpdateItemMock = func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
		return &dynamodb.UpdateItemOutput{}, nil
	}
}

func TestCreateUser(t *testing.T) {
	setup(t)

	svc := mocks.DynamoServiceMock{}
	if err := PutUser(user, svc, logger); err != nil {
		t.Errorf("Failed to create user when it should have been successful: %s", err.Error())
	}

	userBadEmail := &User{
		UserId: "user",
		Email:  "not an email",
	}
	if err := PutUser(userBadEmail, svc, logger); err == nil {
		t.Errorf("Expected to get an error and no err was returned")
	} else if ErrorInvalidEmail != err.Error() {
		t.Errorf("Expected error to be '%s', but was '%s'", ErrorInvalidEmail, err.Error())
	}

	userBadUserId := &User{
		UserId: "",
		Email:  "user@domain.com",
	}
	if err := PutUser(userBadUserId, svc, logger); err == nil {
		t.Errorf("Expected to get an error and no err was returned")
	} else if ErrorInvalidUserId != err.Error() {
		t.Errorf("Expected error to be '%s', but was '%s'", ErrorInvalidUserId, err.Error())
	}

	expectedError := "some error"
	mocks.PutItemMock = func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return nil, fmt.Errorf(expectedError)
	}
	if err := PutUser(user, svc, logger); err == nil {
		t.Errorf("Created user when it should have failed")
	} else if err.Error() != expectedError {
		t.Errorf("Expected error to be '%s', but was '%s'", expectedError, err.Error())
	}
}

func TestGetUserPutInput(t *testing.T) {
	setup(t)

	if input, err := getUserPutInput(user); err != nil {
		t.Errorf("Failed to get input with error '%s'", err.Error())
	} else {
		if input.TableName == nil {
			t.Error("Table name should not be nil")
		} else if *input.TableName != UsersTable {
			t.Errorf("Expected table name to be '%s', but was '%s'", UsersTable, *input.TableName)
		}

		if input.Item == nil {
			t.Error("User should not have generated an empty map")
		}
	}
}

func TestIsEmail(t *testing.T) {
	nonEmail1 := ""
	nonEmail2 := "a@b"
	nonEmail3 := "not an email"
	email1 := "user@domain.com"
	email2 := "a@b.co"

	if isEmail(nonEmail1) {
		t.Errorf("Determined that '%s' is a valid email and it is not", nonEmail1)
	}

	if isEmail(nonEmail2) {
		t.Errorf("Determined that '%s' is a valid email and it is not", nonEmail2)
	}

	if isEmail(nonEmail3) {
		t.Errorf("Determined that '%s' is a valid email and it is not", nonEmail3)
	}

	if !isEmail(email1) {
		t.Errorf("Determined that '%s' is not a valid email, but it is", email1)
	}

	if !isEmail(email2) {
		t.Errorf("Determined that '%s' is not a valid email, but it is", email1)
	}
}

func TestGetUserByKey(t *testing.T) {
	setup(t)

	key := UserKey{UserId: "username"}
	svc := mocks.DynamoServiceMock{}
	if res, err := GetUserByKey(&key, svc, logger); err != nil {
		t.Errorf("Expected to get a user and got the error '%s' instead", err.Error())
	} else if res == nil {
		t.Errorf("Expected to get a user, but got nil")
	} else {
		if user.Email != res.Email {
			t.Errorf("Expected email to be '%s', but was '%s'", user.Email, res.Email)
		}
		// TODO: Check the rest of the fields
	}

	mocks.GetItemMock = func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{},
		}, nil
	}
	if _, err := GetUserByKey(&key, svc, logger); err == nil {
		t.Errorf("Found user when none should have been found")
	} else if err.Error() != "no results found" {
		t.Errorf("Expected error to be 'no results found', but was '%s'", err.Error())
	}

	expectedError := "some error"
	mocks.GetItemMock = func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
		return nil, fmt.Errorf(expectedError)
	}
	if _, err := GetUserByKey(&key, svc, logger); err == nil {
		t.Errorf("Found user when it should have failed")
	} else if err.Error() != expectedError {
		t.Errorf("Expected error to be '%s', but was '%s'", expectedError, err.Error())
	}
}

func TestGetUserGetItemInput(t *testing.T) {
	key := &UserKey{UserId: "username"}
	if input, err := getUserGetItemInput(key); err != nil {
		t.Errorf("Failed to get input with error '%s'", err.Error())
	} else {
		if input.TableName == nil {
			t.Error("Table name should not be nil")
		} else if *input.TableName != UsersTable {
			t.Errorf("Expected table name to be '%s', but was '%s'", UsersTable, *input.TableName)
		}

		if input.Key == nil {
			t.Error("User key should not have generated an empty map")
		} else if input.Key["user_id"].S == nil {
			t.Error("Expected user_id to be a string type")
		} else if *input.Key["user_id"].S != key.UserId {
			t.Errorf("Expected user_id to be '%s', but was '%s'", key.UserId, *input.Key["user_id"].S)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	setup(t)

	key := &UserKey{UserId: "username"}
	svc := mocks.DynamoServiceMock{}
	if err := DeleteUser(key, svc, logger); err != nil {
		t.Errorf("Failed to delete user when it should have been successful: %s", err.Error())
	}

	expectedError := "some error"
	mocks.DeleteItemMock = func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
		return nil, fmt.Errorf(expectedError)
	}
	if err := DeleteUser(key, svc, logger); err == nil {
		t.Errorf("Deleted user when it should have failed")
	} else if err.Error() != expectedError {
		t.Errorf("Expected error to be '%s', but was '%s'", expectedError, err.Error())
	}
}

func TestGetUserDeleteInput(t *testing.T) {
	key := &UserKey{UserId: "username"}
	if input, err := getUserDeleteInput(key); err != nil {
		t.Errorf("Failed to get input with error '%s'", err.Error())
	} else {
		if input.TableName == nil {
			t.Error("Table name should not be nil")
		} else if *input.TableName != UsersTable {
			t.Errorf("Expected table name to be '%s', but was '%s'", UsersTable, *input.TableName)
		}

		if input.Key == nil {
			t.Error("User key should not have generated an empty map")
		} else if input.Key["user_id"].S == nil {
			t.Error("Expected user_id to be a string type")
		} else if *input.Key["user_id"].S != key.UserId {
			t.Errorf("Expected user_id to be '%s', but was '%s'", key.UserId, *input.Key["user_id"].S)
		}
	}
}

/** TEST HELPERS  */

func getValueKey(prefixKey *string, nameKey string, update string) string {
	var keyIdx int
	if prefixKey != nil {
		keyIdx = strings.Index(update, fmt.Sprintf("%s.%s", *prefixKey, nameKey))
	} else {
		keyIdx = strings.Index(update, nameKey)
	}

	if keyIdx == -1 {
		return ""
	}

	startIdx := keyIdx + strings.Index(update[keyIdx:], ":")
	commaIdx := strings.Index(update[startIdx:], ",")
	newLineIdx := strings.Index(update[startIdx:], "\n")

	var endIdx int
	if commaIdx < newLineIdx && commaIdx != -1 {
		endIdx = startIdx + commaIdx
	} else if newLineIdx != -1 {
		endIdx = startIdx + newLineIdx
	} else {
		endIdx = len(update)
	}

	return update[startIdx:endIdx]
}
