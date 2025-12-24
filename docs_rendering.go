package tfdocextras

import (
	"strings"

	"github.com/terraform-docs/terraform-docs/terraform"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TableRowAttribute struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

type RowMetadata struct {
	Attributes    []TableRowAttribute `json:"attributes,omitempty"`
	Enumerations  []string            `json:"enumerations,omitempty"`
	Examples      []TableRowAttribute `json:"examples,omitempty"`
	Links         []TableRowAttribute `json:"links,omitempty"`
	RegexPattern  string              `json:"regex_pattern,omitempty"`
	RegexExamples []string            `json:"regex_examples,omitempty"`
}

type TableRow struct {
	Type         string  `json:"type,omitempty"`
	ComplexType  *string `json:"complex_type,omitempty"`
	Name         string  `json:"name,omitempty"`
	DefaultValue string  `json:"default_value,omitempty"`
	Description  string  `json:"description,omitempty"`
	RowMetadata
}

func (r *TableRow) GetAnchor() string {
	if r.ComplexType == nil {
		return ""
	}

	return strings.ToLower(*r.ComplexType)
}

func (r *TableRow) GetParentType() [2]string {
	if r.ComplexType == nil {
		return [2]string{"", ""}
	}

	values := strings.Split(r.Type, *r.ComplexType)

	if len(values) == 2 {
		return [2]string{values[0], values[1]}
	}

	return [2]string{"", ""}
}

func (r *TableRow) GetMetadata() *RowMetadata {
	return &r.RowMetadata
}

type TableData struct {
	Description string     `json:"description"`
	Rows        []TableRow `json:"rows,omitempty"`
	RowMetadata
}

func (d *TableData) GetMetadata() *RowMetadata {
	return &d.RowMetadata
}

type InputsManifest struct {
	RequiredInputs TableData            `json:"required_inputs,omitempty"`
	OptionalInputs TableData            `json:"optional_inputs,omitempty"`
	NestedInputs   map[string]TableData `json:"nested_inputs,omitempty"`
	ReferenceLinks map[string]string    `json:"reference_links,omitempty"`
}

func newTableData() TableData {
	return TableData{
		Description: "",
		RowMetadata: RowMetadata{
			Attributes:    []TableRowAttribute{},
			Enumerations:  []string{},
			Examples:      []TableRowAttribute{},
			Links:         []TableRowAttribute{},
			RegexPattern:  "",
			RegexExamples: []string{},
		},
		Rows: []TableRow{},
	}
}

func newTableRow(typeStr, name, defaultValue, description string) TableRow {
	return TableRow{
		Type:         typeStr,
		Name:         name,
		DefaultValue: defaultValue,
		Description:  description,
		RowMetadata: RowMetadata{
			Attributes:    []TableRowAttribute{},
			Enumerations:  []string{},
			Examples:      []TableRowAttribute{},
			Links:         []TableRowAttribute{},
			RegexPattern:  "",
			RegexExamples: []string{},
		},
	}
}

func newTemplateData() *InputsManifest {
	return &InputsManifest{
		RequiredInputs: newTableData(),
		OptionalInputs: newTableData(),
		NestedInputs:   make(map[string]TableData),
		ReferenceLinks: make(map[string]string),
	}
}

func processDirectives(directives []DocDirective, manifest *InputsManifest, data *TableData, row *TableRow) {
	var metadata *RowMetadata

	if data != nil {
		metadata = data.GetMetadata()
	} else if row != nil {
		metadata = row.GetMetadata()
	}

	if metadata == nil {
		return
	}

	for _, attr := range directives {
		if (attr.Parsed.Flags & IsInvalid) != 0 {
			continue
		}

		switch attr.Parsed.Type {
		case DirEnum:
			metadata.Enumerations = append(metadata.Enumerations, attr.Parsed.Args...)
		case DirExample:
			metadata.Examples = append(metadata.Examples, TableRowAttribute{
				Name:    attr.Parsed.Args[0],
				Content: getArgOrDefault(attr.Parsed.Args, 1),
			})
		case DirLink:
			if (attr.Parsed.Flags & IsReferenceLink) != 0 {
				manifest.ReferenceLinks[attr.Parsed.Args[0]] = attr.Parsed.Args[1]
			} else if (attr.Parsed.Flags & IsNamedLink) != 0 {
				metadata.Links = append(metadata.Links, TableRowAttribute{
					Name:    attr.Parsed.Args[0],
					Content: getArgOrDefault(attr.Parsed.Args, 1),
				})
			}
		case DirRegex:
			if len(attr.Parsed.Args) >= 1 {
				metadata.RegexPattern = attr.Parsed.Args[0]
				metadata.RegexExamples = append(metadata.RegexExamples, attr.Parsed.Args[1:]...)
			}
		default:
			caser := cases.Title(language.English)
			metadata.Attributes = append(metadata.Attributes, TableRowAttribute{
				Name:    caser.String(attr.Name),
				Content: attr.RawContent,
			})
		}
	}
}

func getArgOrDefault(args []string, index int) string {
	if len(args) > index {
		return args[index]
	}

	return ""
}

func recordNested(group ObjectField, manifest *InputsManifest) {
	if group.NestedDataType == nil {
		return
	}

	if group.Fields != nil && len(group.Fields) > 0 {
		data := newTableData()
		data.Description = strings.Join(group.Documentation.Content, "\n")

		processDirectives(group.Documentation.Directives, manifest, &data, nil)

		for _, field := range group.Fields {
			defaultValue := ""

			if field.DefaultValue != nil {
				defaultValue = *field.DefaultValue
			}

			row := newTableRow(field.DataTypeStr, field.Name, defaultValue, strings.Join(field.Documentation.Content, "\n"))

			if field.NestedDataType != nil {
				row.ComplexType = field.NestedDataType
			}

			processDirectives(field.Documentation.Directives, manifest, nil, &row)

			data.Rows = append(data.Rows, row)
		}

		manifest.NestedInputs[*group.NestedDataType] = data
	}

	for _, field := range group.Fields {
		recordNested(field, manifest)
	}
}

func ParseModuleInputsIntoManifest(inputs []*terraform.Input) *InputsManifest {
	templateData := newTemplateData()

	for _, input := range inputs {
		var extras ObjectGroup
		if input.Type != "" {
			documented, astErr := ParseIntoDocumentedStruct(string(input.Type), input.Name)

			if astErr == nil && documented != nil {
				extras = *documented
			}
		}

		docBlk := parseStringIntoDocBlock(string(input.Description))
		tableRow := newTableRow(string(input.Type), input.Name, input.GetValue(), strings.Join(docBlk.Content, "\n"))

		processDirectives(docBlk.Directives, templateData, nil, &tableRow)

		if extras.ObjectField.NestedDataType != nil {
			tableRow.Type = extras.ObjectField.DataTypeStr
			tableRow.ComplexType = extras.ObjectField.NestedDataType
		}

		if input.Required {
			templateData.RequiredInputs.Rows = append(templateData.RequiredInputs.Rows, tableRow)
		} else {
			templateData.OptionalInputs.Rows = append(templateData.OptionalInputs.Rows, tableRow)
		}

		recordNested(extras.ObjectField, templateData)

		for _, field := range extras.ObjectField.Fields {
			recordNested(field, templateData)
		}
	}

	return templateData
}
