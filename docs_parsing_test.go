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
		{
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
		ObjectField: ObjectField{
			Name: "test_object",
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Optional:       true,
			DefaultValue:   nil,
			DataTypeStr:    "object(TestObject)",
			NestedDataType: strPtr("TestObject"),
			Fields: []ObjectField{
				{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:         "name",
					DataTypeStr:  "string",
					Optional:     false,
					DefaultValue: nil,
					Fields:       nil,
				},
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
		ObjectField: ObjectField{
			Name:           "user_profile",
			DataTypeStr:    "object(UserProfile)",
			NestedDataType: strPtr("UserProfile"),
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Optional: true,
			Fields: []ObjectField{
				{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:         "name",
					DataTypeStr:  "string",
					Optional:     false,
					DefaultValue: nil,
					Fields:       nil,
				},
				{
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					Name:           "address",
					DataTypeStr:    "object(Address)",
					NestedDataType: strPtr("Address"),
					Optional:       false,
					DefaultValue:   nil,
					Fields: []ObjectField{
						{
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							Name:         "street",
							DataTypeStr:  "string",
							Optional:     false,
							DefaultValue: nil,
							Fields:       nil,
						},
						{
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							Name:         "city",
							DataTypeStr:  "string",
							Optional:     false,
							DefaultValue: nil,
							Fields:       nil,
						},
					},
				},
			},
		},
		ParentDataType: strPtr(""),
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
		ObjectField: ObjectField{
			Name:           "instance_config",
			DataTypeStr:    "map(object(InstanceConfig))",
			NestedDataType: strPtr("InstanceConfig"),
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Fields: []ObjectField{
				{
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
					Fields:       nil,
				},
			},
		},
		ParentDataType: strPtr("map"),
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
		ObjectField: ObjectField{
			Name:           "server_config",
			DataTypeStr:    "list(object(ServerConfig))",
			NestedDataType: strPtr("ServerConfig"),
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Fields: []ObjectField{
				{
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
					Fields:       nil,
				},
				{
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
					Fields:       nil,
				},
			},
		},
		ParentDataType: strPtr("list"),
	}

	if !reflect.DeepEqual(*objBlock, expected) {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}

func TestParseObjectWithOptionalNestedObject(t *testing.T) {
	obj := astFunction{
		Name: "object",
		Args: []*astDataType{
			{
				Object: &astObject{
					Pairs: []*astObjectProperty{
						{
							Key: "ssh_keypair_name",
							Value: &astDataType{
								Primitive: strPtr("string"),
							},
						},
						{
							Key: "enable_managed_scaling",
							Value: &astDataType{
								Func: &astFunction{
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
																	Key: "enable_managed_scaling_draining",
																	Value: &astDataType{
																		Primitive: strPtr("bool"),
																	},
																},
																{
																	Key: "enable_scale_in_protection",
																	Value: &astDataType{
																		Primitive: strPtr("bool"),
																	},
																},
																{
																	Key: "target_capacity_percentage",
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
								},
							},
						},
					},
				},
			},
		},
	}

	objBlock := parseObjectFunctionBlock(obj, "root_object")

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	expected := ObjectGroup{
		ObjectField: ObjectField{
			Name:           "root_object",
			DataTypeStr:    "object(RootObject)",
			NestedDataType: strPtr("RootObject"),
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Fields: []ObjectField{
				{
					Name: "ssh_keypair_name",
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					DataTypeStr:  "string",
					Optional:     false,
					DefaultValue: nil,
					Fields:       nil,
				},
				{
					Name: "enable_managed_scaling",
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					DataTypeStr:    "object(EnableManagedScaling)",
					NestedDataType: strPtr("EnableManagedScaling"),
					Optional:       true,
					DefaultValue:   nil,
					Fields: []ObjectField{
						{
							Name: "enable_managed_scaling_draining",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "bool",
							Optional:     false,
							DefaultValue: nil,
							Fields:       nil,
						},
						{
							Name: "enable_scale_in_protection",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "bool",
							Optional:     false,
							DefaultValue: nil,
							Fields:       nil,
						},
						{
							Name: "target_capacity_percentage",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "number",
							Optional:     false,
							DefaultValue: nil,
						},
					},
				},
			},
		},
		ParentDataType: strPtr(""),
	}

	if !reflect.DeepEqual(*objBlock, expected) {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}

func TestParseObjectWithRequiredNestedObject(t *testing.T) {
	obj := astFunction{
		Name: "object",
		Args: []*astDataType{
			{
				Object: &astObject{
					Pairs: []*astObjectProperty{
						{
							Key: "ssh_keypair_name",
							Value: &astDataType{
								Primitive: strPtr("string"),
							},
						},
						{
							Key: "enable_managed_scaling",
							Value: &astDataType{
								Func: &astFunction{
									Name: "object",
									Args: []*astDataType{
										{
											Object: &astObject{
												Pairs: []*astObjectProperty{
													{
														Key: "enable_managed_scaling_draining",
														Value: &astDataType{
															Primitive: strPtr("bool"),
														},
													},
													{
														Key: "enable_scale_in_protection",
														Value: &astDataType{
															Primitive: strPtr("bool"),
														},
													},
													{
														Key: "target_capacity_percentage",
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
					},
				},
			},
		},
	}

	objBlock := parseObjectFunctionBlock(obj, "root_object")

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	expected := ObjectGroup{
		ObjectField: ObjectField{
			Name:           "root_object",
			DataTypeStr:    "object(RootObject)",
			NestedDataType: strPtr("RootObject"),
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Fields: []ObjectField{
				{
					Name: "ssh_keypair_name",
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					DataTypeStr:  "string",
					Optional:     false,
					DefaultValue: nil,
					Fields:       nil,
				},
				{
					Name: "enable_managed_scaling",
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					DataTypeStr:    "object(EnableManagedScaling)",
					NestedDataType: strPtr("EnableManagedScaling"),
					Optional:       false,
					DefaultValue:   nil,
					Fields: []ObjectField{
						{
							Name: "enable_managed_scaling_draining",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "bool",
							Optional:     false,
							DefaultValue: nil,
							Fields:       nil,
						},
						{
							Name: "enable_scale_in_protection",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "bool",
							Optional:     false,
							DefaultValue: nil,
							Fields:       nil,
						},
						{
							Name: "target_capacity_percentage",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "number",
							Optional:     false,
							DefaultValue: nil,
						},
					},
				},
			},
		},
		ParentDataType: strPtr(""),
	}

	if !reflect.DeepEqual(*objBlock, expected) {
		t.Errorf("Expected object group %v, got %v", expected, *objBlock)
	}
}

func TestParseComplexNestedHealthCheckObject(t *testing.T) {
	obj, err := parseAst(`list(object({
	  name            = string
	  health_check = optional(object({
		enabled                    = optional(bool, true)
		invert_health_check_status = optional(bool, false)
	
		calculated_check = optional(object({
		  health_checks_to_monitor = list(string)
		  healthy_threshold        = optional(number, null)
		}), null)
	
		cloudwatch_alarm_check = optional(object({
		  alarm_name   = string
		  alarm_region = optional(string, null)
		}), null)
	
		cloudwatch_alarms = optional(map(object({
		  metric_name = string # HealthCheckPercentageHealthy, HealthCheckStatus, ChildHealthCheckHealthyCount
		  expression  = string # statistic comparison_operator threshold
		})), {})
	
		endpoint_check = optional(object({
		  url                   = string
		  enable_latency_graphs = optional(bool, false)
		}), null)
	  }), null)
	}))`)

	if err != nil {
		t.Fatalf("Failed to parse AST: %v", err)
	}

	if obj.Expr == nil || obj.Expr.Func == nil {
		t.Fatal("Expected non-nil function expression")
	}

	objBlock := parseObjectFunctionBlock(*obj.Expr.Func, "health_check_config")

	if objBlock == nil {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	// Verify the top-level structure
	if objBlock.ObjectField.Name != "health_check_config" {
		t.Errorf("Expected name 'health_check_config', got '%s'", objBlock.ObjectField.Name)
	}

	if objBlock.ObjectField.DataTypeStr != "list(object(HealthCheckConfig))" {
		t.Errorf("Expected DataTypeStr 'list(object(HealthCheckConfig))', got '%s'", objBlock.ObjectField.DataTypeStr)
	}

	if objBlock.ParentDataType == nil || *objBlock.ParentDataType != "list" {
		t.Errorf("Expected ParentDataType 'list', got %v", objBlock.ParentDataType)
	}

	// Verify top-level fields: name and health_check
	if len(objBlock.ObjectField.Fields) != 2 {
		t.Fatalf("Expected 2 top-level fields, got %d", len(objBlock.ObjectField.Fields))
	}

	// Verify 'name' field
	nameField := objBlock.ObjectField.Fields[0]
	if nameField.Name != "name" {
		t.Errorf("Expected first field 'name', got '%s'", nameField.Name)
	}
	if nameField.DataTypeStr != "string" {
		t.Errorf("Expected name DataTypeStr 'string', got '%s'", nameField.DataTypeStr)
	}
	if nameField.Optional {
		t.Error("Expected name to not be optional")
	}

	// Verify 'health_check' field
	healthCheckField := objBlock.ObjectField.Fields[1]
	if healthCheckField.Name != "health_check" {
		t.Errorf("Expected second field 'health_check', got '%s'", healthCheckField.Name)
	}
	if !healthCheckField.Optional {
		t.Error("Expected health_check to be optional")
	}
	if healthCheckField.DefaultValue == nil || *healthCheckField.DefaultValue != "null" {
		t.Errorf("Expected health_check default value 'null', got %v", healthCheckField.DefaultValue)
	}

	// Verify health_check has 6 nested fields
	if len(healthCheckField.Fields) != 6 {
		t.Fatalf("Expected 6 health_check fields, got %d", len(healthCheckField.Fields))
	}

	expectedHealthCheckFields := []string{
		"enabled",
		"invert_health_check_status",
		"calculated_check",
		"cloudwatch_alarm_check",
		"cloudwatch_alarms",
		"endpoint_check",
	}

	for i, expectedName := range expectedHealthCheckFields {
		if healthCheckField.Fields[i].Name != expectedName {
			t.Errorf("Expected health_check field %d to be '%s', got '%s'", i, expectedName, healthCheckField.Fields[i].Name)
		}
	}

	// Verify 'enabled' field
	enabledField := healthCheckField.Fields[0]
	if !enabledField.Optional {
		t.Error("Expected enabled to be optional")
	}
	if enabledField.DataTypeStr != "bool" {
		t.Errorf("Expected enabled DataTypeStr 'bool', got '%s'", enabledField.DataTypeStr)
	}
	if enabledField.DefaultValue == nil || *enabledField.DefaultValue != "true" {
		t.Errorf("Expected enabled default value 'true', got %v", enabledField.DefaultValue)
	}

	// Verify 'calculated_check' nested object
	calculatedCheckField := healthCheckField.Fields[2]
	if !calculatedCheckField.Optional {
		t.Error("Expected calculated_check to be optional")
	}
	if len(calculatedCheckField.Fields) != 2 {
		t.Fatalf("Expected 2 calculated_check fields, got %d", len(calculatedCheckField.Fields))
	}

	// Verify 'health_checks_to_monitor' is a list
	healthChecksToMonitorField := calculatedCheckField.Fields[0]
	if healthChecksToMonitorField.Name != "health_checks_to_monitor" {
		t.Errorf("Expected 'health_checks_to_monitor', got '%s'", healthChecksToMonitorField.Name)
	}
	if healthChecksToMonitorField.DataTypeStr != "list(string)" {
		t.Errorf("Expected DataTypeStr 'list(string)', got '%s'", healthChecksToMonitorField.DataTypeStr)
	}

	// Verify 'cloudwatch_alarms' is a map of objects
	cloudwatchAlarmsField := healthCheckField.Fields[4]
	if !cloudwatchAlarmsField.Optional {
		t.Error("Expected cloudwatch_alarms to be optional")
	}
	if cloudwatchAlarmsField.DataTypeStr != "map(object(CloudwatchAlarms))" {
		t.Errorf("Expected DataTypeStr 'map(object(CloudwatchAlarms))', got '%s'", cloudwatchAlarmsField.DataTypeStr)
	}
	if len(cloudwatchAlarmsField.Fields) != 2 {
		t.Fatalf("Expected 2 cloudwatch_alarms fields, got %d", len(cloudwatchAlarmsField.Fields))
	}
	if cloudwatchAlarmsField.Fields[0].Name != "metric_name" {
		t.Errorf("Expected 'metric_name', got '%s'", cloudwatchAlarmsField.Fields[0].Name)
	}
	if cloudwatchAlarmsField.Fields[1].Name != "expression" {
		t.Errorf("Expected 'expression', got '%s'", cloudwatchAlarmsField.Fields[1].Name)
	}

	// Verify 'endpoint_check' nested object
	endpointCheckField := healthCheckField.Fields[5]
	if !endpointCheckField.Optional {
		t.Error("Expected endpoint_check to be optional")
	}
	if len(endpointCheckField.Fields) != 2 {
		t.Fatalf("Expected 2 endpoint_check fields, got %d", len(endpointCheckField.Fields))
	}
	if endpointCheckField.Fields[0].Name != "url" {
		t.Errorf("Expected 'url', got '%s'", endpointCheckField.Fields[0].Name)
	}
	if endpointCheckField.Fields[1].Name != "enable_latency_graphs" {
		t.Errorf("Expected 'enable_latency_graphs', got '%s'", endpointCheckField.Fields[1].Name)
	}
	if !endpointCheckField.Fields[1].Optional {
		t.Error("Expected enable_latency_graphs to be optional")
	}
	if endpointCheckField.Fields[1].DefaultValue == nil || *endpointCheckField.Fields[1].DefaultValue != "false" {
		t.Errorf("Expected enable_latency_graphs default value 'false', got %v", endpointCheckField.Fields[1].DefaultValue)
	}
}
