package pkg

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	certifications = "Certifications"
)

type Certification struct {
	Name         string `json:"name"`
	DateAchieved string `json:"date_achieved"`
	BadgeLink    string `json:"badge_link,omitempty"`
	DateExpires  string `json:"date_expires,omitempty"`
}

type CertificationKey struct {
	CertificationName string `json:"certification_name"`
}

func compareCertifications(updateBuilder expression.UpdateBuilder, currentCert Certification, updatedCert Certification, idx int) {
	if currentCert.Name != updatedCert.Name {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, certifications, idx, "Name")), expression.Value(updatedCert.Name))
	}

	if currentCert.DateAchieved != updatedCert.DateAchieved {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, certifications, idx, "DateAchieved")), expression.Value(updatedCert.DateAchieved))
	}

	if currentCert.BadgeLink != updatedCert.BadgeLink {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, certifications, idx, "BadgeLink")), expression.Value(updatedCert.BadgeLink))
	}

	if currentCert.DateExpires != updatedCert.DateExpires {
		updateBuilder.Set(expression.Name(fmt.Sprintf(listElementNameFormat, certifications, idx, "DateExpires")), expression.Value(updatedCert.DateExpires))
	}
}
