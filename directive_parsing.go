package tfdocextras

import (
	"regexp"
	"strings"
)

var (
	quoteAndUrlRe   = regexp.MustCompile(`^"([^"]+)"\s+(.+)$`)
	braceAndUrlRe   = regexp.MustCompile(`^\{([^}]+)}\s+(.+)$`)
	enumDelimiterRe = regexp.MustCompile(`\s*\|\s*`)
)

type DirectiveType int

const (
	DirUnsupported DirectiveType = iota
	DirDeprecated
	DirExample
	DirEnum
	DirLink
	DirRegex
	DirSee
	DirSince
	DirType
)

const (
	IsValid byte = 1 << iota
	IsInvalid
	IsNamedLink
	IsReferenceLink
)

type ParsedDirective struct {
	Type DirectiveType
	// Args is the list of arguments for the directive, delimited by spaces while
	// respecting quoted substrings as a single argument.
	Args  []string
	Flags byte
}

func newBasicDirective(dt DirectiveType, content string) ParsedDirective {
	return ParsedDirective{
		Type:  dt,
		Args:  []string{content},
		Flags: IsValid,
	}
}

func newInvalidDirective(dt DirectiveType) ParsedDirective {
	return ParsedDirective{
		Type:  dt,
		Args:  []string{},
		Flags: IsInvalid,
	}
}

func parseExampleDirective(line string) ParsedDirective {
	if matches := quoteAndUrlRe.FindStringSubmatch(line); len(matches) == 3 {
		return ParsedDirective{
			Type:  DirExample,
			Args:  matches[1:],
			Flags: IsValid,
		}
	}

	return newInvalidDirective(DirExample)
}

func parseEnumDirective(line string) ParsedDirective {
	choices := enumDelimiterRe.Split(line, -1)

	return ParsedDirective{
		Type:  DirEnum,
		Args:  choices,
		Flags: IsValid,
	}
}

func parseLinkDirective(line string) ParsedDirective {
	if strings.HasPrefix(line, "\"") {
		if matches := quoteAndUrlRe.FindStringSubmatch(line); len(matches) == 3 {
			return ParsedDirective{
				Type:  DirLink,
				Args:  matches[1:],
				Flags: IsValid | IsNamedLink,
			}
		}
	}

	if strings.HasPrefix(line, "{") {
		if matches := braceAndUrlRe.FindStringSubmatch(line); len(matches) == 3 {
			return ParsedDirective{
				Type:  DirLink,
				Args:  matches[1:],
				Flags: IsValid | IsReferenceLink,
			}
		}
	}

	return newInvalidDirective(DirLink)
}

func parseRegexDirective(line string) ParsedDirective {
	if !strings.HasPrefix(line, "/") {
		return newInvalidDirective(DirRegex)
	}

	closingSlashIdx := -1
	for i := 1; i < len(line); i++ {
		if line[i] == '/' && (i == 1 || line[i-1] != '\\') {
			closingSlashIdx = i
			break
		}
	}

	if closingSlashIdx <= 0 {
		return newInvalidDirective(DirRegex)
	}

	pattern := line[1:closingSlashIdx]
	if _, err := regexp.Compile(pattern); err != nil {
		return newInvalidDirective(DirRegex)
	}

	args := []string{pattern}
	remaining := strings.TrimSpace(line[closingSlashIdx+1:])
	if remaining != "" {
		args = append(args, splitBySpacePreserveQuotes(remaining)...)
	}

	return ParsedDirective{
		Type:  DirRegex,
		Args:  args,
		Flags: IsValid,
	}
}

// splitBySpacePreserveQuotes splits a string by spaces while preserving quoted
// substrings. For example,
//
//	arg1 "hello world" arg3
//
// would be split into:
//
//	["arg1", "hello world", "arg3"]
func splitBySpacePreserveQuotes(s string) []string {
	quoted := false

	return strings.FieldsFunc(s, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}

		return !quoted && r == ' '
	})
}

// parseDirective parses a directive from its name and associated content,
// returning a ParsedDirective encapsulating the details.
func parseDirective(name string, line string) ParsedDirective {
	line = strings.TrimSpace(line)

	switch name {
	case "link":
		return parseLinkDirective(line)
	case "enum":
		return parseEnumDirective(line)
	case "example":
		return parseExampleDirective(line)
	case "regex":
		return parseRegexDirective(line)
	case "type":
		return newBasicDirective(DirType, line)
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
