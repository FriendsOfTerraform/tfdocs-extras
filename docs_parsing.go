package tfdocextras

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DocDirective represents a documentation directive like @since, @param, etc.
type DocDirective struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// FieldDocBlock contains parsed documentation for a field
type FieldDocBlock struct {
	Content    []string       `json:"content"`
	Directives []DocDirective `json:"directives"`
}

// ObjectField represents a field within an object structure
type ObjectField struct {
	Name           string        `json:"name"`
	Documentation  FieldDocBlock `json:"documentation"`
	DataTypeStr    string        `json:"dataType"`
	Optional       bool          `json:"optional"`
	DefaultValue   *string       `json:"defaultValue,omitempty"`
	NestedDataType *string       `json:"nestedDataType,omitempty"`
	Fields         []ObjectField `json:"fields,omitempty"`
}

// ObjectGroup represents a group of related object fields with documentation
type ObjectGroup struct {
	ObjectField `json:",inline"`

	ParentDataType *string `json:"parentDataType,omitempty"`
}

func extractObjectFromArg(arg *astDataType) []ObjectField {
	if isObjectType(*arg) {
		if arg.Func.Args[0].Object != nil {
			return parseObjectBlock(*arg.Func.Args[0].Object)
		}
	} else if arg.Object != nil {
		return parseObjectBlock(*arg.Object)
	}

	return nil
}

func getObjectName(name string) string {
	if name != "" {
		caser := cases.Title(language.English)
		parts := strings.Split(name, "_")

		for i, part := range parts {
			parts[i] = caser.String(part)
		}

		return strings.Join(parts, "")
	}

	return "UnknownObject"
}

func isObjectType(data astDataType) bool {
	return data.Func != nil && data.Func.Name == "object" && len(data.Func.Args) > 0
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

func newObjectField(name string) ObjectField {
	return ObjectField{
		Name:          name,
		Documentation: newEmptyFieldDocBlock(),
		Optional:      false,
		DefaultValue:  nil,
		DataTypeStr:   "",
	}
}

func handleObjectField(field *ObjectField, dataType astDataType) bool {
	if isObjectType(dataType) && dataType.Func.Args[0].Object != nil {
		parseNestedObject(field, dataType.Func.Args[0].Object)

		return true
	}

	return false
}

func parseNestedObject(field *ObjectField, obj *astObject) {
	objectName := getObjectName(field.Name)
	nestedFields := parseObjectBlock(*obj)

	field.DataTypeStr = "object(" + objectName + ")"
	field.NestedDataType = &objectName
	field.Fields = nestedFields
}

func parseOptionalField(field *ObjectField, args []*astDataType) {
	field.Optional = true

	if len(args) >= 1 {
		if handleObjectField(field, *args[0]) {
			return
		}

		if flattened := flattenSimpleTypes(*args[0]); flattened != nil {
			field.DataTypeStr = *flattened
		}

		if len(args) >= 2 {
			if defaultFlattened := flattenSimpleTypes(*args[1]); defaultFlattened != nil {
				field.DefaultValue = defaultFlattened
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
	}
	if data.Number != nil {
		return data.Number
	}
	if data.String != nil {
		return data.String
	}
	if data.Func != nil {
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
		field := newObjectField(pair.Key)

		if pair.Doc != nil {
			field.Documentation = parseDocBlock(*pair.Doc)
		}

		if isOptionalType(*pair.Value) {
			parseOptionalField(&field, pair.Value.Func.Args)
		} else if handleObjectField(&field, *pair.Value) {
			// Object function handled by helper
		} else if flattened := flattenSimpleTypes(*pair.Value); flattened != nil {
			field.DataTypeStr = *flattened
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
		ObjectField:    newObjectField(name),
		ParentDataType: &collectionPrefix,
	}

	var fields []ObjectField

	// Handle optional(object({...}))
	if fxn.Name == "optional" && len(fxn.Args) > 0 {
		objGroup.Optional = true
	}

	if len(fxn.Args) > 0 {
		fields = extractObjectFromArg(fxn.Args[0])
	}

	if fields != nil {
		objGroup.Fields = fields
		objectName := getObjectName(name)
		objectTypeName := "object(" + objectName + ")"

		objGroup.NestedDataType = &objectName

		if collectionPrefix != "" {
			objGroup.DataTypeStr = collectionPrefix + "(" + objectTypeName + ")"
		} else {
			objGroup.DataTypeStr = objectTypeName
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
