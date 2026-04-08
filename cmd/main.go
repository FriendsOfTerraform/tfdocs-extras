package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/FriendsOfTerraform/tfdocs-extras"
	"github.com/terraform-docs/terraform-config-inspect/tfconfig"
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
	jsonOutput := flag.Bool("json", false, "output JSON representation of template data to stdout")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("Module path argument is required")
	}

	modulePath := flag.Arg(0)

	config := print.DefaultConfig()
	config.ModuleRoot = modulePath

	module, err := terraform.LoadWithOptions(config)
	if err != nil {
		log.Fatal(err)
	}

	// tfdocs has a bug where it doesn't set the sensitivity flag on outputs,
	// even if it's set in the raw module config. This is a workaround to read
	// the output sensitivity flags by reloading the module with
	// terraform-config-inspect.
	//
	// see: https://github.com/terraform-docs/terraform-docs/issues/798
	if err := outputSensitivityFallback(modulePath, module.Outputs); err != nil {
		log.Printf("Warning: could not enrich output sensitivity flags: %v", err)
	}

	templateData := tfdocextras.ParseModuleArgsIntoManifest(module.Inputs, module.Outputs)

	// If -json flag is set, output JSON and exit
	if *jsonOutput {
		jsonBytes, err := json.MarshalIndent(templateData, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal JSON: %v", err)
		}
		fmt.Println(string(jsonBytes))
		return
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

func outputSensitivityFallback(modulePath string, outputs []*terraform.Output) error {
	tfModule, diagnostics := tfconfig.LoadModule(modulePath)
	if diagnostics != nil && diagnostics.HasErrors() {
		return diagnostics
	}

	for _, output := range outputs {
		if rawOutput, ok := tfModule.Outputs[output.Name]; ok {
			output.Sensitive = rawOutput.Sensitive
		}
	}

	return nil
}
