package scanner

import (
	"bytes"

	"github.com/policyscript/policyscript/token"

	"github.com/policyscript/policyscript/util"
)

// Scanner lexes tokens from an input string.
type Scanner struct {
	input []rune
	err   ErrorHandler

	// Position
	ch     rune
	line   int
	column int
	offset int

	// Block
	blockInit bool
	block     bool
	addSemi   bool

	ErrorCount int
}

// An ErrorHandler is called with each error the scanner finds.
type ErrorHandler func(msg string, pos util.Position)

// New initializes a new Scanner.
func New(input []byte, err ErrorHandler) *Scanner {
	l := &Scanner{
		input: bytes.Runes(input),
		err:   err,

		// ch will be obtained by calling `next`
		line:   1,
		column: -1,
		offset: -1,

		blockInit: false,
		block:     false,
	}

	l.next()
	return l
}

// Scan the input into a list of tokens.
func (s *Scanner) Scan() []token.Token {
	var tokens []token.Token

	for !s.isAtEnd() {
		tokens = append(tokens, *s.nextToken())
	}

	// EOF token.
	tokens = append(tokens, *s.nextToken())

	return tokens
}

func (s *Scanner) error(msg string, pos *util.Position) {
	if s.err != nil {
		s.err(msg, *pos)
	}
	s.ErrorCount++
}

func (s *Scanner) nextToken() *token.Token {
	// Special scanning for when we are inside a block.
	if s.block {
		return s.nextBlockToken()
	}

	line := s.line
	start := s.skipWhitespaceAndBreaks()

	if s.blockInit {
		// Only continue block init if it's on the same line.
		if s.line == line {
			switch {
			case s.ch == '{':
				// Begin block.
				s.blockInit = false
				s.block = true

				return s.makeSingleRuneToken(token.LBRACE)
			case isAlphaNumeric(s.ch):
				// Get type def identifier.
				start = s.getPosition()
				s.next()
				for isAlphaNumeric(s.ch) {
					s.next()
				}
				return s.makeMultiRuneToken(token.IDENT, start)
			}
		}

		// Otherwise, stop block init.
		s.blockInit = false
	}

	switch s.ch {
	case '_':
		// Only a heading if it is in the first column and followed by a
		// space character.
		if s.column == 0 && s.peek() == ' ' {
			return s.readHeading()
		}
	case '@':
		// Only a block if in first column.
		if s.column == 0 {
			if token := s.eatBlockKeyword(start); token != nil {
				return token
			}
		}
	case '#':
		return s.readComment()
	case 0:
		position := s.getPosition()
		return makeToken(token.EOF, nil, position, position)
	}

	return s.readParagraph(start)
}

func (s *Scanner) nextBlockToken() *token.Token {
	line := s.line
	start := s.getPosition()

	s.skipWhitespaceAndBreaks()

	if s.addSemi && s.line != line && !s.isAtEnd() {
		s.addSemi = false
		return makeToken(token.SEMI, []rune{';'}, start, start)
	}

	s.addSemi = false

	switch s.ch {
	case 0:
		// End block.
		s.block = false
		position := s.getPosition()
		s.error("block does not have closing \"}\"", position)
		return makeToken(token.EOF, nil, position, position)
	case '}':
		// End block.
		s.block = false
		return s.makeSingleRuneToken(token.RBRACE)
	case '`':
		s.addSemi = true
		return s.readText()
	case '|':
		s.addSemi = true
		if isNumeric(s.peek()) {
			return s.readTimeOrDate()
		}
		return s.makeSingleRuneToken(token.ILLEGAL)
	case '(':
		return s.makeSingleRuneToken(token.LPAREN)
	case ')':
		s.addSemi = true
		return s.makeSingleRuneToken(token.RPAREN)
	case ':':
		return s.makeSingleRuneToken(token.COLON)
	case '=':
		return s.makeSingleRuneToken(token.EQ)
	case '!':
		s.next()
		if s.ch == '=' {
			s.next()
			return s.makeMultiRuneToken(token.NOT_EQ, start)
		}
		return s.makeMultiRuneToken(token.ILLEGAL, start)
	case '+':
		return s.makeSingleRuneToken(token.PLUS)
	case '-':
		return s.makeSingleRuneToken(token.MINUS)
	case '*':
		return s.makeSingleRuneToken(token.MULT)
	case '/':
		return s.makeSingleRuneToken(token.DIV)
	case '#':
		return s.readComment()
	case '<':
		s.next()
		if s.ch == '=' {
			s.next()
			return s.makeMultiRuneToken(token.LT_EQ, start)
		}
		return s.makeMultiRuneToken(token.LT, start)
	case '>':
		s.next()
		if s.ch == '=' {
			s.next()
			return s.makeMultiRuneToken(token.GT_EQ, start)
		}
		return s.makeMultiRuneToken(token.GT, start)
	default:
		switch {
		case token.LookupMoney(s.ch) && isNumeric(s.peek()):
			s.addSemi = true
			t := token.MONEY
			return s.readNumber(&t)
		case isAlpha(s.ch):
			// Will add semi after identifier but not keyword.
			return s.readIdentifier()
		case isNumeric(s.ch):
			s.addSemi = true
			return s.readNumber(nil)
		default:
			s.addSemi = true
			return s.makeSingleRuneToken(token.ILLEGAL)
		}
	}
}

