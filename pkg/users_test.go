package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"go.uber.org/zap"
	"strings"
	"testing"
)

var (
	user   *User
	logger *zap.Logger

	deleteItemMock func(*dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
	findItemMock   func(*dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	putItemMock    func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	updateItemMock func(*dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
)

type dynamoServiceMock struct{}

func (d dynamoServiceMock) deleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	return deleteItemMock(input)
}

func (d dynamoServiceMock) findItem(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return findItemMock(input)
}

func (d dynamoServiceMock) putItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return putItemMock(input)
}

func (d dynamoServiceMock) updateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	return updateItemMock(input)
}

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
	}

	deleteItemMock = func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
		return &dynamodb.DeleteItemOutput{}, nil
	}

	attr, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		t.Fatalf("Failed to marshal user int Dynamo attribute map: %s", err.Error())
	}
	findItemMock = func(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		return &dynamodb.ScanOutput{
			Count: aws.Int64(1),
			Items: []map[string]*dynamodb.AttributeValue{
				attr,
			},
		}, nil
	}

	putItemMock = func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{}, nil
	}

	updateItemMock = func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
		return &dynamodb.UpdateItemOutput{}, nil
	}
}

func TestCreateUser(t *testing.T) {
	setup(t)

	svc := dynamoServiceMock{}
	if err := CreateUser(user, svc, logger); err != nil {
		t.Errorf("Failed to create user when it should have been successful: %s", err.Error())
	}

	expectedError := "some error"
	putItemMock = func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return nil, fmt.Errorf(expectedError)
	}
	if err := CreateUser(user, svc, logger); err == nil {
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
		} else if *input.TableName != usersTable {
			t.Errorf("Expected table name to be '%s', but was '%s'", usersTable, *input.TableName)
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

	key := UserKey{Email: "user@domain.com"}
	svc := dynamoServiceMock{}
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

	findItemMock = func(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		return &dynamodb.ScanOutput{
			Count: aws.Int64(0),
			Items: []map[string]*dynamodb.AttributeValue{},
		}, nil
	}
	if _, err := GetUserByKey(&key, svc, logger); err == nil {
		t.Errorf("Found user when none should have been found")
	} else if err.Error() != "no results found" {
		t.Errorf("Expected error to be 'no results found', but was '%s'", err.Error())
	}

	findItemMock = func(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		return &dynamodb.ScanOutput{
			Count: aws.Int64(2),
			Items: []map[string]*dynamodb.AttributeValue{
				{},
				{},
			},
		}, nil
	}
	if _, err := GetUserByKey(&key, svc, logger); err == nil {
		t.Errorf("Found user when too many should have been found")
	} else if err.Error() != "too many results found" {
		t.Errorf("Expected error to be 'too many results found', but was '%s'", err.Error())
	}

	expectedError := "some error"
	findItemMock = func(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		return nil, fmt.Errorf(expectedError)
	}
	if _, err := GetUserByKey(&key, svc, logger); err == nil {
		t.Errorf("Created user when it should have failed")
	} else if err.Error() != expectedError {
		t.Errorf("Expected error to be '%s', but was '%s'", expectedError, err.Error())
	}
}

