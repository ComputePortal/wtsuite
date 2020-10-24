package prototypes

import (
	"fmt"

	"../values"

	"../../context"
)

var Object *ObjectPrototype = allocObjectPrototype()

type ObjectPrototype struct {
	BuiltinPrototype
}

func allocObjectPrototype() *ObjectPrototype {
	return &ObjectPrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

func NewObject(props map[string]values.Value, ctx context.Context) *values.Instance {
	// props == nil -> non-literal object
	return values.NewInstance(Object, values.NewObjectProperties(props, ctx), ctx)
}

func (p *ObjectPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *ObjectPrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*ObjectPrototype); ok {
		if other == p {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (p *ObjectPrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
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
		props := values.AssertObjectProperties(newV.Properties())
		objMap, ok := props.LiteralValues()
		if !ok {
			return nil, ctx.NewError("Error: can't content cast non-literal object")
		}

		if len(typeChildren) != len(objMap) {
			return nil, ctx.NewError(fmt.Sprintf("Error: this Object expects %d type child(ren), got %d", len(objMap), len(typeChildren)))
		}

		newObjMap := make(map[string]values.Value)
		for objKey, objVal := range objMap {
			found := false
			for _, typeChild := range typeChildren {
				if typeChild.Key() == objKey {
					found = true

					newObjVal, err := objVal.Cast(typeChild, ctx)
					if err != nil {
						return nil, err
					}

					newObjMap[objKey] = newObjVal
				}
			}

			if !found {
				errCtx := ctx
				return nil, errCtx.NewError("Error: " + objKey + " not found in type spec")
			}
		}

		return NewObject(newObjMap, ctx), nil
	}
}

func (p *ObjectPrototype) GetIndex(stack values.Stack, this *values.Instance,
	index values.Value, ctx context.Context) (values.Value, error) {
	if !index.IsInstanceOf(String) {
		return nil, ctx.NewError("Error: expected string index")
	}

	if this == nil {
		return values.NewAllNull(ctx), nil
	}

	props := values.AssertObjectProperties(this.Properties())

	if vals, ok := props.LiteralValues(); ok {
		if key, ok := index.LiteralStringValue(); ok {
			if v, hasKey := vals[key]; hasKey {
				return v, nil
			} else {
				return nil, ctx.NewError("Error: LiteralObject doesn't have property '" + key + "'")
			}
		} else {
			return values.NewContextValue(props.GetItem(), ctx), nil
		}
	} else {
		return values.NewContextValue(props.GetItem(), ctx), nil
	}
}

func (p *ObjectPrototype) SetIndex(stack values.Stack, this *values.Instance,
	index values.Value, arg values.Value, ctx context.Context) error {
	if !index.IsInstanceOf(String) {
		return ctx.NewError("Error: expected string index")
	}

	if this == nil {
		return nil
	}

	props := values.AssertObjectProperties(this.Properties())

	if str, ok := index.LiteralStringValue(); ok {
		if items, ok := props.LiteralValues(); ok {
			items[str] = arg
			return nil
		} else {
			props.AppendItem(arg)
		}
	} else {
		props.AppendItem(arg)

		if _, ok := props.LiteralValues(); ok {
			props.RemoveLiteralness()
		}
	}

	return nil
}

func (p *ObjectPrototype) LoopForOf(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	return ctx.NewError("Error: can't loop over Object using 'for of' (hint: use 'for in' instead")
}

func (p *ObjectPrototype) LoopForIn(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	if this == nil {
		return nil
	}

	//props := values.AssertObjectProperties(this.Properties())

	/*if items, ok := props.LiteralValues(); ok {
		for k, v := range items {
			idx := NewLiteralString(k, v.Context())
			if err := fn(idx); err != nil {
				return err
			}
		}

		return nil
	} else {*/
	//v := values.NewContextValue(props.GetItem(), ctx)
	//fmt.Println("NOT LOOPING LITERAL ", reflect.TypeOf(values.UnpackContextValue(v)).String())
	//if _, ok := v.LiteralStringValue(); ok {
	//panic("can't be literal")
	//}
	return fn(NewString(ctx))
	//}
}

func generateObjectPrototype() bool {
	*Object = ObjectPrototype{BuiltinPrototype{
		"Object", nil,
		map[string]BuiltinFunction{
			"assign": NewStaticFunction(&AtLeast{2, &Any{}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					// TODO: take the this properties of the sources and add them to the target

					return args[0].Copy(values.NewCopyCache()), nil
				}),
			"keys": NewStaticFunction(Object,
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {

					// can't convert literalproperties into literalarray because order was lost
					//  (and it's too much effort to create an order preserving map just for this function)
					return NewArray([]values.Value{NewString(ctx)}, ctx), nil
				}),
		},
		NewConstructorGeneratorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if err := CheckInputs(&None{}, args, ctx); err != nil {
				return nil, err
			}

			// props == nil signifies non-literal
			return NewObject(nil, ctx), nil
		},
			func(stack values.Stack, keys []string, args []values.Value,
				ctx context.Context) (values.Value, error) {
				// keys and args can be empty, but must exist, so 'Object<>` is the minimal full type
				if keys == nil || args == nil {
					return nil, ctx.NewError("Error: Object generator expects keyed content types")
				}

				if len(keys) != len(args) {
					panic("lengths must be the same, should've been caught earlier")
				}

				props := make(map[string]values.Value)
				for i, key := range keys {
					arg := args[i]
					props[key] = arg
				}

				return NewObject(props, ctx), nil
			}),
	},
	}

	return true
}

var _ObjectOk = generateObjectPrototype()
