package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

//go:embed templates/resume.typ.template
var embeddedResumeTemplates embed.FS

// Resume represents the complete resume structure
type Resume struct {
	Positions      []string        `json:"positions"`
	Summary        string          `json:"summary"`
	Education      []ResumeEntry   `json:"education"`
	WorkExperience []ResumeEntry   `json:"work_experience"`
	Projects       []ResumeEntry   `json:"projects"`
	Skills         []SkillCategory `json:"skills"`
	Interests      []InterestItem  `json:"interests"`
}

// ResumeEntry represents a single entry with title, location, date, and content
type ResumeEntry struct {
	Title       string `json:"title"`
	Location    string `json:"location,omitempty"`
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
}

// SkillCategory represents a category of skills with a name and list of skills
type SkillCategory struct {
	Name   string      `json:"name"`
	Skills []SkillItem `json:"skills"`
}

// SkillItem represents a single skill with its name and whether it's strong/emphasized
type SkillItem struct {
	Name   string `json:"name"`
	Strong bool   `json:"strong"`
}

// InterestItem represents a single interest with a category and description
type InterestItem struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}

// ParseResumeTypst parses a Typst resume file and returns a Resume struct
func ParseResumeTypst(content string) (*Resume, error) {
	resume := &Resume{}

	// Parse positions
	resume.Positions = parsePositions(content)

	// Parse summary
	resume.Summary = parseSummary(content)

	// Parse education
	resume.Education = parseSection(content, "Education")

	// Parse working experience
	resume.WorkExperience = parseSection(content, "Working Experience")

	// Parse projects
	resume.Projects = parseSection(content, "Projects")

	// Parse skills
	resume.Skills = parseSkills(content)

	// Parse interests
	resume.Interests = parseInterests(content)

	return resume, nil
}

// parsePositions extracts the positions array from the typst file
func parsePositions(content string) []string {
	var positions []string

	// Match the positions tuple in the author section
	positionPattern := regexp.MustCompile(`positions:\s*\(([\s\S]*?)\)`)
	match := positionPattern.FindStringSubmatch(content)
	if len(match) < 2 {
		return positions
	}

	positionsText := match[1]
	// Extract quoted strings
	stringPattern := regexp.MustCompile(`"([^"]*)"`)
	for _, m := range stringPattern.FindAllStringSubmatch(positionsText, -1) {
		if len(m) > 1 {
			positions = append(positions, m[1])
		}
	}

	return positions
}

// parseSummary extracts the summary text
func parseSummary(content string) string {
	// Find the summary section by looking for "= Summary" and "= Education"
	startIdx := strings.Index(content, "= Summary")
	if startIdx == -1 {
		return ""
	}

	// Find the end of summary (next section header)
	endIdx := strings.Index(content[startIdx+9:], "=")
	if endIdx == -1 {
		return ""
	}

	summary := content[startIdx+9 : startIdx+9+endIdx]
	summary = strings.TrimSpace(summary)

	return summary
}

