package main

import (
	"errors"
	"testing"
)

func TestGetIssueID(t *testing.T) {
	tests := []struct {
		Name   string
		Given  string
		Expect string
		Err    error
	}{
		{
			Name:   "Happy Path",
			Given:  "https://user-testing.atlassian.net/browse/TOOL-936",
			Expect: "TOOL-936",
			Err:    nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			result, err := getIssueID(tc.Given)
			if !errors.Is(err, tc.Err) {
				t.Fatalf("expected error to be %v, got: %v", tc.Err, err)
			}
			if result != tc.Expect {
				t.Fatalf("expected result to be %v, got: %v", tc.Expect, result)
			}
		})
	}
}

func TestCreateBranchName(t *testing.T) {
	tests := []struct {
		Name    string
		ID      string
		Summary string
		Expect  string
	}{
		{
			Name:    "Happy Path",
			ID:      "TOOL-936",
			Summary: "abc def",
			Expect:  "TOOL-936--abc-def",
		},
		{
			Name:    "Leading spaces",
			ID:      "TOOL-936",
			Summary: "   Refactor specs for non ut marketing",
			Expect:  "TOOL-936--refactor-specs-for-non-ut-marketing",
		},
		{
			Name:    "Trailing spaces",
			ID:      "TOOL-936",
			Summary: "Refactor specs for non ut marketing ",
			Expect:  "TOOL-936--refactor-specs-for-non-ut-marketing",
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			name := createBranchName(tc.ID, tc.Summary)
			if name != tc.Expect {
				t.Fatalf("expected %v, got %v", tc.Expect, name)
			}
		})
	}
}
