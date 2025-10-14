package main

import (
	"embed"
	"log"
	"os"
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

type TableData struct {
	Rows []TableRow `json:"rows,omitempty"`
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
		Rows: []TableRow{},
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

		tableRow := TableRow{
			Type:         string(input.Type),
			Name:         input.Name,
			DefaultValue: input.GetValue(),
			Description:  string(input.Description),
			Attributes:   []TableRowAttribute{},
		}

		if extras.ObjectField.NestedDataType != nil {
			tableRow.Type = extras.ObjectField.DataTypeStr
			tableRow.ComplexType = extras.ObjectField.NestedDataType
		}

		if extras.Optional {
			templateData.OptionalInputs.Rows = append(templateData.OptionalInputs.Rows, tableRow)
		} else {
			templateData.RequiredInputs.Rows = append(templateData.RequiredInputs.Rows, tableRow)
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

		for _, field := range group.Fields {
			row := TableRow{
				Type:         field.DataTypeStr,
				Name:         field.Name,
				DefaultValue: "",
				Description:  strings.Join(field.Documentation.Content, "\n"),
				Attributes:   []TableRowAttribute{},
			}

			if field.NestedDataType != nil {
				row.Type = *field.NestedDataType
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
	config := print.DefaultConfig()
	config.ModuleRoot = os.Args[1]

	module, err := terraform.LoadWithOptions(config)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFS(inputsTmplContent, "templates/inputs.tmpl")
	if err != nil {
		panic(err)
	}

	templateData := parseModuleInputs(module.Inputs)
	err = tmpl.Execute(os.Stdout, templateData)
	if err != nil {
		panic(err)
	}

	// Print as JSON for easy consumption
	//astModule, _ := json.MarshalIndent(templateData, "", "  ")
	//fmt.Printf("%s\n", astModule)
}
