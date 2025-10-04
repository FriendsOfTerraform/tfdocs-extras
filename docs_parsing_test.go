package main

import (
	"reflect"
	"testing"
)

func optStr(s *string) string {
	if s == nil {
		return "<nil>"
	}

	return *s
}

func strPtr(s string) *string {
	return &s
}

func TestParseDocBlock_WithLineComments(t *testing.T) {
	lines := []AstDocString{
		AstDocString("This is a description"),
		AstDocString("It spans multiple lines"),
		AstDocString("@since 1.0.0"),
		AstDocString("@param name The name parameter"),
	}

	block := AstDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := ParseDocBlock(block)

	expectedContent := []string{
		"This is a description",
		"It spans multiple lines",
	}

	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	if len(result.Directives) != 2 {
		t.Errorf("Expected 2 directives, got %d", len(result.Directives))
	}

	expectedDirectives := []DocDirective{
		{Name: "since", Content: "1.0.0"},
		{Name: "param", Content: "name The name parameter"},
	}

	if !reflect.DeepEqual(result.Directives, expectedDirectives) {
		t.Errorf("Expected directives %v, got %v", expectedDirectives, result.Directives)
	}
}

func TestParseDocBlock_WithBlockComment(t *testing.T) {
	blockContent := AstDocBlockString("This is a block comment\nWith multiple lines\n@deprecated Use new function\n@since 2.0.0")

	block := AstDocBlock{
		Lines: nil,
		Block: &blockContent,
	}

	result := ParseDocBlock(block)

	expectedContent := []string{
		"This is a block comment",
		"With multiple lines",
	}

	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	if len(result.Directives) != 2 {
		t.Errorf("Expected 2 directives, got %d", len(result.Directives))
	}

	expectedDirectives := []DocDirective{
		{Name: "deprecated", Content: "Use new function"},
		{Name: "since", Content: "2.0.0"},
	}

	if !reflect.DeepEqual(result.Directives, expectedDirectives) {
		t.Errorf("Expected directives %v, got %v", expectedDirectives, result.Directives)
	}
}

func TestParseDocBlock_MixedContent(t *testing.T) {
	lines := []AstDocString{
		AstDocString("  Leading whitespace should be trimmed  "),
		AstDocString(""),
		AstDocString("Empty lines are preserved"),
		AstDocString("@example some code example"),
		AstDocString("More content after directive"),
		AstDocString("@returns boolean value"),
	}

	block := AstDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := ParseDocBlock(block)

	// Check content (whitespace should be trimmed)
	expectedContent := []string{
		"Leading whitespace should be trimmed",
		"",
		"Empty lines are preserved",
		"More content after directive",
	}

	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	expectedDirectives := []DocDirective{
		{Name: "example", Content: "some code example"},
		{Name: "returns", Content: "boolean value"},
	}

	if !reflect.DeepEqual(result.Directives, expectedDirectives) {
		t.Errorf("Expected directives %v, got %v", expectedDirectives, result.Directives)
	}
}

func TestParseDocBlock_DirectiveWithoutContent(t *testing.T) {
	lines := []AstDocString{
		AstDocString("Some description"),
		AstDocString("@deprecated"),
		AstDocString("@internal"),
		AstDocString("@final"),
	}

	block := AstDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := ParseDocBlock(block)

	expectedContent := []string{"Some description"}
	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	// Check directives without content
	expectedDirectives := []DocDirective{
		{Name: "deprecated", Content: ""},
		{Name: "internal", Content: ""},
		{Name: "final", Content: ""},
	}

	if !reflect.DeepEqual(result.Directives, expectedDirectives) {
		t.Errorf("Expected directives %v, got %v", expectedDirectives, result.Directives)
	}
}

func TestParseDocBlock_OnlyDirectives(t *testing.T) {
	lines := []AstDocString{
		AstDocString("@since 1.0.0"),
		AstDocString("@author John Doe"),
		AstDocString("@version 2.1.0"),
	}

	block := AstDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := ParseDocBlock(block)

	// Should have no content, only directives
	if len(result.Content) != 0 {
		t.Errorf("Expected no content, got %v", result.Content)
	}

	expectedDirectives := []DocDirective{
		{Name: "since", Content: "1.0.0"},
		{Name: "author", Content: "John Doe"},
		{Name: "version", Content: "2.1.0"},
	}

	if !reflect.DeepEqual(result.Directives, expectedDirectives) {
		t.Errorf("Expected directives %v, got %v", expectedDirectives, result.Directives)
	}
}

