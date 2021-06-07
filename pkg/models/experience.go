package models

type Experience struct {
	Company          string   `json:"company"`
	JobTitle         string   `json:"job_title"`
	StartMonth       string   `json:"start_month"`
	StartYear        int      `json:"start_year"`
	EndMonth         string   `json:"end_month,omitempty"`
	EndYear          int      `json:"end_year,omitempty"`
	Responsibilities []string `json:"responsibilities,omitempty"`
}

type ExperienceKey struct {
	Company  string `json:"company"`
	JobTitle string `json:"job_title"`
}
