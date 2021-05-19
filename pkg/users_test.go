package pkg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestGetUserUpdateBuilder_Matching(t *testing.T) {
	uuid1 := uuid.New()
	user := User{
		Id: &uuid1,
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

	updateBuilder, err := getUserUpdateBuilder(&user, &user)
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
	uuid1 := uuid.New()
	user := User{
		Id: &uuid1,
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

	uuid2 := uuid.New()
	user2 := User{
		Id: &uuid2,
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

	updateBuilder, err := getUserUpdateBuilder(&user, &user2)
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
	uuid1 := uuid.New()
	user := User{
		Id: &uuid1,
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
	user2 := User{
		Id: &uuid1,
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

	updateBuilder, err := getUserUpdateBuilder(&user, &user2)
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
	uuid1 := uuid.New()
	user := User{
		Id: &uuid1,
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
	user2 := User{
		Id:             &uuid1,
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

	updateBuilder, err := getUserUpdateBuilder(&user, &user2)
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
