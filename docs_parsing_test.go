package tfdocextras

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
	lines := []astDocString{
		astDocString("This is a description"),
		astDocString("It spans multiple lines"),
		astDocString("@since 1.0.0"),
		astDocString("@param name The name parameter"),
	}

	block := astDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := parseDocBlock(block)

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
	blockContent := astDocBlockString("This is a block comment\nWith multiple lines\n@deprecated Use new function\n@since 2.0.0")

	block := astDocBlock{
		Lines: nil,
		Block: &blockContent,
	}

	result := parseDocBlock(block)

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
	lines := []astDocString{
		astDocString("  Leading whitespace should be trimmed  "),
		astDocString(""),
		astDocString("Empty lines are preserved"),
		astDocString("@example some code example"),
		astDocString("More content after directive"),
		astDocString("@returns boolean value"),
	}

	block := astDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := parseDocBlock(block)

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
	lines := []astDocString{
		astDocString("Some description"),
		astDocString("@deprecated"),
		astDocString("@internal"),
		astDocString("@final"),
	}

	block := astDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := parseDocBlock(block)

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
	lines := []astDocString{
		astDocString("@since 1.0.0"),
		astDocString("@author John Doe"),
		astDocString("@version 2.1.0"),
	}

	block := astDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := parseDocBlock(block)

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
	lines := []astDocString{
		astDocString("This is just content"),
		astDocString("No directives here"),
		astDocString("Just plain documentation"),
	}

	block := astDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := parseDocBlock(block)

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
	lines := []astDocString{
		astDocString(""),
		astDocString("   "),
		astDocString("\t\t"),
	}

	block := astDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := parseDocBlock(block)

	// All lines should be trimmed to empty strings
	expectedContent := []string{}
	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v, got %v", expectedContent, result.Content)
	}

	if len(result.Directives) != 0 {
		t.Errorf("Expected no directives, got %v", result.Directives)
	}
}

func TestParseDocBlock_ComplexDirective(t *testing.T) {
	lines := []astDocString{
		astDocString("Function description"),
		astDocString("@param name string The user's name"),
		astDocString("@param age number The user's age in years"),
		astDocString("@returns {user: User} The created user object"),
		astDocString("@throws {ValidationError} When validation fails"),
	}

	block := astDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := parseDocBlock(block)

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
	lines := []astDocString{
		astDocString("@ Invalid directive with space"),
		astDocString("@"),
		astDocString("@@double"),
		astDocString("@ "),
		astDocString("@valid content"),
	}

	block := astDocBlock{
		Lines: lines,
		Block: nil,
	}

	result := parseDocBlock(block)

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
	block := astDocBlock{
		Lines: []astDocString{},
		Block: nil,
	}

	result := parseDocBlock(block)

	if len(result.Content) != 0 {
		t.Errorf("Expected no content for empty block, got %v", result.Content)
	}

	if len(result.Directives) != 0 {
		t.Errorf("Expected no directives for empty block, got %v", result.Directives)
	}
}

func TestParseDocBlock_EmptyBlockString(t *testing.T) {
	// Test with empty block string
	emptyBlock := astDocBlockString("")
	block := astDocBlock{
		Lines: nil,
		Block: &emptyBlock,
	}

	result := parseDocBlock(block)

	// Empty block string will result in one empty content line after splitting by newline
	expectedContent := []string{}
	if !reflect.DeepEqual(result.Content, expectedContent) {
		t.Errorf("Expected content %v for empty block string, got %v", expectedContent, result.Content)
	}

	if len(result.Directives) != 0 {
		t.Errorf("Expected no directives for empty block string, got %v", result.Directives)
	}
}

