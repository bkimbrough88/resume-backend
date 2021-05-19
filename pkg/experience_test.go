package pkg

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestCompareExperience_Matching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	experience1 := Experience{
		Company:          "Co",
		JobTitle:         "SRE",
		StartMonth:       "May",
		StartYear:        2020,
		EndMonth: "June",
		EndYear: 2020,
		Responsibilities: []string{
			"foo",
			"bar",
		},
	}

	compareExperience(updateBuilder, experience1, experience1, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
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

	if *expr.Names()["#0"] != "foo" {
		t.Errorf("Expected names to only contain 'foo', but was '%s'", *expr.Names()["#0"])
	}

	if expr.Values()[":0"].S == nil {
		t.Fatal("Expected value be a string, but the string value was null")
	}

	if *expr.Values()[":0"].S != "bar" {
		t.Fatalf("Expected values to only contain 'bar', but was '%s'", *expr.Values()[":0"].S)
	}
}

func TestCompareExperience_NoneMatching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	experience1 := Experience{
		Company:          "Co",
		JobTitle:         "SRE",
		StartMonth:       "May",
		StartYear:        2020,
		EndMonth: "June",
		EndYear: 2020,
		Responsibilities: []string{
			"foo",
			"bar",
		},
	}

	experience2 := Experience{
		Company:          "Com",
		JobTitle:         "Dev",
		StartMonth:       "July",
		StartYear:        2021,
		EndMonth: "October",
		EndYear: 2021,
		Responsibilities: []string{
			"baz",
			"biz",
		},
	}

	compareExperience(updateBuilder, experience1, experience2, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 9 {
		t.Errorf("Expected to have 9 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 9 {
		t.Errorf("Expected to have 9 values, but got %d", len(expr.Values()))
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
		actualValue := expr.Values()[getValueKey(key, *expr.Update())]
		if *name == "Company" {
			if experience2.Company != *actualValue.S {
				t.Errorf("Expected Company to be %s, but was %s", experience2.Company, *actualValue.S)
			}
		} else if *name == "JobTitle" {
			if experience2.JobTitle != *actualValue.S {
				t.Errorf("Expected JobTitle to be %s, but was %s", experience2.JobTitle, *actualValue.S)
			}
		} else if *name == "StartMonth" {
			if experience2.StartMonth != *actualValue.S {
				t.Errorf("Expected StartMonth to be %s, but was %s", experience2.StartMonth, *actualValue.S)
			}
		} else if *name == "StartYear" {
			if actualNumber, err := strconv.Atoi(*actualValue.N); err != nil {
				t.Errorf("Could not parse number from '%s'. Error: %s", *actualValue.N, err.Error())
			} else if experience2.StartYear != actualNumber {
				t.Errorf("Expected StartYear to be %d, but was %d", experience2.StartYear, actualNumber)
			}
		} else if *name == "EndMonth" {
			if experience2.EndMonth != *actualValue.S {
				t.Errorf("Expected EndMonth to be %s, but was %s", experience2.EndMonth, *actualValue.S)
			}
		} else if *name == "EndYear" {
			if actualNumber, err := strconv.Atoi(*actualValue.N); err != nil {
				t.Errorf("Could not parse number from '%s'. Error: %s", *actualValue.N, err.Error())
			} else if experience2.EndYear != actualNumber {
				t.Errorf("Expected EndYear to be %d, but was %d", experience2.EndYear, actualNumber)
			}
		} else if *name == "Responsibilities" {
			found := false
			for _, expectedResponsibility := range experience2.Responsibilities {
				if expectedResponsibility == *actualValue.S {
					found = true
					continue
				}
			}

			if !found {
				t.Errorf("Got the responsibility %s, but it was not found in the updated object", *actualValue.S)
			}
		}
	}
}

func TestCompareExperience_AddResponsibility(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	experience1 := Experience{
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
	}

	experience2 := Experience{
		Company:    "Co",
		JobTitle:   "SRE",
		StartMonth: "May",
		StartYear:  2020,
		EndMonth:   "June",
		EndYear:    2020,
		Responsibilities: []string{
			"foo",
			"bar",
			"baz",
		},
	}

	compareExperience(updateBuilder, experience1, experience2, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 3 {
		t.Errorf("Expected to have 3 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 2 {
		t.Errorf("Expected to have 2 values, but got %d", len(expr.Values()))
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	if !strings.Contains(*expr.Update(), "SET #2 = :1\n") && *expr.Names()["#2"] == "foo" && *expr.Values()[":1"].S == "bar" {
		t.Errorf("Expected update expression to only SET initializer value")
	}

	if !strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Expected update expression to ADD values")
	}

	if strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Did not expect update expression to REMOVE values")
	}

	for key, name := range expr.Names() {
		if *name == "Responsibilities" {
			updateStatements := strings.Split(*expr.Update(), "\n")
			for _, statement := range updateStatements {
				if strings.Contains(statement, "ADD") {
					if !strings.Contains(statement, key) {
						t.Fatal("Expected the responsibility to be added, but was not found in the ADD statement")
					}

					actualValue := expr.Values()[getValueKey(key, statement)]
					if experience2.Responsibilities[2] != *actualValue.S {
						t.Errorf("Expected added responsibility to be '%s', but was '%s'", experience2.Responsibilities[2], *actualValue.S)
					}
				}
			}
		}
	}
}

func TestCompareExperience_ModifyAndAddResponsibility(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	experience1 := Experience{
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
	}

	experience2 := Experience{
		Company:    "Co",
		JobTitle:   "SRE",
		StartMonth: "May",
		StartYear:  2020,
		EndMonth:   "June",
		EndYear:    2020,
		Responsibilities: []string{
			"biz",
			"bar",
			"baz",
		},
	}

	compareExperience(updateBuilder, experience1, experience2, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 3 {
		t.Errorf("Expected to have 3 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 3 {
		t.Errorf("Expected to have 3 values, but got %d", len(expr.Values()))
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

	for key, name := range expr.Names() {
		if *name == "Responsibilities" {
			updateStatements := strings.Split(*expr.Update(), "\n")
			for _, statement := range updateStatements {
				if strings.Contains(statement, "ADD") {
					if !strings.Contains(statement, key) {
						t.Fatal("Expected the responsibility to be added, but was not found in the ADD statement")
					}

					actualValue := expr.Values()[getValueKey(key, statement)]
					if experience2.Responsibilities[2] != *actualValue.S {
						t.Errorf("Expected added responsibility to be '%s', but was '%s'", experience2.Responsibilities[2], *actualValue.S)
					}
				} else if strings.Contains(statement, "SET") {
					if !strings.Contains(statement, key) {
						t.Fatal("Expected the responsibilities to be updated, but was not found in the SET statement")
					}

					actualValue := expr.Values()[getValueKey(key, statement)]
					if experience2.Responsibilities[0] != *actualValue.S {
						t.Errorf("Expected updated responsibility to be '%s', but was '%s'", experience2.Responsibilities[0], *actualValue.S)
					}
				}
			}
		}
	}
}

func TestCompareExperience_RemoveResponsibility(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	experience1 := Experience{
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
	}

	experience2 := Experience{
		Company:    "Co",
		JobTitle:   "SRE",
		StartMonth: "May",
		StartYear:  2020,
		EndMonth:   "June",
		EndYear:    2020,
		Responsibilities: []string{
			"foo",
		},
	}

	compareExperience(updateBuilder, experience1, experience2, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 3 {
		t.Errorf("Expected to have 3 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 1 {
		t.Errorf("Expected to have 1 value, but got %d", len(expr.Values()))
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	if !strings.Contains(*expr.Update(), "SET #0 = :0\n") && *expr.Names()["#0"] == "foo" {
		t.Errorf("Expected update expression to only SET initializer value")
	}

	if strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Did not expect update expression to ADD values")
	}

	if !strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Expected update expression to REMOVE values")
	}

	for key, name := range expr.Names() {
		if *name == "Responsibilities" {
			updateStatements := strings.Split(*expr.Update(), "\n")
			for _, statement := range updateStatements {
				if strings.Contains(statement, "REMOVE") {
					if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 1)) {
						t.Fatal("Expected the responsibility at index 1 to be removed, but was not found in the REMOVE statement")
					}
				}
			}
		}
	}
}

func TestCompareExperience_ModifyAndRemoveResponsibility(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	experience1 := Experience{
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
	}

	experience2 := Experience{
		Company:    "Co",
		JobTitle:   "SRE",
		StartMonth: "May",
		StartYear:  2020,
		EndMonth:   "June",
		EndYear:    2020,
		Responsibilities: []string{
			"baz",
		},
	}

	compareExperience(updateBuilder, experience1, experience2, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 3 {
		t.Errorf("Expected to have 3 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 2 {
		t.Errorf("Expected to have 2 values, but got %d", len(expr.Values()))
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

	for key, name := range expr.Names() {
		if *name == "Responsibilities" {
			updateStatements := strings.Split(*expr.Update(), "\n")
			for _, statement := range updateStatements {
				if strings.Contains(statement, "ADD") {
					if strings.Contains(statement, "REMOVE") {
						if !strings.Contains(statement, fmt.Sprintf("%s[%d]", key, 1)) {
							t.Fatal("Expected the responsibility at index 1 to be removed, but was not found in the REMOVE statement")
						}
					}
				} else if strings.Contains(statement, "SET") {
					if !strings.Contains(statement, key) {
						t.Fatal("Expected the responsibilities to be updated, but was not found in the SET statement")
					}

					actualValue := expr.Values()[getValueKey(key, statement)]
					if experience2.Responsibilities[0] != *actualValue.S {
						t.Errorf("Expected updated responsibility to be '%s', but was '%s'", experience2.Responsibilities[0], *actualValue.S)
					}
				}
			}
		}
	}
}
