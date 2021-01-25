package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

func Contains(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, NewInterface([]string{"container", "key"}, ctx))
  if err != nil {
    return nil, err
  }

	container := args[0]
	switch {
	case tokens.IsNull(container):
		return tokens.NewValueBool(false, ctx), nil
	case tokens.IsList(container):
		lst, err := tokens.AssertList(container)
		if err != nil {
			panic(err)
		}

		ok := false
		if err := lst.Loop(func(i int, val tokens.Token, last bool) error {
			if val.IsSame(args[1]) {
				ok = true
			}

			return nil
		}); err != nil {
			return nil, err
		}

		return tokens.NewValueBool(ok, ctx), nil
	case tokens.IsKeyDict(container):
		d, err := tokens.AssertKeyDict(container)
		if err != nil {
			panic(err)
		}

		_, ok := d.Get(args[1])
		return tokens.NewValueBool(ok, ctx), nil
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: not a list or dict")
	}
}
