package tfdocextras

import "testing"

func TestParseDirective_NamedLinkDirective(t *testing.T) {
	link := "\"Supported AWS service endpoints\" https://docs.aws.amazon.com/general/latest/gr/aws-service-information.html"

	expected := NamedLinkDirective{
		Name: "Supported AWS service endpoints",
		URL:  "https://docs.aws.amazon.com/general/latest/gr/aws-service-information.html",
	}

	actual := ParseDirective("link", link)

	if actual != expected {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_ReferenceLinkDirective(t *testing.T) {
	link := "{route53-routing-policy-failover} https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy-failover.html"

	expected := ReferenceLinkDirective{
		Reference: "route53-routing-policy-failover",
		URL:       "https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy-failover.html",
	}

	actual := ParseDirective("link", link)

	if actual != expected {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}

func TestParseDirective_ExampleDirective(t *testing.T) {
	link := "\"Failover Routing Policy Example\" #failover-routing-policy"

	expected := ExampleDirective{
		Name:      "Failover Routing Policy Example",
		Reference: "#failover-routing-policy",
	}

	actual := ParseDirective("example", link)

	if actual != expected {
		t.Errorf("Expected %+v, but got %+v", expected, actual)
	}
}
