package lexer

import (
	"github.com/policyscript/policyscript/token"

	"github.com/policyscript/policyscript/util"
)

// Lexer lexes tokens from an input string.
type Lexer struct {
	input  []rune
	ch     rune
	line   int
	column int
	offset int
}

// New initializes a new Lexer.
func New(input string) *Lexer {
	l := &Lexer{
		input:  []rune(input),
		line:   1,
		column: -1,
		offset: -1,
	}

	l.advance()
	return l
}

// Scan the input into a list of tokens.
func (l *Lexer) Scan() []*token.Token {
	var tokens []*token.Token

	for !l.isAtEnd() {
		tokens = append(tokens, l.nextToken())
	}

	// EOF token.
	tokens = append(tokens, l.nextToken())

	return tokens
}

func (l *Lexer) nextToken() *token.Token {
	l.consumeWhitespaceAndEmptyLines()

	switch l.ch {
	case '_':
		if l.column == 0 && l.peek() == ' ' {
			return l.readHeading()
		}
		return l.readParagraph()
	case '#':
		return l.readComment()
	default:
		return l.readParagraph()
	}
}

func (l *Lexer) readHeading() *token.Token {
	return l.readUntilDoubleLineBreak(token.HEADING)
}

func (l *Lexer) readParagraph() *token.Token {
	return l.readUntilDoubleLineBreak(token.PARAGRAPH)
}

func (l *Lexer) readComment() *token.Token {
	var (
		start   util.Position = l.getPosition()
		end     util.Position = start
		literal []rune        = nil
	)

	for {
		// Consume the leading '#'.
		l.advance()

		startOffset := l.offset

		// Advance until new line.
		for !l.isAtEnd() && l.ch != '\n' {
			l.advance()
		}

		end = l.getPosition()
		literal = append(literal, l.input[startOffset:end.Offset]...)

		// If not at end, advance to consume '\n'.
		if !l.isAtEnd() {
			l.advance()
		}

		// If at end, return token.
		if l.isAtEnd() {
			return makeToken(token.COMMENT, literal, start, end)
		}

		nextLineIsComment := l.consumeWhitespaceTillChar('#')
		if !nextLineIsComment {
			return makeToken(token.COMMENT, literal, start, end)
		}

		// Otherwise, add the line break we skipped.
		literal = append(literal, '\n')
	}
}

// readUntilDoubleLineBreak will read a token until it hits a double line break,
// EOF or a comment.
func (l *Lexer) readUntilDoubleLineBreak(tokenType token.Type) *token.Token {
	var (
		start   util.Position = l.getPosition()
		end     util.Position = start
		literal []rune        = nil
	)

	for {
		var (
			startOffset    = l.offset
			whitespaceOnly = true
		)

		// Return early if new line begins with a comment.
		if l.ch == '#' {
			return makeToken(tokenType, literal, start, end)
		}

		// Advance until new line.
		for !l.isAtEnd() && l.ch != '\n' {
			// Set whitespace only to false if line contains a non-
			// whitespace character.
			if !isWhitespace(l.ch) {
				whitespaceOnly = false
			}

			l.advance()
		}

		// Update position if non-whitespace only line.
		if !whitespaceOnly {
			end = l.getPosition()

			// If literal is not nil (has been through a cycle),
			// append the line break.
			if literal != nil {
				literal = append(literal, '\n')
			} else {
				literal = []rune{}
			}

			literal = append(literal, l.input[startOffset:end.Offset]...)
		}

		// If at end, return token. If line was whitespace only, return
		// since it is a "double line break".
		if l.isAtEnd() || whitespaceOnly {
			if !l.isAtEnd() {
				l.advance()
			}

			return makeToken(tokenType, literal, start, end)
		}

		l.advance()
	}
}

func (l *Lexer) advance() {
	if l.ch == '\n' {
		l.line++
		l.column = -1
	}

	if l.offset+1 >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.offset+1]
	}

	l.offset++
	l.column++
}

func (l *Lexer) peek() rune {
	if l.offset+1 >= len(l.input) {
		return 0
	}
	return l.input[l.offset+1]
}

// consumeWhitespaceTillChar returns false and resets position if no match.
func (l *Lexer) consumeWhitespaceTillChar(char rune) bool {
	// Memoize values.
	var (
		ch     = l.ch
		line   = l.line
		column = l.column
		offset = l.offset
	)

	l.skipWhitespace()

	if l.ch == char {
		return true
	}

	// Else, restore to memo values and return.
	l.ch = ch
	l.line = line
	l.column = column
	l.offset = offset
	return false
}

func (l *Lexer) consumeWhitespaceAndEmptyLines() {
	// Memoize values.
	var (
		ch     = l.ch
		line   = l.line
		column = l.column
		offset = l.offset
	)

	for {
		l.skipWhitespace()

		if l.ch == '\n' {
			// If new line, consume then store current values.
			l.advance()
			ch = l.ch
			line = l.line
			column = l.column
			offset = l.offset
		} else {
			// Else, restore to memo values and return.
			l.ch = ch
			l.line = line
			l.column = column
			l.offset = offset
			return
		}
	}
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.advance()
	}
}

func (l *Lexer) skipWhitespaceAndBreaks() {
	for l.ch == '\n' || isWhitespace(l.ch) {
		l.advance()
	}
}

func (l *Lexer) isAtEnd() bool {
	return l.offset >= len(l.input)
}

func (l *Lexer) getPosition() util.Position {
	return util.Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.offset,
	}
}

func isWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\r'
}

func makeToken(tokenType token.Type, literal []rune, start, end util.Position) *token.Token {
	return &token.Token{
		Type:    tokenType,
		Literal: string(literal),
		Range: util.Range{
			Start: start,
			End:   end,
		},
	}
}
