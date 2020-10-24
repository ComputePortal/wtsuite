package prototypes

import (
	"../values"
)

var Int8Array *BuiltinPrototype = allocBuiltinPrototype()
var Int16Array *BuiltinPrototype = allocBuiltinPrototype()
var Int32Array *BuiltinPrototype = allocBuiltinPrototype()

var Uint8Array *BuiltinPrototype = allocBuiltinPrototype()
var Uint8ClampedArray *BuiltinPrototype = allocBuiltinPrototype()
var Uint16Array *BuiltinPrototype = allocBuiltinPrototype()
var Uint32Array *BuiltinPrototype = allocBuiltinPrototype()

func genIntArrayConstructor(proto values.Prototype) EvalConstructorType {
	return genTypedArrayConstructor(proto, Int)
}

func genIntArrayMemberMap(proto values.Prototype) map[string]BuiltinFunction {
	return genTypedArrayMemberMap(proto, Int)
}

func generateIntArrayPrototypes() bool {
	*Int8Array = BuiltinPrototype{
		"Int8Array", TypedArray,
		genIntArrayMemberMap(Int8Array),
		NewConstructorFunctionNoContentGenerator(
			genIntArrayConstructor(Int8Array),
			Int8Array,
		),
	}

	*Int16Array = BuiltinPrototype{
		"Int16Array", TypedArray,
		genIntArrayMemberMap(Int16Array),
		NewConstructorFunctionNoContentGenerator(
			genIntArrayConstructor(Int16Array),
			Int16Array,
		),
	}

	*Int32Array = BuiltinPrototype{
		"Int32Array", TypedArray,
		genIntArrayMemberMap(Int32Array),
		NewConstructorFunctionNoContentGenerator(
			genIntArrayConstructor(Int32Array),
			Int32Array,
		),
	}

	*Uint8Array = BuiltinPrototype{
		"Uint8Array", TypedArray,
		genIntArrayMemberMap(Uint8Array),
		NewConstructorFunctionNoContentGenerator(
			genIntArrayConstructor(Uint8Array),
			Uint8Array,
		),
	}

	*Uint8ClampedArray = BuiltinPrototype{
		"Uint8ClampedArray", Uint8Array,
		genIntArrayMemberMap(Uint8ClampedArray),
		NewConstructorFunctionNoContentGenerator(
			genIntArrayConstructor(Uint8ClampedArray),
			Uint8ClampedArray,
		),
	}

	*Uint16Array = BuiltinPrototype{
		"Uint16Array", TypedArray,
		genIntArrayMemberMap(Uint16Array),
		NewConstructorFunctionNoContentGenerator(
			genIntArrayConstructor(Uint16Array),
			Uint16Array,
		),
	}

	*Uint32Array = BuiltinPrototype{
		"Uint32Array", TypedArray,
		genIntArrayMemberMap(Uint32Array),
		NewConstructorFunctionNoContentGenerator(
			genIntArrayConstructor(Uint32Array),
			Uint32Array,
		),
	}

	return true
}

var _IntArraysOk = generateIntArrayPrototypes()
