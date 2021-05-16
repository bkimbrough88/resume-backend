package pkg

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	certifications = "Certifications"
)

type Certification struct {
	Name        string `json:"name"`
	BadgeLink   string `json:"badge_link,omitempty"`
	Description string `json:"description,omitempty"`
}

type CertificationKey struct {
	CertificationName string `json:"certification_name"`
}

type Certifications struct {
	Certifications []Certification `json:"certifications"`
}

func compareCertifications(builder expression.UpdateBuilder, currentCert Certification, updatedCert Certification, idx int) {
	if currentCert.Name != updatedCert.Name {
		builder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, certifications, idx, "Name")), expression.Value(updatedCert.Name))
	}

	if currentCert.BadgeLink != updatedCert.BadgeLink {
		builder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, certifications, idx, "BadgeLink")), expression.Value(updatedCert.BadgeLink))
	}

	if currentCert.Description != updatedCert.Description {
		builder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, certifications, idx, "Description")), expression.Value(updatedCert.Description))
	}
}
