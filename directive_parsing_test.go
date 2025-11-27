package tfdocextras

import "testing"

func TestParseDirective_NamedLinkDirective(t *testing.T) {
	link := "\"Supported AWS service endpoints\" https://docs.aws.amazon.com/general/latest/gr/aws-service-information.html"

	expected := ParsedDirective{
		Type:   DirLink,
		First:  "Supported AWS service endpoints",
		Second: "https://docs.aws.amazon.com/general/latest/gr/aws-service-information.html",
		Flags:  IsValid | IsNamedLink,
	}

	actual := ParseDirective("link", link)

	if actual != expected {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_ReferenceLinkDirective(t *testing.T) {
	link := "{route53-routing-policy-failover} https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy-failover.html"

	expected := ParsedDirective{
		Type:   DirLink,
		First:  "route53-routing-policy-failover",
		Second: "https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy-failover.html",
		Flags:  IsValid | IsReferenceLink,
	}

	actual := ParseDirective("link", link)

	if actual != expected {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_ExampleDirective(t *testing.T) {
	link := "\"Failover Routing Policy Example\" #failover-routing-policy"

	expected := ParsedDirective{
		Type:   DirExample,
		First:  "Failover Routing Policy Example",
		Second: "#failover-routing-policy",
		Flags:  IsValid,
	}

	actual := ParseDirective("example", link)

	if actual != expected {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}
