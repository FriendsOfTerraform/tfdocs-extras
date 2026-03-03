package tfdocextras

import (
	"strings"

	"github.com/terraform-docs/terraform-docs/terraform"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MetadataAttribute struct {
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
}

type ArgumentMetadata struct {
	Attributes    []MetadataAttribute `json:"attributes,omitempty"`
	Enumerations  []string            `json:"enumerations,omitempty"`
	Examples      []MetadataAttribute `json:"examples,omitempty"`
	Links         []MetadataAttribute `json:"links,omitempty"`
	RegexPattern  string              `json:"regex_pattern,omitempty"`
	RegexExamples []string            `json:"regex_examples,omitempty"`
}

type Argument struct {
	Type         string  `json:"type,omitempty"`
	ComplexType  *string `json:"complex_type,omitempty"`
	Name         string  `json:"name,omitempty"`
	DefaultValue *string `json:"default_value,omitempty"`
	Required     bool    `json:"required,omitempty"`
	Sensitive    bool    `json:"sensitive,omitempty"`
	Description  string  `json:"description,omitempty"`
	ArgumentMetadata
}

func (r *Argument) GetAnchor() string {
	if r.ComplexType == nil {
		return ""
	}

	return strings.ToLower(*r.ComplexType)
}

func (r *Argument) GetParentType() [2]string {
	if r.ComplexType == nil {
		return [2]string{"", ""}
	}

	values := strings.Split(r.Type, *r.ComplexType)

	if len(values) == 2 {
		return [2]string{values[0], values[1]}
	}

	return [2]string{"", ""}
}

func (r *Argument) GetMetadata() *ArgumentMetadata {
	return &r.ArgumentMetadata
}

func (r *Argument) HasDefaultValue() bool {
	if r.Sensitive || r.Required {
		return false
	}

	return r.DefaultValue != nil
}

type ArgGroupType int

const (
	Anonymous ArgGroupType = iota
	RequiredInput
	OptionalInput
	Output
)

type ArgumentGroup struct {
	Type        ArgGroupType `json:"-"`
	Description string       `json:"description"`
	Rows        []Argument   `json:"rows,omitempty"`
	ArgumentMetadata
}

func (d *ArgumentGroup) GetMetadata() *ArgumentMetadata {
	return &d.ArgumentMetadata
}

type ModuleManifest struct {
	RequiredInputs ArgumentGroup            `json:"required_inputs,omitempty"`
	OptionalInputs ArgumentGroup            `json:"optional_inputs,omitempty"`
	Outputs        ArgumentGroup            `json:"outputs,omitempty"`
	Objects        map[string]ArgumentGroup `json:"objects,omitempty"`
	ReferenceLinks map[string]string        `json:"reference_links,omitempty"`
}

func newArgumentGroup(grpType ArgGroupType) ArgumentGroup {
	return ArgumentGroup{
		Type:        grpType,
		Description: "",
		ArgumentMetadata: ArgumentMetadata{
			Attributes:    []MetadataAttribute{},
			Enumerations:  []string{},
			Examples:      []MetadataAttribute{},
			Links:         []MetadataAttribute{},
			RegexPattern:  "",
			RegexExamples: []string{},
		},
		Rows: []Argument{},
	}
}

func newArgument(typeStr, name, description string) Argument {
	return Argument{
		Type:         typeStr,
		Name:         name,
		DefaultValue: nil,
		Required:     false,
		Sensitive:    false,
		Description:  description,
		ArgumentMetadata: ArgumentMetadata{
			Attributes:    []MetadataAttribute{},
			Enumerations:  []string{},
			Examples:      []MetadataAttribute{},
			Links:         []MetadataAttribute{},
			RegexPattern:  "",
			RegexExamples: []string{},
		},
	}
}

func newModuleManifest() *ModuleManifest {
	return &ModuleManifest{
		RequiredInputs: newArgumentGroup(RequiredInput),
		OptionalInputs: newArgumentGroup(OptionalInput),
		Objects:        make(map[string]ArgumentGroup),
		ReferenceLinks: make(map[string]string),
		Outputs:        newArgumentGroup(Output),
	}
}

func processDirectives(directives []DocDirective, manifest *ModuleManifest, data *ArgumentGroup, row *Argument) {
	var metadata *ArgumentMetadata

	if data != nil {
		metadata = data.GetMetadata()
	} else if row != nil {
		metadata = row.GetMetadata()
	}

	if metadata == nil {
		return
	}

	for _, attr := range directives {
		if (attr.Parsed.Flags&IsInvalid) != 0 || attr.Parsed.Type == DirType {
			continue
		}

		switch attr.Parsed.Type {
		case DirEnum:
			metadata.Enumerations = append(metadata.Enumerations, attr.Parsed.Args...)
		case DirExample:
			metadata.Examples = append(metadata.Examples, MetadataAttribute{
				Name:    attr.Parsed.Args[0],
				Content: getArgOrDefault(attr.Parsed.Args, 1),
			})
		case DirLink:
			if (attr.Parsed.Flags & IsReferenceLink) != 0 {
				manifest.ReferenceLinks[attr.Parsed.Args[0]] = attr.Parsed.Args[1]
			} else if (attr.Parsed.Flags & IsNamedLink) != 0 {
				metadata.Links = append(metadata.Links, MetadataAttribute{
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
			metadata.Attributes = append(metadata.Attributes, MetadataAttribute{
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

func recordNested(group StructProperty, manifest *ModuleManifest) {
	if group.NestedDataType == nil {
		return
	}

	if group.Properties != nil && len(group.Properties) > 0 {
		data := newArgumentGroup(Anonymous)
		data.Description = strings.Join(group.Documentation.Content, "\n")

		processDirectives(group.Documentation.Directives, manifest, &data, nil)

		for _, field := range group.Properties {
			row := newArgument(field.DataTypeStr, field.Name, strings.Join(field.Documentation.Content, "\n"))

			if field.DefaultValue != nil {
				row.DefaultValue = field.DefaultValue
			}

			if field.NestedDataType != nil {
				row.ComplexType = field.NestedDataType
			}

			processDirectives(field.Documentation.Directives, manifest, nil, &row)

			data.Rows = append(data.Rows, row)
		}

		manifest.Objects[*group.NestedDataType] = data
	}

	for _, field := range group.Properties {
		recordNested(field, manifest)
	}
}

func processStructAndNested(extras *DocumentedStruct, manifest *ModuleManifest) {
	if extras == nil {
		return
	}

	recordNested(extras.StructProperty, manifest)
}

func createArgumentFromDocBlock(name string, typeStr string, docBlk PropertyDocBlock, extras *DocumentedStruct, manifest *ModuleManifest) Argument {
	tableRow := newArgument(typeStr, name, strings.Join(docBlk.Content, "\n"))

	processDirectives(docBlk.Directives, manifest, nil, &tableRow)

	if extras != nil && extras.StructProperty.NestedDataType != nil {
		tableRow.Type = extras.StructProperty.DataTypeStr
		tableRow.ComplexType = extras.StructProperty.NestedDataType
	}

	return tableRow
}

func findTypeDirective(directives []DocDirective) *DocDirective {
	for _, directive := range directives {
		if directive.Name == "type" {
			return &directive
		}
	}

	return nil
}

func ParseModuleArgsIntoManifest(inputs []*terraform.Input, outputs []*terraform.Output) *ModuleManifest {
	templateData := newModuleManifest()

	for _, input := range inputs {
		var extras *DocumentedStruct
		if input.Type != "" {
			if documented, err := ParseIntoDocumentedStruct(string(input.Type), input.Name); err == nil && documented != nil {
				extras = documented
			}
		}

		docBlk := parseStringIntoDocBlock(string(input.Description))
		defaultValue := input.GetValue()
		tableRow := createArgumentFromDocBlock(input.Name, string(input.Type), docBlk, extras, templateData)
		tableRow.DefaultValue = &defaultValue
		tableRow.Required = input.Required

		if tableRow.Required {
			templateData.RequiredInputs.Rows = append(templateData.RequiredInputs.Rows, tableRow)
		} else {
			templateData.OptionalInputs.Rows = append(templateData.OptionalInputs.Rows, tableRow)
		}

		processStructAndNested(extras, templateData)
	}

	for _, output := range outputs {
		outputType := "unknown"

		if output.Description == "" {
			tableRow := newArgument(outputType, output.Name, "")
			tableRow.Sensitive = output.Sensitive
			templateData.Outputs.Rows = append(templateData.Outputs.Rows, tableRow)

			continue
		}

		docBlk := parseStringIntoDocBlock(string(output.Description))
		typeDef := findTypeDirective(docBlk.Directives)

		var extras *DocumentedStruct
		if typeDef != nil && len(typeDef.Parsed.Args) > 0 {
			outputType = typeDef.Parsed.Args[0]
			extras, _ = ParseIntoDocumentedStruct(outputType, output.Name)
		}

		tableRow := createArgumentFromDocBlock(output.Name, outputType, docBlk, extras, templateData)
		tableRow.Sensitive = output.Sensitive
		templateData.Outputs.Rows = append(templateData.Outputs.Rows, tableRow)

		processStructAndNested(extras, templateData)
	}

	return templateData
}
