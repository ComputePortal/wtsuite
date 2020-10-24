package prototypes

import (
	"../values"

	"../../context"
)

var Image *BuiltinPrototype = allocBuiltinPrototype()

func generateImagePrototype() bool {
	*Image = BuiltinPrototype{
		"Image", nil,
		map[string]BuiltinFunction{},
		NewConstructorFunctionNoContentGenerator(func(stack values.Stack, args []values.Value, ctx context.Context) (values.Value, error) {
			if err := CheckInputs(&And{Number, Number}, args, ctx); err != nil {
				return nil, err
			}

			return NewInstance(HTMLImageElement, ctx), nil
		}, HTMLImageElement),
	}

	return true
}

var _ImageOk = generateImagePrototype()
