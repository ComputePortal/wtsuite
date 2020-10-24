package prototypes

import (
	"../values"

	"../../context"
)

var RegExp *BuiltinPrototype = allocBuiltinPrototype()

func generateRegExpPrototype() bool {
	*RegExp = BuiltinPrototype{
		"RegExp", nil,
		map[string]BuiltinFunction{
			"exec": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewAltArray(RegExpArray, []values.Value{NewString(ctx)}, ctx), nil
				}),
			"global":     NewGetter(Boolean),
			"ignoreCase": NewGetter(Boolean),
			"lastIndex":  NewGetter(Int),
			"multiline":  NewGetter(Boolean),
			"source":     NewGetter(String),
			"test":       NewNormal(String, Boolean),
		},
		NewConstructor(&And{String, &Opt{String}}, RegExp),
	}

	return true
}

var _RegExpOk = generateRegExpPrototype()
