package prototypes

import (
	"../values"

	"../../context"
)

var NodeJS_process *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_processprototype() bool {
	*NodeJS_process = BuiltinPrototype{
		"process", nil,
		map[string]BuiltinFunction{
			"argv": NewStaticGetterFunction(
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					content := NewString(ctx)

					return NewArray([]values.Value{content}, ctx), nil
				}),
			"argv0":    NewStaticGetter(String),
			"execPath": NewStaticGetter(String),
			"exit":     NewStatic(Int, nil),
		},
		nil,
	}

	return true
}

var _NodeJS_processOk = generateNodeJS_processprototype()