func TestParseDocBlock_OnlyContent(t *testing.T) {
	lines := []AstDocString{
		AstDocString("This is just content"),
		AstDocString("No directives here"),
		AstDocString("Just plain documentation"),
	}

	block := AstDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := ParseDocBlock(block)

	expectedContent := []string{
		"This is just content",
		"No directives here",
		"Just plain documentation",
	}

	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	// Should have no directives
	if len(result.Directives) != 0 {
		t.Errorf("Expected no directives, got %v", result.Directives)
	}
}

func TestParseDocBlock_EmptyLines(t *testing.T) {
	lines := []AstDocString{
		AstDocString(""),
		AstDocString("   "),
		AstDocString("\t\t"),
	}

	block := AstDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := ParseDocBlock(block)

	// All lines should be trimmed to empty strings
	expectedContent := []string{"", "", ""}
	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	if len(result.Directives) != 0 {
		t.Errorf("Expected no directives, got %v", result.Directives)
	}
}

func TestParseDocBlock_ComplexDirective(t *testing.T) {
	lines := []AstDocString{
		AstDocString("Function description"),
		AstDocString("@param name string The user's name"),
		AstDocString("@param age number The user's age in years"),
		AstDocString("@returns {user: User} The created user object"),
		AstDocString("@throws {ValidationError} When validation fails"),
	}

	block := AstDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := ParseDocBlock(block)

	expectedContent := []string{"Function description"}
	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	expectedDirectives := []DocDirective{
		{Name: "param", Content: "name string The user's name"},
		{Name: "param", Content: "age number The user's age in years"},
		{Name: "returns", Content: "{user: User} The created user object"},
		{Name: "throws", Content: "{ValidationError} When validation fails"},
	}

	if !reflect.DeepEqual(result.Directives, expectedDirectives) {
		t.Errorf("Expected directives %v, got %v", expectedDirectives, result.Directives)
	}
}

func TestParseDocBlock_EdgeCases(t *testing.T) {
	lines := []AstDocString{
		AstDocString("@ Invalid directive with space"),
		AstDocString("@"),
		AstDocString("@@double"),
		AstDocString("@ "),
		AstDocString("@valid content"),
	}

	block := AstDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := ParseDocBlock(block)

	// The function treats ALL lines starting with @ as directives
	// So there should be no content lines
	if len(result.Content) != 0 {
		t.Errorf("Expected no content, got %v", result.Content)
	}

	// All @ lines should be parsed as directives
	expectedDirectives := []DocDirective{
		{Name: "", Content: "Invalid directive with space"}, // "@ Invalid..." -> name="", content="Invalid..."
		{Name: "", Content: ""},                             // "@" -> name="", content=""
		{Name: "@double", Content: ""},                      // "@@double" -> name="@double", content=""
		{Name: "", Content: ""},                             // "@ " -> name="", content=""
		{Name: "valid", Content: "content"},                 // "@valid content" -> name="valid", content="content"
	}

	if !reflect.DeepEqual(result.Directives, expectedDirectives) {
		t.Errorf("Expected directives %v, got %v", expectedDirectives, result.Directives)
	}
}

func TestParseDocBlock_EmptyBlock(t *testing.T) {
	// Test with empty Lines slice
	block := AstDocBlock{
		Lines: []AstDocString{},
		Block: nil,
	}

	result := ParseDocBlock(block)

	if len(result.Content) != 0 {
		t.Errorf("Expected no content for empty block, got %v", result.Content)
	}

	if len(result.Directives) != 0 {
		t.Errorf("Expected no directives for empty block, got %v", result.Directives)
	}
}

func TestParseDocBlock_EmptyBlockString(t *testing.T) {
	// Test with empty block string
	emptyBlock := AstDocBlockString("")
	block := AstDocBlock{
		Lines: nil,
		Block: &emptyBlock,
	}

	result := ParseDocBlock(block)

	// Empty block string will result in one empty content line after splitting by newline
	expectedContent := []string{""}
	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v for empty block string, got %v", expectedContent, result.Content)
	}

	if len(result.Directives) != 0 {
		t.Errorf("Expected no directives for empty block string, got %v", result.Directives)
	}
}

