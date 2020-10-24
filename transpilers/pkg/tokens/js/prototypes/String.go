package prototypes

import (
	"../values"

	"../../context"
)

var String *StringPrototype = allocStringPrototype()

type StringPrototype struct {
	BuiltinPrototype
}

func allocStringPrototype() *StringPrototype {
	return &StringPrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

// exactly the same as BuiltinPrototype.Check, but with *StringPrototype receiver
func (p *StringPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *StringPrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*StringPrototype); ok {
		if other == p {
			return true
		} else {
			return false
		}
	} else {
		return p == other_
	}
}

func (p *StringPrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	if typeChildren != nil {
		return nil, ctx.NewError("Error: " + p.Name() + " can't have content types")
	}

	newV, ok := v.ChangeInstanceInterface(p, false, true)
	if !ok {
		return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + p.Name())
	}

	return newV, nil
}

func NewString(ctx context.Context) values.Value {
	return values.NewInstance(String, values.NewStringProperties(false, "", ctx), ctx)
}

func NewLiteralString(value string, ctx context.Context) values.Value {
	return values.NewInstance(String, values.NewStringProperties(true, value, ctx), ctx)
}

func (p *StringPrototype) GetIndex(stack values.Stack, this *values.Instance,
	index values.Value, ctx context.Context) (values.Value, error) {
	str, okS := this.LiteralStringValue()
	idx, okI := this.LiteralIntValue()

	if okS && okI {
		if idx < 0 || idx > len(str)-1 {
			return nil, ctx.NewError("Error: index out of range")
		}

		return NewLiteralString(str[idx:idx+1], ctx), nil
	} else {
		return NewString(ctx), nil
	}
}

func (p *StringPrototype) SetIndex(stack values.Stack, this *values.Instance,
	index values.Value, arg values.Value, ctx context.Context) error {
	return ctx.NewError("Error: can't set string character via indexing")
}

func (p *StringPrototype) LoopForIn(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	return ctx.NewError("Error: can't iterate over String using 'in' (hint: use regular 'for', or 'for of' if you are only interested in the chars")
}

func (p *StringPrototype) LoopForOf(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	if this == nil {
		return fn(NewString(ctx))
	}

	props := values.AssertStringProperties(this.Properties())

	if str, ok := props.LiteralValue(); ok {
		for i := 0; i < len(str); i++ {
			if err := fn(NewLiteralString(str[i:i+1], ctx)); err != nil {
				return err
			}
		}

		return nil
	} else {
		return fn(NewString(ctx))
	}
}

func generateStringPrototype() bool {
	*String = StringPrototype{BuiltinPrototype{
		"String", nil,
		map[string]BuiltinFunction{
			"charAt":        NewNormal(Int, String),
			"charCodeAt":    NewNormal(Int, Int),
			"codePointAt":   NewNormal(Int, Int),
			"concat":        NewNormal(&Rest{String}, String),
			"endsWith":      NewNormal(&And{String, &Opt{Int}}, Boolean),
			"fromCharCode":  NewStatic(&Rest{Int}, String),
			"fromCodePoint": NewStatic(&Rest{Int}, String),
			"includes":      NewNormal(&And{String, &Opt{Int}}, String),
			"indexOf":       NewNormal(&And{String, &Opt{Int}}, Int),
			"lastIndexOf":   NewNormal(&And{String, &Opt{Int}}, Int),
			"length":        NewGetter(Int),
			"localeCompare": NewNormal(&Rest{&Any{}}, Int),
			"match": NewNormalFunction(RegExp,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					content := NewString(ctx)
					return NewArray([]values.Value{content}, ctx), nil
				}),
			"normalize": NewNormal(&Opt{String}, String),
			"padEnd":    NewNormal(&And{Int, &Opt{String}}, String),
			"padStart":  NewNormal(&And{Int, &Opt{String}}, String),
			"repeat":    NewNormal(Int, String),
			"replace":   NewNormal(&And{&Or{RegExp, String}, &Or{String, &Function{}}}, String),
			"search":    NewNormal(&Or{RegExp, String}, Int),
			"slice":     NewNormal(&And{Int, &Opt{Int}}, String),
			"split": NewNormalFunction(&And{&Or{RegExp, String}, &Opt{Int}},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {

					content := NewString(ctx)
					return NewArray([]values.Value{content}, ctx), nil
				}),
			"startsWith":        NewNormal(&And{String, &Opt{Int}}, Boolean),
			"substring":         NewNormal(&And{Int, &Opt{Int}}, String),
			"toLocaleLowerCase": NewNormal(&Rest{&Any{}}, String),
			"toLocaleUpperCase": NewNormal(&Rest{&Any{}}, String),
			"toLowerCase":       NewNormal(&None{}, String),
			"toString":          NewNormal(&None{}, String),
			"toUpperCase":       NewNormal(&None{}, String),
			"trim":              NewNormal(&None{}, String),
			"trimLeft":          NewNormal(&None{}, String),
			"trimRight":         NewNormal(&None{}, String),
		},
		NewConstructorGenerator(&Any{}, String,
			func(stack values.Stack, keys []string, args []values.Value,
				ctx context.Context) (values.Value, error) {
				if keys != nil || args != nil {
					return nil, ctx.NewError("Error: unexpected content types")
				}
				return NewString(ctx), nil
			}),
	},
	}

	return true
}

var _StringOk = generateStringPrototype()
