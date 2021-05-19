package pkg

import (
	"fmt"
	"strconv"
	"strings"
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

	if !strings.Contains(*expr.Update(), "SET") {
		t.Errorf("Expected update expression to SET values")
	}

	if strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Did not expect update expression to ADD values")
	}

	if strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Did not expect update expression to REMOVE values")
	}

	// Exit if the counts are off
	if t.Failed() {
		t.FailNow()
	}

	validateSkill(skill2, expr, t, 0)
}

func validateSkill(updatedSkill Skill, expr expression.Expression, t *testing.T, idx int) {
	var skillssKey string
	for key, name := range expr.Names() {
		if *name == skills {
			skillssKey = fmt.Sprintf("%s[%d]", key, idx)
		}
	}

	if skillssKey == "" {
		t.Fatalf("Expected to find '%s' in the names list, but it was not there", skills)
	}

	for key, name := range expr.Names() {
		actualValue := expr.Values()[getValueKey(&skillssKey, key, *expr.Update())]
		if *name == "Name" {
			if updatedSkill.Name != *actualValue.S {
				t.Errorf("Expected Name to be %s, but was %s", updatedSkill.Name, *actualValue.S)
			}
		} else if *name == "YearsOfExperience" {
			if actualNumber, err := strconv.Atoi(*actualValue.N); err != nil {
				t.Errorf("Could not parse number from '%s'. Error: %s", *actualValue.N, err.Error())
			} else if updatedSkill.YearsOfExperience != actualNumber {
				t.Errorf("Expected YearsOfExperience to be %d, but was %d", updatedSkill.YearsOfExperience, actualNumber)
			}
		}
	}
}
