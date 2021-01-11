package lexer_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/policyscript/policyscript/lexer"
	"github.com/policyscript/policyscript/token"
)

var _ = Describe("Lexer", func() {
	each([][2]string{
		{"\n_ A", "_ A"},
		{"\n_ A\n", "_ A"},
		{"\n_ A\n\n", "_ A"},
		{"\n_ A\nB", "_ A\nB"},
		{"\n_ A\nB#C", "_ A\nB#C"},
		{"\n_ A\n# B", "_ A"},
		{"\n_ A\n # B", "_ A\n # B"},
		{"\n_ A\n\t# B", "_ A\n\t# B"},
		{"\n_ A\nB\n", "_ A\nB"},
		{"\n_ A\nB\n\n", "_ A\nB"},
		{"\n_ A\nB\n\nC", "_ A\nB"},
		{"\n_ A\n\tB", "_ A\n\tB"},
		{"\n_ A\n  B", "_ A\n  B"},
		{"\n_ A", "_ A"},
		{"\n_ _A", "_ _A"},
	}, "should find secondary title tokens", func(input, expects string) {
		l := lexer.New(input)
		tokens := l.Scan()
		t := tokens[0]

		Expect(t.Type).To(Equal(token.HEADING))
		Expect(t.Literal).To(Equal(expects))
	})

	each([][2]string{
		{"\nA", "A"},
		{"\n - A\n", " - A"},
		{"\nA\n", "A"},
		{"\nA\n\n", "A"},
		{"\nA\nB", "A\nB"},
		{"\nA\n# B", "A"},
		{"\nA\n # B", "A\n # B"},
		{"\nA\n\t# B", "A\n\t# B"},
		{"\nA\nB\n", "A\nB"},
		{"\nA\nB\n\n", "A\nB"},
		{"\nA\nB\n\nC", "A\nB"},
		{"\nA\n\tB", "A\n\tB"},
		{"\nA\n  B", "A\n  B"},
		{"\n # A\n # B\n", " # A\n # B"},
		{"\n\t# A\n\t# B\n", "\t# A\n\t# B"},
	}, "should find paragraph tokens", func(input, expects string) {
		l := lexer.New(input)
		tokens := l.Scan()
		t := tokens[0]

		Expect(t.Type).To(Equal(token.PARAGRAPH))
		Expect(t.Literal).To(Equal(expects))
	})

	each([][2]string{
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
	}, "should find comment tokens", func(input, expects string) {
		l := lexer.New(input)
		tokens := l.Scan()
		t := tokens[0]

		Expect(t.Type).To(Equal(token.COMMENT))
		Expect(t.Literal).To(Equal(expects))
	})
})

func each(items [][2]string, title string, handler func(input, expects string)) {
	for _, pair := range items {
		input := pair[0]
		expects := pair[1]

		text := fmt.Sprintf("%s (input: %q | expects: %q)", title, input, expects)
		It(text, func() {
			handler(input, expects)
		})
	}
}
