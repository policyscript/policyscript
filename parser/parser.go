package parser

import (
	"fmt"

	"github.com/policyscript/policyscript/ast"
	"github.com/policyscript/policyscript/scanner"
	"github.com/policyscript/policyscript/token"
	"github.com/policyscript/policyscript/util"
)

type (
	prefixParseFn func() (right ast.Expr)
	infixParseFn  func(left ast.Expr) (right ast.Expr)
)

const (
	LOWEST      int = iota + 1
	ASSIGN          // : for type assignment
	OR              // or
	AND             // and
	EQUALS          // =
	LESSGREATER     // > or < or >= or <=
	SUM             // + or -
	PRODUCT         // * or /
	PREFIX          // -X
)

var precedences = map[token.Type]int{
	token.COLON:  ASSIGN,
	token.OR:     OR,
	token.AND:    AND,
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.LT_EQ:  LESSGREATER,
	token.GT_EQ:  LESSGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.DIV:    PRODUCT,
	token.MULT:   PRODUCT,
}

type Parser struct {
	s      scanner.Scanner
	errors util.ErrorList

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(s scanner.Scanner) *Parser {
	p := &Parser{s: s}

	// Read 2 tokens to set cur and peek token.
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.prefixParseFns[token.IDENT] = p.parseIdentifier

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.infixParseFns[token.COLON] = p.parseDeclare

	return p
}

func (p *Parser) Errors() util.ErrorList { return p.errors }

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for !p.curTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.curToken.Type {
	case token.COMMENT:
		return p.parseComment()
	case token.IF:
		return p.parseIfStatement()
	case token.ELSE:
		return p.parseElseStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.SET:
		return p.parseSetStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseComment() ast.Stmt {
	// TODO: implement
	return nil
}

func (p *Parser) parseIfStatement() ast.Stmt {
	// TODO: implement
	return nil
}

func (p *Parser) parseElseStatement() ast.Stmt {
	// TODO: implement
	return nil
}

func (p *Parser) parseForStatement() ast.Stmt {
	// TODO: implement
	return nil
}

func (p *Parser) parseSetStatement() ast.Stmt {
	// TODO: implement
	return nil
}

func (p *Parser) parseExpressionStatement() ast.Stmt {
	// TODO: implement
	return nil
}

func (p *Parser) parseIdentifier() ast.Expr {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseDeclare(left ast.Expr) ast.Expr {
	name, ok := left.(*ast.Identifier)
	if !ok {
		p.errors.Add("left side of : must be a variable", left.Range())
		return nil
	}
	exp := &ast.DeclareExpression{Token: p.curToken, Ident: name}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) parseExpression(precedence int) ast.Expr {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors.Add(fmt.Sprintf("no prefix parse function for %s", p.curToken.Type),
			&p.curToken.Range)
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMI) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = *p.s.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}
