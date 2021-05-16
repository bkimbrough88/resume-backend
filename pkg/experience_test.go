package pkg

import (
	"strconv"
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

	if len(expr.Names()) > 1 && *expr.Names()["#0"] == "foo" {
		t.Errorf("Expected to have 0 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) > 1 && *expr.Values()[":0"].S != "bar" {
		t.Errorf("Expected to have 0 values, but got %d", len(expr.Values()))
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
		}

		//TODO: Check for responsibilities
	}
}
