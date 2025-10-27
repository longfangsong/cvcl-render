package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/coverletter.typ.template
var embeddedTemplates embed.FS

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

// RenderCoverLetter renders the cover letter template with the provided data
func RenderCoverLetter(templatePath string, data CoverLetterData) (string, error) {
	// Read template file (from embedded or filesystem)
	templateContent, err := getTemplateContent(templatePath)
	if err != nil {
		return "", err
	}

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
func RenderCoverLetterFromJSON(templatePath string, jsonData string) (string, error) {
	var data CoverLetterData
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON data: %w", err)
	}

	return RenderCoverLetter(templatePath, data)
}

// RenderCoverLetterFromJSONFile renders the template using a JSON file
func RenderCoverLetterFromJSONFile(templatePath string, jsonFilePath string) (string, error) {
	jsonContent, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read JSON file: %w", err)
	}

	return RenderCoverLetterFromJSON(templatePath, string(jsonContent))
}

// CompileTypstToPDF compiles a Typst file to PDF using the typst compile command
func CompileTypstToPDF(typstFilePath string, pdfOutputPath string) error {
	cmd := exec.Command("typst", "compile", typstFilePath, pdfOutputPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to compile Typst file: %w", err)
	}
	return nil
}

// RenderAndCompile renders the template and optionally compiles to PDF
func RenderAndCompile(templatePath string, data CoverLetterData, outputDir string, skipPDF bool) (typstFile string, pdfFile string, err error) {
	// Render template
	result, err := RenderCoverLetter(templatePath, data)
	if err != nil {
		return "", "", err
	}

	// Generate output filenames
	pos := strings.ReplaceAll(data.Position, " ", "_")
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

// handleRender handles the /render POST endpoint
func handleRender(templatePath string, outputDir string, skipPDF bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(RenderResponse{
				Success: false,
				Error:   "Method not allowed. Please use POST.",
			})
			return
		}

		var req RenderRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(RenderResponse{
				Success: false,
				Error:   fmt.Sprintf("Invalid JSON: %v", err),
			})
			return
		}

		data := CoverLetterData(req)

		// Render and compile
		_, pdfFile, err := RenderAndCompile(templatePath, data, outputDir, skipPDF)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(RenderResponse{
				Success: false,
				Error:   fmt.Sprintf("Rendering failed: %v", err),
			})
			return
		}

		// Return the PDF file
		pdfContent, err := os.ReadFile(pdfFile)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(RenderResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to read PDF file: %v", err),
			})
			return
		}

		// Get filename for Content-Disposition header
		fileName := filepath.Base(pdfFile)

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfContent)))
		w.WriteHeader(http.StatusOK)
		w.Write(pdfContent)
	}
}

// handleHealth handles the /health GET endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	// Define flags
	templatePath := flag.String("template", "templates/coverletter.typ.template", "Path to the Typst template file")
	outputDir := flag.String("output-dir", ".", "Output directory for the generated files")
	port := flag.String("port", "8080", "Port to listen on for HTTP server")
	cliMode := flag.Bool("cli", false, "Run in CLI mode instead of HTTP server")
	jsonFile := flag.String("data", "", "Path to JSON file containing the data (CLI mode only)")
	jsonString := flag.String("json", "", "JSON string containing the data (CLI mode only)")
	skipPDF := flag.Bool("skip-pdf", false, "Skip PDF compilation and only output the rendered Typst file")

	flag.Parse()

	// CLI mode
	if *cliMode {
		runCLI(*templatePath, *outputDir, *jsonFile, *jsonString, *skipPDF)
	} else {
		// HTTP server mode (default)
		runHTTPServer(*templatePath, *outputDir, *port, *skipPDF)
	}
}

// runCLI runs the program in CLI mode
func runCLI(templatePath, outputDir, jsonFile, jsonString string, skipPDF bool) {
	var result string
	var data CoverLetterData
	var err error

	// Parse data based on input
	if jsonFile != "" {
		jsonContent, err := os.ReadFile(jsonFile)
		if err != nil {
			log.Fatalf("Failed to read JSON file: %v", err)
		}
		result, err = RenderCoverLetterFromJSON(templatePath, string(jsonContent))
		if err != nil {
			log.Fatalf("Error rendering template: %v", err)
		}
		err = json.Unmarshal(jsonContent, &data)
		if err != nil {
			log.Fatalf("Error parsing JSON data: %v", err)
		}
	} else if jsonString != "" {
		result, err = RenderCoverLetterFromJSON(templatePath, jsonString)
		if err != nil {
			log.Fatalf("Error rendering template: %v", err)
		}
		err = json.Unmarshal([]byte(jsonString), &data)
		if err != nil {
			log.Fatalf("Error parsing JSON data: %v", err)
		}
	} else {
		log.Fatal("Please provide either -data or -json flag with the input data")
	}

	// Generate output filenames
	pos := strings.ReplaceAll(data.Position, " ", "_")
	typstFileName := fmt.Sprintf("Cover_Letter_%s_%s.typ", data.FirstName, pos)
	pdfFileName := fmt.Sprintf("Cover_Letter_%s_%s.pdf", data.FirstName, pos)

	typstFilePath := filepath.Join(outputDir, typstFileName)
	pdfFilePath := filepath.Join(outputDir, pdfFileName)

	// Write the rendered Typst file
	err = os.WriteFile(typstFilePath, []byte(result), 0644)
	if err != nil {
		log.Fatalf("Failed to write Typst file: %v", err)
	}
	fmt.Printf("Rendered Typst file saved to: %s\n", typstFilePath)

	// Compile to PDF if not skipped
	if !skipPDF {
		err = CompileTypstToPDF(typstFilePath, pdfFilePath)
		if err != nil {
			log.Fatalf("Error compiling PDF: %v", err)
		}
		fmt.Printf("PDF compiled successfully: %s\n", pdfFilePath)
	}
}

// runHTTPServer runs the program as an HTTP server
func runHTTPServer(templatePath, outputDir, port string, skipPDF bool) {
	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Register HTTP handlers
	http.HandleFunc("/render", handleRender(templatePath, outputDir, skipPDF))
	http.HandleFunc("/health", handleHealth)

	// Start server
	addr := ":" + port
	log.Printf("Starting HTTP server on http://localhost:%s", port)
	log.Printf("Endpoints:")
	log.Printf("  POST /render - Render cover letter from JSON")
	log.Printf("  GET /health - Health check")
	log.Printf("  Template: %s", templatePath)
	log.Printf("  Output directory: %s", outputDir)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
