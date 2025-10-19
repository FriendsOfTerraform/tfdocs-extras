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

type TableRowAttribute struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

type TableRow struct {
	Type         string              `json:"type,omitempty"`
	ComplexType  *string             `json:"complex_type,omitempty"`
	Name         string              `json:"name,omitempty"`
	DefaultValue string              `json:"default_value,omitempty"`
	Description  string              `json:"description,omitempty"`
	Attributes   []TableRowAttribute `json:"attributes,omitempty"`
}

func (row *TableRow) GetAnchor() string {
	if row.ComplexType == nil {
		return ""
	}

	return strings.ToLower(*row.ComplexType)
}

func (row *TableRow) GetParentType() [2]string {
	if row.ComplexType == nil {
		return [2]string{"", ""}
	}

	values := strings.Split(row.Type, *row.ComplexType)

	if len(values) == 2 {
		return [2]string{values[0], values[1]}
	}

	return [2]string{"", ""}
}

type TableData struct {
	Description string              `json:"description"`
	Attributes  []TableRowAttribute `json:"attributes"`
	Rows        []TableRow          `json:"rows,omitempty"`
}

type TemplateData struct {
	RequiredInputs TableData            `json:"required_inputs,omitempty"`
	OptionalInputs TableData            `json:"optional_inputs,omitempty"`
	NestedInputs   map[string]TableData `json:"nested_inputs,omitempty"`
}

//go:embed templates/inputs.tmpl
var inputsTmplContent embed.FS

func newTableData() TableData {
	return TableData{
		Description: "",
		Attributes:  []TableRowAttribute{},
		Rows:        []TableRow{},
	}
}

func newTemplateData() *TemplateData {
	return &TemplateData{
		RequiredInputs: newTableData(),
		OptionalInputs: newTableData(),
		NestedInputs:   make(map[string]TableData),
	}
}

func parseModuleInputs(inputs []*terraform.Input) *TemplateData {
	templateData := newTemplateData()

	for _, input := range inputs {
		var extras tfdocextras.ObjectGroup
		if input.Type != "" {
			documented, astErr := tfdocextras.ParseIntoDocumentedStruct(string(input.Type), input.Name)

			if astErr == nil && documented != nil {
				extras = *documented
			}
		}

		docBlk := tfdocextras.ParseStringIntoDocBlock(string(input.Description))

		tableRow := TableRow{
			Type:         string(input.Type),
			Name:         input.Name,
			DefaultValue: input.GetValue(),
			Description:  strings.Join(docBlk.Content, "\n"),
			Attributes:   []TableRowAttribute{},
		}

		for _, attr := range docBlk.Directives {
			attribute := TableRowAttribute{
				Name:    attr.Name,
				Content: attr.Content,
			}

			tableRow.Attributes = append(tableRow.Attributes, attribute)
		}

		if extras.ObjectField.NestedDataType != nil {
			tableRow.Type = extras.ObjectField.DataTypeStr
			tableRow.ComplexType = extras.ObjectField.NestedDataType
		}

		if input.Required {
			templateData.RequiredInputs.Rows = append(templateData.RequiredInputs.Rows, tableRow)
		} else {
			templateData.OptionalInputs.Rows = append(templateData.OptionalInputs.Rows, tableRow)
		}

		recordNested(extras.ObjectField, templateData.NestedInputs)

		for _, field := range extras.ObjectField.Fields {
			recordNested(field, templateData.NestedInputs)
		}
	}

	return templateData
}

func recordNested(group tfdocextras.ObjectField, record map[string]TableData) {
	if group.NestedDataType == nil {
		return
	}

	if group.Fields != nil && len(group.Fields) > 0 {
		data := newTableData()
		data.Description = strings.Join(group.Documentation.Content, "\n")

		for _, attr := range group.Documentation.Directives {
			attribute := TableRowAttribute{
				Name:    attr.Name,
				Content: attr.Content,
			}

			data.Attributes = append(data.Attributes, attribute)
		}

		for _, field := range group.Fields {
			row := TableRow{
				Type:         field.DataTypeStr,
				Name:         field.Name,
				DefaultValue: "",
				Description:  strings.Join(field.Documentation.Content, "\n"),
				Attributes:   []TableRowAttribute{},
			}

			if field.NestedDataType != nil {
				row.ComplexType = field.NestedDataType
			}

			if field.DefaultValue != nil {
				row.DefaultValue = *field.DefaultValue
			}

			for _, attr := range field.Documentation.Directives {
				attribute := TableRowAttribute{
					Name:    attr.Name,
					Content: attr.Content,
				}

				row.Attributes = append(row.Attributes, attribute)
			}

			data.Rows = append(data.Rows, row)
		}

		record[*group.NestedDataType] = data
	}

	for _, field := range group.Fields {
		recordNested(field, record)
	}
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

	templateData := parseModuleInputs(module.Inputs)
	var templateOutput bytes.Buffer
	err = tmpl.Execute(&templateOutput, templateData)
	if err != nil {
		panic(err)
	}

	// Replace content between HTML comments
	updatedContent := replaceContentBetweenMarkers(
		string(readmeContent),
		"<!-- TFDOCS_EXTRAS_START -->",
		"<!-- TFDOCS_EXTRAS_END -->",
		templateOutput.String(),
	)

	// Write the updated content back to README.md
	err = os.WriteFile(readmePath, []byte(updatedContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write README.md: %v", err)
	}

	fmt.Println("README.md updated successfully")
}

// replaceContentBetweenMarkers replaces content between startMarker and endMarker
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
