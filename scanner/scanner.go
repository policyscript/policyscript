package scanner

import (
	"bytes"

	"github.com/policyscript/policyscript/token"

	"github.com/policyscript/policyscript/util"
)

// Scanner lexes tokens from an input string.
type Scanner struct {
	input  []rune
	ch     rune
	line   int
	column int
	offset int
}

// New initializes a new Scanner.
func New(input []byte) *Scanner {
	l := &Scanner{
		input:  bytes.Runes(input),
		line:   1,
		column: -1,
		offset: -1,
	}

	l.next()
	return l
}

// Scan the input into a list of tokens.
func (l *Scanner) Scan() []*token.Token {
	var tokens []*token.Token

	for !l.isAtEnd() {
		tokens = append(tokens, l.nextToken())
	}

	// EOF token.
	tokens = append(tokens, l.nextToken())

	return tokens
}

func (l *Scanner) nextToken() *token.Token {
	start := l.skipWhitespaceAndBreaks()

	switch l.ch {
	case '_':
		// Only a heading if it is in the first column and followed by a
		// space character.
		if l.column == 0 && l.peek() == ' ' {
			return l.readHeading()
		}
		return l.readParagraph(start)
	case '#':
		return l.readComment()
	case 0:
		position := l.getPosition()
		return makeToken(token.EOF, nil, position, position)
	default:
		return l.readParagraph(start)
	}
}

func (l *Scanner) readHeading() *token.Token {
	return l.readUntilDoubleLineBreak(token.HEADING, nil)
}

func (l *Scanner) readParagraph(start *util.Position) *token.Token {
	return l.readUntilDoubleLineBreak(token.PARAGRAPH, start)
}

func (l *Scanner) readComment() *token.Token {
	var (
		start   *util.Position = l.getPosition()
		end     *util.Position = start
		literal []rune         = nil
	)

	for {
		// Consume the leading '#'.
		l.next()

		startOffset := l.offset

		// Advance until new line.
		for !l.isAtEnd() && l.ch != '\n' {
			l.next()
		}

		end = l.getPosition()
		literal = append(literal, l.input[startOffset:end.Offset]...)

		// If not at end, advance to consume '\n'.
		if !l.isAtEnd() {
			l.next()
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
func (l *Scanner) readUntilDoubleLineBreak(tokenType token.Type, prevStart *util.Position) *token.Token {
	var (
		start   *util.Position = l.getPosition()
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
		for !l.isAtEnd() && l.ch != '\n' {
			// Set whitespace only to false if line contains a non-
			// whitespace character.
			if !isWhitespace(l.ch) {
				// Return early if new line begins with a comment.
				if whitespaceOnly && l.ch == '#' {
					return makeToken(tokenType, literal, start, end)
				}

				whitespaceOnly = false
			}

			l.next()
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
				l.next()
			}

			return makeToken(tokenType, literal, start, end)
		}

		l.next()
		startOffset = l.offset
	}
}

// This will return false and reset position if no match. It does not consume
// the character.
func (l *Scanner) consumeWhitespaceTillChar(char rune) bool {
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

	// Else, reset position and return false.
	l.ch = ch
	l.line = line
	l.column = column
	l.offset = offset
	return false
}

// This will skip through empty lines, but will return the position at the start
// of the first non-empty line.
func (l *Scanner) skipWhitespaceAndBreaks() *util.Position {
	position := &util.Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.offset,
	}

	for {
		l.skipWhitespace()

		if l.ch == '\n' {
			// If new line, consume then store current values.
			l.next()
			position.Line = l.line
			position.Column = l.column
			position.Offset = l.offset
		} else {
			return position
		}
	}
}

func (l *Scanner) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.next()
	}
}

func (l *Scanner) next() {
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

func (l *Scanner) peek() rune {
	if l.offset+1 >= len(l.input) {
		return 0
	}
	return l.input[l.offset+1]
}

func (l *Scanner) isAtEnd() bool {
	return l.offset >= len(l.input)
}

func (l *Scanner) getPosition() *util.Position {
	return &util.Position{
		Line:   l.line,
		Column: l.column,
		Offset: l.offset,
	}
}

func isWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\r'
}

func makeToken(tokenType token.Type, literal []rune, start, end *util.Position) *token.Token {
	return &token.Token{
		Type:    tokenType,
		Literal: string(literal),
		Range: util.Range{
			Start: start,
			End:   end,
		},
	}
}
