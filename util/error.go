package util

// Error describes an error and the position of the error.
type Error struct {
	Msg string
	Rng Range
}

func (e Error) Error() string {
	return e.Rng.String() + ": " + e.Msg
}

// ErrorList is a list of pointers to errors.
type ErrorList []*Error

// Add appends an error to the error list.
func (p *ErrorList) Add(msg string, rng *Range) {
	*p = append(*p, &Error{Msg: msg, Rng: *rng})
}
