package prototypes

import (
	"../values"

	"../../context"
)

var Math *BuiltinPrototype = allocBuiltinPrototype()

func generateMathPrototype() bool {
	minMaxFn := func(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) (values.Value, error) {
		isInt := true
		for _, arg := range args {
			if !arg.IsInstanceOf(Number) {
				panic("expected number")
			}

			if !arg.IsInstanceOf(Int) {
				isInt = false
			}
		}

		if isInt {
			return NewInstance(Int, ctx), nil
		} else {
			return NewInstance(Number, ctx), nil
		}
	}

	*Math = BuiltinPrototype{
		"Math", nil,
		map[string]BuiltinFunction{
			"E":       NewStaticGetter(Number),
			"LN2":     NewStaticGetter(Number),
			"LN10":    NewStaticGetter(Number),
			"LOG2E":   NewStaticGetter(Number),
			"LOG10E":  NewStaticGetter(Number),
			"PI":      NewStaticGetter(Number),
			"SQRT1_2": NewStaticGetter(Number),
			"SQRT2":   NewStaticGetter(Number),

			"abs":    NewStatic(Number, Number),
			"acos":   NewStatic(Number, Number),
			"acosh":  NewStatic(Number, Number),
			"asin":   NewStatic(Number, Number),
			"asinh":  NewStatic(Number, Number),
			"atan":   NewStatic(Number, Number),
			"atan2":  NewStatic(&And{Number, Number}, Number),
			"atanh":  NewStatic(Number, Number),
			"cbrt":   NewStatic(Number, Number),
			"ceil":   NewStatic(Number, Int),
			"clz32":  NewStatic(Int, Int),
			"cos":    NewStatic(Number, Number),
			"cosh":   NewStatic(Number, Number),
			"exp":    NewStatic(Number, Number),
			"expm1":  NewStatic(Number, Number),
			"floor":  NewStatic(Number, Int),
			"fround": NewStatic(Number, Number),
			"hypot":  NewStatic(&Rest{Number}, Number),
			"imul":   NewStatic(&And{Int, Int}, Int),
			"log":    NewStatic(Number, Number),
			"log10":  NewStatic(Number, Number),
			"log1p":  NewStatic(Number, Number),
			"log2":   NewStatic(Number, Number),
			"max":    NewStaticFunction(&AtLeast{2, Number}, minMaxFn),
			"min":    NewStaticFunction(&AtLeast{2, Number}, minMaxFn),
			"pow": NewStaticFunction(&Or{&And{Number, Number}, &And{Int, Int}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {

					if args[0].IsInstanceOf(Int) && args[1].IsInstanceOf(Int) {
						return NewInstance(Int, ctx), nil
					} else {
						return NewInstance(Number, ctx), nil
					}
				}),
			"random": NewStatic(&None{}, Number),
			"round":  NewStatic(Number, Int),
			"sign":   NewStatic(Number, Int),
			"sin":    NewStatic(Number, Number),
			"sinh":   NewStatic(Number, Number),
			"sqrt":   NewStatic(Number, Number),
			"tan":    NewStatic(Number, Number),
			"tanh":   NewStatic(Number, Number),
			"trunc":  NewStatic(Number, Int),
		},
		nil,
	}

	return true
}

var _MathOk = generateMathPrototype()