func TestParseDocBlock_BlockWithNewlines(t *testing.T) {
	// Test block comment with various newline patterns
	blockContent := AstDocBlockString("First line\n\nSecond line after empty line\n@since 1.0.0\n\n@param test\n")

	block := AstDocBlock{
		Lines: nil,
		Block: &blockContent,
	}

	result := ParseDocBlock(block)

	expectedContent := []string{
		"First line",
		"",
		"Second line after empty line",
		"",
		"",
	}

	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	expectedDirectives := []DocDirective{
		{Name: "since", Content: "1.0.0"},
		{Name: "param", Content: "test"},
	}

	if !reflect.DeepEqual(result.Directives, expectedDirectives) {
		t.Errorf("Expected directives %v, got %v", expectedDirectives, result.Directives)
	}
}

func TestFlattenSimpleTypes_PrimitiveTypes(t *testing.T) {
	primitiveStr := "string"
	primitiveObj := AstDataType{
		Primitive: &primitiveStr,
	}

	primitiveFlattened := FlattenSimpleTypes(primitiveObj)
	if primitiveFlattened == nil || *primitiveFlattened != "string" {
		t.Errorf("Expected 'string', got %v", primitiveFlattened)
	}

	numberStr := "30"
	numberObj := AstDataType{
		Number: &numberStr,
	}
	numberFlattened := FlattenSimpleTypes(numberObj)
	if numberFlattened == nil || *numberFlattened != "30" {
		t.Errorf("Expected '30', got %v", numberFlattened)
	}

	stringStr := "\"hello\""
	stringObj := AstDataType{
		String: &stringStr,
	}
	stringFlattened := FlattenSimpleTypes(stringObj)
	if stringFlattened == nil || *stringFlattened != "\"hello\"" {
		t.Errorf("Expected '\"hello\"', got %v", stringFlattened)
	}
}

func TestFlattenSimpleTypes_Functions(t *testing.T) {
	genericType := "string"
	defaultValue := "\"default\""

	listObj := AstDataType{
		Func: &AstFunction{
			Name: "list",
			Args: []*AstDataType{
				{
					Primitive: &genericType,
				},
			},
		},
	}

	fxnFlattened := FlattenSimpleTypes(listObj)
	if fxnFlattened == nil || *fxnFlattened != "list(string)" {
		t.Errorf("Expected \"list(string)\", got \"%s\"", optStr(fxnFlattened))
	}

	optionalObj := AstDataType{
		Func: &AstFunction{
			Name: "optional",
			Args: []*AstDataType{
				{
					Primitive: &genericType,
				},
				{
					String: &defaultValue,
				},
			},
		},
	}

	optionalFlattened := FlattenSimpleTypes(optionalObj)
	if optionalFlattened == nil || *optionalFlattened != "optional(string, \"default\")" {
		t.Errorf("Expected \"optional(string, \"default\")\", got \"%s\"", optStr(optionalFlattened))
	}
}

func TestParseObjectBlock_AstObject(t *testing.T) {
	obj := AstObject{
		Pairs: []*AstObjectProperty{
			{
				Doc: &AstDocBlock{
					Lines: []AstDocString{
						AstDocString("The name of the user"),
						AstDocString(""),
						AstDocString("@since 1.0.0"),
					},
				},
				Key: "enable_managed_scaling_draining",
				Value: &AstDataType{
					Func: &AstFunction{
						Name: "optional",
						Args: []*AstDataType{
							{
								Primitive: strPtr("bool"),
							},
							{
								Primitive: strPtr("true"),
							},
						},
					},
				},
			},
			{
				Doc: &AstDocBlock{
					Lines: []AstDocString{
						AstDocString("The age of the user"),
						AstDocString(""),
						AstDocString("@since 1.0.0"),
					},
				},
				Key: "enable_scale_in_protection",
				Value: &AstDataType{
					Primitive: strPtr("number"),
				},
			},
		},
	}

	actual := ParseObjectBlock(obj)

	if len(actual) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(actual))
	}

	expected := []ObjectField{
		{
			VariableMetadata: VariableMetadata{
				Name: "enable_managed_scaling_draining",
				Documentation: FieldDocBlock{
					Content: []string{"The name of the user", ""},
					Directives: []DocDirective{
						{Name: "since", Content: "1.0.0"},
					},
				},
				Optional:     true,
				DefaultValue: strPtr("true"),
			},
			PrimitiveDataType: "bool",
		},
		{
			VariableMetadata: VariableMetadata{
				Name: "enable_scale_in_protection",
				Documentation: FieldDocBlock{
					Content: []string{"The age of the user", ""},
					Directives: []DocDirective{
						{Name: "since", Content: "1.0.0"},
					},
				},
				Optional:     false,
				DefaultValue: nil,
			},
			PrimitiveDataType: "number",
		},
	}

	if reflect.DeepEqual(actual, expected) != true {
		t.Errorf("Expected fields %v, got %v", expected, actual)
	}
}

