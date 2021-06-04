package models

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestCompareDegrees_Matching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	degree := Degree{
		Degree:    "BS",
		Major:     "CS",
		School:    "University",
		StartYear: 2017,
		EndYear:   2021,
	}

	compareDegrees(updateBuilder, degree, degree, 0)

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

func TestCompareDegrees_NoneMatching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	degree := Degree{
		Degree:    "BS",
		Major:     "CS",
		School:    "University",
		StartYear: 2017,
		EndYear:   2021,
	}

	degree2 := Degree{
		Degree:    "BA",
		Major:     "CA",
		School:    "College",
		StartYear: 2018,
		EndYear:   2020,
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

	if !strings.Contains(*expr.Update(), "SET") {
		t.Errorf("Expected update expression to SET values")
	}

	if strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Did not expect update expression to ADD values")
	}

	if strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Did not expect update expression to REMOVE values")
	}

	validateDegree(degree2, expr, t, 0)
}

func validateDegree(updatedDegree Degree, expr expression.Expression, t *testing.T, idx int) {
	var degreesKey string
	for key, name := range expr.Names() {
		if *name == degrees {
			degreesKey = fmt.Sprintf("%s[%d]", key, idx)
		}
	}

	if degreesKey == "" {
		t.Fatalf("Expected to find '%s' in the names list, but it was not there", degrees)
	}

	for key, name := range expr.Names() {
		actualValue := expr.Values()[getValueKey(&degreesKey, key, *expr.Update())]
		if *name == "Degree" {
			if updatedDegree.Degree != *actualValue.S {
				t.Errorf("Expected Degree to be %s, but was %s", updatedDegree.Degree, *actualValue.S)
			}
		} else if *name == "Major" {
			if updatedDegree.Major != *actualValue.S {
				t.Errorf("Expected Major to be %s, but was %s", updatedDegree.Major, *actualValue.S)
			}
		} else if *name == "School" {
			if updatedDegree.School != *actualValue.S {
				t.Errorf("Expected School to be %s, but was %s", updatedDegree.School, *actualValue.S)
			}
		} else if *name == "StartYear" {
			if actualNumber, err := strconv.Atoi(*actualValue.N); err != nil {
				t.Errorf("Could not parse number from '%s'. Error: %s", *actualValue.N, err.Error())
			} else if updatedDegree.StartYear != actualNumber {
				t.Errorf("Expected StartYear to be %d, but was %d", updatedDegree.StartYear, actualNumber)
			}
		} else if *name == "EndYear" {
			if actualNumber, err := strconv.Atoi(*actualValue.N); err != nil {
				t.Errorf("Could not parse number from '%s'. Error: %s", *actualValue.N, err.Error())
			} else if updatedDegree.EndYear != actualNumber {
				t.Errorf("Expected EndYear to be %d, but was %d", updatedDegree.EndYear, actualNumber)
			}
		}
	}
}
