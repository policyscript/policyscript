package util

type (

	// Position represents a position within a file.
	Position struct {

		// Line is a 1-indexed line number.
		Line int

		// Column is a 0-indexed column number.
		Column int

		// Offset is a 0-indexed position in the file string.
		Offset int
	}

	// Range represents a range of two positions within a file.
	Range struct {

		// Start is the beginning position.
		Start Position

		// End is the end position.
		End Position
	}
)