func TestParseDocBlock_BlockWithNewlines(t *testing.T) {
	// Test block comment with various newline patterns
	blockContent := astDocBlockString("First line\n\nSecond line after empty line\n@since 1.0.0\n\n@param test\n")

	block := astDocBlock{
		Lines: nil,
		Block: &blockContent,
	}

	result := parseDocBlock(block)

	expectedContent := []string{
		"First line",
		"",
		"Second line after empty line",
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
	primitiveObj := astDataType{
		Primitive: &primitiveStr,
	}

	primitiveFlattened := flattenSimpleTypes(primitiveObj)
	if primitiveFlattened == nil || *primitiveFlattened != "string" {
		t.Errorf("Expected 'string', got %v", primitiveFlattened)
	}

	numberStr := "30"
	numberObj := astDataType{
		Number: &numberStr,
	}
	numberFlattened := flattenSimpleTypes(numberObj)
	if numberFlattened == nil || *numberFlattened != "30" {
		t.Errorf("Expected '30', got %v", numberFlattened)
	}

	stringStr := "\"hello\""
	stringObj := astDataType{
		String: &stringStr,
	}
	stringFlattened := flattenSimpleTypes(stringObj)
	if stringFlattened == nil || *stringFlattened != "\"hello\"" {
		t.Errorf("Expected '\"hello\"', got %v", stringFlattened)
	}
}

func TestFlattenSimpleTypes_Functions(t *testing.T) {
	genericType := "string"
	defaultValue := "\"default\""

	listObj := astDataType{
		Func: &astFunction{
			Name: "list",
			Args: []*astDataType{
				{
					Primitive: &genericType,
				},
			},
		},
	}

	fxnFlattened := flattenSimpleTypes(listObj)
	if fxnFlattened == nil || *fxnFlattened != "list(string)" {
		t.Errorf("Expected \"list(string)\", got \"%s\"", optStr(fxnFlattened))
	}

	optionalObj := astDataType{
		Func: &astFunction{
			Name: "optional",
			Args: []*astDataType{
				{
					Primitive: &genericType,
				},
				{
					String: &defaultValue,
				},
			},
		},
	}

	optionalFlattened := flattenSimpleTypes(optionalObj)
	if optionalFlattened == nil || *optionalFlattened != "optional(string, \"default\")" {
		t.Errorf("Expected \"optional(string, \"default\")\", got \"%s\"", optStr(optionalFlattened))
	}
}

func TestParseObjectBlock_AstObject(t *testing.T) {
	obj := astObject{
		Pairs: []*astObjectProperty{
			{
				Doc: &astDocBlock{
					Lines: []astDocString{
						astDocString("The name of the user"),
						astDocString(""),
						astDocString("@since 1.0.0"),
					},
				},
				Key: "enable_managed_scaling_draining",
				Value: &astDataType{
					Func: &astFunction{
						Name: "optional",
						Args: []*astDataType{
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
				Doc: &astDocBlock{
					Lines: []astDocString{
						astDocString("The age of the user"),
						astDocString(""),
						astDocString("@since 1.0.0"),
					},
				},
				Key: "enable_scale_in_protection",
				Value: &astDataType{
					Primitive: strPtr("number"),
				},
			},
		},
	}

	actual := parseObjectBlock(obj)

	if len(actual) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(actual))
	}

	expected := []ObjectField{
		{
			VariableMetadata: VariableMetadata{
				Name: "enable_managed_scaling_draining",
				Documentation: FieldDocBlock{
					Content: []string{"The name of the user"},
					Directives: []DocDirective{
						{Name: "since", Content: "1.0.0"},
					},
				},
				DataTypeStr:  "bool",
				Optional:     true,
				DefaultValue: strPtr("true"),
			},
		},
		{
			VariableMetadata: VariableMetadata{
				Name: "enable_scale_in_protection",
				Documentation: FieldDocBlock{
					Content: []string{"The age of the user"},
					Directives: []DocDirective{
						{Name: "since", Content: "1.0.0"},
					},
				},
				DataTypeStr:  "number",
				Optional:     false,
				DefaultValue: nil,
			},
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected fields %v, got %v", expected, actual)
	}
}

func TestParseObjectFunctionBlock_OptionalObject(t *testing.T) {
	obj := astFunction{
		Name: "optional",
		Args: []*astDataType{
			{
				Func: &astFunction{
					Name: "object",
					Args: []*astDataType{
						{
							Object: &astObject{
								Pairs: []*astObjectProperty{
									{
										Key: "name",
										Value: &astDataType{
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

	objBlock := parseObjectFunctionBlock(obj, "test_object")

	expected := ObjectGroup{
		VariableMetadata: VariableMetadata{
			Name: "test_object",
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Optional:     true,
			DefaultValue: nil,
			DataTypeStr:  "object(TestObject)",
		},
		Fields: []ObjectField{
			{
				VariableMetadata: VariableMetadata{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:         "name",
					DataTypeStr:  "string",
					Optional:     false,
					DefaultValue: nil,
				},
				NestedDataType: nil,
			},
		},
		ParentDataType: strPtr(""),
	}

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	if !reflect.DeepEqual(*objBlock, expected) {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}

func TestParseObjectFunctionBlock_NestedObject(t *testing.T) {
	obj := astFunction{
		Name: "optional",
		Args: []*astDataType{
			{
				Func: &astFunction{
					Name: "object",
					Args: []*astDataType{
						{
							Object: &astObject{
								Pairs: []*astObjectProperty{
									{
										Key: "name",
										Value: &astDataType{
											Primitive: strPtr("string"),
										},
									},
									{
										Key: "address",
										Value: &astDataType{
											Object: &astObject{
												Pairs: []*astObjectProperty{
													{
														Key: "street",
														Value: &astDataType{
															Primitive: strPtr("string"),
														},
													},
													{
														Key: "city",
														Value: &astDataType{
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

	objBlock := parseObjectFunctionBlock(obj, "user_profile")

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	expected := ObjectGroup{
		VariableMetadata: VariableMetadata{
			Name:        "user_profile",
			DataTypeStr: "object(UserProfile)",
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Optional: true,
		},
		ParentDataType: strPtr(""),
		Fields: []ObjectField{
			{
				VariableMetadata: VariableMetadata{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:         "name",
					DataTypeStr:  "string",
					Optional:     false,
					DefaultValue: nil,
				},
				NestedDataType: nil,
			},
			{
				VariableMetadata: VariableMetadata{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:         "address",
					DataTypeStr:  "object(Address)",
					Optional:     false,
					DefaultValue: nil,
				},
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
								DataTypeStr:  "string",
								Optional:     false,
								DefaultValue: nil,
							},
							NestedDataType: nil,
						},
						{
							VariableMetadata: VariableMetadata{
								Documentation: FieldDocBlock{
									Content:    []string{},
									Directives: []DocDirective{},
								},
								Name:         "city",
								DataTypeStr:  "string",
								Optional:     false,
								DefaultValue: nil,
							},
							NestedDataType: nil,
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(*objBlock, expected) {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}

func TestParseMapFunctionBlock_MapObject(t *testing.T) {
	obj := astFunction{
		Name: "map",
		Args: []*astDataType{
			{
				Func: &astFunction{
					Name: "object",
					Args: []*astDataType{
						{
							Object: &astObject{
								Pairs: []*astObjectProperty{
									{
										Doc: &astDocBlock{
											Lines: []astDocString{
												astDocString("Specify the number of EC2 instances that should be running in the group"),
												astDocString(""),
												astDocString("@since 1.0.0"),
											},
										},
										Key: "desired_instances",
										Value: &astDataType{
											Primitive: strPtr("number"),
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

	objBlock := parseObjectFunctionBlock(obj, "instance_config")

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	expected := ObjectGroup{
		VariableMetadata: VariableMetadata{
			Name:        "instance_config",
			DataTypeStr: "map(object(InstanceConfig))",
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
		},
		ParentDataType: strPtr("map"),
		Fields: []ObjectField{
			{
				VariableMetadata: VariableMetadata{
					Name: "desired_instances",
					Documentation: FieldDocBlock{
						Content: []string{
							"Specify the number of EC2 instances that should be running in the group",
						},
						Directives: []DocDirective{
							{Name: "since", Content: "1.0.0"},
						},
					},
					DataTypeStr:  "number",
					Optional:     false,
					DefaultValue: nil,
				},
				NestedDataType: nil,
			},
		},
	}

	if !reflect.DeepEqual(*objBlock, expected) {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}

func TestParseListFunctionBlock_ListObject(t *testing.T) {
	obj := astFunction{
		Name: "list",
		Args: []*astDataType{
			{
				Func: &astFunction{
					Name: "object",
					Args: []*astDataType{
						{
							Object: &astObject{
								Pairs: []*astObjectProperty{
									{
										Doc: &astDocBlock{
											Lines: []astDocString{
												astDocString("The name of the server"),
												astDocString("@required true"),
											},
										},
										Key: "server_name",
										Value: &astDataType{
											Primitive: strPtr("string"),
										},
									},
									{
										Doc: &astDocBlock{
											Lines: []astDocString{
												astDocString("The port number for the server"),
												astDocString("@default 80"),
											},
										},
										Key: "port",
										Value: &astDataType{
											Func: &astFunction{
												Name: "optional",
												Args: []*astDataType{
													{
														Primitive: strPtr("number"),
													},
													{
														Number: strPtr("80"),
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

	objBlock := parseObjectFunctionBlock(obj, "server_config")

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	expected := ObjectGroup{
		VariableMetadata: VariableMetadata{
			Name:        "server_config",
			DataTypeStr: "list(object(ServerConfig))",
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
		},
		ParentDataType: strPtr("list"),
		Fields: []ObjectField{
			{
				VariableMetadata: VariableMetadata{
					Name: "server_name",
					Documentation: FieldDocBlock{
						Content: []string{
							"The name of the server",
						},
						Directives: []DocDirective{
							{Name: "required", Content: "true"},
						},
					},
					DataTypeStr:  "string",
					Optional:     false,
					DefaultValue: nil,
				},
				NestedDataType: nil,
			},
			{
				VariableMetadata: VariableMetadata{
					Name: "port",
					Documentation: FieldDocBlock{
						Content: []string{
							"The port number for the server",
						},
						Directives: []DocDirective{
							{Name: "default", Content: "80"},
						},
					},
					DataTypeStr:  "number",
					Optional:     true,
					DefaultValue: strPtr("80"),
				},
				NestedDataType: nil,
			},
		},
	}

	if !reflect.DeepEqual(*objBlock, expected) {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}
