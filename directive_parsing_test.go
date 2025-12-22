package tfdocextras

import (
	"testing"

	"github.com/go-test/deep"
)

func TestParseDirective_NamedLinkDirective(t *testing.T) {
	link := "\"Supported AWS service endpoints\" https://docs.aws.amazon.com/general/latest/gr/aws-service-information.html"

	expected := ParsedDirective{
		Type: DirLink,
		Args: []string{
			"Supported AWS service endpoints",
			"https://docs.aws.amazon.com/general/latest/gr/aws-service-information.html",
		},
		Flags: IsValid | IsNamedLink,
	}

	actual := ParseDirective("link", link)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_ReferenceLinkDirective(t *testing.T) {
	link := "{route53-routing-policy-failover} https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy-failover.html"

	expected := ParsedDirective{
		Type: DirLink,
		Args: []string{
			"route53-routing-policy-failover",
			"https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy-failover.html",
		},
		Flags: IsValid | IsReferenceLink,
	}

	actual := ParseDirective("link", link)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_ExampleDirective(t *testing.T) {
	link := "\"Failover Routing Policy Example\" #failover-routing-policy"

	expected := ParsedDirective{
		Type: DirExample,
		Args: []string{

			"Failover Routing Policy Example",
			"#failover-routing-policy",
		},
		Flags: IsValid,
	}

	actual := ParseDirective("example", link)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_EnumDirective(t *testing.T) {
	raw := "value1|value2|value3"

	expected := ParsedDirective{
		Type: DirEnum,
		Args: []string{
			"value1",
			"value2",
			"value3",
		},
		Flags: IsValid,
	}

	actual := ParseDirective("enum", raw)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_EnumDirectiveWithSpaces(t *testing.T) {
	raw := "value1 | value2 | value3"

	expected := ParsedDirective{
		Type: DirEnum,
		Args: []string{
			"value1",
			"value2",
			"value3",
		},
		Flags: IsValid,
	}

	actual := ParseDirective("enum", raw)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_RegexDirective(t *testing.T) {
	raw := "/^[a-zA-Z0-9_-]{5}$/ abcd4 efgh_ ijkl-"

	expected := ParsedDirective{
		Type: DirRegex,
		Args: []string{
			"^[a-zA-Z0-9_-]{5}$",
			"abcd4",
			"efgh_",
			"ijkl-",
		},
		Flags: IsValid,
	}

	actual := ParseDirective("regex", raw)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_RegexDirectiveWithSpaces(t *testing.T) {
	raw := "/\\w+ \\w+/ \"hello world\" \"foo bar\""

	expected := ParsedDirective{
		Type: DirRegex,
		Args: []string{
			"\\w+ \\w+",
			"\"hello world\"",
			"\"foo bar\"",
		},
		Flags: IsValid,
	}

	actual := ParseDirective("regex", raw)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_RegexWithSlash(t *testing.T) {
	raw := "/^https?:\\/\\// https://example.com http://test.com"

	expected := ParsedDirective{
		Type: DirRegex,
		Args: []string{
			"^https?:\\/\\/",
			"https://example.com",
			"http://test.com",
		},
		Flags: IsValid,
	}

	actual := ParseDirective("regex", raw)

	if diff := deep.Equal(expected, actual); diff != nil {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}
