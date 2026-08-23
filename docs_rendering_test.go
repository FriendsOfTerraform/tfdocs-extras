package tfdocsextras

import (
	"sort"
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
		ParentType:  Output,
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
		ParentType:  Output,
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
		ParentType:  Output,
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

	if _, exists := manifest.Objects["output.nat_gateways"]; !exists {
		t.Error("Expected nat_gateways nested object to be recorded")
	}

	complexTypeStr := "output.nat_gateways"
	expectedOutput := Argument{
		ParentType:  Output,
		Type:        "map(object(output.nat_gateways))",
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

	natGatewaysObj := manifest.Objects["output.nat_gateways"]
	if diff := deep.Equal(len(natGatewaysObj.Rows), 2); diff != nil {
		t.Errorf("nat_gateways field count mismatch:\n%v", diff)
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
			ParentType:  Output,
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
			ParentType:  Output,
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

	expectedObjects := []string{"output.config", "output.config.server", "output.config.database"}
	for _, objName := range expectedObjects {
		if _, exists := manifest.Objects[objName]; !exists {
			t.Errorf("Expected nested object '%s' to be recorded", objName)
		}
	}

	configObj := manifest.Objects["output.config"]
	if diff := deep.Equal(len(configObj.Rows), 2); diff != nil {
		t.Errorf("config field count mismatch:\n%v", diff)
	}

	serverComplexType := "output.config.server"
	expectedServerField := Argument{
		Type:        "object(output.config.server)",
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

	databaseComplexType := "output.config.database"
	expectedDatabaseField := Argument{
		Type:        "object(output.config.database)",
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

func TestParseModuleArgsIntoManifest_OutputSensitiveSummary(t *testing.T) {
	outputs := []*terraform.Output{
		{
			Name:        "sensitive_summary",
			Description: "Sensitive summary value for output handling tests.\n\n@type string",
			Sensitive:   true,
		},
	}

	manifest := ParseModuleArgsIntoManifest([]*terraform.Input{}, outputs)

	if diff := deep.Equal(len(manifest.Outputs.Rows), 1); diff != nil {
		t.Fatalf("Output rows count mismatch:\n%v", diff)
	}

	output := manifest.Outputs.Rows[0]

	if !output.Sensitive {
		t.Fatalf("Expected output 'sensitive_summary' to be marked sensitive")
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

	arg := createArgumentFromDocBlock(RequiredInput, "test_arg", "string", docBlk, nil, "test_arg", manifest)

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

// TestParseModuleArgsIntoManifest_NestedObjectsWithTheSameName covers two
// inputs that each declare a nested object called `details`. Keyed by the
// field name alone they collapsed into one entry, so the second silently
// overwrote the first and the first input's table linked to the wrong type.
func TestParseModuleArgsIntoManifest_NestedObjectsWithTheSameName(t *testing.T) {
	inputs := []*terraform.Input{
		{
			Name:        "alpha",
			Description: "Alpha.",
			Required:    true,
			Type: `object({
  /// Alpha's details.
  details = optional(object({
    alpha_only = string
  }), null)
})`,
		},
		{
			Name:        "beta",
			Description: "Beta.",
			Required:    true,
			Type: `object({
  /// Beta's details.
  details = optional(object({
    beta_only = number
  }), null)
})`,
		},
	}

	manifest := ParseModuleArgsIntoManifest(inputs, []*terraform.Output{})

	expected := []string{"alpha", "alpha.details", "beta", "beta.details"}
	if diff := deep.Equal(objectTypeNames(manifest), expected); diff != nil {
		t.Fatalf("Object type names mismatch:\n%v", diff)
	}

	// Each nested object keeps its own fields.
	if diff := deep.Equal(fieldNames(manifest.Objects["alpha.details"]), []string{"alpha_only"}); diff != nil {
		t.Errorf("alpha.details field mismatch:\n%v", diff)
	}
	if diff := deep.Equal(fieldNames(manifest.Objects["beta.details"]), []string{"beta_only"}); diff != nil {
		t.Errorf("beta.details field mismatch:\n%v", diff)
	}

	// And each parent links to its own, in both the type and the anchor.
	details := manifest.Objects["alpha"].Rows[0]
	if diff := deep.Equal(*details.ComplexType, "alpha.details"); diff != nil {
		t.Errorf("alpha's details complex type mismatch:\n%v", diff)
	}
	if diff := deep.Equal(details.Type, "object(alpha.details)"); diff != nil {
		t.Errorf("alpha's details type mismatch:\n%v", diff)
	}
	if diff := deep.Equal(details.GetAnchor(), "alphadetails"); diff != nil {
		t.Errorf("alpha's details anchor mismatch:\n%v", diff)
	}
	if diff := deep.Equal(details.GetParentType(), [2]string{"object(", ")"}); diff != nil {
		t.Errorf("alpha's details parent type mismatch:\n%v", diff)
	}
}

// TestParseModuleArgsIntoManifest_InputAndOutputWithTheSameName covers a
// module that passes an object input straight back out under the same name.
// Both object types claimed the same entry, and the input's was the one lost.
func TestParseModuleArgsIntoManifest_InputAndOutputWithTheSameName(t *testing.T) {
	inputs := []*terraform.Input{
		{
			Name:        "listener",
			Description: "Listener input.",
			Required:    true,
			Type: `object({
  /// Input side port.
  port = number
})`,
		},
	}

	outputs := []*terraform.Output{
		{
			Name: "listener",
			Description: `Listener output.

@type object({
  /// Output side hostname.
  hostname = string
})`,
		},
	}

	manifest := ParseModuleArgsIntoManifest(inputs, outputs)

	expected := []string{"listener", "output.listener"}
	if diff := deep.Equal(objectTypeNames(manifest), expected); diff != nil {
		t.Fatalf("Object type names mismatch:\n%v", diff)
	}

	if diff := deep.Equal(fieldNames(manifest.Objects["listener"]), []string{"port"}); diff != nil {
		t.Errorf("Input object type mismatch:\n%v", diff)
	}
	if diff := deep.Equal(fieldNames(manifest.Objects["output.listener"]), []string{"hostname"}); diff != nil {
		t.Errorf("Output object type mismatch:\n%v", diff)
	}

	if diff := deep.Equal(*manifest.Outputs.Rows[0].ComplexType, "output.listener"); diff != nil {
		t.Errorf("Output complex type mismatch:\n%v", diff)
	}
	if diff := deep.Equal(*manifest.RequiredInputs.Rows[0].ComplexType, "listener"); diff != nil {
		t.Errorf("Input complex type mismatch:\n%v", diff)
	}
}

// TestAnchorFor pins the anchor to what a markdown renderer derives from the
// heading of the same name, which is what the object tables link to.
func TestAnchorFor(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"health_check", "health_check"},
		{"listener.health_check", "listenerhealth_check"},
		{"output.Config", "outputconfig"},
		{"deeply.nested.object", "deeplynestedobject"},
		{"with space", "with-space"},
		{"keeps-dashes", "keeps-dashes"},
	}

	for _, tt := range tests {
		if diff := deep.Equal(anchorFor(tt.name), tt.expected); diff != nil {
			t.Errorf("Anchor for %q mismatch:\n%v", tt.name, diff)
		}
	}
}

// objectTypeNames returns the recorded object type names, sorted, which is the
// order the Objects section renders them in.
func objectTypeNames(manifest *ModuleManifest) []string {
	names := make([]string, 0, len(manifest.Objects))
	for name := range manifest.Objects {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

func fieldNames(group ArgumentGroup) []string {
	names := make([]string, 0, len(group.Rows))
	for _, row := range group.Rows {
		names = append(names, row.Name)
	}

	return names
}
