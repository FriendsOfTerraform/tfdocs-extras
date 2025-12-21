package tfdocextras

import (
	"log"
	"strings"

	"github.com/terraform-docs/terraform-docs/terraform"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const ExtrasMarkerStart = "<!-- TFDOCS_EXTRAS_START -->"
const ExtrasMarkerEnd = "<!-- TFDOCS_EXTRAS_END -->"

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
	Examples     []TableRowAttribute `json:"examples,omitempty"`
	Links        []TableRowAttribute `json:"links,omitempty"`
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
	Examples    []TableRowAttribute `json:"examples,omitempty"`
	Links       []TableRowAttribute `json:"links,omitempty"`
	Rows        []TableRow          `json:"rows,omitempty"`
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
		Attributes:  []TableRowAttribute{},
		Rows:        []TableRow{},
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

// processDirectives applies directive logic to either a TableData or TableRow
func processDirectives(directives []DocDirective, manifest *InputsManifest, data *TableData, row *TableRow) {
	for _, attr := range directives {
		if (attr.Parsed.Flags & IsInvalid) != 0 {
			continue
		}

		rowAttr := TableRowAttribute{
			Name:    attr.Parsed.Args[0],
			Content: attr.Parsed.Args[1],
		}

		switch attr.Parsed.Type {
		case DirExample:
			if data != nil {
				data.Examples = append(data.Examples, rowAttr)
			} else if row != nil {
				row.Examples = append(row.Examples, rowAttr)
			}
		case DirLink:
			if (attr.Parsed.Flags & IsReferenceLink) != 0 {
				manifest.ReferenceLinks[attr.Parsed.Args[0]] = attr.Parsed.Args[1]
			} else if (attr.Parsed.Flags & IsNamedLink) != 0 {
				if data != nil {
					data.Links = append(data.Links, rowAttr)
				} else if row != nil {
					row.Links = append(row.Links, rowAttr)
				}
			}
		default:
			caser := cases.Title(language.English)

			attribute := TableRowAttribute{
				Name:    caser.String(attr.Name),
				Content: attr.RawContent,
			}

			if data != nil {
				data.Attributes = append(data.Attributes, attribute)
			} else if row != nil {
				row.Attributes = append(row.Attributes, attribute)
			}
		}
	}
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
			row := TableRow{
				Type:         field.DataTypeStr,
				Name:         field.Name,
				DefaultValue: "",
				Description:  strings.Join(field.Documentation.Content, "\n"),
				Attributes:   []TableRowAttribute{},
				Examples:     []TableRowAttribute{},
				Links:        []TableRowAttribute{},
			}

			if field.NestedDataType != nil {
				row.ComplexType = field.NestedDataType
			}

			if field.DefaultValue != nil {
				row.DefaultValue = *field.DefaultValue
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

		tableRow := TableRow{
			Type:         string(input.Type),
			Name:         input.Name,
			DefaultValue: input.GetValue(),
			Description:  strings.Join(docBlk.Content, "\n"),
			Attributes:   []TableRowAttribute{},
		}

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

// ReplaceContentBetweenMarkers replaces content between startMarker and endMarker
// Both markers must exist on their own lines
func ReplaceContentBetweenMarkers(content, startMarker, endMarker, newContent string) string {
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
