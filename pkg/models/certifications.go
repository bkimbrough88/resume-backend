package models

type Certification struct {
	Name         string `json:"name"`
	DateAchieved string `json:"date_achieved"`
	BadgeLink    string `json:"badge_link,omitempty"`
	DateExpires  string `json:"date_expires,omitempty"`
}

type CertificationKey struct {
	CertificationName string `json:"certification_name"`
}
