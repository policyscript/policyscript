package ast

import (
	"github.com/policyscript/policyscript/token"
	"github.com/policyscript/policyscript/util"
)

type (
	// Node is implemented by all ast nodes.
	Node interface {
		Range() *util.Range
	}

	// Stmt is implemented by all statement nodes.
	Stmt interface {
		Node
		statementNode()
	}

	// Expr is implemented by all expression nodes.
	Expr interface {
		Node
		expressionNode()
	}
)

// The Program node is the high level node containing the entire program.
type Program struct {
	Stmts []Stmt
}

func (p *Program) Range() *util.Range {
	if len(p.Stmts) == 0 {
		return &util.Range{}
	}
	return &util.Range{
		Start: p.Stmts[0].Range().Start,
		End:   p.Stmts[len(p.Stmts)-1].Range().End,
	}
}

/* --- Statements --- */

// The HeadingStatement node.
type ExpressionStatement struct {
	Token token.Token
	Expr  Expr
}

func (s *ExpressionStatement) statementNode()     {}
func (s *ExpressionStatement) Range() *util.Range { return &s.Token.Range }

// The HeadingStatement node.
type HeadingStatement struct {
	Token token.Token
	Depth int
	Value string
}

func (s *HeadingStatement) statementNode()     {}
func (s *HeadingStatement) Range() *util.Range { return &s.Token.Range }

// The ParagraphStatement node.
type ParagraphStatement struct {
	Token token.Token
	Value string
}

func (s *ParagraphStatement) statementNode()     {}
func (s *ParagraphStatement) Range() *util.Range { return &s.Token.Range }

// The CommentStatement node.
type CommentStatement struct {
	Token token.Token
	Value string
}

func (s *CommentStatement) statementNode()     {}
func (s *CommentStatement) Range() *util.Range { return &s.Token.Range }

// The BlockStatement node.
type BlockStatement struct {
	Token token.Token
	Ident *Identifier // nil unless @define
	end   token.Token
	Stmts []Stmt
}

func (s *BlockStatement) statementNode() {}
func (s *BlockStatement) Range() *util.Range {
	return &util.Range{Start: s.Token.Range.Start, End: s.end.Range.End}
}

// The ScopeStatement node.
type ScopeStatement struct {
	Stmts []Stmt
}

func (s *ScopeStatement) statementNode() {}
func (s *ScopeStatement) Range() *util.Range {
	if len(s.Stmts) == 0 {
		return &util.Range{}
	}
	return &util.Range{
		Start: s.Stmts[0].Range().Start,
		End:   s.Stmts[len(s.Stmts)-1].Range().End,
	}
}

// The IfStatement node.
type IfStatement struct {
	Token     token.Token
	Condition Expr
	Block     *ScopeStatement
}

func (s *IfStatement) statementNode() {}
func (s *IfStatement) Range() *util.Range {
	return &util.Range{Start: s.Token.Range.Start, End: s.Block.Range().End}
}

// The ElseStatement node.
type ElseStatement struct {
	Token     token.Token
	Condition Expr
	Block     *ScopeStatement
}

func (s *ElseStatement) statementNode() {}
func (s *ElseStatement) Range() *util.Range {
	return &util.Range{Start: s.Token.Range.Start, End: s.Block.Range().End}
}

// The ForStatement node.
type ForStatement struct {
	Token token.Token
	Ident *Identifier
	Iter  Expr
	Block *ScopeStatement
}

func (s *ForStatement) statementNode() {}
func (s *ForStatement) Range() *util.Range {
	return &util.Range{Start: s.Token.Range.Start, End: s.Block.Range().End}
}

/* --- Expressions -- */

// The Identifier node.
type Identifier struct {
	Token token.Token
	Value string
}

func (e *Identifier) expressionNode()    {}
func (e *Identifier) Range() *util.Range { return &e.Token.Range }

// The PrefixExpression node.
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expr
}

func (e *PrefixExpression) expressionNode() {}
func (e *PrefixExpression) Range() *util.Range {
	return &util.Range{Start: e.Token.Range.Start, End: e.Right.Range().End}
}

// The InfixExpression node.
type InfixExpression struct {
	token    token.Token
	Left     Expr
	Operator string
	Right    Expr
}

func (e *InfixExpression) expressionNode() {}
func (e *InfixExpression) Range() *util.Range {
	return &util.Range{Start: e.Left.Range().Start, End: e.Right.Range().End}
}

// The DeclareExpression node.
type DeclareExpression struct {
	Token token.Token
	Ident *Identifier
	Value Expr
}

func (e *DeclareExpression) expressionNode() {}
func (e *DeclareExpression) Range() *util.Range {
	return &util.Range{Start: e.Ident.Range().Start, End: e.Value.Range().End}
}

// The SetExpression node.
type SetExpression struct {
	Token token.Token
	Ident *Identifier
	Value Expr
}

func (e *SetExpression) expressionNode() {}
func (e *SetExpression) Range() *util.Range {
	return &util.Range{Start: e.Token.Range.Start, End: e.Value.Range().End}
}

// The Condition node.
type Condition struct {
	Token token.Token
	Value bool
}

func (e *Condition) expressionNode()    {}
func (e *Condition) Range() *util.Range { return &e.Token.Range }

// The TextLiteral node.
type TextLiteral struct {
	Token token.Token
	Value string
}

func (e *TextLiteral) expressionNode()    {}
func (e *TextLiteral) Range() *util.Range { return &e.Token.Range }

// The IntegerLiteral node.
type IntegerLiteral struct {
	Token token.Token
	Value int
}

func (e *IntegerLiteral) expressionNode()    {}
func (e *IntegerLiteral) Range() *util.Range { return &e.Token.Range }

// The DecimalLiteral node.
type DecimalLiteral struct {
	Token token.Token
	Value float32
}

func (e *DecimalLiteral) expressionNode()    {}
func (e *DecimalLiteral) Range() *util.Range { return &e.Token.Range }

// The MoneyLiteral node.
type MoneyLiteral struct {
	Token  token.Token
	Value  float32
	Symbol string
}

func (e *MoneyLiteral) expressionNode()    {}
func (e *MoneyLiteral) Range() *util.Range { return &e.Token.Range }

// The PercentLiteral node.
type PercentLiteral struct {
	Token token.Token
	Value float32
}

func (e *PercentLiteral) expressionNode()    {}
func (e *PercentLiteral) Range() *util.Range { return &e.Token.Range }

// The PeriodLiteral node.
type PeriodLiteral struct {
	Token  token.Token
	Value  int
	Symbol string
}

func (e *PeriodLiteral) expressionNode()    {}
func (e *PeriodLiteral) Range() *util.Range { return &e.Token.Range }

// The DateLiteral node.
type DateLiteral struct {
	Token token.Token
	Year  int
	Month int
	Day   int
}

func (e *DateLiteral) expressionNode()    {}
func (e *DateLiteral) Range() *util.Range { return &e.Token.Range }

// The TimeLiteral node.
type TimeLiteral struct {
	Token   token.Token
	Hours   int
	Minutes int
	Seconds int
}

func (e *TimeLiteral) expressionNode()    {}
func (e *TimeLiteral) Range() *util.Range { return &e.Token.Range }
