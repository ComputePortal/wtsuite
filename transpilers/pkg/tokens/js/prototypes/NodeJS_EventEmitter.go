package prototypes

import (
	"../values"

	"../../context"
)

var NodeJS_EventEmitter *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_EventEmitterPrototype() bool {
	*NodeJS_EventEmitter = BuiltinPrototype{
		"EventEmitter", nil,
		map[string]BuiltinFunction{
			"addListener": NewNormalFunction(&And{String, &Function{}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					if err := args[1].EvalMethod(stack.Parent(), []values.Value{}, ctx); err != nil {
						return nil, err
					}

					return nil, nil
				}),
		},
		nil,
	}

	return true
}

var _NodeJS_EventEmitterOk = generateNodeJS_EventEmitterPrototype()
