package prototypes

import (
	"fmt"

	"../values"

	"../../context"
)

type SetPrototype struct {
	BuiltinPrototype
}

var Set *SetPrototype = allocSetPrototype()

func allocSetPrototype() *SetPrototype {
	return &SetPrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

func (p *SetPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *SetPrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*SetPrototype); ok {
		if other == p {
			return true
		} else {
			return false
		}
	} else {
		_, ok = other_.IsImplementedBy(p)
		return ok
	}
}

func (p *SetPrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	newV_, ok := v.ChangeInstanceInterface(p, false, true)
	if !ok {
		return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + p.Name())
	}

	newV, ok := newV_.(*values.Instance)
	if !ok {
		panic("unexpected")
	}

	if typeChildren == nil {
		return newV, nil
	} else {
		if len(typeChildren) != 1 {
			return nil, ctx.NewError(fmt.Sprintf("Error: Set expects 1 type child, got %d", len(typeChildren)))
		}

		typeChild := typeChildren[0]

		// now cast all the items
		props := values.AssertSetProperties(newV.Properties())
		items := props.GetItems()
		newVProps := values.NewSetProperties(ctx)
		newV = values.NewInstance(Set, newVProps, ctx)
		for _, item := range items {
			newItem, err := item.Cast(typeChild, ctx)
			if err != nil {
				return nil, err
			}

			newVProps.AppendItem(newItem)
		}

		return newV, nil
	}
}

func NewSet(items []values.Value, ctx context.Context) values.Value {
	props := values.NewSetProperties(ctx)

	for _, item := range items {
		props.AppendItem(item)
	}

	return values.NewInstance(Set, props, ctx)
}

func (p *SetPrototype) LoopForOf(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	if this == nil {
		return nil
	}

	props := values.AssertSetProperties(this.Properties())

	items := props.GetItems()

	return fn(values.NewMulti(items, ctx))
}

func generateSetPrototype() bool {
	hasDeleteFn := func(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) (values.Value, error) {
		return NewBoolean(ctx), nil
	}

	*Set = SetPrototype{BuiltinPrototype{
		"Set", nil,
		map[string]BuiltinFunction{
			"add": NewMethodLikeNormalFunction(&Any{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {

					props := values.AssertSetProperties(this.Properties())

					props.AppendItem(args[0])

					return this, nil
				}),
			"delete": NewMethodLikeNormalFunction(&Any{}, hasDeleteFn),
			"has":    NewNormalFunction(&Any{}, hasDeleteFn),
			"size":   NewGetter(Int),
		},
		NewConstructorGeneratorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if len(args) != 0 {
				return nil, ctx.NewError("Error: expected 0 arguments")
			}

			return values.NewInstance(Set, values.NewSetProperties(ctx), ctx), nil
		}, func(stack values.Stack, keys []string, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if keys != nil {
				return nil, ctx.NewError("Error: unexpected keyed content type for Set")
			}

			if args == nil {
				return NewSet([]values.Value{}, ctx), nil
			} else if len(args) == 1 {
				return NewSet(args, ctx), nil
			} else {
				return nil, ctx.NewError("Error: expected single content typed entry for Set")
			}
		}),
	}}

	return true
}

var _SetOk = generateSetPrototype()
