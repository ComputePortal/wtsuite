package functions

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func joinStringList(l *tokens.List, sep string, ctx context.Context) (tokens.Token, error) {
	var b strings.Builder

	if err := l.Loop(func(i int, item tokens.Token, last bool) error {
		if i > 0 {
			b.WriteString(sep)
		}

		switch s := item.(type) {
		case *tokens.String:
			b.WriteString(s.Value())
			return nil
		default:
			errCtx := context.MergeContexts(ctx, item.Context())
			return errCtx.NewError("Error: expected a string")
		}
	}); err != nil {
		return nil, err
	}

	return tokens.NewString(b.String(), ctx)
}

func joinStrings(a *tokens.String, b *tokens.String, ctx context.Context) (tokens.Token, error) {
	return tokens.NewString(a.Value()+b.Value(), ctx)
}

func Join(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	switch len(args) {
	case 1:
		switch a := args[0].(type) {
		case *tokens.List:
			return joinStringList(a, "", ctx)
		default:
			return nil, ctx.NewError("Error: expected list")
		}
	case 2:
		switch a := args[0].(type) {
		case *tokens.List:
			switch b := args[1].(type) {
			case *tokens.String:
				return joinStringList(a, b.Value(), ctx)
			default:
				return nil, ctx.NewError("Error: expected separator string")
			}
		case *tokens.String:
			switch b := args[1].(type) {
			case *tokens.String:
				return joinStrings(a, b, ctx)
			default:
				return nil, ctx.NewError("Error: expected second string")
			}
		default:
			return nil, ctx.NewError("Error: expected string or list")
		}
	default:
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}
}
