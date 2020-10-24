package functions

import (
	"strings"

	"../tokens/context"
	tokens "../tokens/html"
)

func Error(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	var b strings.Builder
	b.WriteString("User Error: ")

	for i, arg := range args {
		s, err := tokens.AssertString(arg)
		if err != nil {
			return nil, err
		}

		b.WriteString(s.Value())

		if i < len(args)-1 {
			b.WriteString(" ")
		}
	}

	b.WriteString("\n")

	errCtx := ctx
	return nil, errCtx.NewError(b.String())
}
