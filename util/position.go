package util

import "fmt"

// Position represents a position within a file.
type Position struct {

	// Filename where the token is located.
	Filename string

	// Line is a 1-indexed line number.
	Line int

	// Column is a 0-indexed column number.
	Column int

	// Offset is a 0-indexed position in the file string.
	Offset int
}

func (pos Position) String() string {
	s := pos.Filename
	if s != "" {
		s += ":"
	}
	s += fmt.Sprintf("%d", pos.Line)
	s += ":"
	s += fmt.Sprintf("%d", pos.Column)
	return s
}

// Range represents a range of two positions within a file.
type Range struct {

	// Start is the beginning position.
	Start Position

	// End is the end position.
	End Position
}

func (r Range) String() string {
	return r.Start.String() + "-" + r.End.String()
}
