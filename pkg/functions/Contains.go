package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

func Contains(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
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
