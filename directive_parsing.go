package tfdocextras

import (
	"regexp"
	"strings"
)

// Compile regex patterns once at package initialization for better performance
var (
	quoteAndUrlPattern = regexp.MustCompile(`^"([^"]+)"\s+(.+)$`)
	braceAndUrlPattern = regexp.MustCompile(`^\{([^}]+)\}\s+(.+)$`)
)

type ExampleDirective struct {
	Name      string
	Reference string
}

type NamedLinkDirective struct {
	Name string
	URL  string
}

type ReferenceLinkDirective struct {
	Reference string
	URL       string
}

func ParseDirective(name string, line string) interface{} {
	line = strings.TrimSpace(line)

	switch name {
	case "link":
		return parseLinkDirective(line)
	case "example":
		return parseExampleDirective(line)
	default:
		return nil
	}
}

func parseLinkDirective(line string) interface{} {
	if strings.HasPrefix(line, "\"") {
		if matches := quoteAndUrlPattern.FindStringSubmatch(line); len(matches) == 3 {
			return NamedLinkDirective{
				Name: matches[1],
				URL:  matches[2],
			}
		}
	}

	if strings.HasPrefix(line, "{") {
		if matches := braceAndUrlPattern.FindStringSubmatch(line); len(matches) == 3 {
			return ReferenceLinkDirective{
				Reference: matches[1],
				URL:       matches[2],
			}
		}
	}

	return nil
}

func parseExampleDirective(line string) interface{} {
	if matches := quoteAndUrlPattern.FindStringSubmatch(line); len(matches) == 3 {
		return ExampleDirective{
			Name:      matches[1],
			Reference: matches[2],
		}
	}
	return nil
}
