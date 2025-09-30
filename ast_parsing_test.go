package main

import (
	"testing"
)

func TestParse_SimpleFunction(t *testing.T) {
	input := `map(string)`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Expr.Func == nil {
		t.Fatal("Expected AstFunction, got nil")
	}

	if result.Expr.Func.Name != "map" {
		t.Errorf("Expected AstFunction name 'map', got '%s'", result.Expr.Func.Name)
	}

	if len(result.Expr.Func.Args) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(result.Expr.Func.Args))
	}

	if result.Expr.Func.Args[0].Primitive == nil || *result.Expr.Func.Args[0].Primitive != "string" {
		t.Error("Expected primitive argument 'string'")
	}
}

func TestParse_NestedFunction(t *testing.T) {
	input := `map(AstObject({ name = string age = number }))`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Check outer AstFunction
	if result.Expr.Func == nil || result.Expr.Func.Name != "map" {
		t.Fatal("Expected outer AstFunction 'map'")
	}

	// Check inner AstFunction
	innerFunc := result.Expr.Func.Args[0].Func
	if innerFunc == nil || innerFunc.Name != "AstObject" {
		t.Fatal("Expected inner AstFunction 'AstObject'")
	}

	// Check AstObject argument
	obj := innerFunc.Args[0].Object
	if obj == nil {
		t.Fatal("Expected AstObject argument")
	}

	if len(obj.Pairs) != 2 {
		t.Errorf("Expected 2 AstObject properties, got %d", len(obj.Pairs))
	}

	// Check first property
	if obj.Pairs[0].Key != "name" {
		t.Errorf("Expected first property key 'name', got '%s'", obj.Pairs[0].Key)
	}
	if obj.Pairs[0].Value.Primitive == nil || *obj.Pairs[0].Value.Primitive != "string" {
		t.Error("Expected first property value 'string'")
	}

	// Check second property
	if obj.Pairs[1].Key != "age" {
		t.Errorf("Expected second property key 'age', got '%s'", obj.Pairs[1].Key)
	}
	if obj.Pairs[1].Value.Primitive == nil || *obj.Pairs[1].Value.Primitive != "number" {
		t.Error("Expected second property value 'number'")
	}
}

func TestParse_Object(t *testing.T) {
	input := `{ name = string age = number }`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Expr.Object == nil {
		t.Fatal("Expected AstObject, got nil")
	}

	obj := result.Expr.Object
	if len(obj.Pairs) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(obj.Pairs))
	}

	if obj.Pairs[0].Key != "name" || obj.Pairs[1].Key != "age" {
		t.Error("Expected properties 'name' and 'age'")
	}
}

func TestParse_WithDocLineComments(t *testing.T) {
	input := `{
		/// This is a name field
		/// It stores the user's name
		name = string
		age = number
	}`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	obj := result.Expr.Object
	if obj == nil || len(obj.Pairs) != 2 {
		t.Fatal("Expected AstObject with 2 properties")
	}

	// Check documentation
	doc := obj.Pairs[0].Doc
	if doc == nil {
		t.Fatal("Expected documentation on first property")
	}

	if len(doc.Lines) != 2 {
		t.Errorf("Expected 2 doc lines, got %d", len(doc.Lines))
	}

	expectedLines := []string{
		"This is a name field",
		"It stores the user's name",
	}

	for i, expected := range expectedLines {
		if string(doc.Lines[i]) != expected {
			t.Errorf("Expected doc line %d to be '%s', got '%s'", i, expected, string(doc.Lines[i]))
		}
	}
}

func TestParse_WithDocBlockComment(t *testing.T) {
	input := `{
		/**
		 * This is a block comment
		 * It describes the name field
		 * @since 1.0.0
		 */
		name = string
	}`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	obj := result.Expr.Object
	if obj == nil || len(obj.Pairs) != 1 {
		t.Fatal("Expected AstObject with 1 property")
	}

	doc := obj.Pairs[0].Doc
	if doc == nil || doc.Block == nil {
		t.Fatal("Expected block documentation")
	}

	expectedContent := "This is a block comment\nIt describes the name field\n@since 1.0.0"
	if string(*doc.Block) != expectedContent {
		t.Errorf("Expected block content '%s', got '%s'", expectedContent, string(*doc.Block))
	}
}

