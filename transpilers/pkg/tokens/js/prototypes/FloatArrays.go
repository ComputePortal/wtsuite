package prototypes

import (
	"../values"
)

var Float32Array *BuiltinPrototype = allocBuiltinPrototype()
var Float64Array *BuiltinPrototype = allocBuiltinPrototype()

func genFloatArrayConstructor(proto values.Prototype) EvalConstructorType {
	return genTypedArrayConstructor(proto, Number)
}

func genFloatArrayMemberMap(proto values.Prototype) map[string]BuiltinFunction {
	return genTypedArrayMemberMap(proto, Number)
}

func generateFloatArrayPrototypes() bool {
	*Float32Array = BuiltinPrototype{
		"Float32Array", TypedArray,
		genFloatArrayMemberMap(Float32Array),
		NewConstructorFunctionNoContentGenerator(
			genFloatArrayConstructor(Float32Array),
			Float32Array,
		),
	}

	*Float64Array = BuiltinPrototype{
		"Float64Array", TypedArray,
		genFloatArrayMemberMap(Float64Array),
		NewConstructorFunctionNoContentGenerator(
			genFloatArrayConstructor(Float64Array),
			Float64Array,
		),
	}

	return true
}

var _FloatArraysOk = generateFloatArrayPrototypes()
