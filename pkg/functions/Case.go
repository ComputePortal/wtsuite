package functions

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func strUpperRange(fn func(string) string, s *tokens.String, start *tokens.Int, stop *tokens.Int, ctx context.Context) (tokens.Token, error) {
	if start.Value() > s.Len()-1 {
		return s, nil
	}

	if start.Value() < 0 {
		errCtx := ctx
		return nil, errCtx.NewError("Error: negative start index")
	}

	if stop.Value() < 0 {
		errCtx := ctx
		return nil, errCtx.NewError("Error: negative stop index")
	}

	n := s.Len()

	if stop.Value() < n {
		n = stop.Value()
	}

	snew := fn(s.Value()[start.Value():n])
	sold := s.Value()[n:]

	return tokens.NewString(snew+sold, ctx)
}

func lowerUpper(fn func(string) string, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) == 0 {
		return nil, ctx.NewError("Error: expected at least 1 argument")
	}

	s, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	n := len(s.Value())
	switch len(args) {
	case 1:
		return strUpperRange(fn, s, tokens.NewValueInt(0, ctx), tokens.NewValueInt(n, ctx), ctx)
	case 2:
		stop, err := tokens.AssertInt(args[1])
		if err != nil {
			return nil, err
		}

		return strUpperRange(fn, s, tokens.NewValueInt(0, ctx), stop, ctx)
	case 3:
		start, err := tokens.AssertInt(args[1])
		if err != nil {
			return nil, err
		}

		stop, err := tokens.AssertInt(args[2])
		if err != nil {
			return nil, err
		}

		return strUpperRange(fn, s, start, stop, ctx)
	default:
		return nil, ctx.NewError("Error: expected 1, 2 or 3 arguments")
	}
}

func Upper(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return lowerUpper(strings.ToUpper, args, ctx)
}

func Lower(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return lowerUpper(strings.ToLower, args, ctx)
}

// smart capitaloization of titles
// ref:https://capitalizemytitle.com/#capitalizationrules
func Caps(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	s, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(s.Value())

	result := make([]string, len(fields))

	for i, field := range fields {
		switch field {
		case "of", "and":
			result[i] = field
		default:
			result[i] = strings.Title(field)
		}
	}

	return tokens.NewString(strings.Join(result, " "), ctx)
}