func TestParse_PrimitiveTypes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		field    string
	}{
		{"string", "string", "Primitive"},
		{"number", "number", "Primitive"},
		{"bool", "bool", "Primitive"},
		{`"hello world"`, `"hello world"`, "String"},
		{"42", "42", "Number"},
		{"-3.14", "-3.14", "Number"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, err := Parse(test.input)
			if err != nil {
				t.Fatalf("Parse failed for '%s': %v", test.input, err)
			}

			switch test.field {
			case "Primitive":
				if result.Expr.Primitive == nil || *result.Expr.Primitive != test.expected {
					t.Errorf("Expected primitive '%s', got %v", test.expected, result.Expr.Primitive)
				}
			case "String":
				if result.Expr.String == nil || *result.Expr.String != test.expected {
					t.Errorf("Expected string '%s', got %v", test.expected, result.Expr.String)
				}
			case "Number":
				if result.Expr.Number == nil || *result.Expr.Number != test.expected {
					t.Errorf("Expected number '%s', got %v", test.expected, result.Expr.Number)
				}
			}
		})
	}
}

func TestParse_ComplexFunction(t *testing.T) {
	input := `optional(number, 0)`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	fn := result.Expr.Func
	if fn == nil || fn.Name != "optional" {
		t.Fatal("Expected AstFunction 'optional'")
	}

	if len(fn.Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(fn.Args))
	}

	// First arg should be primitive "number"
	if fn.Args[0].Primitive == nil || *fn.Args[0].Primitive != "number" {
		t.Error("Expected first argument to be primitive 'number'")
	}

	// Second arg should be number "0"
	if fn.Args[1].Number == nil || *fn.Args[1].Number != "0" {
		t.Error("Expected second argument to be number '0'")
	}
}

func TestParse_EmptyObject(t *testing.T) {
	input := `{}`

	result, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Expr.Object == nil {
		t.Fatal("Expected AstObject")
	}

	if len(result.Expr.Object.Pairs) != 0 {
		t.Errorf("Expected empty AstObject, got %d pairs", len(result.Expr.Object.Pairs))
	}
}

func TestParse_InvalidSyntax(t *testing.T) {
	invalidInputs := []string{
		`map(`,              // Unclosed AstFunction
		`{ name = }`,        // Missing value
		`{ = string }`,      // Missing key
		`map AstObject({})`, // Missing parentheses
		`{,}`,               // Invalid comma placement
	}

	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			_, err := Parse(input)
			if err == nil {
				t.Errorf("Expected error for invalid input '%s', but got none", input)
			}
		})
	}
}

func TestDocString_Capture(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/// Simple comment", "Simple comment"},
		{"///No space", "No space"},
		{"///   Multiple spaces", "Multiple spaces"},
		{"///\t\tTab indented", "Tab indented"},
		{"///", ""},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			var docStr AstDocString
			err := docStr.Capture([]string{test.input})
			if err != nil {
				t.Fatalf("Capture failed: %v", err)
			}

			if string(docStr) != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, string(docStr))
			}
		})
	}
}

func TestDocBlockString_Capture(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple block",
			input:    "/** Simple block comment */",
			expected: "Simple block comment",
		},
		{
			name: "Multi-line block",
			input: `/**
 * First line
 * Second line
 * Third line
 */`,
			expected: "First line\nSecond line\nThird line",
		},
		{
			name: "Block with mixed formatting",
			input: `/**
  First line without asterisk
 * Second line with asterisk
   * Third line with indented asterisk
 */`,
			expected: "First line without asterisk\nSecond line with asterisk\nThird line with indented asterisk",
		},
		{
			name:     "Empty block",
			input:    "/** */",
			expected: "",
		},
		{
			name: "Block with empty lines",
			input: `/**
 * First line
 *
 * Third line
 */`,
			expected: "First line\n\nThird line",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var docBlock AstDocBlockString
			err := docBlock.Capture([]string{test.input})
			if err != nil {
				t.Fatalf("Capture failed: %v", err)
			}

			if string(docBlock) != test.expected {
				t.Errorf("Expected '%s', got '%s'", test.expected, string(docBlock))
			}
		})
	}
}

func TestDocString_Capture_EmptyInput(t *testing.T) {
	var docStr AstDocString
	err := docStr.Capture([]string{})
	if err != nil {
		t.Errorf("Expected no error for empty input, got %v", err)
	}

	if string(docStr) != "" {
		t.Errorf("Expected empty string, got '%s'", string(docStr))
	}
}

func TestDocBlockString_Capture_EmptyInput(t *testing.T) {
	var docBlock AstDocBlockString
	err := docBlock.Capture([]string{})
	if err != nil {
		t.Errorf("Expected no error for empty input, got %v", err)
	}

	if string(docBlock) != "" {
		t.Errorf("Expected empty string, got '%s'", string(docBlock))
	}
}
