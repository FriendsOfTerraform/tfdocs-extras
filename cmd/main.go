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
	updatedContent := tfdocextras.ReplaceContentBetweenMarkers(
		string(readmeContent),
		tfdocextras.ExtrasMarkerStart,
		tfdocextras.ExtrasMarkerEnd,
		templateOutput.String(),
	)

	// Write the updated content back to README.md
	err = os.WriteFile(readmePath, []byte(updatedContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write README.md: %v", err)
	}

	fmt.Println("README.md updated successfully")
}
