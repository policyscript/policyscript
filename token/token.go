package token

import "github.com/policyscript/policyscript/util"

type (

	// Type is the identifier for the token.
	Type string

	// Token is a unit lexed from source code.
	Token struct {

		// Type is the Type.
		Type Type

		// Literal is the string literal of the token.
		Literal string

		// Range is the positional range of the token within the source.
		Range *util.Range
	}
)

const (
	ILLEGAL Type = "illegal"
	COMMENT Type = "comment"
	EOF     Type = "EOF"

	// Identifiers

	IDENT Type = "identifier"

	// Literals.

	INTEGER Type = "integer"
	DECIMAL Type = "decimal"
	MONEY   Type = "money"
	PERIOD  Type = "period"
	PERCENT Type = "percent"
	TEXT    Type = "text"
	DATE    Type = "date"
	TIME    Type = "time"

	// Documentation literals.

	HEADING   Type = "heading"
	PARAGRAPH Type = "paragraph"

	// Operators.

	EQ     Type = "="
	NOT_EQ Type = "!="
	PLUS   Type = "+"
	MINUS  Type = "-"
	MULT   Type = "*"
	DIV    Type = "/"
	LT     Type = "<"
	GT     Type = ">"
	GT_EQ  Type = ">="
	LT_EQ  Type = "<="

	// Delimiters.

	LPAREN Type = "("
	RPAREN Type = ")"
	LBRACE Type = "{"
	RBRACE Type = "}"
	COLON  Type = ":"

	// Keywords.
	IF    Type = "if"
	FOR   Type = "for"
	IN    Type = "in"
	SET   Type = "set"
	TO    Type = "to"
	TRUE  Type = "true"
	FALSE Type = "false"
	AND   Type = "and"
	OR    Type = "or"
	LIST  Type = "list"

	// Block keywords.

	META    Type = "@meta"
	DEFINE  Type = "@define"
	ENUM    Type = "@enum"
	INPUTS  Type = "@inputs"
	OUTPUTS Type = "@outputs"
	LOCALS  Type = "@locals"
	CODE    Type = "@code"

	// Controls.

	SCOPE_START Type = "start"
	SCOPE_END   Type = "end"
	LINE_END    Type = ";"
)

var blockKeywords = map[string]Type{
	"@meta":    META,
	"@define":  DEFINE,
	"@enum":    ENUM,
	"@inputs":  INPUTS,
	"@outputs": OUTPUTS,
	"@locals":  LOCALS,
	"@code":    CODE,
}

// LookupBlockKeyword will return the block keyword and true, or ILLEGAL and
// false if not found.
func LookupBlockKeyword(input []rune) (Type, bool) {
	if kw, ok := blockKeywords[string(input)]; ok {
		return kw, true
	}
	return ILLEGAL, false
}

// Valid periods
var validPeriods = map[string]bool{
	"year":    true,
	"years":   true,
	"month":   true,
	"months":  true,
	"day":     true,
	"days":    true,
	"hour":    true,
	"hours":   true,
	"minute":  true,
	"minutes": true,
	"second":  true,
	"seconds": true,
}

// LookupPeriodKeyword will return true if a period keyword is valid.
func LookupPeriodKeyword(input []rune) bool {
	if _, ok := validPeriods[string(input)]; ok {
		return true
	}
	return false
}

var keywords = map[string]Type{
	"if":    IF,
	"for":   FOR,
	"in":    IN,
	"set":   SET,
	"to":    TO,
	"true":  TRUE,
	"false": FALSE,
	"and":   AND,
	"or":    OR,
	"list":  LIST,
}

// LookupIdent will return the keyword, or IDENT which is any other alpha-
// numeric string.
func LookupIdent(input []rune) Type {
	if kw, ok := keywords[string(input)]; ok {
		return kw
	}
	return IDENT
}
