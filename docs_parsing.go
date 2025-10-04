package tfdocextras

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DocDirective represents a documentation directive like @since, @param, etc.
type DocDirective struct {
	Name    string
	Content string
}

// FieldDocBlock contains parsed documentation for a field
type FieldDocBlock struct {
	Content    []string
	Directives []DocDirective
}

// VariableMetadata contains metadata about a variable or field
type VariableMetadata struct {
	Name          string
	Documentation FieldDocBlock
	DataTypeStr   string
	Optional      bool
	DefaultValue  *string
}

// ObjectField represents a field within an object structure
type ObjectField struct {
	VariableMetadata

	NestedDataType *ObjectGroup
}

// ObjectGroup represents a group of related object fields with documentation
type ObjectGroup struct {
	VariableMetadata

	Fields         []ObjectField
	ParentDataType *string
}

// GetObjectName returns the CamelCase name for this object group
func (o *ObjectGroup) GetObjectName() string {
	if o.Name != "" {
		// Convert to CamelCase
		caser := cases.Title(language.English)
		parts := strings.Split(o.Name, "_")

		for i, part := range parts {
			parts[i] = caser.String(part)
		}

		return strings.Join(parts, "")
	}

	return "UnknownObject"
}

func extractObjectFromArg(arg *astDataType) []ObjectField {
	if arg.Func != nil && arg.Func.Name == "object" && len(arg.Func.Args) > 0 {
		if arg.Func.Args[0].Object != nil {
			return parseObjectBlock(*arg.Func.Args[0].Object)
		}
	} else if arg.Object != nil {
		return parseObjectBlock(*arg.Object)
	}

	return nil
}

func isOptionalType(data astDataType) bool {
	return data.Func != nil && data.Func.Name == "optional"
}

func iterateDocLines(block astDocBlock, fn func(line string)) {
	if block.Block != nil {
		lines := strings.Split(string(*block.Block), "\n")
		for _, line := range lines {
			fn(line)
		}
	} else {
		for _, line := range block.Lines {
			fn(string(line))
		}
	}
}

func newEmptyFieldDocBlock() FieldDocBlock {
	return FieldDocBlock{
		Content:    []string{},
		Directives: []DocDirective{},
	}
}

func newVariableMetadata(name string) VariableMetadata {
	return VariableMetadata{
		Name:          name,
		Documentation: newEmptyFieldDocBlock(),
		Optional:      false,
		DefaultValue:  nil,
		DataTypeStr:   "",
	}
}

func parseNestedObject(field *ObjectField, obj *astObject) {
	nestedFields := parseObjectBlock(*obj)
	nestedObjectGroup := &ObjectGroup{
		VariableMetadata: VariableMetadata{
			Name: field.VariableMetadata.Name,
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Optional:     false,
			DefaultValue: nil,
		},
		Fields: nestedFields,
	}

	field.VariableMetadata.DataTypeStr = "object(" + nestedObjectGroup.GetObjectName() + ")"
	field.NestedDataType = nestedObjectGroup
}

func parseOptionalField(field *ObjectField, args []*astDataType) {
	field.VariableMetadata.Optional = true

	if len(args) >= 1 {
		if flattened := flattenSimpleTypes(*args[0]); flattened != nil {
			field.VariableMetadata.DataTypeStr = *flattened
		}

		if len(args) >= 2 {
			if defaultFlattened := flattenSimpleTypes(*args[1]); defaultFlattened != nil {
				field.VariableMetadata.DefaultValue = defaultFlattened
			}
		}
	}
}

func trimEmptyLines(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	// Find first non-empty line
	start := 0
	for start < len(lines) && lines[start] == "" {
		start++
	}

	// Find last non-empty line
	end := len(lines) - 1
	for end >= start && lines[end] == "" {
		end--
	}

	// If all lines are empty, return empty slice
	if start > end {
		return []string{}
	}

	return lines[start : end+1]
}

