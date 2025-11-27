package tfdocextras

import (
	"testing"

	"github.com/go-test/deep"
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

	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Errorf("RawContent mismatch:\n%v", diff)
	}

	if len(result.Directives) != 2 {
		t.Errorf("Expected 2 directives, got %d", len(result.Directives))
	}

	expectedDirectives := []DocDirective{
		{Name: "since", RawContent: "1.0.0", Parsed: ParsedDirective{Type: DirSince, First: "1.0.0", Flags: IsValid}},
		{Name: "param", RawContent: "name The name parameter", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
	}

	if diff := deep.Equal(result.Directives, expectedDirectives); diff != nil {
		t.Errorf("Directives mismatch:\n%v", diff)
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

	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Errorf("RawContent mismatch:\n%v", diff)
	}

	if len(result.Directives) != 2 {
		t.Errorf("Expected 2 directives, got %d", len(result.Directives))
	}

	expectedDirectives := []DocDirective{
		{Name: "deprecated", RawContent: "Use new function", Parsed: ParsedDirective{Type: DirDeprecated, First: "Use new function", Flags: IsValid}},
		{Name: "since", RawContent: "2.0.0", Parsed: ParsedDirective{Type: DirSince, First: "2.0.0", Flags: IsValid}},
	}

	if diff := deep.Equal(result.Directives, expectedDirectives); diff != nil {
		t.Errorf("Directives mismatch:\n%v", diff)
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

	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Errorf("RawContent mismatch:\n%v", diff)
	}

	expectedDirectives := []DocDirective{
		{Name: "example", RawContent: "some code example", Parsed: ParsedDirective{Type: DirExample, Flags: IsInvalid}},
		{Name: "returns", RawContent: "boolean value", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
	}

	if diff := deep.Equal(result.Directives, expectedDirectives); diff != nil {
		t.Errorf("Directives mismatch:\n%v", diff)
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
	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Errorf("RawContent mismatch:\n%v", diff)
	}

	// Check directives without content
	expectedDirectives := []DocDirective{
		{Name: "deprecated", RawContent: "", Parsed: ParsedDirective{Type: DirDeprecated, First: "", Flags: IsValid}},
		{Name: "internal", RawContent: "", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
		{Name: "final", RawContent: "", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
	}

	if diff := deep.Equal(result.Directives, expectedDirectives); diff != nil {
		t.Errorf("Directives mismatch:\n%v", diff)
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
		{Name: "since", RawContent: "1.0.0", Parsed: ParsedDirective{Type: DirSince, First: "1.0.0", Flags: IsValid}},
		{Name: "author", RawContent: "John Doe", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
		{Name: "version", RawContent: "2.1.0", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
	}

	if diff := deep.Equal(result.Directives, expectedDirectives); diff != nil {
		t.Errorf("Directives mismatch:\n%v", diff)
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

	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Errorf("RawContent mismatch:\n%v", diff)
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
	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Errorf("RawContent mismatch:\n%v", diff)
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
	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Errorf("RawContent mismatch:\n%v", diff)
	}

	expectedDirectives := []DocDirective{
		{Name: "param", RawContent: "name string The user's name", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
		{Name: "param", RawContent: "age number The user's age in years", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
		{Name: "returns", RawContent: "{user: User} The created user object", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
		{Name: "throws", RawContent: "{ValidationError} When validation fails", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
	}

	if diff := deep.Equal(result.Directives, expectedDirectives); diff != nil {
		t.Errorf("Directives mismatch:\n%v", diff)
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
		{Name: "", RawContent: "Invalid directive with space", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}}, // "@ Invalid..." -> name="", content="Invalid..."
		{Name: "", RawContent: "", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},                             // "@" -> name="", content=""
		{Name: "@double", RawContent: "", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},                      // "@@double" -> name="@double", content=""
		{Name: "", RawContent: "", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},                             // "@ " -> name="", content=""
	}

	if diff := deep.Equal(result.Directives[:4], expectedDirectives); diff != nil {
		t.Errorf("Directives mismatch:\n%v", diff)
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
	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Error(diff)
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

	if diff := deep.Equal(result.Content, expectedContent); diff != nil {
		t.Errorf("RawContent mismatch:\n%v", diff)
	}

	expectedDirectives := []DocDirective{
		{Name: "since", RawContent: "1.0.0", Parsed: ParsedDirective{Type: DirSince, First: "1.0.0", Flags: IsValid}},
		{Name: "param", RawContent: "test", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
	}

	if diff := deep.Equal(result.Directives, expectedDirectives); diff != nil {
		t.Errorf("Directives mismatch:\n%v", diff)
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

func TestParseIntoDocumentedStruct_ObjectWithDocStringComments(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`object({
		/// The name of the user
		///
		/// @since 1.0.0
		enable_managed_scaling_draining = optional(bool, true)
		/// The age of the user
		///
		/// @since 1.0.0
		enable_scale_in_protection = number
	})`, "test_object"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
		t.Fatalf("Failed to parse: %v", err)
	}

	if len(objBlock.Fields) != 2 {
		t.Fatalf("Expected 2 fields, got %d", len(objBlock.Fields))
	}

	expected := ObjectGroup{
		ObjectField: ObjectField{
			Name:           "test_object",
			DataTypeStr:    "object(TestObject)",
			NestedDataType: strPtr("TestObject"),
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Fields: []ObjectField{
				{
					Name: "enable_managed_scaling_draining",
					Documentation: FieldDocBlock{
						Content: []string{"The name of the user"},
						Directives: []DocDirective{
							{Name: "since", RawContent: "1.0.0", Parsed: ParsedDirective{Type: DirSince, First: "1.0.0", Flags: IsValid}},
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
							{Name: "since", RawContent: "1.0.0", Parsed: ParsedDirective{Type: DirSince, First: "1.0.0", Flags: IsValid}},
						},
					},
					DataTypeStr:  "number",
					Optional:     false,
					DefaultValue: nil,
				},
			},
		},
		ParentDataType: strPtr(""),
	}

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Errorf("ObjectGroup mismatch:\n%v", diff)
	}
}

func TestParseIntoDocumentedStruct_OptionalObject(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`optional(object({
		name = string
	}))`, "test_object"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
		t.Fatalf("Failed to parse: %v", err)
	}

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

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Errorf("ObjectGroup mismatch:\n%v", diff)
	}
}

func TestParseIntoDocumentedStruct_OptionalObjectWithObject(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`optional(object({
	  name = string
	  address = object({
	    street = string
	    city   = string
	  })
	}))`, "user_profile"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
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

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Error(diff)
	}
}

func TestParseIntoDocumentedStruct_MapOfObjects(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`map(object({
	  /// Specify the number of EC2 instances that should be running in the group
	  ///
	  /// @since 1.0.0
	  desired_instances = number
	}))`, "instance_config"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
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
							{Name: "since", RawContent: "1.0.0", Parsed: ParsedDirective{Type: DirSince, First: "1.0.0", Flags: IsValid}},
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

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Error(diff)
	}
}

func TestParseIntoDocumentedStruct_ListOfObjects(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`list(object({
	  /// The name of the server
	  /// @required true
	  server_name = string
	  /// The port number for the server
	  /// @default 80
	  port = optional(number, 80)
	}))`, "server_config"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
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
							{Name: "required", RawContent: "true", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
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
							{Name: "default", RawContent: "80", Parsed: ParsedDirective{Type: DirUnsupported, Flags: IsInvalid}},
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

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Error(diff)
	}
}

func TestParseIntoDocumentedStruct_ObjectWithOptionalObject(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`object({
	  ssh_keypair_name = string
	  enable_managed_scaling = optional(object({
	    enable_managed_scaling_draining = bool
	    enable_scale_in_protection      = bool
	    target_capacity_percentage      = number
	  }))
	})`, "root_object"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
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

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Error(diff)
	}
}

func TestParseIntoDocumentedStruct_ObjectWithNestedObject(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`object({
	  ssh_keypair_name = string
	  enable_managed_scaling = object({
	    enable_managed_scaling_draining = bool
	    enable_scale_in_protection      = bool
	    target_capacity_percentage      = number
	  })
	})`, "root_object"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
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

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Error(diff)
	}
}

func TestParseIntoDocumentedStruct_DoubleNestedWithinOptional(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`list(object({
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
	}))`, "health_check_config"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	expected := ObjectGroup{
		ObjectField: ObjectField{
			Name:           "health_check_config",
			DataTypeStr:    "list(object(HealthCheckConfig))",
			NestedDataType: strPtr("HealthCheckConfig"),
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Fields: []ObjectField{
				{
					Name: "name",
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					DataTypeStr:  "string",
					Optional:     false,
					DefaultValue: nil,
				},
				{
					Name: "health_check",
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					DataTypeStr:    "object(HealthCheck)",
					NestedDataType: strPtr("HealthCheck"),
					Optional:       true,
					DefaultValue:   strPtr("null"),
					Fields: []ObjectField{
						{
							Name: "enabled",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "bool",
							Optional:     true,
							DefaultValue: strPtr("true"),
						},
						{
							Name: "invert_health_check_status",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "bool",
							Optional:     true,
							DefaultValue: strPtr("false"),
						},
						{
							Name: "calculated_check",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:    "object(CalculatedCheck)",
							NestedDataType: strPtr("CalculatedCheck"),
							Optional:       true,
							DefaultValue:   strPtr("null"),
							Fields: []ObjectField{
								{
									Name: "health_checks_to_monitor",
									Documentation: FieldDocBlock{
										Content:    []string{},
										Directives: []DocDirective{},
									},
									DataTypeStr:  "list(string)",
									Optional:     false,
									DefaultValue: nil,
								},
								{
									Name: "healthy_threshold",
									Documentation: FieldDocBlock{
										Content:    []string{},
										Directives: []DocDirective{},
									},
									DataTypeStr:  "number",
									Optional:     true,
									DefaultValue: strPtr("null"),
								},
							},
						},
						{
							Name: "cloudwatch_alarm_check",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:    "object(CloudwatchAlarmCheck)",
							NestedDataType: strPtr("CloudwatchAlarmCheck"),
							Optional:       true,
							DefaultValue:   strPtr("null"),
							Fields: []ObjectField{
								{
									Name: "alarm_name",
									Documentation: FieldDocBlock{
										Content:    []string{},
										Directives: []DocDirective{},
									},
									DataTypeStr:  "string",
									Optional:     false,
									DefaultValue: nil,
								},
								{
									Name: "alarm_region",
									Documentation: FieldDocBlock{
										Content:    []string{},
										Directives: []DocDirective{},
									},
									DataTypeStr:  "string",
									Optional:     true,
									DefaultValue: strPtr("null"),
								},
							},
						},
						{
							Name: "cloudwatch_alarms",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:    "map(object(CloudwatchAlarms))",
							NestedDataType: strPtr("CloudwatchAlarms"),
							Optional:       true,
							DefaultValue:   strPtr("{}"),
							Fields: []ObjectField{
								{
									Name: "metric_name",
									Documentation: FieldDocBlock{
										Content:    []string{},
										Directives: []DocDirective{},
									},
									DataTypeStr:  "string",
									Optional:     false,
									DefaultValue: nil,
								},
								{
									Name: "expression",
									Documentation: FieldDocBlock{
										Content:    []string{},
										Directives: []DocDirective{},
									},
									DataTypeStr:  "string",
									Optional:     false,
									DefaultValue: nil,
								},
							},
						},
						{
							Name: "endpoint_check",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:    "object(EndpointCheck)",
							NestedDataType: strPtr("EndpointCheck"),
							Optional:       true,
							DefaultValue:   strPtr("null"),
							Fields: []ObjectField{
								{
									Name: "url",
									Documentation: FieldDocBlock{
										Content:    []string{},
										Directives: []DocDirective{},
									},
									DataTypeStr:  "string",
									Optional:     false,
									DefaultValue: nil,
								},
								{
									Name: "enable_latency_graphs",
									Documentation: FieldDocBlock{
										Content:    []string{},
										Directives: []DocDirective{},
									},
									DataTypeStr:  "bool",
									Optional:     true,
									DefaultValue: strPtr("false"),
								},
							},
						},
					},
				},
			},
		},
		ParentDataType: strPtr("list"),
	}

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Error(diff)
	}
}

func TestParseIntoDocumentedStruct_ObjectWithMapOfObjects(t *testing.T) {
	var objBlock *ObjectGroup

	if parsed, err := ParseIntoDocumentedStruct(`object({
	  configurations = map(object({
	    setting_a = string
	    setting_b = number
	  }))
	})`, "config_object"); err == nil && parsed != nil {
		objBlock = parsed
	} else {
		t.Fatal("Expected non-nil ObjectGroup")
	}

	expected := ObjectGroup{
		ObjectField: ObjectField{
			Name:           "config_object",
			DataTypeStr:    "object(ConfigObject)",
			NestedDataType: strPtr("ConfigObject"),
			Documentation: FieldDocBlock{
				Content:    []string{},
				Directives: []DocDirective{},
			},
			Fields: []ObjectField{
				{
					Name: "configurations",
					Documentation: FieldDocBlock{
						Content:    []string{},
						Directives: []DocDirective{},
					},
					DataTypeStr:    "map(object(Configurations))",
					NestedDataType: strPtr("Configurations"),
					Optional:       false,
					DefaultValue:   nil,
					Fields: []ObjectField{
						{
							Name: "setting_a",
							Documentation: FieldDocBlock{
								Content:    []string{},
								Directives: []DocDirective{},
							},
							DataTypeStr:  "string",
							Optional:     false,
							DefaultValue: nil,
						},
						{
							Name: "setting_b",
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

	if diff := deep.Equal(*objBlock, expected); diff != nil {
		t.Error(diff)
	}
}
