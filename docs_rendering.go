package tfdocsextras

import (
	"strings"

	"github.com/terraform-docs/terraform-docs/terraform"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type MetadataAttribute struct {
	Name    string `json:"name"`
	Content string `json:"content"`
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
	ParentType   ArgGroupType `json:"-"`
	Type         string       `json:"type"`
	ComplexType  *string      `json:"complex_type,omitempty"`
	Name         string       `json:"name"`
	DefaultValue *string      `json:"default_value"`
	Required     bool         `json:"required"`
	Sensitive    bool         `json:"sensitive"`
	Description  string       `json:"description"`
	ArgumentMetadata
}

func (r *Argument) GetAnchor() string {
	if r.ComplexType == nil {
		return ""
	}

	return anchorFor(*r.ComplexType)
}

// objectPathSeparator joins an object type name to the one it is nested in.
// A dot cannot appear in an HCL identifier, so a path built out of dots is
// unambiguous, and it reads the way the attribute itself would be referenced.
const objectPathSeparator = "."

// outputObjectPrefix namespaces the object types reached from outputs. Input
// names are unique among inputs, and output names among outputs, but a module
// may well have an output named after the input it passes through - and both
// would otherwise claim the same object type name, with the second silently
// overwriting the first.
const outputObjectPrefix = "output" + objectPathSeparator

// joinObjectPath qualifies a nested object type name with the path of the one
// it is declared in, so that two objects with the same field name - say a
// `details` inside each of two variables - stay distinct.
func joinObjectPath(parent string, name string) string {
	if parent == "" {
		return name
	}

	return parent + objectPathSeparator + name
}

// requalifyObjectType rewrites the `object(name)` inside a type expression to
// carry the nested type's full path, keeping the rendered type and the anchor
// it links to in agreement.
func requalifyObjectType(dataType string, name string, path string) string {
	return strings.Replace(dataType, "object("+name+")", "object("+path+")", 1)
}

// anchorFor derives the fragment identifier a markdown renderer gives to a
// heading: lowercased, spaces turned into dashes, and anything that is not a
// letter, digit, dash or underscore dropped. Object type names are paths now,
// and the dots in them do not survive into the anchor.
func anchorFor(name string) string {
	var anchor strings.Builder

	for _, char := range strings.ToLower(name) {
		switch {
		case char >= 'a' && char <= 'z', char >= '0' && char <= '9', char == '-', char == '_':
			anchor.WriteRune(char)
		case char == ' ':
			anchor.WriteRune('-')
		}
	}

	return anchor.String()
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

	return r.DefaultValue != nil && *r.DefaultValue != ""
}

func (r *Argument) HasNonEmptyObjectDefault() bool {
	if r.ComplexType == nil || !r.HasDefaultValue() {
		return false
	}

	if !strings.HasPrefix(strings.TrimSpace(r.Type), "object(") {
		return false
	}

	defaultValue := strings.TrimSpace(*r.DefaultValue)

	return defaultValue != "{}" && defaultValue != "null"
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
	Required    bool         `json:"required"`
	Description string       `json:"description"`
	Rows        []Argument   `json:"rows"`
	ArgumentMetadata
}

func (d *ArgumentGroup) GetMetadata() *ArgumentMetadata {
	return &d.ArgumentMetadata
}

type ModuleManifest struct {
	RequiredInputs ArgumentGroup            `json:"required_inputs"`
	OptionalInputs ArgumentGroup            `json:"optional_inputs"`
	Outputs        ArgumentGroup            `json:"outputs"`
	Objects        map[string]ArgumentGroup `json:"objects"`
	ReferenceLinks map[string]string        `json:"reference_links"`
}

func newArgumentGroup(grpType ArgGroupType) ArgumentGroup {
	return ArgumentGroup{
		Type:        grpType,
		Description: "",
		Required:    false,
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

func newArgument(parentType ArgGroupType, typeStr, name, description string) Argument {
	return Argument{
		ParentType:   parentType,
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

// recordNested records the object type declared by group under 'path', its
// fully qualified name, and recurses into the object types declared inside it.
func recordNested(group StructProperty, path string, manifest *ModuleManifest) {
	if group.NestedDataType == nil {
		return
	}

	if len(group.Properties) > 0 {
		data := newArgumentGroup(Anonymous)
		data.Required = !group.Optional
		data.Description = strings.Join(group.Documentation.Content, "\n")

		processDirectives(group.Documentation.Directives, manifest, &data, nil)

		for _, field := range group.Properties {
			row := newArgument(data.Type, field.DataTypeStr, field.Name, strings.Join(field.Documentation.Content, "\n"))

			if field.DefaultValue != nil {
				row.DefaultValue = field.DefaultValue
			}

			if field.NestedDataType != nil {
				nested := joinObjectPath(path, *field.NestedDataType)
				row.ComplexType = &nested
				row.Type = requalifyObjectType(field.DataTypeStr, *field.NestedDataType, nested)
			}

			row.Required = !field.Optional

			processDirectives(field.Documentation.Directives, manifest, nil, &row)

			data.Rows = append(data.Rows, row)
		}

		manifest.Objects[path] = data
	}

	for _, field := range group.Properties {
		if field.NestedDataType == nil {
			continue
		}

		recordNested(field, joinObjectPath(path, *field.NestedDataType), manifest)
	}
}

func processStructAndNested(extras *DocumentedStruct, path string, manifest *ModuleManifest) {
	if extras == nil {
		return
	}

	recordNested(extras.StructProperty, path, manifest)
}

func createArgumentFromDocBlock(argType ArgGroupType, name string, typeStr string, docBlk PropertyDocBlock, extras *DocumentedStruct, path string, manifest *ModuleManifest) Argument {
	tableRow := newArgument(argType, typeStr, name, strings.Join(docBlk.Content, "\n"))

	processDirectives(docBlk.Directives, manifest, nil, &tableRow)

	if extras != nil && extras.StructProperty.NestedDataType != nil {
		tableRow.Type = requalifyObjectType(extras.StructProperty.DataTypeStr, *extras.StructProperty.NestedDataType, path)
		tableRow.ComplexType = &path
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

		argType := OptionalInput
		if input.Required {
			argType = RequiredInput
		}

		// An input's object types are rooted at the input name, which is
		// already unique among inputs.
		path := input.Name

		tableRow := createArgumentFromDocBlock(argType, input.Name, string(input.Type), docBlk, extras, path, templateData)
		tableRow.DefaultValue = &defaultValue
		tableRow.Required = input.Required

		if tableRow.Required {
			templateData.RequiredInputs.Rows = append(templateData.RequiredInputs.Rows, tableRow)
		} else {
			templateData.OptionalInputs.Rows = append(templateData.OptionalInputs.Rows, tableRow)
		}

		processStructAndNested(extras, path, templateData)
	}

	for _, output := range outputs {
		outputType := "unknown"

		if output.Description == "" {
			tableRow := newArgument(Output, outputType, output.Name, "")
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

		// An output may share its name with the input it passes through, so
		// its object types are namespaced to keep the two apart.
		path := outputObjectPrefix + output.Name

		tableRow := createArgumentFromDocBlock(Output, output.Name, outputType, docBlk, extras, path, templateData)
		tableRow.Sensitive = output.Sensitive
		templateData.Outputs.Rows = append(templateData.Outputs.Rows, tableRow)

		processStructAndNested(extras, path, templateData)
	}

	return templateData
}
