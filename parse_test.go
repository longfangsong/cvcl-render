package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseResumeTypst(t *testing.T) {
	// Read the resume file
	content, err := os.ReadFile("templates/resume.typ")
	print(content)
	if err != nil {
		t.Fatalf("Failed to read resume file: %v", err)
	}

	// Parse it
	resume, err := ParseResumeTypst(string(content))
	if err != nil {
		t.Fatalf("Failed to parse resume: %v", err)
	}

	// Basic validations
	if len(resume.Positions) == 0 {
		t.Errorf("Expected positions, got %d", len(resume.Positions))
	}

	if resume.Summary == "" {
		t.Error("Expected summary, got empty string")
	}

	if len(resume.Education) == 0 {
		t.Errorf("Expected education entries, got %d", len(resume.Education))
	}

	if len(resume.WorkExperience) == 0 {
		t.Errorf("Expected work experience entries, got %d", len(resume.WorkExperience))
	}

	if len(resume.Projects) == 0 {
		t.Errorf("Expected projects, got %d", len(resume.Projects))
	}

	if len(resume.Skills) == 0 {
		t.Errorf("Expected skills categories, got %d", len(resume.Skills))
	}

	if len(resume.Interests) == 0 {
		t.Errorf("Expected interests, got %d", len(resume.Interests))
	}

	// Print the parsed data as JSON for inspection
	jsonData, _ := json.MarshalIndent(resume, "", "  ")
	t.Logf("Parsed Resume:\n%s", string(jsonData))
}

// func TestParsePositions(t *testing.T) {
// 	content, _ := os.ReadFile("templates/resume.typ")
// 	positions := parsePositions(string(content))

// 	expectedPositions := []string{
// 		"Software Engineer",
// 		"Embedded developer",
// 		"Fullstack developer",
// 		"Cloud developer",
// 		"System software developer",
// 	}

// 	if len(positions) != len(expectedPositions) {
// 		t.Errorf("Expected %d positions, got %d", len(expectedPositions), len(positions))
// 	}

// 	for i, expected := range expectedPositions {
// 		if i >= len(positions) {
// 			t.Errorf("Position %d missing", i)
// 			continue
// 		}
// 		if positions[i] != expected {
// 			t.Errorf("Position %d: expected %q, got %q", i, expected, positions[i])
// 		}
// 	}
// }

// func TestParseSummary(t *testing.T) {
// 	content, _ := os.ReadFile("templates/resume.typ")
// 	summary := parseSummary(string(content))

// 	if summary == "" {
// 		t.Error("Expected non-empty summary")
// 	}

// 	if !contains(summary, "software developer") && !contains(summary, "Software developer") {
// 		t.Errorf("Summary does not contain expected content: %s", summary)
// 	}
// }

// func TestParseEducation(t *testing.T) {
// 	content, _ := os.ReadFile("templates/resume.typ")
// 	education := parseSection(string(content), "Education")

// 	if len(education) < 2 {
// 		t.Errorf("Expected at least 2 education entries, got %d", len(education))
// 	}

// 	// Check first entry
// 	if len(education) > 0 {
// 		first := education[0]
// 		if !contains(first.Title, "Chalmers") && !contains(first.Title, "Technology") {
// 			t.Errorf("First education title unexpected: %s", first.Title)
// 		}
// 		if len(first.Items) == 0 {
// 			t.Error("First education entry has no items")
// 		}
// 	}
// }

// func TestParseSkills(t *testing.T) {
// 	content, _ := os.ReadFile("templates/resume.typ")
// 	skills := parseSkills(string(content))

// 	if len(skills) == 0 {
// 		t.Error("Expected skill categories, got 0")
// 	}

// 	// Check if Languages category exists
// 	found := false
// 	for _, skill := range skills {
// 		if skill.Name == "Languages" {
// 			found = true
// 			if len(skill.Skills) == 0 {
// 				t.Error("Languages category has no skills")
// 			}
// 			break
// 		}
// 	}
// 	if !found {
// 		t.Error("Languages skill category not found")
// 	}
// }

// func TestParseInterests(t *testing.T) {
// 	content, _ := os.ReadFile("templates/resume.typ")
// 	interests := parseInterests(string(content))

// 	if len(interests) == 0 {
// 		t.Error("Expected interest items, got 0")
// 	}

// 	// Check if Technical Writing is present
// 	found := false
// 	for _, interest := range interests {
// 		if interest.Category == "Technical Writing" {
// 			found = true
// 			if interest.Description == "" {
// 				t.Error("Technical Writing has empty description")
// 			}
// 			break
// 		}
// 	}
// 	if !found {
// 		t.Error("Technical Writing interest not found")
// 	}
// }

// // Helper function
// func contains(s, substr string) bool {
// 	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && s[:len(substr)] == substr) || stringContainsSubstring(s, substr))
// }

// func stringContainsSubstring(s, substr string) bool {
// 	for i := 0; i <= len(s)-len(substr); i++ {
// 		if s[i:i+len(substr)] == substr {
// 			return true
// 		}
// 	}
// 	return false
// }
