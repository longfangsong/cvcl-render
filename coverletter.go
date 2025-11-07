package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/coverletter.typ.template
var embeddedTemplates embed.FS

// getTemplateContent reads template from embedded files or filesystem
func getTemplateContent(templatePath string) (string, error) {
	// Try to read from embedded templates first
	if strings.HasSuffix(templatePath, "coverletter.typ.template") || templatePath == "templates/coverletter.typ.template" {
		data, err := embeddedTemplates.ReadFile("templates/coverletter.typ.template")
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

// CoverLetterData represents the data to be rendered in the template
type CoverLetterData struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Homepage   string `json:"homepage,omitempty"`
	Phone      string `json:"phone,omitempty"`
	GitHub     string `json:"github,omitempty"`
	LinkedIn   string `json:"linkedin,omitempty"`
	Position   string `json:"position"`
	Addressee  string `json:"addressee"`
	Opening    string `json:"opening"`
	AboutMe    string `json:"about_me"`
	WhyMe      string `json:"why_me"`
	WhyCompany string `json:"why_company"`
}

// RenderRequest is the JSON request body for the render endpoint
type RenderRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Homepage   string `json:"homepage,omitempty"`
	Phone      string `json:"phone,omitempty"`
	GitHub     string `json:"github,omitempty"`
	LinkedIn   string `json:"linkedin,omitempty"`
	Position   string `json:"position"`
	Addressee  string `json:"addressee"`
	Opening    string `json:"opening"`
	AboutMe    string `json:"about_me"`
	WhyMe      string `json:"why_me"`
	WhyCompany string `json:"why_company"`
}

// RenderResponse is the JSON response for the render endpoint
type RenderResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	TypstFile string `json:"typst_file,omitempty"`
	PDFFile   string `json:"pdf_file,omitempty"`
	Error     string `json:"error,omitempty"`
}

// RenderCoverLetter renders the cover letter template with the provided data
func RenderCoverLetter(templateContent string, data CoverLetterData) (string, error) {
	// Parse template
	tmpl, err := template.New("coverletter").Parse(templateContent)
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

// RenderCoverLetterFromJSON renders the template using JSON data
func RenderCoverLetterFromJSON(templateContent string, jsonData string) (string, error) {
	var data CoverLetterData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON data: %w", err)
	}

	return RenderCoverLetter(templateContent, data)
}

// RenderCoverLetterFromJSONFile renders the template using a JSON file
func RenderCoverLetterFromJSONFile(templateContent string, jsonFilePath string) (string, error) {
	jsonContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read JSON file: %w", err)
	}

	return RenderCoverLetterFromJSON(templateContent, string(jsonContent))
}

// CompileTypstToPDF compiles a Typst file to PDF using the typst compile command
func CompileTypstToPDF(typstFilePath string, pdfOutputPath string) error {
	cmd := exec.Command("typst", "compile", typstFilePath, pdfOutputPath)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("typst compile failed: %s\nOutput:\n%s\n", err, string(stdoutStderr))
		return fmt.Errorf("failed to compile Typst file: %w", err)
	}
	return nil
}

// RenderAndCompileCoverLetter renders the template and optionally compiles to PDF
func RenderAndCompileCoverLetter(templateContent string, data CoverLetterData, outputDir string, skipPDF bool) (typstFile string, pdfFile string, err error) {
	// Render template
	result, err := RenderCoverLetter(templateContent, data)
	if err != nil {
		return "", "", err
	}

	// Generate output filenames
	pos := data.Position
	pos = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, pos)
	typstFileName := fmt.Sprintf("Cover_Letter_%s_%s.typ", data.FirstName, pos)
	pdfFileName := fmt.Sprintf("Cover_Letter_%s_%s.pdf", data.FirstName, pos)

	typstFilePath := filepath.Join(outputDir, typstFileName)
	pdfFilePath := filepath.Join(outputDir, pdfFileName)

	// Write the rendered Typst file
	err = os.WriteFile(typstFilePath, []byte(result), 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to write Typst file: %w", err)
	}

	// Compile to PDF if not skipped
	if !skipPDF {
		err = CompileTypstToPDF(typstFilePath, pdfFilePath)
		if err != nil {
			return "", "", err
		}
	}

	return typstFilePath, pdfFilePath, nil
}
