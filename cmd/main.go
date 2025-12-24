package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/FriendsOfTerraform/tfdocs-extras"
	"github.com/terraform-docs/terraform-docs/print"
	"github.com/terraform-docs/terraform-docs/terraform"
)

//go:embed templates/inputs.tmpl
var inputsTmplContent embed.FS

const ExtrasMarkerStart = "<!-- TFDOCS_EXTRAS_START -->"
const ExtrasMarkerEnd = "<!-- TFDOCS_EXTRAS_END -->"

// ReplaceContentBetweenMarkers replaces content between startMarker and endMarker
// Both markers must exist on their own lines
func replaceContentBetweenMarkers(content, startMarker, endMarker, newContent string) string {
	lines := strings.Split(content, "\n")
	var result []string
	insideMarkers := false
	foundStart := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == startMarker {
			result = append(result, line)
			result = append(result, newContent)
			insideMarkers = true
			foundStart = true
			continue
		}

		if trimmedLine == endMarker {
			result = append(result, line)
			insideMarkers = false
			continue
		}

		if !insideMarkers {
			result = append(result, line)
		}
	}

	if !foundStart {
		log.Fatal("Could not find start marker in README.md")
	}

	return strings.Join(result, "\n")
}

func main() {
	modulePath := os.Args[1]

	if modulePath == "" {
		log.Fatal("Module path argument is required")
	}

	config := print.DefaultConfig()
	config.ModuleRoot = modulePath

	module, err := terraform.LoadWithOptions(config)
	if err != nil {
		log.Fatal(err)
	}

	// Read the README.md file
	readmePath := filepath.Join(modulePath, "README.md")
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		log.Fatalf("Failed to read README.md: %v", err)
	}

	// Generate the template output
	tmpl, err := template.New("inputs.tmpl").Funcs(template.FuncMap{
		"indent": func(spaces int, str string) string {
			return "\n" + strings.Repeat("  ", spaces)
		},
	}).ParseFS(inputsTmplContent, "templates/inputs.tmpl")
	if err != nil {
		panic(err)
	}

	templateData := tfdocextras.ParseModuleInputsIntoManifest(module.Inputs)
	var templateOutput bytes.Buffer
	err = tmpl.Execute(&templateOutput, templateData)
	if err != nil {
		panic(err)
	}

	// Replace content between HTML comments
	updatedContent := replaceContentBetweenMarkers(
		string(readmeContent),
		ExtrasMarkerStart,
		ExtrasMarkerEnd,
		templateOutput.String(),
	)

	// Write the updated content back to README.md
	err = os.WriteFile(readmePath, []byte(updatedContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write README.md: %v", err)
	}

	fmt.Println("README.md updated successfully")
}
