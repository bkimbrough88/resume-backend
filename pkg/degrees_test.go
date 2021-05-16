package pkg

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestCompareDegrees_Matching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	degree := Degree{
		Degree: "BS",
		Major: "CS",
		School: "University",
		StartYear: 2017,
		EndYear: 2021,
	}

	compareDegrees(updateBuilder, degree, degree, 0)

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

func TestCompareDegrees_NoneMatching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	degree := Degree{
		Degree: "BS",
		Major: "CS",
		School: "University",
		StartYear: 2017,
		EndYear: 2021,
	}

	degree2 := Degree{
		Degree: "BA",
		Major: "CA",
		School: "College",
		StartYear: 2018,
		EndYear: 2020,
	}

	compareDegrees(updateBuilder, degree, degree2, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 7 {
		t.Errorf("Expected to have 7 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 6 {
		t.Errorf("Expected to have 6 values, but got %d", len(expr.Values()))
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	for key, name := range expr.Names() {
		actualValue := expr.Values()[getValueKey(key, *expr.Update())]
		if *name == "Degree" {
			if degree2.Degree != *actualValue.S {
				t.Errorf("Expected Degree to be %s, but was %s", degree2.Degree, *actualValue.S)
			}
		} else if *name == "Major" {
			if degree2.Major != *actualValue.S {
				t.Errorf("Expected Major to be %s, but was %s", degree2.Major, *actualValue.S)
			}
		} else if *name == "School" {
			if degree2.School != *actualValue.S {
				t.Errorf("Expected School to be %s, but was %s", degree2.School, *actualValue.S)
			}
		} else if *name == "StartYear" {
			if actualNumber, err := strconv.Atoi(*actualValue.N); err != nil {
				t.Errorf("Could not parse number from '%s'. Error: %s", *actualValue.N, err.Error())
			} else if degree2.StartYear != actualNumber {
				t.Errorf("Expected StartYear to be %d, but was %d", degree2.StartYear, actualNumber)
			}
		} else if *name == "EndYear" {
			if actualNumber, err := strconv.Atoi(*actualValue.N); err != nil {
				t.Errorf("Could not parse number from '%s'. Error: %s", *actualValue.N, err.Error())
			} else if degree2.EndYear != actualNumber {
				t.Errorf("Expected EndYear to be %d, but was %d", degree2.EndYear, actualNumber)
			}
		}
	}
}
