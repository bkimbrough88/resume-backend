package main

import (
	"encoding/json"
	"fmt"
	resume "github.com/bkimbrough88/resume-backend/pkg"
	"net/http"
	"testing"
)

func TestConstructErrorResponse(t *testing.T) {
	validateErrorResponse(fmt.Errorf("invalid email"), http.StatusBadRequest, t)
	validateErrorResponse(fmt.Errorf("no results found"), http.StatusNotFound, t)
	validateErrorResponse(fmt.Errorf("too many results found"), http.StatusNotFound, t)
	validateErrorResponse(fmt.Errorf("anything else"), http.StatusInternalServerError, t)
}

func validateErrorResponse(providerError error, expectedStatusCode int, t *testing.T) {
	if err := constructErrorResponse(providerError); err == nil {
		t.Error("Expected an error response, but did not get one")
	} else {
		responseStr := err.Error()
		response := &EventResponse{}
		if jsonErr := json.Unmarshal([]byte(responseStr), response); jsonErr != nil {
			t.Errorf("Failed to unmarshal response JSON: %s", jsonErr.Error())
		} else {
			if response.StatusCode != expectedStatusCode {
				t.Errorf("Expected status code for '%s' to be %d, but was %d", providerError.Error(), expectedStatusCode, response.StatusCode)
			}

			if response.Error != providerError.Error() {
				t.Errorf("Expected error to be '%s', but was '%s'", providerError.Error(), response.Error)
			}

			if response.User != nil {
				t.Errorf("User should not have been populated")
			}
		}
	}
}

func TestGetSuccessResponse(t *testing.T) {
	user := &resume.User{
		Certifications: []resume.Certification{
			{
				Name:         "Some Cert",
				BadgeLink:    "https://example.com",
				DateAchieved: "10-28-2019",
				DateExpires:  "10-28-2022",
			},
		},
		Degrees: []resume.Degree{
			{
				Degree:    "BS",
				Major:     "CS",
				School:    "University",
				StartYear: 2017,
				EndYear:   2021,
			},
		},
		Email: "user@domain.com",
		Experience: []resume.Experience{
			{
				Company:    "Co",
				JobTitle:   "SRE",
				StartMonth: "May",
				StartYear:  2020,
				EndMonth:   "June",
				EndYear:    2020,
				Responsibilities: []string{
					"foo",
					"bar",
				},
			},
		},
		Github:      "https://github.com/user",
		GivenName:   "John",
		Location:    "Place, State",
		Linkedin:    "https://www.linkedin.com/in/user",
		PhoneNumber: "999-999-9999",
		Skills: []resume.Skill{
			{
				Name:              "Go",
				YearsOfExperience: 2,
			},
		},
		Summary: "My awesome summary",
		SurName: "Doe",
	}

	if res, err := getSuccessResponse(http.StatusOK, user); err != nil {
		t.Errorf("Did not expect error response: %s", err.Error())
	} else {
		response := &EventResponse{}
		if jsonErr := json.Unmarshal([]byte(res), response); jsonErr != nil {
			t.Errorf("Failed to unmarshal response JSON: %s", jsonErr.Error())
		} else {
			if response.StatusCode != http.StatusOK {
				t.Errorf("Expected status code to be %d, but was %d", http.StatusOK, response.StatusCode)
			}

			if response.Error != "" {
				t.Errorf("Expected error to be empty, but was '%s'", response.Error)
			}

			if response.User == nil {
				t.Error("Expected user to be set, but was nil")
			}
		}
	}

	if res, err := getSuccessResponse(http.StatusAccepted, nil); err != nil {
		t.Errorf("Did not expect error response: %s", err.Error())
	} else {
		response := &EventResponse{}
		if jsonErr := json.Unmarshal([]byte(res), response); jsonErr != nil {
			t.Errorf("Failed to unmarshal response JSON: %s", jsonErr.Error())
		} else {
			if response.StatusCode != http.StatusAccepted {
				t.Errorf("Expected status code to be %d, but was %d", http.StatusAccepted, response.StatusCode)
			}

			if response.Error != "" {
				t.Errorf("Expected error to be empty, but was '%s'", response.Error)
			}

			if response.User != nil {
				t.Error("Expected user to be nil, but was set")
			}
		}
	}
}