func (s *Scanner) readText() *token.Token {
	start := s.getPosition()
	s.next()
	textStart := s.offset

	for s.ch != '`' && !s.isAtEnd() {
		s.next()
	}

	literal := s.input[textStart:s.offset]

	if s.isAtEnd() {
		s.error("text does not have closing \"`\"", start)
	} else {
		s.next()
	}

	return makeToken(token.TEXT, literal, start, s.getPosition())
}

func (s *Scanner) readIdentifier() *token.Token {
	start := s.getPosition()
	for isAlphaNumeric(s.ch) {
		s.next()
	}
	var (
		end                  = s.getPosition()
		literal              = s.input[start.Offset:end.Offset]
		tokenType, isKeyword = token.LookupIdent(literal)
	)
	// Don't end line after a keyword.
	if !isKeyword {
		s.addSemi = true
	}
	return makeToken(tokenType, literal, start, end)
}

func (s *Scanner) readNumber(tokenType *token.Type) *token.Token {
	t := token.INTEGER
	start := s.getPosition()
	isInteger := true

	// Eat first char, already scanned.
	s.next()

	// Before decimal point can have arbitrary underscores.
	for isNumericUnderscore(s.ch) {
		s.next()
	}

	// If decimal point and next is numeric, parse up until numeric ends.
	if s.ch == '.' && isNumeric(s.peek()) {
		isInteger = false
		t = token.DECIMAL
		s.next()
		for isNumeric(s.ch) {
			s.next()
		}
	}

	// If no provided token type, can optionally "search" for percentage or
	// period.
	if tokenType == nil {
		if s.ch == '%' {
			t = token.PERCENT
			s.next()
		} else if isInteger && s.checkPeriod() {
			// Can only be period if integer value.
			t = token.PERIOD
		}
	} else {
		t = *tokenType
	}
	return s.makeMultiRuneToken(t, start)
}

func (s *Scanner) readTimeOrDate() *token.Token {
	start := s.getPosition()

	// Eat `|`.
	s.next()

	for isNumeric(s.ch) {
		s.next()
	}

	if s.ch == ':' {
		s.readTime(start)
		return s.makeMultiRuneToken(token.TIME, start)
	}

	// Default to date even if `/` is not present.
	s.readDate(start)
	return s.makeMultiRuneToken(token.DATE, start)
}

func (s *Scanner) readDate(start *util.Position) {
	for isNumeric(s.ch) || s.ch == '/' {
		s.next()
	}
	if s.ch == '|' {
		s.next()
	} else {
		s.error("invalid date", start)
	}
}

func (s *Scanner) readTime(start *util.Position) {
	for isNumeric(s.ch) || s.ch == ':' {
		s.next()
	}
	if s.ch == '|' {
		s.next()
	} else {
		s.error("invalid time", start)
	}
}

// This will return false and reset position if no match.
func (s *Scanner) checkPeriod() bool {
	// Memoize values.
	var (
		ch     = s.ch
		line   = s.line
		column = s.column
		offset = s.offset
	)

	s.skipWhitespace()

	if !isAlpha(s.ch) {
		s.resetPosition(ch, line, column, offset)
		return false
	}

	start := s.offset
	for isAlpha(s.ch) {
		s.next()
	}
	end := s.offset

	literal := s.input[start:end]
	if ok := token.LookupPeriodKeyword(literal); ok {
		return true
	}
	s.resetPosition(ch, line, column, offset)
	return false
}

func (s *Scanner) readHeading() *token.Token {
	return s.readUntilDoubleLineBreak(token.HEADING, nil)
}

func (s *Scanner) readParagraph(start *util.Position) *token.Token {
	return s.readUntilDoubleLineBreak(token.PARAGRAPH, start)
}

func (s *Scanner) readComment() *token.Token {
	var (
		start   *util.Position = s.getPosition()
		end     *util.Position = start
		literal []rune         = nil
	)

	for {
		// Consume the leading '#'.
		s.next()

		startOffset := s.offset

		// Advance until new line.
		for !s.isAtEnd() && s.ch != '\n' {
			s.next()
		}

		end = s.getPosition()
		literal = append(literal, s.input[startOffset:end.Offset]...)

		// If not at end, advance to consume '\n'.
		if !s.isAtEnd() {
			s.next()
		}

		// If at end, return token.
		if s.isAtEnd() {
			return makeToken(token.COMMENT, literal, start, end)
		}

		nextLineIsComment := s.consumeWhitespaceTillChar('#')
		if !nextLineIsComment {
			return makeToken(token.COMMENT, literal, start, end)
		}

		// Otherwise, add the line break we skipped.
		literal = append(literal, '\n')
	}
}

