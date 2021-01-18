package scanner_test

import (
	"fmt"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/policyscript/policyscript/scanner"
	"github.com/policyscript/policyscript/token"
	"github.com/policyscript/policyscript/util"
)

var _ = Describe("Scanner", func() {
	util.Each("can scan title", [][2]string{
		{"\n_ A", "_ A"},
		{"\n_ A\n", "_ A"},
		{"\n_ A\n\n", "_ A"},
		{"\n_ A\nB", "_ A\nB"},
		{"\n_ A\nB#C", "_ A\nB#C"},
		{"\n_ A\n# B", "_ A"},
		{"\n_ A\n # B", "_ A"},
		{"\n_ A\n\t# B", "_ A"},
		{"\n_ A\nB\n", "_ A\nB"},
		{"\n_ A\nB\n\n", "_ A\nB"},
		{"\n_ A\nB\n\nC", "_ A\nB"},
		{"\n_ A\n\tB", "_ A\n\tB"},
		{"\n_ A\n  B", "_ A\n  B"},
		{"\n_ A", "_ A"},
		{"\n_ _A", "_ _A"},
	}, func(input, expects string) {
		l := scanner.New([]byte(input), nil)
		tokens := l.Scan()
		t := tokens[0]

		Expect(t.Type).To(Equal(token.HEADING))
		Expect(t.Literal).To(Equal(expects))
	})

	util.Each("can scan paragraph", [][2]string{
		{"\nA", "A"},
		{"\n - A\n", " - A"},
		{"\nA\n", "A"},
		{"\nA\n\n", "A"},
		{"\nA\nB", "A\nB"},
		{"\nA\n# B", "A"},
		{"\nA\n # B", "A"},
		{"\nA\n\t# B", "A"},
		{"\nA\nB\n", "A\nB"},
		{"\nA\nB\n\n", "A\nB"},
		{"\nA\nB\n\nC", "A\nB"},
		{"\nA\n\tB", "A\n\tB"},
		{"\nA\n  B", "A\n  B"},
	}, func(input, expects string) {
		l := scanner.New([]byte(input), nil)
		tokens := l.Scan()
		t := tokens[0]

		Expect(t.Type).To(Equal(token.PARAGRAPH))
		Expect(t.Literal).To(Equal(expects))
	})

	util.Each("can scan comment", [][2]string{
		{"\n#A", "A"},
		{"\n#A\n#B", "A\nB"},
		{"\n# A", " A"},
		{"\n# A\n", " A"},
		{"\n# A\n\n", " A"},
		{"\n# A\nC", " A"},
		{"\n# A\n# B", " A\n B"},
		{"\n# A\n# B\n", " A\n B"},
		{"\n# A\n # B\n", " A\n B"},
		{"\n# A\n# B\n\n", " A\n B"},
		{"\n# A\n# B\n\n# C", " A\n B"},
		{"\n# A\n# \tB", " A\n \tB"},
		{"\n# A\n#   B", " A\n   B"},
		{"\n # A\n # B\n", " A\n B"},
		{"\n\t# A\n\t# B\n", " A\n B"},
	}, func(input, expects string) {
		l := scanner.New([]byte(input), nil)
		tokens := l.Scan()
		t := tokens[0]

		Expect(t.Type).To(Equal(token.COMMENT))
		Expect(t.Literal).To(Equal(expects))
	})

	eachTokens("can scan program", []inputAndTokens{
		{input: "_ Heading", expects: []token.Token{
			{Type: token.HEADING, Literal: "_ Heading", Range: &util.Range{
				Start: &util.Position{Line: 1, Column: 0, Offset: 0},
				End:   &util.Position{Line: 1, Column: 9, Offset: 9},
			}},
			{Type: token.EOF, Literal: "", Range: &util.Range{
				Start: &util.Position{Line: 1, Column: 9, Offset: 9},
				End:   &util.Position{Line: 1, Column: 9, Offset: 9},
			}},
		}},
		{input: "_ Heading\n # Comment", expects: []token.Token{
			{Type: token.HEADING, Literal: "_ Heading", Range: &util.Range{
				Start: &util.Position{Line: 1, Column: 0, Offset: 0},
				End:   &util.Position{Line: 1, Column: 9, Offset: 9},
			}},
			{Type: token.COMMENT, Literal: " Comment", Range: &util.Range{
				Start: &util.Position{Line: 2, Column: 1, Offset: 11},
				End:   &util.Position{Line: 2, Column: 10, Offset: 20},
			}},
			{Type: token.EOF, Literal: "", Range: &util.Range{
				Start: &util.Position{Line: 2, Column: 10, Offset: 20},
				End:   &util.Position{Line: 2, Column: 10, Offset: 20},
			}},
		}},
		{input: "_ Heading\n\n  Paragraph\n continued", expects: []token.Token{
			{Type: token.HEADING, Literal: "_ Heading", Range: &util.Range{
				Start: &util.Position{Line: 1, Column: 0, Offset: 0},
				End:   &util.Position{Line: 1, Column: 9, Offset: 9},
			}},
			{Type: token.PARAGRAPH, Literal: "  Paragraph\n continued", Range: &util.Range{
				Start: &util.Position{Line: 3, Column: 0, Offset: 11},
				End:   &util.Position{Line: 4, Column: 10, Offset: 33},
			}},
			{Type: token.EOF, Literal: "", Range: &util.Range{
				Start: &util.Position{Line: 4, Column: 10, Offset: 33},
				End:   &util.Position{Line: 4, Column: 10, Offset: 33},
			}},
		}},
	})
})

type inputAndTokens struct {
	input   string
	expects []token.Token
}

// Each runs every item in a list of inputs and expects.
func eachTokens(title string, items []inputAndTokens) {
	for i, pair := range items {
		var (
			input   = pair.input
			expects = pair.expects
			text    = fmt.Sprintf("%s (index: %d | input: %q | expects: %v)",
				title, i, input, expects)
		)
		ginkgo.It(text, func() {
			l := scanner.New([]byte(input), nil)
			tokens := l.Scan()

			Expect(len(tokens)).To(Equal(len(expects)), "number of tokens incorrect")

			for i, t := range tokens {
				Expect(t).To(Equal(expects[i]), fmt.Sprintf("token index %d mismatch", i))
			}
		})
	}
}