func TestGetUserScanInput(t *testing.T) {
	filter := expression.Name("foo").Equal(expression.Value("bar"))
	if input, err := getUserScanInput(&filter); err != nil {
		t.Errorf("Failed to get input with error '%s'", err.Error())
	} else {
		if input.TableName == nil {
			t.Error("Table name should not be nil")
		} else if *input.TableName != usersTable {
			t.Errorf("Expected table name to be '%s', but was '%s'", usersTable, *input.TableName)
		}

		if input.FilterExpression == nil {
			t.Error("Filter should not have generated an empty map")
		} else if *input.FilterExpression != "#0 = :0" {
			t.Errorf("Expected filter expression to be '#0 = :0', but ws '%s'", *input.FilterExpression)
		}

		if len(input.ExpressionAttributeNames) != 1 {
			t.Errorf("Expected to have 1 attribute name in expression, but got %d", len(input.ExpressionAttributeNames))
		} else if *input.ExpressionAttributeNames["#0"] != "foo" {
			t.Errorf("Expected expression attribute name to be 'foo', but was '%s'", *input.ExpressionAttributeNames["#0"])
		}

		if len(input.ExpressionAttributeValues) != 1 {
			t.Errorf("Expected to have 1 attribute value in expression, but got %d", len(input.ExpressionAttributeValues))
		} else if input.ExpressionAttributeValues[":0"].S == nil {
			t.Error("Expected expression attribute value to be a string type and it was not")
		} else if *input.ExpressionAttributeValues[":0"].S != "bar" {
			t.Errorf("Expected expression attribute value to be 'foo', but was '%s'", *input.ExpressionAttributeValues[":0"].S)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	setup(t)

	key := &UserKey{Email: "user@domain.com"}
	svc := dynamoServiceMock{}
	if err := UpdateUser(key, user, svc, logger); err != nil {
		t.Errorf("Failed to update user when it should have been successful: %s", err.Error())
	}

	newKey := &UserKey{Email: "anotheruser@domain.com"}
	if err := UpdateUser(newKey, user, svc, logger); err != nil {
		t.Errorf("Failed to update user when it should have been successful: %s", err.Error())
	}

	badUser := &User{Email: "not an email"}
	if err := UpdateUser(key, badUser, svc, logger); err == nil {
		t.Errorf("Updated user when it should have failed")
	} else if err.Error() != "invalid email" {
		t.Errorf("Expected error to be 'invalid email', but was '%s'", err.Error())
	}

	expectedError := "some error"
	updateItemMock = func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
		return nil, fmt.Errorf(expectedError)
	}
	if err := UpdateUser(key, user, svc, logger); err == nil {
		t.Errorf("Updated user when it should have failed")
	} else if err.Error() != expectedError {
		t.Errorf("Expected error to be '%s', but was '%s'", expectedError, err.Error())
	}
}

func TestGetUserUpdateBuilder_Matching(t *testing.T) {
	setup(t)

	updateBuilder, err := getUserUpdateBuilder(user, user)
	if err != nil {
		t.Fatalf("Did not get update builder. Error: %s", err.Error())
	}

	expr, err := expression.NewBuilder().WithUpdate(*updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 1 {
		t.Errorf("Expected to have 1 name, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 1 {
		t.Errorf("Expected to have 1 value, but got %d", len(expr.Values()))
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	if !strings.Contains(*expr.Update(), "SET") {
		t.Errorf("Expected update expression to SET values")
	}

	if strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Did not expect update expression to ADD values")
	}

	if strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Did not expect update expression to REMOVE values")
	}

	if *expr.Names()["#0"] != "LastUpdated" {
		t.Errorf("Expected to names only contain 'LastUpdated', but was '%s'", *expr.Names()["#0"])
	}
}

func TestGetUserUpdateBuilder_NoneMatch(t *testing.T) {
	setup(t)

	user2 := &User{
		Certifications: []Certification{
			{
				Name:         "Another Cert",
				BadgeLink:    "https://domain.com",
				DateAchieved: "12-31-2021",
				DateExpires:  "12-31-2022",
			},
		},
		Degrees: []Degree{
			{
				Degree:    "BA",
				Major:     "CA",
				School:    "College",
				StartYear: 2018,
				EndYear:   2020,
			},
		},
		Email: "different-user@test.com",
		Experience: []Experience{
			{
				Company:    "Com",
				JobTitle:   "Dev",
				StartMonth: "July",
				StartYear:  2021,
				EndMonth:   "October",
				EndYear:    2021,
				Responsibilities: []string{
					"baz",
					"biz",
				},
			},
		},
		Github:      "https://github.com/different-user",
		GivenName:   "Jane",
		Location:    "Another Place, Another State",
		Linkedin:    "https://www.linkedin.com/in/different-user",
		PhoneNumber: "111-111-1111",
		Skills: []Skill{
			{
				Name:              "Java",
				YearsOfExperience: 10,
			},
		},
		Summary: "A better summary",
		SurName: "Smith",
	}

	updateBuilder, err := getUserUpdateBuilder(user, user2)
	if err != nil {
		t.Fatalf("Did not get update builder. Error: %s", err.Error())
	}

	expr, err := expression.NewBuilder().WithUpdate(*updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 28 {
		t.Errorf("Expected to have 28 name, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 28 {
		t.Errorf("Expected to have 28 value, but got %d", len(expr.Values()))
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	if !strings.Contains(*expr.Update(), "SET") {
		t.Errorf("Expected update expression to SET values")
	}

	if strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Did not expect update expression to ADD values")
	}

	if strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Did not expect update expression to REMOVE values")
	}

	for key, name := range expr.Names() {
		actualValue := expr.Values()[getValueKey(nil, key, *expr.Update())]
		if *name == "Email" {
			if user2.Email != *actualValue.S {
				t.Errorf("Expected Email to be %s, but was %s", user2.Email, *actualValue.S)
			}
		} else if *name == "Github" {
			if user2.Github != *actualValue.S {
				t.Errorf("Expected Github to be %s, but was %s", user2.Github, *actualValue.S)
			}
		} else if *name == "GivenName" {
			if user2.GivenName != *actualValue.S {
				t.Errorf("Expected GivenName to be %s, but was %s", user2.GivenName, *actualValue.S)
			}
		} else if *name == "Location" {
			if user2.Location != *actualValue.S {
				t.Errorf("Expected Location to be %s, but was %s", user2.Location, *actualValue.S)
			}
		} else if *name == "Linkedin" {
			if user2.Linkedin != *actualValue.S {
				t.Errorf("Expected Linkedin to be %s, but was %s", user2.Linkedin, *actualValue.S)
			}
		} else if *name == "PhoneNumber" {
			if user2.PhoneNumber != *actualValue.S {
				t.Errorf("Expected PhoneNumber to be %s, but was %s", user2.PhoneNumber, *actualValue.S)
			}
		} else if *name == "Summary" {
			if user2.Summary != *actualValue.S {
				t.Errorf("Expected Summary to be %s, but was %s", user2.Summary, *actualValue.S)
			}
		} else if *name == "SurName" {
			if user2.SurName != *actualValue.S {
				t.Errorf("Expected SurName to be %s, but was %s", user2.SurName, *actualValue.S)
			}
		} else if *name == "Certifications" {
			validateCert(user2.Certifications[0], expr, t, 0)
		} else if *name == "Degrees" {
			validateDegree(user2.Degrees[0], expr, t, 0)
		} else if *name == "Experience" {
			validateExperience(user2.Experience[0], expr, t, 0)
		} else if *name == "Skills" {
			validateSkill(user2.Skills[0], expr, t, 0)
		}
	}
}

func TestGetUserUpdateBuilder_AddLists(t *testing.T) {
	setup(t)

	user2 := &User{
		Certifications: []Certification{
			{
				Name:         "Some Cert",
				BadgeLink:    "https://example.com",
				DateAchieved: "10-28-2019",
				DateExpires:  "10-28-2022",
			},
			{
				Name:         "Another Cert",
				BadgeLink:    "https://domain.com",
				DateAchieved: "12-31-2021",
				DateExpires:  "12-31-2022",
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
			{
				Degree:    "BA",
				Major:     "CA",
				School:    "College",
				StartYear: 2018,
				EndYear:   2020,
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
			{
				Company:    "Com",
				JobTitle:   "Dev",
				StartMonth: "July",
				StartYear:  2021,
				EndMonth:   "October",
				EndYear:    2021,
				Responsibilities: []string{
					"baz",
					"biz",
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
			{
				Name:              "Java",
				YearsOfExperience: 10,
			},
		},
		Summary: "My awesome summary",
		SurName: "Doe",
	}

	updateBuilder, err := getUserUpdateBuilder(user, user2)
	if err != nil {
		t.Fatalf("Did not get update builder. Error: %s", err.Error())
	}

	expr, err := expression.NewBuilder().WithUpdate(*updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 4 {
		t.Errorf("Expected to have 4 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 5 {
		t.Errorf("Expected to have 5 values, but got %d", len(expr.Values()))
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	if !strings.Contains(*expr.Update(), "SET") {
		t.Errorf("Expected update expression to SET values")
	}

	if !strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Expected update expression to ADD values")
	}

	if strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Did not expect update expression to REMOVE values")
	}

	updateStatements := strings.Split(*expr.Update(), "\n")
	for key, name := range expr.Names() {
		if *name == "Certifications" {
			for _, statement := range updateStatements {
				if strings.Contains(statement, "ADD") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 1)) {
						t.Fatal("Expected the certification at index 1 to be added, but was not found in the ADD statement")
					}
				}
			}
		} else if *name == "Degrees" {
			for _, statement := range updateStatements {
				if strings.Contains(statement, "ADD") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 1)) {
						t.Fatal("Expected the degree at index 1 to be added, but was not found in the ADD statement")
					}
				}
			}
		} else if *name == "Experience" {
			for _, statement := range updateStatements {
				if strings.Contains(statement, "ADD") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 1)) {
						t.Fatal("Expected the experience at index 1 to be added, but was not found in the ADD statement")
					}
				}
			}
		} else if *name == "Skills" {
			for _, statement := range updateStatements {
				if strings.Contains(statement, "ADD") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 1)) {
						t.Fatal("Expected the skill at index 1 to be added, but was not found in the ADD statement")
					}
				}
			}
		}
	}
}

