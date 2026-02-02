package tfdocextras

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/terraform-docs/terraform-docs/terraform"
)

func TestParseModuleArgsIntoManifest_OutputsWithoutDescription(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name:        "vpc_id",
			Description: "",
			Sensitive:   false,
		},
		{
			Name:        "api_key",
			Description: "",
			Sensitive:   true,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 2); diff != nil {
		t.Errorf("Output rows count mismatch:\n%v", diff)
	}

	expectedFirstOutput := Argument{
		Type:        "unknown",
		Name:        "vpc_id",
		Description: "",
		Sensitive:   false,
		ArgumentMetadata: ArgumentMetadata{
			Attributes:    []MetadataAttribute{},
			Enumerations:  []string{},
			Examples:      []MetadataAttribute{},
			Links:         []MetadataAttribute{},
			RegexPattern:  "",
			RegexExamples: []string{},
		},
	}

	if diff := deep.Equal(manifest.Outputs.Rows[0], expectedFirstOutput); diff != nil {
		t.Errorf("First output mismatch:\n%v", diff)
	}

	expectedSecondOutput := Argument{
		Type:        "unknown",
		Name:        "api_key",
		Description: "",
		Sensitive:   true,
		ArgumentMetadata: ArgumentMetadata{
			Attributes:    []MetadataAttribute{},
			Enumerations:  []string{},
			Examples:      []MetadataAttribute{},
			Links:         []MetadataAttribute{},
			RegexPattern:  "",
			RegexExamples: []string{},
		},
	}

	if diff := deep.Equal(manifest.Outputs.Rows[1], expectedSecondOutput); diff != nil {
		t.Errorf("Second output mismatch:\n%v", diff)
	}
}

func TestParseModuleArgsIntoManifest_OutputsWithTypeDirective(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name: "vpc_arn",
			Description: `The ARN of the VPC
			
@since 1.0.0
@type string`,
			Sensitive: false,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 1); diff != nil {
		t.Fatalf("Output rows count mismatch:\n%v", diff)
	}

	if len(manifest.Outputs.Rows[0].ArgumentMetadata.Attributes) == 0 {
		t.Error("Expected @since directive to be processed into attributes")
	}

	expected := Argument{
		Type:        "string",
		Name:        "vpc_arn",
		Description: "The ARN of the VPC",
		Sensitive:   false,
		ArgumentMetadata: ArgumentMetadata{
			Attributes:    manifest.Outputs.Rows[0].ArgumentMetadata.Attributes, // Use actual value
			Enumerations:  []string{},
			Examples:      []MetadataAttribute{},
			Links:         []MetadataAttribute{},
			RegexPattern:  "",
			RegexExamples: []string{},
		},
	}

	if diff := deep.Equal(manifest.Outputs.Rows[0], expected); diff != nil {
		t.Errorf("Output mismatch:\n%v", diff)
	}
}

func TestParseModuleArgsIntoManifest_OutputsWithComplexType(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name: "nat_gateways",
			Description: `Map of NAT gateways

@since 1.0.0
@type map(object({
  /// The availability zone
  availability_zone = string
  
  /// The association ID
  association_id = string
}))`,
			Sensitive: false,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 1); diff != nil {
		t.Fatalf("Output rows count mismatch:\n%v", diff)
	}

	if _, exists := manifest.Objects["NatGateways"]; !exists {
		t.Error("Expected NatGateways nested object to be recorded")
	}

	complexTypeStr := "NatGateways"
	expectedOutput := Argument{
		Type:        "map(object(NatGateways))",
		ComplexType: &complexTypeStr,
		Name:        "nat_gateways",
		Description: "Map of NAT gateways",
		Sensitive:   false,
		ArgumentMetadata: ArgumentMetadata{
			Attributes:    manifest.Outputs.Rows[0].ArgumentMetadata.Attributes, // Use actual value
			Enumerations:  []string{},
			Examples:      []MetadataAttribute{},
			Links:         []MetadataAttribute{},
			RegexPattern:  "",
			RegexExamples: []string{},
		},
	}

	if diff := deep.Equal(manifest.Outputs.Rows[0], expectedOutput); diff != nil {
		t.Errorf("Output mismatch:\n%v", diff)
	}

	natGatewaysObj := manifest.Objects["NatGateways"]
	if diff := deep.Equal(len(natGatewaysObj.Rows), 2); diff != nil {
		t.Errorf("NatGateways field count mismatch:\n%v", diff)
	}

	expectedFields := []struct {
		name string
		typ  string
	}{
		{"availability_zone", "string"},
		{"association_id", "string"},
	}

	for i, expected := range expectedFields {
		actualField := Argument{
			Type: natGatewaysObj.Rows[i].Type,
			Name: natGatewaysObj.Rows[i].Name,
		}
		expectedField := Argument{
			Type: expected.typ,
			Name: expected.name,
		}
		if diff := deep.Equal(actualField, expectedField); diff != nil {
			t.Errorf("Field %d mismatch:\n%v", i, diff)
		}
	}
}