// parseSection extracts all resume entries for a given section (Education, Working Experience, Projects)
func parseSection(content string, sectionName string) []ResumeEntry {
	var entries []ResumeEntry

	// Find the section by its header
	sectionMarker := "= " + sectionName
	startIdx := strings.Index(content, sectionMarker)
	if startIdx == -1 {
		return entries
	}

	// Find the next section marker (line starting with "=")
	nextSectionIdx := startIdx + len(sectionMarker)
	sectionLinePattern := regexp.MustCompile(`\n=\s+\w+`)
	match := sectionLinePattern.FindStringIndex(content[nextSectionIdx:])
	var endIdx int
	if match == nil {
		endIdx = len(content)
	} else {
		endIdx = nextSectionIdx + match[0]
	}

	sectionContent := content[nextSectionIdx:endIdx]

	// Parse resume-entry and resume-item blocks
	pos := 0
	for pos < len(sectionContent) {
		// Look for #resume-entry(
		entryMarker := "#resume-entry("
		entryIdx := strings.Index(sectionContent[pos:], entryMarker)
		if entryIdx == -1 {
			break
		}
		entryIdx += pos

		// Find the end of this entry (stops at #resume-item[)
		itemMarker := "#resume-item["
		itemIdx := strings.Index(sectionContent[entryIdx+len(entryMarker):], itemMarker)

		var entryEndIdx int
		if itemIdx != -1 {
			entryEndIdx = entryIdx + len(entryMarker) + itemIdx
		} else {
			entryEndIdx = len(sectionContent)
		}

		// Extract and parse entry
		entryContent := sectionContent[entryIdx+len(entryMarker) : entryEndIdx]
		entry := parseResumeEntry(entryContent)

		// Look for resume-item immediately after this entry
		itemStart := strings.Index(sectionContent[entryEndIdx:], itemMarker)
		if itemStart != -1 {
			itemStart += entryEndIdx
			// Find the end of this item (stops at #resume-entry or line starting with =)
			itemContent := sectionContent[itemStart+len(itemMarker):]

			// Look for next #resume-entry or line starting with =
			nextEntryIdx := strings.Index(itemContent, "#resume-entry(")
			nextLineMarkerIdx := strings.Index(itemContent, "\n=")

			var itemEndOffset int
			if nextLineMarkerIdx != -1 && (nextEntryIdx == -1 || nextLineMarkerIdx < nextEntryIdx) {
				itemEndOffset = nextLineMarkerIdx
			} else if nextEntryIdx != -1 {
				itemEndOffset = nextEntryIdx
			} else {
				itemEndOffset = len(itemContent)
			}

			// Find and remove the last ] before the end
			contentToProcess := itemContent[:itemEndOffset]
			lastBracketIdx := strings.LastIndex(contentToProcess, "]")
			if lastBracketIdx != -1 {
				contentToProcess = contentToProcess[:lastBracketIdx]
			}

			entry.Content = parseResumeItems(contentToProcess)
			pos = itemStart + len(itemMarker) + itemEndOffset
		} else {
			pos = entryEndIdx
		}

		entries = append(entries, entry)
	}

	return entries
}

// parseResumeEntry extracts fields from a resume-entry block
func parseResumeEntry(content string) ResumeEntry {
	entry := ResumeEntry{}

	// Extract title - can be either [bracketed] or "quoted"
	// First try bracketed format
	titleBracketPattern := regexp.MustCompile(`title:\s*\[([^\[]*(?:\[[^\]]*\][^\[]*)*)\]`)
	if titleMatch := titleBracketPattern.FindStringSubmatch(content); len(titleMatch) > 1 {
		title := titleMatch[1]
		entry.Title = strings.TrimSpace(title)
	} else {
		// Try quoted format
		titleQuotePattern := regexp.MustCompile(`title:\s*"([^"]*)"`)
		if titleMatch := titleQuotePattern.FindStringSubmatch(content); len(titleMatch) > 1 {
			entry.Title = titleMatch[1]
		}
	}

	// Extract location
	locationPattern := regexp.MustCompile(`location:\s*"([^"]*)"`)
	if locationMatch := locationPattern.FindStringSubmatch(content); len(locationMatch) > 1 {
		entry.Location = locationMatch[1]
	} else {
		// Try matching github-link or other function patterns like: location: github-link("org/repo")
		githubLinkPattern := regexp.MustCompile(`location:\s*(?:\[?\s*#)?github-link\("([^"]+)"\)`)
		if githubMatch := githubLinkPattern.FindStringSubmatch(content); len(githubMatch) > 1 {
			entry.Location = githubMatch[1]
		} else {
			// Try matching bracketed location like: location: [something]
			bracketPattern := regexp.MustCompile(`location:\s*\[([^\]]+)\]`)
			if bracketMatch := bracketPattern.FindStringSubmatch(content); len(bracketMatch) > 1 {
				entry.Location = strings.TrimSpace(bracketMatch[1])
			}
		}
	}

	// Extract date
	datePattern := regexp.MustCompile(`date:\s*"([^"]*)"`)
	if dateMatch := datePattern.FindStringSubmatch(content); len(dateMatch) > 1 {
		entry.Date = dateMatch[1]
	}

	// Extract description - can be either [bracketed] or "quoted"
	// First try bracketed format
	descBracketPattern := regexp.MustCompile(`description:\s*\[([^\[]*(?:\[[^\]]*\][^\[]*)*)\]`)
	if descMatch := descBracketPattern.FindStringSubmatch(content); len(descMatch) > 1 {
		description := descMatch[1]
		entry.Description = strings.TrimSpace(description)
	} else {
		// Try quoted format
		descQuotePattern := regexp.MustCompile(`description:\s*"([^"]*)"`)
		if descMatch := descQuotePattern.FindStringSubmatch(content); len(descMatch) > 1 {
			entry.Description = descMatch[1]
		}
	}

	return entry
}

