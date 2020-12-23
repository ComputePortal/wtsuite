package functions

import (
	"reflect"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func castIntToString(t *tokens.Int, ctx context.Context) (*tokens.String, error) {
	return tokens.NewString(t.Write(), ctx)
}

func castFloatToString(t *tokens.Float, ctx context.Context) (*tokens.String, error) {
	return tokens.NewString(t.Write(), ctx)
}

func castStringToString(t *tokens.String, ctx context.Context) (*tokens.String, error) {
	return tokens.NewString(t.Write(), ctx)
}

func castColorToString(t *tokens.Color, ctx context.Context) (*tokens.String, error) {
	return tokens.NewString(t.Write(), ctx)
}

func joinListAsString(lst *tokens.List, sepStr string, ctx context.Context) (*tokens.String, error) {
	var b strings.Builder

	if err := lst.Loop(func(i int, value tokens.Token, last bool) error {
		subStr, err := strInternal([]tokens.Token{value}, ctx)
		if err != nil {
			return err
		}

		b.WriteString(subStr.Value())

		if !last {
			b.WriteString(sepStr)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return tokens.NewString(b.String(), ctx)
}

func strInternal(args []tokens.Token, ctx context.Context) (*tokens.String, error) {
	switch t := args[0].(type) {
	case *tokens.Int:
		if len(args) != 1 {
			return nil, ctx.NewError("Error: expected 1 argument")
		}
		return castIntToString(t, ctx)
	case *tokens.Float:
		if len(args) != 1 {
			return nil, ctx.NewError("Error: expected 1 argument")
		}
		return castFloatToString(t, ctx)
	case *tokens.String:
		if len(args) != 1 {
			return nil, ctx.NewError("Error: expected 1 argument")
		}
		return castStringToString(t, ctx)
	case *tokens.Color:
		if len(args) != 1 {
			return nil, ctx.NewError("Error: expected 1 argument")
		}
		return castColorToString(t, ctx)
	case *tokens.List:
		if !(len(args) == 1 || len(args) == 2) {
			return nil, ctx.NewError("Error: expected 1 or 2 arguments")
		}
		sepStr := ""
		if len(args) == 2 {
			sep, err := tokens.AssertString(args[1])
			if err != nil {
				return nil, err
			}

			sepStr = sep.Value()
		}
		return joinListAsString(t, sepStr, ctx)
	default:
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected primitive, got " + reflect.TypeOf(t).String())
		err.AppendContextString("Info: called here", ctx)
		return nil, err
	}
}

func Str(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	return strInternal(args, ctx)
}
