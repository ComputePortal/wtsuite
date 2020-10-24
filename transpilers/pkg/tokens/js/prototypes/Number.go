package prototypes

import (
	"../values"

	"../../context"
)

var Number *BuiltinPrototype = allocBuiltinPrototype()

func NewNumber(ctx context.Context) values.Value {
	return values.NewInstance(Number, values.NewNumberProperties(false, 0.0, ctx), ctx)
}

func NewLiteralNumber(v float64, ctx context.Context) values.Value {
	return values.NewInstance(Number, values.NewNumberProperties(true, v, ctx), ctx)
}

func generateNumberPrototype() bool {
	*Number = BuiltinPrototype{
		"Number", nil,
		map[string]BuiltinFunction{
			"EPSILON":           NewStaticGetter(Number),
			"MAX_SAFE_INTEGER":  NewStaticGetter(Int),
			"MAX_VALUE":         NewStaticGetter(Number),
			"MIN_SAFE_INTEGER":  NewStaticGetter(Int),
			"MIN_VALUE":         NewStaticGetter(Number),
			"NaN":               NewStaticGetter(Number),
			"NEGATIVE_INFINITY": NewStaticGetter(Number),
			"POSITIVE_INFINITY": NewStaticGetter(Number),

			"isFinite":      NewStatic(&Any{}, Boolean),
			"isInteger":     NewStatic(&Any{}, Boolean),
			"isNaN":         NewStatic(&Any{}, Boolean),
			"isSafeInteger": NewStatic(&Any{}, Boolean),
			"parseFloat":    NewStatic(&Or{Number, String}, Number),
			"parseInt":      NewStatic(&And{&Or{Number, String}, &Opt{Int}}, Int),
			"toExponential": NewNormal(&Opt{Int}, String),
			"toFixed":       NewNormal(&Opt{Int}, String),
			"toLocalString": NewNormal(&Rest{&Any{}}, String),
			"toPrecision":   NewNormal(&Opt{Int}, String),
			"toString":      NewNormal(&Opt{Int}, String),
		},
		NewConstructorGenerator(&Or{Number, String}, Number,
			func(stack values.Stack, keys []string, args []values.Value,
				ctx context.Context) (values.Value, error) {
				if keys != nil || args != nil {
					return nil, ctx.NewError("Error: unexpected content types")
				}
				return NewInstance(Number, ctx), nil
			}),
	}

	return true
}

var _NumberOk = generateNumberPrototype()