// parseResumeItems extracts content from a resume-item block and returns as a single string
func parseResumeItems(content string) string {
	content = strings.TrimSpace(content)

	// Remove backslash line continuations
	content = strings.ReplaceAll(content, `\`, "")

	// Clean up extra whitespace but preserve line structure
	lines := strings.Split(content, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	return strings.Join(cleanedLines, "\n")
}

// parseSkills extracts skill categories and items
func parseSkills(content string) []SkillCategory {
	var skills []SkillCategory

	// Find skills section by string search
	skillsMarker := "= Skills"
	startIdx := strings.Index(content, skillsMarker)
	if startIdx == -1 {
		return skills
	}

	// Find the Interests section or end of content
	endMarker := "= Interests"
	endIdx := strings.Index(content[startIdx:], endMarker)
	if endIdx == -1 {
		endIdx = len(content)
	} else {
		endIdx = startIdx + endIdx
	}

	skillsContent := content[startIdx+len(skillsMarker) : endIdx]

	// Find all resume-skill-item blocks
	skillItemPattern := regexp.MustCompile(`#resume-skill-item\(\s*"([^"]*)"\s*,\s*\(([\s\S]*?)\)\s*\)`)
	skillMatches := skillItemPattern.FindAllStringSubmatch(skillsContent, -1)

	for _, skillMatch := range skillMatches {
		if len(skillMatch) < 3 {
			continue
		}

		category := SkillCategory{
			Name: skillMatch[1],
		}

		// Parse skill items
		skillListContent := skillMatch[2]
		category.Skills = parseSkillItems(skillListContent)

		skills = append(skills, category)
	}

	return skills
}

// parseSkillItems extracts individual skills from a skill category
func parseSkillItems(content string) []SkillItem {
	var items []SkillItem

	// Match both strong("...") and plain "..." patterns
	strongPattern := regexp.MustCompile(`strong\("([^"]*)"\)`)
	plainStringPattern := regexp.MustCompile(`"([^"]*)"`)

	// Create map of strong items for later lookup
	strongMatches := strongPattern.FindAllStringSubmatchIndex(content, -1)

	// Create map of strong items for later lookup
	strongItems := make(map[string]bool)
	for _, match := range strongMatches {
		if len(match) >= 4 {
			skillName := content[match[2]:match[3]]
			strongItems[skillName] = true
		}
	}

	// Now find all string matches
	stringMatches := plainStringPattern.FindAllStringSubmatch(content, -1)
	for _, match := range stringMatches {
		if len(match) > 1 {
			skillName := match[1]
			isStrong := strongItems[skillName]
			items = append(items, SkillItem{
				Name:   skillName,
				Strong: isStrong,
			})
		}
	}

	return items
}

// parseInterests extracts interest items
func parseInterests(content string) []InterestItem {
	var interests []InterestItem

	// Find interests section by looking for "= Interests"
	startIdx := strings.Index(content, "= Interests")
	if startIdx == -1 {
		return interests
	}

	interestsContent := content[startIdx:]

	// Parse resume-skill-item blocks manually to handle multi-line content properly
	pos := 0
	for pos < len(interestsContent) {
		// Look for #resume-skill-item(
		itemMarker := "#resume-skill-item("
		itemIdx := strings.Index(interestsContent[pos:], itemMarker)
		if itemIdx == -1 {
			break
		}
		itemIdx += pos

		// Find the category name (first quoted string)
		categoryStart := strings.Index(interestsContent[itemIdx+len(itemMarker):], "\"")
		if categoryStart == -1 {
			pos = itemIdx + len(itemMarker)
			continue
		}
		categoryStart += itemIdx + len(itemMarker)

		categoryEnd := strings.Index(interestsContent[categoryStart+1:], "\"")
		if categoryEnd == -1 {
			pos = itemIdx + len(itemMarker)
			continue
		}
		categoryEnd += categoryStart + 1

		category := interestsContent[categoryStart+1 : categoryEnd]

		// Find the next #resume-skill-item( or end of file
		nextItemIdx := strings.Index(interestsContent[itemIdx+len(itemMarker):], itemMarker)
		var blockEndIdx int
		if nextItemIdx == -1 {
			blockEndIdx = len(interestsContent)
		} else {
			blockEndIdx = itemIdx + len(itemMarker) + nextItemIdx
		}

		// Extract the block content
		blockContent := interestsContent[itemIdx:blockEndIdx]

		// Find the opening and closing brackets [ ... ]
		descStart := strings.Index(blockContent, "[")
		descEnd := strings.LastIndex(blockContent, "]")
		if descStart != -1 && descEnd != -1 && descEnd > descStart {
			descriptionRaw := blockContent[descStart+1 : descEnd]
			description := strings.TrimSpace(descriptionRaw)

			interests = append(interests, InterestItem{
				Category:    category,
				Description: description,
			})
		}

		pos = blockEndIdx
	}

	return interests
}

// Author represents the author information for the resume

// Author represents the author information for the resume
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Homepage  string `json:"homepage,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Github    string `json:"github,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Birth     string `json:"birth,omitempty"`
	Linkedin  string `json:"linkedin,omitempty"`
}

