package tfdocextras

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// DocDirective represents a documentation directive like @since, @param, etc.
type DocDirective struct {
	Name       string          `json:"name"`
	Parsed     ParsedDirective `json:"parsed"`
	RawContent string          `json:"rawContent"`
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

func isCollectionType(data astDataType) bool {
	return data.Func != nil && (data.Func.Name == "map" || data.Func.Name == "list")
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
	field.Fields = parseObjectBlock(*obj)
	field.DataTypeStr = "object(" + objectName + ")"
	field.NestedDataType = &objectName
}

func parseOptionalField(field *ObjectField, args []*astDataType) {
	field.Optional = true

	if len(args) == 0 {
		return
	}

	// Try to parse the first argument (the type)
	parseOptionalFieldType(field, args[0])

	// Handle default value if present (second argument)
	if len(args) >= 2 {
		setDefaultValue(field, args[1])
	}
}

// parseOptionalFieldType handles parsing the type argument of an optional() call
func parseOptionalFieldType(field *ObjectField, arg *astDataType) {
	// Handle map(object({...})) or list(object({...}))
	if isCollectionType(*arg) {
		if parseCollectionOfObjects(field, arg.Func) {
			return
		}
	}

	// Handle object({...}) directly
	if handleObjectField(field, *arg) {
		return
	}

	// Handle primitive types and other functions
	if flattened := flattenSimpleTypes(*arg); flattened != nil {
		field.DataTypeStr = *flattened
	}
}

// parseCollectionOfObjects handles map(object({...})) or list(object({...})) patterns
func parseCollectionOfObjects(field *ObjectField, fn *astFunction) bool {
	if len(fn.Args) == 0 {
		return false
	}

	// Check if the argument is object({...})
	if !isObjectType(*fn.Args[0]) || fn.Args[0].Func.Args[0].Object == nil {
		return false
	}

	objectName := getObjectName(field.Name)
	field.Fields = parseObjectBlock(*fn.Args[0].Func.Args[0].Object)
	field.DataTypeStr = fn.Name + "(object(" + objectName + "))"
	field.NestedDataType = &objectName

	return true
}

// setDefaultValue sets the default value for a field from an astDataType
func setDefaultValue(field *ObjectField, defaultArg *astDataType) {
	// Try to flatten primitives, numbers, strings, and function calls
	if defaultFlattened := flattenSimpleTypes(*defaultArg); defaultFlattened != nil {
		field.DefaultValue = defaultFlattened
		return
	}

	// Handle empty object literal: {}
	if defaultArg.Object != nil {
		emptyObj := "{}"
		field.DefaultValue = &emptyObj
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
	lineNo := 1
	indentation := ""
	doc := FieldDocBlock{}

	iterateDocLines(block, func(line string) {
		trimmed := strings.TrimLeft(line, " \t")

		if lineNo == 1 {
			indentation = line[:len(line)-len(trimmed)]
		}

		if strings.HasPrefix(line, indentation) {
			line = line[len(indentation):]
		}

		if strings.HasPrefix(line, "@") {
			name, content, _ := strings.Cut(line[1:], " ")
			doc.Directives = append(doc.Directives, DocDirective{
				Name:       name,
				Parsed:     ParseDirective(name, content),
				RawContent: content,
			})
		} else {
			doc.Content = append(doc.Content, strings.TrimSpace(line))
		}

		lineNo++
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

		parseFieldType(&field, pair.Value)
		fields = append(fields, field)
	}

	return fields
}

// parseFieldType determines and sets the type information for a field
func parseFieldType(field *ObjectField, value *astDataType) {
	switch {
	case isOptionalType(*value):
		parseOptionalField(field, value.Func.Args)
		return
	case isCollectionType(*value):
		// Handle map(object({...})) or list(object({...}))
		if parseCollectionOfObjects(field, value.Func) {
			return
		}
		// Fall back to flattening if not a collection of objects
	case handleObjectField(field, *value):
		return // Object function handled by helper
	case value.Object != nil:
		parseNestedObject(field, value.Object)
		return
	}

	// Fallback: try to flatten the type into a simple string
	if flattened := flattenSimpleTypes(*value); flattened != nil {
		field.DataTypeStr = *flattened
	}
}

// buildObjectTypeName constructs an object type string with optional prefix (map, list, etc.)
func buildObjectTypeName(name, prefix string) string {
	objectName := getObjectName(name)
	objectType := "object(" + objectName + ")"

	if prefix != "" {
		return prefix + "(" + objectType + ")"
	}
	return objectType
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

	// Handle optional(object({...}))
	if fxn.Name == "optional" {
		objGroup.Optional = true
	}

	if len(fxn.Args) > 0 {
		if fields := extractObjectFromArg(fxn.Args[0]); fields != nil {
			objGroup.Fields = fields
			objectName := getObjectName(name)
			objGroup.NestedDataType = &objectName
			objGroup.DataTypeStr = buildObjectTypeName(name, collectionPrefix)
		}
	}

	return objGroup
}

func parseStringIntoDocBlock(input string) FieldDocBlock {
	str := astDocBlockString(input)
	blk := astDocBlock{
		Block: &str,
	}

	return parseDocBlock(blk)
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