// readUntilDoubleLineBreak will read a token until it hits a double line break,
// EOF or a comment.
func (s *Scanner) readUntilDoubleLineBreak(tokenType token.Type, prevStart *util.Position) *token.Token {
	var (
		start   *util.Position = s.getPosition()
		end     *util.Position = start
		literal []rune         = nil
	)

	if prevStart != nil {
		start = prevStart
	}

	var (
		startOffset    = start.Offset
		whitespaceOnly bool
	)

	for {
		whitespaceOnly = true

		// Advance until new line.
		for !s.isAtEnd() && s.ch != '\n' {
			// Set whitespace only to false if line contains a non-
			// whitespace character.
			if !isWhitespace(s.ch) {
				// Return early if new line begins with a comment.
				if whitespaceOnly && s.ch == '#' {
					return makeToken(tokenType, literal, start, end)
				}

				whitespaceOnly = false
			}

			s.next()
		}

		// Update position if non-whitespace only line.
		if !whitespaceOnly {
			end = s.getPosition()

			// If literal is not nil (has been through a cycle),
			// append the line break.
			if literal != nil {
				literal = append(literal, '\n')
			} else {
				literal = []rune{}
			}

			literal = append(literal, s.input[startOffset:end.Offset]...)
		}

		// If at end, return token. If line was whitespace only, return
		// since it is a "double line break".
		if s.isAtEnd() || whitespaceOnly {
			if !s.isAtEnd() {
				s.next()
			}

			return makeToken(tokenType, literal, start, end)
		}

		s.next()
		startOffset = s.offset
	}
}

// This will return false and reset position if no match. It does not consume
// the character.
func (s *Scanner) consumeWhitespaceTillChar(char rune) bool {
	// Memoize values.
	var (
		ch     = s.ch
		line   = s.line
		column = s.column
		offset = s.offset
	)

	s.skipWhitespace()

	if s.ch == char {
		return true
	}

	// Else, reset position and return false.
	s.resetPosition(ch, line, column, offset)
	return false
}

func (s *Scanner) resetPosition(ch rune, line int, column int, offset int) {
	s.ch = ch
	s.line = line
	s.column = column
	s.offset = offset
}

// This will skip through empty lines, but will return the position at the start
// of the first non-empty line.
func (s *Scanner) skipWhitespaceAndBreaks() *util.Position {
	position := &util.Position{
		Line:   s.line,
		Column: s.column,
		Offset: s.offset,
	}

	for {
		s.skipWhitespace()

		if s.ch == '\n' {
			// If new line, consume then store current values.
			s.next()
			position.Line = s.line
			position.Column = s.column
			position.Offset = s.offset
		} else {
			return position
		}
	}
}

func (s *Scanner) skipWhitespace() {
	for isWhitespace(s.ch) {
		s.next()
	}
}

func (s *Scanner) next() {
	if s.ch == '\n' {
		s.line++
		s.column = -1
	}

	if s.offset+1 >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = s.input[s.offset+1]
	}

	s.offset++
	s.column++
}

func (s *Scanner) eatBlockKeyword(start *util.Position) *token.Token {
	// Eat '@'.
	s.next()

	for isAlpha(s.ch) {
		s.next()
	}

	keyword := s.input[start.Offset:s.offset]

	tokenType, ok := token.LookupBlockKeyword(keyword)
	if !ok {
		return nil
	}

	// If block keyword, set pending to true.
	s.blockInit = true
	return makeToken(tokenType, keyword, start, s.getPosition())
}

func (s *Scanner) peek() rune {
	if s.offset+1 >= len(s.input) {
		return 0
	}
	return s.input[s.offset+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.offset >= len(s.input)
}

func (s *Scanner) getPosition() *util.Position {
	return &util.Position{
		Line:   s.line,
		Column: s.column,
		Offset: s.offset,
	}
}

func (s *Scanner) makeSingleRuneToken(tokenType token.Type) *token.Token {
	start := s.getPosition()
	s.next()
	end := s.getPosition()
	return makeToken(tokenType, s.input[start.Offset:end.Offset], start, end)
}

func (s *Scanner) makeMultiRuneToken(tokenType token.Type, start *util.Position) *token.Token {
	end := s.getPosition()
	literal := s.input[start.Offset:end.Offset]
	return makeToken(tokenType, literal, start, end)
}

func isWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\r'
}

func makeToken(tokenType token.Type, literal []rune, start, end *util.Position) *token.Token {
	return &token.Token{
		Type:    tokenType,
		Literal: string(literal),
		Range: &util.Range{
			Start: start,
			End:   end,
		},
	}
}

func isAlpha(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumeric(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isNumericUnderscore(ch rune) bool {
	return isNumeric(ch) || ch == '_'
}

func isAlphaNumeric(ch rune) bool {
	return isAlpha(ch) || isNumeric(ch)
}
