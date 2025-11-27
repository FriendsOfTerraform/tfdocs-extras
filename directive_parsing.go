package tfdocextras

import (
	"regexp"
	"strings"
)

var (
	quoteAndUrlPattern = regexp.MustCompile(`^"([^"]+)"\s+(.+)$`)
	braceAndUrlPattern = regexp.MustCompile(`^\{([^}]+)}\s+(.+)$`)
)

type DirectiveType int

const (
	DirUnsupported DirectiveType = iota
	DirDeprecated
	DirExample
	DirLink
	DirRegex
	DirSee
	DirSince
)

const (
	IsValid byte = 1 << iota
	IsInvalid
	IsNamedLink
	IsReferenceLink
)

type ParsedDirective struct {
	Type   DirectiveType
	First  string
	Second string
	Flags  byte
}

func ParseDirective(name string, line string) ParsedDirective {
	line = strings.TrimSpace(line)

	switch name {
	case "link":
		return parseLinkDirective(line)
	case "example":
		return parseExampleDirective(line)
	case "regex":
		return parseRegexDirective(line)
	case "deprecated":
		return newBasicDirective(DirDeprecated, line)
	case "see":
		return newBasicDirective(DirSee, line)
	case "since":
		return newBasicDirective(DirSince, line)
	default:
		return newInvalidDirective(DirUnsupported)
	}
}

func newBasicDirective(dt DirectiveType, content string) ParsedDirective {
	return ParsedDirective{
		Type:   dt,
		First:  content,
		Second: "",
		Flags:  IsValid,
	}
}

func newInvalidDirective(dt DirectiveType) ParsedDirective {
	return ParsedDirective{
		Type:   dt,
		First:  "",
		Second: "",
		Flags:  IsInvalid,
	}
}

func parseExampleDirective(line string) ParsedDirective {
	if matches := quoteAndUrlPattern.FindStringSubmatch(line); len(matches) == 3 {
		return ParsedDirective{
			Type:   DirExample,
			First:  matches[1],
			Second: matches[2],
			Flags:  IsValid,
		}
	}

	return newInvalidDirective(DirExample)
}

func parseLinkDirective(line string) ParsedDirective {
	if strings.HasPrefix(line, "\"") {
		if matches := quoteAndUrlPattern.FindStringSubmatch(line); len(matches) == 3 {
			return ParsedDirective{
				Type:   DirLink,
				First:  matches[1],
				Second: matches[2],
				Flags:  IsValid | IsNamedLink,
			}
		}
	}

	if strings.HasPrefix(line, "{") {
		if matches := braceAndUrlPattern.FindStringSubmatch(line); len(matches) == 3 {
			return ParsedDirective{
				Type:   DirLink,
				First:  matches[1],
				Second: matches[2],
				Flags:  IsValid | IsReferenceLink,
			}
		}
	}

	return newInvalidDirective(DirLink)
}

func parseRegexDirective(line string) ParsedDirective {
	if strings.HasPrefix(line, "/") && strings.HasSuffix(line, "/") {
		pattern := line[1 : len(line)-1]

		if _, err := regexp.Compile(pattern); err == nil {
			return ParsedDirective{
				Type:   DirRegex,
				First:  pattern,
				Second: "",
				Flags:  IsValid,
			}
		}
	}

	return newInvalidDirective(DirRegex)
}
