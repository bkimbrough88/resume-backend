package pkg

import (
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestCompareSkills_Matching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	skill := Skill{
		Name:              "Go",
		YearsOfExperience: 2,
	}

	compareSkills(updateBuilder, skill, skill, 0)

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

func TestCompareSkills_NoneMatching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	skill := Skill{
		Name:              "Go",
		YearsOfExperience: 2,
	}

	skill2 := Skill{
		Name:              "Java",
		YearsOfExperience: 10,
	}

	compareSkills(updateBuilder, skill, skill2, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 4 {
		t.Errorf("Expected to have 4 names, but got %d", len(expr.Names()))
	}

	if len(expr.Values()) != 3 {
		t.Errorf("Expected to have 3 values, but got %d", len(expr.Values()))
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	for key, name := range expr.Names() {
		actualValue := expr.Values()[getValueKey(key, *expr.Update())]
		if *name == "Name" {
			if skill2.Name != *actualValue.S {
				t.Errorf("Expected Name to be %s, but was %s", skill2.Name, *actualValue.S)
			}
		} else if *name == "YearsOfExperience" {
			if actualNumber, err := strconv.Atoi(*actualValue.N); err != nil {
				t.Errorf("Could not parse number from '%s'. Error: %s", *actualValue.N, err.Error())
			} else if skill2.YearsOfExperience != actualNumber {
				t.Errorf("Expected YearsOfExperience to be %d, but was %d", skill2.YearsOfExperience, actualNumber)
			}
		}
	}
}
