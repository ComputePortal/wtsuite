package prototypes

import (
	"../values"

	"../../context"
)

var TextEncoder *BuiltinPrototype = allocBuiltinPrototype()

func generateTextEncoder() bool {
	*TextEncoder = BuiltinPrototype{
		"TextEncoder", nil,
		map[string]BuiltinFunction{
			// illegal to create empty array
			"encode": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewAltArray(Uint8Array, []values.Value{NewInt(ctx)}, ctx), nil
				}),
		},
		NewConstructor(&None{}, TextEncoder),
	}

	return true
}

var _TextEncoderOk = generateTextEncoder()
