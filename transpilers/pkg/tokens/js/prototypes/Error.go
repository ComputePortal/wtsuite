package prototypes

import (
	"../values"

	"../../context"
)

var Error *BuiltinPrototype = allocBuiltinPrototype()

func NewError(ctx context.Context) *values.Instance {
	return NewInstance(Error, ctx)
}

func generateErrorPrototype() bool {
	*Error = BuiltinPrototype{
		"Error", nil,
		map[string]BuiltinFunction{
			"message": NewGetter(String),
		},
		NewConstructorGenerator(&Opt{String}, Error,
			func(stack values.Stack, keys []string, args []values.Value,
				ctx context.Context) (values.Value, error) {
				if keys != nil || args != nil {
					return nil, ctx.NewError("Error: unexpected content types")
				}
				return NewError(ctx), nil
			}),
	}

	return true
}

var _ErrorOk = generateErrorPrototype()
