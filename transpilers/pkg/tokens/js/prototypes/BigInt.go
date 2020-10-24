package prototypes

import (
	"../values"

	"../../context"
)

var BigInt *BuiltinPrototype = allocBuiltinPrototype()

func NewBigInt(ctx context.Context) values.Value {
	return values.NewInstance(BigInt, values.NewIntProperties(false, 0, ctx), ctx)
}

func generateBigIntPrototype() bool {
	*BigInt = BuiltinPrototype{
		"BigInt", Int,
		map[string]BuiltinFunction{},
		NewConstructor(&Any{}, BigInt),
	}

	return true
}

var _BigIntOk = generateBigIntPrototype()
