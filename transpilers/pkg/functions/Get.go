package functions

import (
	"reflect"

	"../tokens/context"
	tokens "../tokens/html"
)

func Get(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) < 2 {
		return nil, ctx.NewError("Error: expected at least 2 arguments")
	}

	container := args[0]
	switch {
	case tokens.IsList(container):
		if len(args) != 2 {
			return nil, ctx.NewError("Error: expected 2 arguments")
		}

		lst, err := tokens.AssertList(container)
		if err != nil {
			panic(err)
		}

		index, err := tokens.AssertInt(args[1])
		if err != nil {
			return nil, err
		}

		value, err := lst.Get(index)
		if err != nil {
			errCtx := index.Context()
			return nil, errCtx.NewError("Error: " + err.Error())
		}

		return value, nil
	case tokens.IsKeyDict(container):
		if len(args) > 3 {
			return nil, ctx.NewError("Error: expected 2 or 3 arguments")
		}

		d, err := tokens.AssertKeyDict(container)
		if err != nil {
			panic(err)
		}

		// can be int or string, or perhaps another key-type
		key := args[1]

		if !(tokens.IsInt(key) || tokens.IsString(key)) {
			errCtx := key.Context()
			err := errCtx.NewError("Error: expected int or string")
			err.AppendContextString("Info: needed here", ctx)
			return nil, err
		}

		if value, ok := d.Get(key); !ok {
			if len(args) == 3 {
				return args[2], nil
			} else {
				errCtx := key.Context()
				err := errCtx.NewError("Error: key not found in dict (" + key.Dump("") + ")")
				err.AppendContextString("Info: used here", ctx)
				return nil, err
			}
		} else {
			return value, nil
		}
	default:
		errCtx := container.Context()
		err := errCtx.NewError("Error: not a container (" + reflect.TypeOf(container).String() + ")")
		return nil, err
	}
}
