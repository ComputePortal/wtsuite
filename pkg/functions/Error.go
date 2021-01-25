package functions

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Error(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

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
