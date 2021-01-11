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
		Range util.Range
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

	TRUE  Type = "true"
	FALSE Type = "false"
	IF    Type = "if"
	FOR   Type = "for"
	IN    Type = "in"
	SET   Type = "set"
	TO    Type = "to"

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