func flattenSimpleTypes(data astDataType) *string {
	if data.Primitive != nil {
		return data.Primitive
	} else if data.Number != nil {
		return data.Number
	} else if data.String != nil {
		return data.String
	} else if data.Func != nil {
		fxnName := data.Func.Name

		// Build the function signature: "functionName(arg1, arg2, ...)"
		var args []string
		for _, arg := range data.Func.Args {
			if argFlattened := flattenSimpleTypes(*arg); argFlattened != nil {
				args = append(args, *argFlattened)
			}
		}

		result := fxnName + "(" + strings.Join(args, ", ") + ")"
		return &result
	}

	return nil
}

func parseDocBlock(block astDocBlock) FieldDocBlock {
	doc := FieldDocBlock{}

	iterateDocLines(block, func(line string) {
		if strings.HasPrefix(line, "@") {
			parts := strings.SplitN(line[1:], " ", 2)
			if len(parts) == 2 {
				doc.Directives = append(doc.Directives, DocDirective{
					Name:    parts[0],
					Content: parts[1],
				})
			} else if len(parts) == 1 {
				doc.Directives = append(doc.Directives, DocDirective{
					Name:    parts[0],
					Content: "",
				})
			}
		} else {
			doc.Content = append(doc.Content, strings.TrimSpace(line))
		}
	})

	doc.Content = trimEmptyLines(doc.Content)

	return doc
}

func parseObjectBlock(obj astObject) []ObjectField {
	var fields []ObjectField

	for _, pair := range obj.Pairs {
		field := ObjectField{
			VariableMetadata: newVariableMetadata(pair.Key),
		}

		if pair.Doc != nil {
			field.VariableMetadata.Documentation = parseDocBlock(*pair.Doc)
		}

		if isOptionalType(*pair.Value) {
			parseOptionalField(&field, pair.Value.Func.Args)
		} else if flattened := flattenSimpleTypes(*pair.Value); flattened != nil {
			field.VariableMetadata.DataTypeStr = *flattened
		} else if pair.Value.Object != nil {
			parseNestedObject(&field, pair.Value.Object)
		}

		fields = append(fields, field)
	}

	return fields
}

func parseObjectFunctionBlock(fxn astFunction, name string) *ObjectGroup {
	collectionPrefix := ""

	switch fxn.Name {
	case "map":
		collectionPrefix = "map"
	case "list":
		collectionPrefix = "list"
	case "optional", "object":
		// No collection prefix needed
	default:
		return nil
	}

	objGroup := &ObjectGroup{
		VariableMetadata: newVariableMetadata(name),
		ParentDataType:   &collectionPrefix,
	}

	var fields []ObjectField

	// Handle optional(object({...}))
	if fxn.Name == "optional" && len(fxn.Args) > 0 {
		objGroup.VariableMetadata.Optional = true
	}

	if len(fxn.Args) > 0 {
		fields = extractObjectFromArg(fxn.Args[0])
	}

	if fields != nil {
		objGroup.Fields = fields
		objectTypeName := "object(" + objGroup.GetObjectName() + ")"

		if collectionPrefix != "" {
			objGroup.VariableMetadata.DataTypeStr = collectionPrefix + "(" + objectTypeName + ")"
		} else {
			objGroup.VariableMetadata.DataTypeStr = objectTypeName
		}
	}

	return objGroup
}

// ParseIntoDocumentedStruct parses a Terraform type definition string into a documented object group.
// This is the main entry point for the library.
//
// Example usage:
//
//	group, err := ParseIntoDocumentedStruct(`optional(object({
//	  /// The user's name
//	  /// @since 1.0.0
//	  name = string
//	  age = optional(number, 18)
//	}))`, "user_config")
func ParseIntoDocumentedStruct(input string, name string) (*ObjectGroup, error) {
	root, err := parseAst(input)
	if err != nil {
		return nil, err
	}

	value := root.Expr

	if value.Func != nil {
		return parseObjectFunctionBlock(*value.Func, name), nil
	}

	return nil, nil
}
