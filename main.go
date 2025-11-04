package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

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

		// Get template content
		templateContent, err := getTemplateContent(templatePath)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(RenderResponse{
				Success: false,
				Error:   fmt.Sprintf("Failed to read template: %v", err),
			})
			return
		}

		// Render and compile
		_, pdfFile, err := RenderAndCompileCoverLetter(templateContent, data, outputDir, skipPDF)
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

// handleRenderResume handles the /render-resume POST endpoint
func handleRenderResume(templatePath string, outputDir string, skipPDF bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Method not allowed. Please use POST.",
			})
			return
		}

		var data ResumeData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Invalid JSON: %v", err),
			})
			return
		}

		// Get template content
		templateContent, err := getResumeTemplateContent(templatePath)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to read template: %v", err),
			})
			return
		}

		// Render and compile
		_, pdfFile, err := RenderAndCompileResume(templateContent, data, outputDir, skipPDF)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Rendering failed: %v", err),
			})
			return
		}

		// Return the PDF file
		pdfContent, err := os.ReadFile(pdfFile)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to read PDF file: %v", err),
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
	var data CoverLetterData
	var err error

	// Get template content
	templateContent, err := getTemplateContent(templatePath)
	if err != nil {
		log.Fatalf("Failed to read template: %v", err)
	}

	// Parse data based on input
	if jsonFile != "" {
		jsonContent, err := os.ReadFile(jsonFile)
		if err != nil {
			log.Fatalf("Failed to read JSON file: %v", err)
		}
		_, err = RenderCoverLetterFromJSON(templateContent, string(jsonContent))
		if err != nil {
			log.Fatalf("Error rendering template: %v", err)
		}
		err = json.Unmarshal(jsonContent, &data)
		if err != nil {
			log.Fatalf("Error parsing JSON data: %v", err)
		}
	} else if jsonString != "" {
		_, err = RenderCoverLetterFromJSON(templateContent, jsonString)
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

	// Render and compile
	typstFilePath, pdfFilePath, err := RenderAndCompileCoverLetter(templateContent, data, outputDir, skipPDF)
	if err != nil {
		log.Fatalf("Error rendering/compiling: %v", err)
	}

	fmt.Printf("Rendered Typst file saved to: %s\n", typstFilePath)

	// Print PDF path if compiled
	if !skipPDF {
		fmt.Printf("PDF compiled successfully: %s\n", pdfFilePath)
	}
}

// handleParseResume handles resume parsing requests
func handleParseResume() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Method not allowed. Please use POST.",
			})
			return
		}

		// Read the Typst file path from query parameter or request body
		var req struct {
			FilePath string `json:"file_path"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			req.FilePath = r.URL.Query().Get("file_path")
		}

		if req.FilePath == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Missing required parameter: file_path",
			})
			return
		}

		// Read the Typst file
		content, err := os.ReadFile(req.FilePath)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to read file: %v", err),
			})
			return
		}

		// Parse the resume
		resume, err := ParseResumeTypst(string(content))
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("Failed to parse resume: %v", err),
			})
			return
		}

		// Return the parsed resume as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"resume":  resume,
		})
	}
}

// runHTTPServer runs the program as an HTTP server
func runHTTPServer(templatePath, outputDir, port string, skipPDF bool) {
	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Resume template path (default to templates/resume.typ.template)
	resumeTemplatePath := "templates/resume.typ.template"

	// Register HTTP handlers
	http.HandleFunc("/render", handleRender(templatePath, outputDir, skipPDF))
	http.HandleFunc("/render-resume", handleRenderResume(resumeTemplatePath, outputDir, skipPDF))
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/parse-resume", handleParseResume())

	// Start server
	addr := ":" + port
	log.Printf("Starting HTTP server on http://localhost:%s", port)
	log.Printf("Endpoints:")
	log.Printf("  POST /render - Render cover letter from JSON")
	log.Printf("  POST /render-resume - Render resume from JSON")
	log.Printf("  POST /parse-resume - Parse resume from Typst file")
	log.Printf("  GET /health - Health check")
	log.Printf("  Template: %s", templatePath)
	log.Printf("  Resume Template: %s", resumeTemplatePath)
	log.Printf("  Output directory: %s", outputDir)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
