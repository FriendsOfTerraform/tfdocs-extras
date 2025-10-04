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
	Optional      bool
	DefaultValue  *string
}

type ObjectField struct {
	VariableMetadata

	PrimitiveDataType string
	NestedDataType    *ObjectField
}

type ObjectGroup struct {
	VariableMetadata

	Fields []ObjectField
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

func isOptionalType(data AstDataType) bool {
	if data.Func != nil && data.Func.Name == "optional" {
		return true
	}

	return false
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

	return doc
}

func ParseObjectBlock(obj AstObject) []ObjectField {
	var fields []ObjectField

	for _, pair := range obj.Pairs {
		field := ObjectField{
			VariableMetadata: VariableMetadata{
				Name: pair.Key,
				Documentation: FieldDocBlock{
					Content:    []string{},
					Directives: []DocDirective{},
				},
			},
		}

		if pair.Doc != nil {
			field.VariableMetadata.Documentation = ParseDocBlock(*pair.Doc)
		}

		if isOptionalType(*pair.Value) {
			field.VariableMetadata.Optional = true
			if len(pair.Value.Func.Args) >= 1 {
				if flattened := FlattenSimpleTypes(*pair.Value.Func.Args[0]); flattened != nil {
					field.PrimitiveDataType = *flattened
				}

				if len(pair.Value.Func.Args) >= 2 {
					if defaultFlattened := FlattenSimpleTypes(*pair.Value.Func.Args[1]); defaultFlattened != nil {
						field.VariableMetadata.DefaultValue = defaultFlattened
					}
				}
			}
		} else {
			if flattened := FlattenSimpleTypes(*pair.Value); flattened != nil {
				field.PrimitiveDataType = *flattened
			} else if pair.Value.Object != nil {
				nestedFields := ParseObjectBlock(*pair.Value.Object)
				if len(nestedFields) == 1 {
					field.NestedDataType = &nestedFields[0]
				} else if len(nestedFields) > 1 {
					field.NestedDataType = &ObjectField{
						VariableMetadata: VariableMetadata{
							Name: "",
						},
						NestedDataType: nil,
					}
				}
			}
		}

		fields = append(fields, field)
	}

	return fields
}

func ParseFunctionBlock(fxn AstFunction, name string) *ObjectGroup {
	objGroup := &ObjectGroup{
		VariableMetadata: VariableMetadata{
			Name: name,
		},
	}

	// Handle optional(object({...}))
	if fxn.Name == "optional" && len(fxn.Args) > 0 {
		firstArg := fxn.Args[0]

		// Check if the first argument is an object function
		if firstArg.Func != nil && firstArg.Func.Name == "object" && len(firstArg.Func.Args) > 0 {
			objectArg := firstArg.Func.Args[0]
			if objectArg.Object != nil {
				objGroup.Fields = ParseObjectBlock(*objectArg.Object)
			}
		}
	} else if fxn.Name == "object" && len(fxn.Args) > 0 {
		// Handle direct object({...})
		objectArg := fxn.Args[0]

		if objectArg.Object != nil {
			objGroup.Fields = ParseObjectBlock(*objectArg.Object)
		}
	} else {
		return nil
	}

	return objGroup
}
