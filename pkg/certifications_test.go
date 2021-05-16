package pkg

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestCompareCertifications_Matching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	cert := Certification{
		Name:        "Some Cert",
		BadgeLink:   "https://example.com",
		DateAchieved: "10-28-2019",
		DateExpires: "10-28-2022",
	}

	compareCertifications(updateBuilder, cert, cert, 0)

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

func TestCompareCertifications_NoneMatching(t *testing.T) {
	updateBuilder := expression.Set(expression.Name("foo"), expression.Value("bar"))

	cert := Certification{
		Name:        "Some Cert",
		BadgeLink:   "https://example.com",
		DateAchieved: "10-28-2019",
		DateExpires: "10-28-2022",
	}

	cert2 := Certification{
		Name:        "Another Cert",
		BadgeLink:   "https://domain.com",
		DateAchieved: "12-31-2021",
		DateExpires: "12-31-2022",
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

	for key, name := range expr.Names() {
		actualValue := expr.Values()[getValueKey(key, *expr.Update())]
		if *name == "Name" {
			if cert2.Name != *actualValue.S {
				t.Errorf("Expected Name to be %s, but was %s", cert2.Name, *actualValue.S)
			}
		} else if *name == "DateAchieved" {
			if cert2.DateAchieved != *actualValue.S {
				t.Errorf("Expected DateAchieved to be %s, but was %s", cert2.DateAchieved, *actualValue.S)
			}
		} else if *name == "BadgeLink" {
			if cert2.BadgeLink != *actualValue.S {
				t.Errorf("Expected BadgeLink to be %s, but was %s", cert2.BadgeLink, *actualValue.S)
			}
		} else if *name == "DateExpires" {
			if cert2.DateExpires != *actualValue.S {
				t.Errorf("Expected DateExpires to be %s, but was %s", cert2.DateExpires, *actualValue.S)
			}
		}
	}
}
