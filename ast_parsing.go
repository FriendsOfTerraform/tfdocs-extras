package main

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// AstDocString represents a documentation line with the `///` prefix stripped
type AstDocString string

func (d *AstDocString) Capture(values []string) error {
	if len(values) > 0 {
		stripped := strings.TrimPrefix(values[0], "///")
		stripped = strings.TrimLeft(stripped, " \t")
		*d = AstDocString(stripped)
	}
	return nil
}

// AstDocBlockString represents a documentation block with `/**`, `*/`, and `*` prefixes stripped
type AstDocBlockString string

func (d *AstDocBlockString) Capture(values []string) error {
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

		*d = AstDocBlockString(result)
	}
	return nil
}

type AstDocBlock struct {
	Lines []AstDocString     `  @DocLine+`
	Block *AstDocBlockString `| @DocBlock`
}

type AstDataType struct {
	Func      *AstFunction `  @@`
	Object    *AstObject   `| @@`
	Primitive *string      `| @Ident`
	Number    *string      `| @Number`
	String    *string      `| @String`
}

type AstObjectProperty struct {
	Doc   *AstDocBlock `@@?`
	Key   string       `@Ident`
	Value *AstDataType `"=" @@`
}

type AstObject struct {
	Pairs []*AstObjectProperty `"{" @@* "}"`
}

type AstFunction struct {
	Name string         `@Ident`
	Args []*AstDataType `"(" ( @@ ( "," @@ )* )? ")"`
}

type AstRoot struct {
	Expr *AstDataType `@@`
}

func ParseAst(str string) (*AstRoot, error) {
	parser, err := participle.Build[AstRoot](
		participle.Lexer(lexer.MustSimple([]lexer.SimpleRule{
			{"DocBlock", `/\*\*([^*]|\*+[^*/])*\*+/`},
			{"DocLine", `///[^\n]*`},
			{"String", `"([^"\\]|\\.)*"`},
			{"Number", `-?\d+(\.\d+)?`},
			{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`},
			{"Punct", `[\(\)\{\}=,]`},
			{"Whitespace", `[ \t\r\n]+`},
		})),
		participle.Elide("Whitespace"),
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
