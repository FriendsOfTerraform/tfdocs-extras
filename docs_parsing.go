package main

import (
	"strings"
)

type DocDirective struct {
	Name    string
	Content string
}

type FieldDocBlock struct {
	Content    []string
	Directives []DocDirective
}

type ObjectField struct {
	Name              string
	Documentation     FieldDocBlock
	PrimitiveDataType string
	NestedDataType    *ObjectField
	Optional          bool
	DefaultValue      *string
}

type VariableDocument struct {
	Parent []string
	Fields []ObjectField
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

func ParseObjectBlock(obj AstObject) []ObjectField {
	var fields []ObjectField

	for _, pair := range obj.Pairs {
		field := ObjectField{
			Name: pair.Key,
		}

		if pair.Doc != nil {
			field.Documentation = ParseDocBlock(*pair.Doc)
		}

		if isOptionalType(*pair.Value) {
			field.Optional = true
			if len(pair.Value.Func.Args) >= 1 {
				if flattened := FlattenSimpleTypes(*pair.Value.Func.Args[0]); flattened != nil {
					field.PrimitiveDataType = *flattened
				}

				if len(pair.Value.Func.Args) >= 2 {
					if defaultFlattened := FlattenSimpleTypes(*pair.Value.Func.Args[1]); defaultFlattened != nil {
						field.DefaultValue = defaultFlattened
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
						Name:           "",
						NestedDataType: nil,
					}
				}
			}
		}

		fields = append(fields, field)
	}

	return fields
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

func ParseAstRoot(root AstRoot) ([]VariableDocument, error) {
	var doc []VariableDocument

	return doc, nil
}
