package models

type Skill struct {
	Name              string `json:"name"`
	YearsOfExperience int    `json:"years_of_experience,omitempty"`
}

type SkillKey struct {
	Name string `json:"name"`
}
