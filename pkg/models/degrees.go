package models

type Degree struct {
	Degree    string `json:"degree"`
	Major     string `json:"major"`
	School    string `json:"school"`
	StartYear int    `json:"start_year"`
	EndYear   int    `json:"end_year,omitempty"`
}

type DegreeKey struct {
	Degree string `json:"degree"`
	Major  string `json:"major"`
	School string `json:"school"`
}