// ResumeData represents the complete data structure for rendering a resume
type ResumeData struct {
	Author         Author          `json:"author"`
	Positions      []string        `json:"positions"`
	Summary        string          `json:"summary"`
	Education      []ResumeEntry   `json:"education"`
	WorkExperience []ResumeEntry   `json:"work_experience"`
	Projects       []ResumeEntry   `json:"projects"`
	Skills         []SkillCategory `json:"skills"`
	Interests      []InterestItem  `json:"interests"`
}

// getResumeTemplateContent reads resume template from embedded files or filesystem
func getResumeTemplateContent(templatePath string) (string, error) {
	// Try to read from embedded templates first
	if strings.HasSuffix(templatePath, "resume.typ.template") || templatePath == "templates/resume.typ.template" {
		data, err := embeddedResumeTemplates.ReadFile("templates/resume.typ.template")
		if err == nil {
			return string(data), nil
		}
	}

	// Fall back to filesystem
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}
	return string(templateContent), nil
}

// RenderResume renders the resume template with the provided data
func RenderResume(templateContent string, data ResumeData) (string, error) {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"contains": strings.Contains,
	}

	// Parse template
	tmpl, err := template.New("resume").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// RenderResumeFromJSON renders the template using JSON data
func RenderResumeFromJSON(templateContent string, jsonData string) (string, error) {
	var data ResumeData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON data: %w", err)
	}

	return RenderResume(templateContent, data)
}

// RenderAndCompileResume renders the resume template and optionally compiles to PDF
func RenderAndCompileResume(templateContent string, data ResumeData, outputDir string, skipPDF bool) (typstFile string, pdfFile string, err error) {
	// Render template
	result, err := RenderResume(templateContent, data)
	if err != nil {
		return "", "", err
	}

	// Generate output filenames
	typstFileName := fmt.Sprintf("Resume_%s_%s.typ", data.Author.Firstname, data.Author.Lastname)
	pdfFileName := fmt.Sprintf("Resume_%s_%s.pdf", data.Author.Firstname, data.Author.Lastname)

	typstFilePath := filepath.Join(outputDir, typstFileName)
	pdfFilePath := filepath.Join(outputDir, pdfFileName)

	// Write the rendered Typst file
	err = os.WriteFile(typstFilePath, []byte(result), 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to write Typst file: %w", err)
	}

	// Compile to PDF if not skipped
	if !skipPDF {
		cmd := exec.Command("typst", "compile", typstFilePath, pdfFilePath)
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("typst compile failed: %s\nOutput:\n%s\n", err, string(stdoutStderr))
			return "", "", fmt.Errorf("failed to compile Typst file: %w", err)
		}
	}

	return typstFilePath, pdfFilePath, nil
}
