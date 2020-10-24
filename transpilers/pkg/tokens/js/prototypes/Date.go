package prototypes

import (
	"../values"

	"../../context"
)

var Date *BuiltinPrototype = allocBuiltinPrototype()

func generateDatePrototype() bool {
	*Date = BuiltinPrototype{
		"Date", nil,
		map[string]BuiltinFunction{
			"getTime":        NewNormal(&None{}, Int),
			"setTime":        NewMethodLikeNormal(Int, Int),
			"toGMTString":    NewNormal(&None{}, String),
			"toLocaleString": NewNormal(String, String),
		},
		NewConstructorGenerator(&Opt{&Or{Number, String}}, Date,
			func(stack values.Stack, keys []string, args []values.Value,
				ctx context.Context) (values.Value, error) {
				if keys != nil || args != nil {
					return nil, ctx.NewError("Error: unexpected content types")
				}

				return NewInstance(Date, ctx), nil
			}),
	}

	return true
}

var _DateOk = generateDatePrototype()
