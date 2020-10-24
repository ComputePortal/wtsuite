package prototypes

import (
	"fmt"

	"../values"

	"../../context"
)

var Map *MapPrototype = allocMapPrototype()

type MapPrototype struct {
	BuiltinPrototype
}

func allocMapPrototype() *MapPrototype {
	return &MapPrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

func (p *MapPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *MapPrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*MapPrototype); ok {
		if other == p {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (p *MapPrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
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
		if len(typeChildren) != 2 {
			return nil, ctx.NewError(fmt.Sprintf("Error: Map expects 2 type children, got %d", len(typeChildren)))
		}

		keyType := typeChildren[0]
		itemType := typeChildren[1]

		// now cast all the items
		props := values.AssertMapProperties(newV.Properties())

		newVProps := values.NewMapProperties(ctx)
		newV = values.NewInstance(Map, newVProps, ctx)

		keys := props.GetKeys()
		for _, key := range keys {
			newKey, err := key.Cast(keyType, ctx)
			if err != nil {
				return nil, err
			}

			newVProps.AppendKey(newKey)
		}

		items := props.GetItems()
		for _, item := range items {
			newItem, err := item.Cast(itemType, ctx)
			if err != nil {
				return nil, err
			}

			newVProps.AppendItem(newItem)
		}

		return newV, nil
	}
}

func (p *MapPrototype) LoopForIn(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	return ctx.NewError("Error: for in doesnt work for map")
	/*props := values.AssertMapProperties(this.Properties())

	key := props.GetKey()

	return fn(key)*/
}

func (p *MapPrototype) LoopForOf(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	props := values.AssertMapProperties(this.Properties())

	key := props.GetKey()
	item := props.GetItem()

	return fn(NewLiteralArray([]values.Value{key, item}, ctx))
}

func NewMap(keys []values.Value, items []values.Value, ctx context.Context) values.Value {
  props := values.NewMapProperties(ctx)

  for _, key := range keys {
    props.AppendKey(key)
  }

  for _, item := range items {
    props.AppendItem(item)
  }

  return values.NewInstance(Map, props, ctx)
}

func generateMapPrototype() bool {
	getHasFn := func(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) (values.Value, values.Value, error) {
		props := values.AssertMapProperties(this.Properties())
		return props.GetKey(), props.GetItem(), nil
	}

	*Map = MapPrototype{BuiltinPrototype{
		"Map", nil,
		map[string]BuiltinFunction{
			"clear":  NewNormal(&None{}, nil),
			"delete": NewMethodLikeNormal(&Any{}, Boolean),
			"get": NewNormalFunction(&Any{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					_, item, err := getHasFn(stack, this, args, ctx)
					if err != nil {
						return nil, err
					}

					return item.Copy(values.NewCopyCache()), nil
				}),
			"set": NewNormalFunction(&And{&Any{}, &Any{}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {

					props := values.AssertMapProperties(this.Properties())

					props.AppendKey(args[0])
					props.AppendItem(args[1])

					return nil, nil
				}),

			"has": NewNormalFunction(&Any{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					if _, _, err := getHasFn(stack, this, args, ctx); err != nil {
						return nil, err
					}

					return NewInstance(Boolean, ctx), nil
				}),
			"size": NewGetter(Int),
		},
		NewConstructorGeneratorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {

			if len(args) != 0 {
				return nil, ctx.NewError("Error: expected 0 arguments")
			}

			return NewInstance(Map, ctx), nil
		}, func(stack values.Stack, keys []string, args []values.Value,
      ctx context.Context) (values.Value, error) {
      if keys != nil {
        return nil, ctx.NewError("Error: unexpected keyed content type for Map")
      }

      if args == nil {
        return NewMap([]values.Value{}, []values.Value{}, ctx), nil
      } else if len(args) == 2 {
        return NewMap([]values.Value{args[0]}, []values.Value{args[1]}, ctx), nil
      } else {
        return nil, ctx.NewError("Error: expected 2 type arguments for Map")
      }
    }),
	},
	}

	return true
}

var _MapOk = generateMapPrototype()
