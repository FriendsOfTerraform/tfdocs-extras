package tfdocextras

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// astDocString represents a documentation line with the `///` prefix stripped
type astDocString string

func (d *astDocString) Capture(values []string) error {
	if len(values) > 0 {
		stripped := strings.TrimPrefix(values[0], "///")
		stripped = strings.TrimLeft(stripped, " \t")
		*d = astDocString(stripped)
	}
	return nil
}

// astDocBlockString represents a documentation block with `/**`, `*/`, and `*` prefixes stripped
type astDocBlockString string

func (d *astDocBlockString) Capture(values []string) error {
	if len(values) > 0 {
		content := values[0]
		content = strings.TrimPrefix(content, "/**")
		content = strings.TrimSuffix(content, "*/")

		lines := strings.Split(content, "\n")
		var cleanLines []string

		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "*") {
				line = strings.TrimPrefix(line, "*")
				line = strings.TrimLeft(line, " \t")
			}

			cleanLines = append(cleanLines, line)
		}

		result := strings.Join(cleanLines, "\n")
		result = strings.Trim(result, "\n")

		*d = astDocBlockString(result)
	}
	return nil
}

type astDocBlock struct {
	Lines []astDocString     `  @DocLine+`
	Block *astDocBlockString `| @DocBlock`
}

type astDataType struct {
	Func      *astFunction `  @@`
	Object    *astObject   `| @@`
	Primitive *string      `| @Ident`
	Number    *string      `| @Number`
	String    *string      `| @String`
}

type astObjectProperty struct {
	Doc   *astDocBlock `@@?`
	Key   string       `@Ident`
	Value *astDataType `"=" @@`
}

type astObject struct {
	Pairs []*astObjectProperty `"{" @@* "}"`
}

type astFunction struct {
	Name string         `@Ident`
	Args []*astDataType `"(" ( @@ ( "," @@ )* )? ")"`
}

type astRoot struct {
	Expr *astDataType `@@`
}

func parseAst(str string) (*astRoot, error) {
	parser, err := participle.Build[astRoot](
		participle.Lexer(lexer.MustSimple([]lexer.SimpleRule{
			{"DocBlock", `/\*\*([^*]|\*+[^*/])*\*+/`},
			{"DocLine", `///[^\n]*`},
			{"Comment", `#[^\n]*`},
			{"String", `"([^"\\]|\\.)*"`},
			{"Number", `-?\d+(\.\d+)?`},
			{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`},
			{"Punct", `[\(\)\{\}=,]`},
			{"Whitespace", `[ \t\r\n]+`},
		})),
		participle.Elide("Whitespace", "Comment"),
	)
	if err != nil {
		return nil, err
	}

	ast, err := parser.ParseString("", str)
	if err != nil {
		return nil, err
	}

	return ast, nil
}