func TestParseModuleArgsIntoManifest_OutputsSensitiveField(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name:        "regular_output",
			Description: "A regular output",
			Sensitive:   false,
		},
		{
			Name:        "sensitive_output",
			Description: "A sensitive output",
			Sensitive:   true,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 2); diff != nil {
		t.Fatalf("Output rows count mismatch:\n%v", diff)
	}

	expectedOutputs := []Argument{
		{
			Type:        "unknown",
			Name:        "regular_output",
			Description: "A regular output",
			Sensitive:   false,
			ArgumentMetadata: ArgumentMetadata{
				Attributes:    []MetadataAttribute{},
				Enumerations:  []string{},
				Examples:      []MetadataAttribute{},
				Links:         []MetadataAttribute{},
				RegexPattern:  "",
				RegexExamples: []string{},
			},
		},
		{
			Type:        "unknown",
			Name:        "sensitive_output",
			Description: "A sensitive output",
			Sensitive:   true,
			ArgumentMetadata: ArgumentMetadata{
				Attributes:    []MetadataAttribute{},
				Enumerations:  []string{},
				Examples:      []MetadataAttribute{},
				Links:         []MetadataAttribute{},
				RegexPattern:  "",
				RegexExamples: []string{},
			},
		},
	}

	if diff := deep.Equal(manifest.Outputs.Rows, expectedOutputs); diff != nil {
		t.Errorf("Outputs mismatch:\n%v", diff)
	}
}

func TestParseModuleArgsIntoManifest_OutputsWithNestedObjects(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name: "config",
			Description: `Configuration settings

@type object({
  /// Server configuration
  server = object({
    /// Server host
    host = string
    
    /// Server port
    port = number
  })
  
  /// Database configuration
  database = object({
    /// Database name
    name = string
  })
})`,
			Sensitive: false,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 1); diff != nil {
		t.Fatalf("Output rows count mismatch:\n%v", diff)
	}

	expectedObjects := []string{"Config", "Server", "Database"}
	for _, objName := range expectedObjects {
		if _, exists := manifest.Objects[objName]; !exists {
			t.Errorf("Expected nested object '%s' to be recorded", objName)
		}
	}

	configObj := manifest.Objects["Config"]
	if diff := deep.Equal(len(configObj.Rows), 2); diff != nil {
		t.Errorf("Config field count mismatch:\n%v", diff)
	}

	serverComplexType := "Server"
	expectedServerField := Argument{
		Type:        "object(Server)",
		ComplexType: &serverComplexType,
		Name:        "server",
	}
	actualServerField := Argument{
		Type:        configObj.Rows[0].Type,
		ComplexType: configObj.Rows[0].ComplexType,
		Name:        configObj.Rows[0].Name,
	}
	if diff := deep.Equal(actualServerField, expectedServerField); diff != nil {
		t.Errorf("Server field mismatch:\n%v", diff)
	}

	databaseComplexType := "Database"
	expectedDatabaseField := Argument{
		Type:        "object(Database)",
		ComplexType: &databaseComplexType,
		Name:        "database",
	}
	actualDatabaseField := Argument{
		Type:        configObj.Rows[1].Type,
		ComplexType: configObj.Rows[1].ComplexType,
		Name:        configObj.Rows[1].Name,
	}
	if diff := deep.Equal(actualDatabaseField, expectedDatabaseField); diff != nil {
		t.Errorf("Database field mismatch:\n%v", diff)
	}
}

func TestParseModuleArgsIntoManifest_OutputsWithSensitiveFlag(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name:        "vpc_id",
			Description: "The ID of the VPC\n\n@type string",
			Sensitive:   false,
		},
		{
			Name:        "private_key",
			Description: "The private key\n\n@type string",
			Sensitive:   true,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 2); diff != nil {
		t.Errorf("Output rows count mismatch:\n%v", diff)
	}

	// Extract key fields for comparison
	type outputCheck struct {
		Name      string
		Sensitive bool
	}

	actualOutputs := []outputCheck{
		{Name: manifest.Outputs.Rows[0].Name, Sensitive: manifest.Outputs.Rows[0].Sensitive},
		{Name: manifest.Outputs.Rows[1].Name, Sensitive: manifest.Outputs.Rows[1].Sensitive},
	}

	expectedOutputs := []outputCheck{
		{Name: "vpc_id", Sensitive: false},
		{Name: "private_key", Sensitive: true},
	}

	if diff := deep.Equal(actualOutputs, expectedOutputs); diff != nil {
		t.Errorf("Outputs sensitive flags mismatch:\n%v", diff)
	}
}

