package prototypes

import (
	"../values"

	"../../context"
)

var Response *BuiltinPrototype = allocBuiltinPrototype()

func generateResponsePrototype() bool {
	*Response = BuiltinPrototype{
		"Response", nil,
		map[string]BuiltinFunction{
			"ok":         NewGetter(Boolean),
			"status":     NewGetter(Int),
			"statusText": NewGetter(String),
			"blob": NewNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewResolvedPromise(NewInstance(Blob, ctx), ctx)
				}),
			"json": NewNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewResolvedPromise(NewObject(nil, ctx), ctx)
				}),
		},
		nil,
	}

	return true
}

var _ResponseOk = generateResponsePrototype()
