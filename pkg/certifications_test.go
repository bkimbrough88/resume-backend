package pkg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestCompareCertifications_Matching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	cert := Certification{
		Name:         "Some Cert",
		BadgeLink:    "https://example.com",
		DateAchieved: "10-28-2019",
		DateExpires:  "10-28-2022",
	}

	compareCertifications(updateBuilder, cert, cert, 0)

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

func TestCompareCertifications_NoneMatching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	cert := Certification{
		Name:         "Some Cert",
		BadgeLink:    "https://example.com",
		DateAchieved: "10-28-2019",
		DateExpires:  "10-28-2022",
	}

	cert2 := Certification{
		Name:         "Another Cert",
		BadgeLink:    "https://domain.com",
		DateAchieved: "12-31-2021",
		DateExpires:  "12-31-2022",
	}

	compareCertifications(updateBuilder, cert, cert2, 0)

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		t.Fatalf("Could not build expression with resulting updateBuilder. Error: %s", err.Error())
	}

	if len(expr.Names()) != 6 {
		t.Errorf("Expected to have 6 names, but got %d", len(expr.Names()))
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

	if strings.Contains(*expr.Update(), "ADD") {
		t.Errorf("Did not expect update expression to ADD values")
	}

	if strings.Contains(*expr.Update(), "REMOVE") {
		t.Errorf("Did not expect update expression to REMOVE values")
	}

	validateCert(cert2, expr, t, 0)
}

func validateCert(updatedCert Certification, expr expression.Expression, t *testing.T, idx int) {
	var certsKey string
	for key, name := range expr.Names() {
		if *name == certifications {
			certsKey = fmt.Sprintf("%s[%d]", key, idx)
		}
	}

	if certsKey == "" {
		t.Fatalf("Expected to find '%s' in the names list, but it was not there", certifications)
	}

	for key, name := range expr.Names() {
		actualValue := expr.Values()[getValueKey(&certsKey, key, *expr.Update())]
		if *name == "Name" {
			if updatedCert.Name != *actualValue.S {
				t.Errorf("Expected Name to be %s, but was %s", updatedCert.Name, *actualValue.S)
			}
		} else if *name == "DateAchieved" {
			if updatedCert.DateAchieved != *actualValue.S {
				t.Errorf("Expected DateAchieved to be %s, but was %s", updatedCert.DateAchieved, *actualValue.S)
			}
		} else if *name == "BadgeLink" {
			if updatedCert.BadgeLink != *actualValue.S {
				t.Errorf("Expected BadgeLink to be %s, but was %s", updatedCert.BadgeLink, *actualValue.S)
			}
		} else if *name == "DateExpires" {
			if updatedCert.DateExpires != *actualValue.S {
				t.Errorf("Expected DateExpires to be %s, but was %s", updatedCert.DateExpires, *actualValue.S)
			}
		}
	}
}
