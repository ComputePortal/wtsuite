package prototypes

import (
	"../values"

	"../../context"
)

var Boolean *BuiltinPrototype = allocBuiltinPrototype()

func NewBoolean(ctx context.Context) values.Value {
	return values.NewInstance(Boolean, values.NewBooleanProperties(false, false, ctx), ctx)
}

func NewLiteralBoolean(v bool, ctx context.Context) values.Value {
	return values.NewInstance(Boolean, values.NewBooleanProperties(true, v, ctx), ctx)
}

func generateBooleanPrototype() bool {
	*Boolean = BuiltinPrototype{
		"Boolean", nil,
		map[string]BuiltinFunction{},
		NewConstructorNoContentGenerator(Number, Boolean),
	}

	return true
}

var _BooleanOk = generateBooleanPrototype()