func TestGetUserUpdateBuilder_RemoveLists(t *testing.T) {
	setup(t)

	user2 := &User{
		Certifications: []Certification{},
		Degrees:        []Degree{},
		Email:          "user@domain.com",
		Experience:     []Experience{},
		Github:         "https://github.com/user",
		GivenName:      "John",
		Location:       "Place, State",
		Linkedin:       "https://www.linkedin.com/in/user",
		PhoneNumber:    "999-999-9999",
		Skills:         []Skill{},
		Summary:        "My awesome summary",
		SurName:        "Doe",
	}

	updateBuilder, err := getUserUpdateBuilder(user, user2)
	if err != nil {
		t.Fatalf("Did not get update builder. Error: %s", err.Error())
	}

	expr, err := expression.NewBuilder().WithUpdate(*updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 4 {
		t.Errorf("Expected to have 4 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 1 {
		t.Errorf("Expected to have 1 value, but got %d", len(expr.Values()))
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	if !strings.Contains(*expr.Update(), "SET") {
		t.Errorf("Expected update expression to SET values")
	}

	if strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Did not expect update expression to ADD values")
	}

	if !strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Expected update expression to REMOVE values")
	}

	updateStatements := strings.Split(*expr.Update(), "\n")
	for key, name := range expr.Names() {
		if *name == "Certifications" {
			for _, statement := range updateStatements {
				if strings.Contains(statement, "REMOVE") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 0)) {
						t.Fatal("Expected the certification at index 0 to be removed, but was not found in the REMOVE statement")
					}
				}
			}
		} else if *name == "Degrees" {
			for _, statement := range updateStatements {
				if strings.Contains(statement, "REMOVE") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 0)) {
						t.Fatal("Expected the degree at index 0 to be removed, but was not found in the REMOVE statement")
					}
				}
			}
		} else if *name == "Experience" {
			for _, statement := range updateStatements {
				if strings.Contains(statement, "REMOVE") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 0)) {
						t.Fatal("Expected the experience at index 0 to be removed, but was not found in the REMOVE statement")
					}
				}
			}
		} else if *name == "Skills" {
			for _, statement := range updateStatements {
				if strings.Contains(statement, "REMOVE") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 0)) {
						t.Fatal("Expected the skill at index 0 to be removed, but was not found in the REMOVE statement")
					}
				}
			}
		}
	}
}