func TestParseFunctionBlock_OptionalObject(t *testing.T) {
	obj := AstFunction{
		Name: "optional",
		Args: []*AstDataType{
			{
				Func: &AstFunction{
					Name: "object",
					Args: []*AstDataType{
						{
							Object: &AstObject{
								Pairs: []*AstObjectProperty{
									{
										Key: "name",
										Value: &AstDataType{
											Primitive: strPtr("string"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	objBlock := ParseFunctionBlock(obj, "test_object")

	expected := ObjectGroup{
		VariableMetadata: VariableMetadata{
			Name: "test_object",
		},
		Fields: []ObjectField{
			{
				VariableMetadata: VariableMetadata{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:         "name",
					Optional:     false,
					DefaultValue: nil,
				},
				PrimitiveDataType: "string",
				NestedDataType:    nil,
			},
		},
	}

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	if reflect.DeepEqual(*objBlock, expected) != true {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}

func TestParseFunctionBlock_NestedObject(t *testing.T) {
	obj := AstFunction{
		Name: "optional",
		Args: []*AstDataType{
			{
				Func: &AstFunction{
					Name: "object",
					Args: []*AstDataType{
						{
							Object: &AstObject{
								Pairs: []*AstObjectProperty{
									{
										Key: "name",
										Value: &AstDataType{
											Primitive: strPtr("string"),
										},
									},
									{
										Key: "address",
										Value: &AstDataType{
											Object: &AstObject{
												Pairs: []*AstObjectProperty{
													{
														Key: "street",
														Value: &AstDataType{
															Primitive: strPtr("string"),
														},
													},
													{
														Key: "city",
														Value: &AstDataType{
															Primitive: strPtr("string"),
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	objBlock := ParseFunctionBlock(obj, "user_profile")

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	expected := ObjectGroup{
		VariableMetadata: VariableMetadata{
			Name: "user_profile",
		},
		Fields: []ObjectField{
			{
				VariableMetadata: VariableMetadata{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:         "name",
					Optional:     false,
					DefaultValue: nil,
				},
				PrimitiveDataType: "string",
				NestedDataType:    nil,
			},
			{
				VariableMetadata: VariableMetadata{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:         "address",
					Optional:     false,
					DefaultValue: nil,
				},
				PrimitiveDataType: "object(Address)",
				NestedDataType: &ObjectGroup{
					VariableMetadata: VariableMetadata{
						Name: "address",
						Documentation: FieldDocBlock{
							Content:    []string{},
							Directives: []DocDirective{},
						},
						Optional:     false,
						DefaultValue: nil,
					},
					Fields: []ObjectField{
						{
							VariableMetadata: VariableMetadata{
								Documentation: FieldDocBlock{
									Content:    []string{},
									Directives: []DocDirective{},
								},
								Name:         "street",
								Optional:     false,
								DefaultValue: nil,
							},
							PrimitiveDataType: "string",
							NestedDataType:    nil,
						},
						{
							VariableMetadata: VariableMetadata{
								Documentation: FieldDocBlock{
									Content:    []string{},
									Directives: []DocDirective{},
								},
								Name:         "city",
								Optional:     false,
								DefaultValue: nil,
							},
							PrimitiveDataType: "string",
							NestedDataType:    nil,
						},
					},
				},
			},
		},
	}

	if reflect.DeepEqual(*objBlock, expected) != true {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}
