package prototypes

import (
  "../values"

  "../../context"
)

var FontFaceSet *BuiltinPrototype = allocBuiltinPrototype()

func generateFontFaceSetPrototype() bool {
	*FontFaceSet = BuiltinPrototype{
		"FontFaceSet", nil,
		map[string]BuiltinFunction{
			"ready":   NewGetterFunction(func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
        return NewVoidPromise(ctx)
      }),
		},
		nil,
	}

	return true
}

var _FontFaceSetOk = generateFontFaceSetPrototype()

