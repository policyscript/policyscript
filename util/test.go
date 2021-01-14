package util

import (
	"fmt"

	"github.com/onsi/ginkgo"
)

// Each runs every item in a list of inputs and expects.
func Each(title string, items [][2]string, handler func(input, expects string)) {
	for i, pair := range items {
		var (
			input   = pair[0]
			expects = pair[1]
			text    = fmt.Sprintf("%s (index: %d | input: %q | expects: %q)",
				title, i, input, expects)
		)
		ginkgo.It(text, func() { handler(input, expects) })
	}
}
