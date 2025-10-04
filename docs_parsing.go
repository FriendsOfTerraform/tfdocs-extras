package main

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type DocDirective struct {
	Name    string
	Content string
}

type FieldDocBlock struct {
	Content    []string
	Directives []DocDirective
}

type VariableMetadata struct {
	Name          string
	Documentation FieldDocBlock
	DataTypeStr   string
	Optional      bool
	DefaultValue  *string
}

type ObjectField struct {
	VariableMetadata

	NestedDataType *ObjectGroup
}

type ObjectGroup struct {
	VariableMetadata

	Fields         []ObjectField
	ParentDataType *string
}

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

func extractObjectFromArg(arg *AstDataType) []ObjectField {
	if arg.Func != nil && arg.Func.Name == "object" && len(arg.Func.Args) > 0 {
		if arg.Func.Args[0].Object != nil {
			return ParseObjectBlock(*arg.Func.Args[0].Object)
		}
	} else if arg.Object != nil {
		return ParseObjectBlock(*arg.Object)
	}

	return nil
}

func isOptionalType(data AstDataType) bool {
	return data.Func != nil && data.Func.Name == "optional"
}

func iterateDocLines(block AstDocBlock, fn func(line string)) {
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

func parseNestedObject(field *ObjectField, obj *AstObject) {
	nestedFields := ParseObjectBlock(*obj)
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

func parseOptionalField(field *ObjectField, args []*AstDataType) {
	field.VariableMetadata.Optional = true

	if len(args) >= 1 {
		if flattened := FlattenSimpleTypes(*args[0]); flattened != nil {
			field.VariableMetadata.DataTypeStr = *flattened
		}

		if len(args) >= 2 {
			if defaultFlattened := FlattenSimpleTypes(*args[1]); defaultFlattened != nil {
				field.VariableMetadata.DefaultValue = defaultFlattened
			}
		}
	}
}

// trimEmptyLines removes empty lines from the beginning and end of a slice
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

func FlattenSimpleTypes(data AstDataType) *string {
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
			if argFlattened := FlattenSimpleTypes(*arg); argFlattened != nil {
				args = append(args, *argFlattened)
			}
		}

		result := fxnName + "(" + strings.Join(args, ", ") + ")"
		return &result
	}

	return nil
}

func ParseDocBlock(block AstDocBlock) FieldDocBlock {
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

func ParseObjectBlock(obj AstObject) []ObjectField {
	var fields []ObjectField

	for _, pair := range obj.Pairs {
		field := ObjectField{
			VariableMetadata: newVariableMetadata(pair.Key),
		}

		if pair.Doc != nil {
			field.VariableMetadata.Documentation = ParseDocBlock(*pair.Doc)
		}

		if isOptionalType(*pair.Value) {
			parseOptionalField(&field, pair.Value.Func.Args)
		} else if flattened := FlattenSimpleTypes(*pair.Value); flattened != nil {
			field.VariableMetadata.DataTypeStr = *flattened
		} else if pair.Value.Object != nil {
			parseNestedObject(&field, pair.Value.Object)
		}

		fields = append(fields, field)
	}

	return fields
}

// ParseObjectFunctionBlock handles parsing of object-containing functions including:
// - optional(object({...}))
// - object({...})
// - map(object({...}))
// - list(object({...}))
func ParseObjectFunctionBlock(fxn AstFunction, name string) *ObjectGroup {
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

func ParseIntoDocumentedGroup(root AstRoot, name string) *ObjectGroup {
	value := root.Expr

	if value.Func != nil {
		return ParseObjectFunctionBlock(*value.Func, name)
	}

	return nil
}