func TestDeleteUser(t *testing.T) {
	setup(t)

	key := UserKey{Email: "user@domain.com"}
	svc := dynamoServiceMock{}
	if err := DeleteUser(&key, svc, logger); err != nil {
		t.Errorf("Failed to delete user when it should have been successful: %s", err.Error())
	}

	expectedError := "some error"
	deleteItemMock = func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
		return nil, fmt.Errorf(expectedError)
	}
	if err := DeleteUser(&key, svc, logger); err == nil {
		t.Errorf("Deleted user when it should have failed")
	} else if err.Error() != expectedError {
		t.Errorf("Expected error to be '%s', but was '%s'", expectedError, err.Error())
	}
}

func TestGetUserDeleteInput(t *testing.T) {
	email := "user@domain.com"
	key := &UserKey{Email: email}
	if input, err := getUserDeleteInput(key); err != nil {
		t.Errorf("Failed to get input with error '%s'", err.Error())
	} else {
		if input.TableName == nil {
			t.Error("Table name should not be nil")
		} else if *input.TableName != usersTable {
			t.Errorf("Expected table name to be '%s', but was '%s'", usersTable, *input.TableName)
		}

		if input.Key == nil {
			t.Error("User key should not have generated an empty map")
		} else if input.Key["email"].S == nil {
			t.Error("Expected email to be a string type")
		} else if *input.Key["email"].S != email {
			t.Errorf("Expected email to be '%s', but was '%s'", email, *input.Key["email"].S)
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