func TestParseModuleArgsIntoManifest_OutputWithoutTypeDirective(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name:        "result",
			Description: "Some result value\n\n@since 1.0.0",
			Sensitive:   false,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 1); diff != nil {
		t.Fatalf("Output rows count mismatch:\n%v", diff)
	}

	output := manifest.Outputs.Rows[0]

	if len(output.ArgumentMetadata.Attributes) == 0 {
		t.Error("Expected @since directive to be processed")
	}

	type outputCheck struct {
		Type string
		Name string
	}

	actual := outputCheck{Type: output.Type, Name: output.Name}
	expected := outputCheck{Type: "unknown", Name: "result"}

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Errorf("Output mismatch:\n%v", diff)
	}
}

func TestParseModuleArgsIntoManifest_OutputWithInvalidTypeDirective(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name:        "bad_output",
			Description: "Bad type definition\n\n@type map(object({",
			Sensitive:   false,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 1); diff != nil {
		t.Fatalf("Output rows count mismatch:\n%v", diff)
	}

	output := manifest.Outputs.Rows[0]

	type outputCheck struct {
		Name            string
		HasNonEmptyType bool
	}

	actual := outputCheck{
		Name:            output.Name,
		HasNonEmptyType: output.Type != "",
	}
	expected := outputCheck{
		Name:            "bad_output",
		HasNonEmptyType: true,
	}

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Errorf("Output mismatch:\n%v", diff)
	}
}

func TestCreateArgumentFromDocBlock(t *testing.T) {
	manifest := newModuleManifest()

	docBlk := PropertyDocBlock{
		Content: []string{"Test description", "Second line"},
		Directives: []DocDirective{
			{
				Name:       "since",
				RawContent: "1.0.0",
				Parsed:     ParsedDirective{Type: DirSince, Args: []string{"1.0.0"}, Flags: IsValid},
			},
		},
	}

	arg := createArgumentFromDocBlock("test_arg", "string", docBlk, nil, manifest)

	if len(arg.ArgumentMetadata.Attributes) == 0 {
		t.Error("Expected @since directive to be processed into attributes")
	}

	type argCheck struct {
		Name        string
		Type        string
		Description string
	}

	actual := argCheck{
		Name:        arg.Name,
		Type:        arg.Type,
		Description: arg.Description,
	}
	expected := argCheck{
		Name:        "test_arg",
		Type:        "string",
		Description: "Test description\nSecond line",
	}

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Errorf("Argument mismatch:\n%v", diff)
	}
}

func TestFindTypeDirective(t *testing.T) {
	directives := []DocDirective{
		{Name: "since", RawContent: "1.0.0"},
		{Name: "type", RawContent: "string"},
		{Name: "deprecated", RawContent: "Use something else"},
	}

	result := findTypeDirective(directives)

	if result == nil {
		t.Fatal("Expected to find @type directive")
	}

	// Compare key fields
	type directiveCheck struct {
		Name       string
		RawContent string
	}

	actual := directiveCheck{
		Name:       result.Name,
		RawContent: result.RawContent,
	}
	expected := directiveCheck{
		Name:       "type",
		RawContent: "string",
	}

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Errorf("Directive mismatch:\n%v", diff)
	}
}

func TestFindTypeDirective_NotFound(t *testing.T) {
	directives := []DocDirective{
		{Name: "since", RawContent: "1.0.0"},
		{Name: "deprecated", RawContent: "Use something else"},
	}

	result := findTypeDirective(directives)

	var nilDirective *DocDirective
	if diff := deep.Equal(result, nilDirective); diff != nil {
		t.Errorf("Expected nil when @type directive not found:\n%v", diff)
	}
}

func TestFindTypeDirective_EmptyList(t *testing.T) {
	directives := []DocDirective{}

	result := findTypeDirective(directives)

	var nilDirective *DocDirective
	if diff := deep.Equal(result, nilDirective); diff != nil {
		t.Errorf("Expected nil for empty directives list:\n%v", diff)
	}
}
