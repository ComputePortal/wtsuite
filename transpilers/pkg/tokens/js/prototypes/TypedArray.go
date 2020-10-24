package prototypes

import (
	"../values"

	"../../context"
)

var TypedArray *BuiltinPrototype = allocBuiltinPrototype()

func genTypedArrayConstructor(proto values.Prototype, itemProto values.Prototype) EvalConstructorType {
	return func(stack values.Stack, args []values.Value,
		ctx context.Context) (values.Value, error) {

		if err := CheckInputs(&Or{Int, &Or{Array, ArrayBuffer}}, args, ctx); err != nil {
			return nil, err
		}

		if args[0].IsInstanceOf(Array) {
			if item, err := args[0].GetIndex(stack, NewInt(ctx), ctx); err == nil {
				if !item.IsInstanceOf(itemProto) {
					return nil, ctx.NewError("Error: input array doesn't contain just Ints")
				}
			} else {
				panic("unexpected")
			}
		}

		content := NewInstance(itemProto, ctx)
		return NewAltArray(proto, []values.Value{content}, ctx), nil
	}
}

func genTypedArrayMemberMap(proto values.Prototype, itemProto values.Prototype) map[string]BuiltinFunction {
	return map[string]BuiltinFunction{
		"slice": NewNormalFunction(&And{&Opt{Int}, &Opt{Int}},
			func(stack values.Stack, this *values.Instance,
				args []values.Value, ctx context.Context) (values.Value, error) {
				item, err := this.GetIndex(stack, NewInt(ctx), ctx)
				if err != nil {
					panic("unexpected")
				}

				return NewAltArray(proto, []values.Value{item}, ctx), nil
			}),
		"from": NewStaticFunction(&Or{Set, Array},
			func(stack values.Stack, this *values.Instance, args []values.Value,
				ctx context.Context) (values.Value, error) {
				args_ := values.UnpackMulti(args)
				arg, ok := args_[0].(*values.Instance)
				if !ok {
					return nil, ctx.NewError("Error: unable to create array from something with unknown content")
				}

				if arg.IsInstanceOf(Set) {
					props := values.AssertSetProperties(arg.Properties())
					items := props.GetItems()
					for _, item := range items {
						if !item.IsInstanceOf(itemProto) {
							return nil, ctx.NewError("Error: expected Set<" + itemProto.Name() + ">" +
								", got Set <" + item.TypeName() + ">")
						}
					}

					return NewAltArray(proto, items, ctx), nil
				} else if arg.IsInstanceOf(Array) {
					props := values.AssertArrayProperties(arg.Properties())

					item := props.GetItem()
					if !item.IsInstanceOf(itemProto) {
						return nil, ctx.NewError("Error: expected Array<" + itemProto.Name() + ">" +
							", got Array<" + item.TypeName() + ">")
					}

					return NewAltArray(proto, []values.Value{item}, ctx), nil
				} else {
					panic("expected Set or Array")
				}
			}),
	}
}

func generateTypedArrayPrototype() bool {
	*TypedArray = BuiltinPrototype{
		"TypedArray", Array,
		map[string]BuiltinFunction{
			// TODO: check that content is same type
			"buffer": NewGetter(ArrayBuffer),
			"set":    NewNormal(&And{Array, &Opt{Int}}, nil),
			"subarray": NewNormalFunction(&And{&Opt{Int}, &Opt{Int}}, // faster than slice
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					dummyIdx := NewInt(ctx)
					item, err := this.GetIndex(stack, dummyIdx, ctx)
					if err != nil {
						panic(err)
					}

					proto, ok := this.GetInstancePrototype()
					if !ok {
						panic("unexpected")
					}
					return NewAltArray(proto, []values.Value{item}, ctx), nil
				}),
			"BYTES_PER_ELEMENT": NewGetter(Int),
		},
		nil,
	}

	return true
}

var _TypedArrayOk = generateTypedArrayPrototype()
