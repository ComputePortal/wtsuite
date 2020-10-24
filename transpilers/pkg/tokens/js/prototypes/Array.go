package prototypes

import (
	"fmt"

	"../values"

	"../../context"
)

type ArrayPrototype struct {
	BuiltinPrototype
}

var Array *ArrayPrototype = allocArrayPrototype()

func allocArrayPrototype() *ArrayPrototype {
	return &ArrayPrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

func (p *ArrayPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *ArrayPrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*ArrayPrototype); ok {
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

func (p *ArrayPrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
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
			return nil, ctx.NewError(fmt.Sprintf("Error: Array expects 1 type child, got %d", len(typeChildren)))
		}

		typeChild := typeChildren[0]

		// now cast all the items
		props := values.AssertArrayProperties(newV.Properties())
		if items, ok := props.LiteralValues(); ok {
			newItems := make([]values.Value, len(items))
			for i, item := range items {
				var err error
				newItems[i], err = item.Cast(typeChild, ctx)
				if err != nil {
					return nil, err
				}
			}

			return values.NewInstance(Array, values.NewArrayProperties(true, newItems, ctx), ctx), nil
		} else {
			item := props.GetItem()
			newItem, err := item.Cast(typeChild, ctx)
			if err != nil {
				return nil, err
			}

			return values.NewInstance(Array, values.NewArrayProperties(false, []values.Value{newItem}, ctx), ctx), nil
		}
	}
}

func NewArray(items []values.Value, ctx context.Context) values.Value {
	return values.NewInstance(Array, values.NewArrayProperties(false, items, ctx), ctx)
}

func NewAltArray(proto values.Prototype, items []values.Value, ctx context.Context) values.Value {
	return values.NewInstance(proto, values.NewArrayProperties(false, items, ctx), ctx)
}

func NewLiteralArray(items []values.Value, ctx context.Context) values.Value {
	return values.NewInstance(Array, values.NewArrayProperties(true, items, ctx), ctx)
}

func (p *ArrayPrototype) GetIndex(stack values.Stack, this *values.Instance,
	index values.Value, ctx context.Context) (values.Value, error) {
	if index == nil {
		panic("index shouldnt be nil")
	}

	if index == nil || !index.IsInstanceOf(Int) {
		return nil, ctx.NewError("Error: expected Int index")
	}

	if this == nil {
		//errCtx := ctx
		//return nil, errCtx.NewError("Error: array content unknown")
		//panic("array content unknown")
		return values.NewAllNull(ctx), nil
	}

	props := values.AssertArrayProperties(this.Properties())

	idx, okI := index.LiteralIntValue()

	if items, ok := props.LiteralValues(); ok && okI {
		if idx < 0 || idx > len(items)-1 {
			if idx >= 0 && len(items) == 0 {
				return values.NewAllNull(ctx), nil
			}
			return nil, ctx.NewError(fmt.Sprintf("Error: index out of range (i=%d, n=%d)\n", idx, len(items)))
		}

		return items[idx], nil
	} else if items, ok := props.LiteralValues(); ok && len(items) == 0 {
		// return All, we can't do anything else
		return values.NewAllNull(ctx), nil
	} else {
		return values.NewContextValue(props.GetItem(), ctx), nil
	}
}

func (p *ArrayPrototype) SetIndex(stack values.Stack, this *values.Instance,
	index values.Value, arg values.Value, ctx context.Context) error {
	if !index.IsInstanceOf(Int) {
		return ctx.NewError("Error: expected Int index")
	}

	if this == nil {
		panic("setting index of nil array")
		return nil
	}

	props := values.AssertArrayProperties(this.Properties())

	props.AppendItem(arg)

	// TODO: can literal index and literal array be used? (what about in for loops?)

	return nil
}

func (p *ArrayPrototype) LoopForIn(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	return ctx.NewError("Error: can't iterate over Array using 'in' (hint: use regular 'for', or 'for of' if your are only interested in the items")
}

func (p *ArrayPrototype) LoopForOf(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	if this == nil {
		return nil
	}

	props := values.AssertArrayProperties(this.Properties())

	/*if items, ok := props.LiteralValues(); ok {
		if VERBOSITY >= 3 || (len(items) > 3 && VERBOSITY >= 1) {
			errCtx := ctx
			err := errCtx.NewError("Warning: evaluating literal loop (hint: use regular for wrap length with Int() to avoid this)")
			fmt.Fprintf(os.Stderr, err.Error())
		}

		for i := 0; i < len(items); i++ {
			if err := fn(items[i]); err != nil {
				return err
			}
		}

		return nil
	} else {*/
	item := props.GetItem()

	return fn(item)
	//}
}

func generateArrayPrototype() bool {
	popShiftFn := func(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) (values.Value, error) {
		dummyIdx := NewInt(ctx)
		if v, err := this.GetIndex(stack, dummyIdx, ctx); err == nil {
			return v, nil
		} else {
			panic("unexpected")
		}
	}

	pushUnshiftFn := func(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) (values.Value, error) {
		dummyIdx := NewInt(ctx)
		if err := this.SetIndex(stack, dummyIdx, args[0], ctx); err != nil {
			return nil, err
		}

		// returns new length
		return NewInt(ctx), nil
	}

	filterFn := func(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) (values.Value, error) {
		dummyIdx := NewInt(ctx)
		item, err := this.GetIndex(stack, dummyIdx, ctx)
		if err != nil {
			panic("unexpected")
		}

		result, err := args[0].EvalFunction(stack.Parent(), []values.Value{item},
			args[0].Context())
		if err != nil {
			return nil, err
		}

		if !result.IsInstanceOf(Boolean) {
			errCtx := args[0].Context()
			return nil, errCtx.NewError("Error: " +
				"function doesn't return a Boolean but a " + result.TypeName())
		}

		return NewBoolean(ctx), nil
	}

	*Array = ArrayPrototype{BuiltinPrototype{
		"Array", nil,
		map[string]BuiltinFunction{
			"concat": NewNormalFunction(&Rest{Array},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					vals := make([]values.Value, 0)

					dummyIdx := NewInt(ctx)
					if val, err := this.GetIndex(stack, dummyIdx, ctx); err == nil {
						vals = append(vals, val)
					} else {
						panic(err)
					}

					for _, arg := range args {
						if val, err := arg.GetIndex(stack, dummyIdx, ctx); err == nil {
							vals = append(vals, val)
						} else {
							panic(err)
						}
					}

					// values can be a list of MultiValues
					return NewArray(vals, ctx), nil
				}),
			"copyWithin": NewNormal(&And{Int, &And{Int, &Opt{Int}}}, nil),
			"every":      NewNormalFunction(&Function{}, filterFn),
			"fill": NewMethodLikeNormalFunction(&And{&Any{}, &Opt{&And{Int, &Opt{Int}}}},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					props := values.AssertArrayProperties(this.Properties())
					props.AppendItem(args[0])
					return this, nil
				}),
			"find": NewNormalFunction(&Function{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return nil, ctx.NewError("Error: use findIndex instead (find returns undefined when nothing is found, which is dumb, should be null)")
				}),
			"findIndex": NewNormalFunction(&Function{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					dummyIdx := NewInt(ctx)
					item, err := this.GetIndex(stack, dummyIdx, ctx)
					if err != nil {
						panic("unexpected")
					}

					result, err := args[0].EvalFunction(stack.Parent(), []values.Value{item}, ctx)
					if err != nil {
						return nil, err
					}

					if !result.IsInstanceOf(Boolean) {
						errCtx := args[0].Context()
						return nil, errCtx.NewError("Error: function doesn't return a bool")
					}

					return NewInt(ctx), nil
				}),
			"filter": NewNormalFunction(&Function{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					dummyIdx := NewInt(ctx)
					item, err := this.GetIndex(stack, dummyIdx, ctx)
					if err != nil {
						panic("unexpected")
					}

					result, err := args[0].EvalFunction(stack.Parent(), []values.Value{item}, ctx)
					if err != nil {
						return nil, err
					}

					if !result.IsInstanceOf(Boolean) {
						errCtx := args[0].Context()
						return nil, errCtx.NewError("Error: function doesn't return a bool")
					}

					return NewArray([]values.Value{item}, ctx), nil
				}),
			"forEach": NewNormalFunction(&Function{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					dummyIdx := NewInt(ctx)
					item, err := this.GetIndex(stack, dummyIdx, ctx)
					if err != nil {
						panic("unexpected")
					}

					if err := args[0].EvalMethod(stack.Parent(), []values.Value{item},
						ctx); err != nil {
						return nil, err
					}

					return nil, nil
				}),
			"from": NewStaticFunction(Set,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					args_ := values.UnpackMulti(args)
					arg, ok := args_[0].(*values.Instance)
					if !ok {
						return nil, ctx.NewError("Error: unable to create array from something with unknown content")
					}

					props := values.AssertSetProperties(arg.Properties())

					return NewArray(props.GetItems(), ctx), nil
				}),
			"indexOf": NewNormalFunction(&And{&Any{}, &Opt{Int}},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					// TODO: return error if type is not at all available in array
					return NewInt(ctx), nil
				}),
			"join": NewNormalFunction(&Opt{String},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					dummyIdx := NewInt(ctx)
					item, err := this.GetIndex(stack, dummyIdx, ctx)
					if err != nil {
						panic("unexpected")
					}

					if !item.IsInstanceOf(String, Number, Boolean) {
						return nil, ctx.NewError("Error: array contents not stringable (" + item.TypeName() + ")")
					}

					return NewString(ctx), nil
				}),
			"length": NewGetterFunction(
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					if this == nil {
						return NewInt(ctx), nil
					}

					props := values.AssertArrayProperties(this.Properties())

					if n, ok := props.Length(); ok {
						return NewLiteralInt(n, ctx), nil
					} else {
						return NewInt(ctx), nil
					}
				}),
			"map": NewNormalFunction(&Function{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					dummyIdx := NewInt(ctx)
					item, err := this.GetIndex(stack, dummyIdx, ctx)
					if err != nil {
						panic("unexpected")
					}

					var result values.Value

					fnArg_ := values.UnpackContextValue(args[0])
					if fnArg, ok := fnArg_.(*values.Function); ok && fnArg.Length() == 2 {
						result, err = fnArg.EvalFunction(stack.Parent(), []values.Value{item, NewInt(ctx)}, ctx)
						if err != nil {
							return nil, err
						}
					} else {
						result, err = args[0].EvalFunction(stack.Parent(), []values.Value{item}, ctx)
						if err != nil {
							return nil, err
						}
					}

					return NewArray([]values.Value{result}, ctx), nil
				}),
			"pop":  NewNormalFunction(&None{}, popShiftFn),
			"push": NewMethodLikeNormalFunction(&Any{}, pushUnshiftFn),
			"reduce": NewNormalFunction(&And{&Function{}, &Opt{&Any{}}},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					item, err := this.GetIndex(stack, NewInt(ctx), ctx)
					if err != nil {
						panic("unexpected")
					}

					if len(args) == 2 {
						if _, err := args[0].EvalFunction(stack.Parent(), []values.Value{args[1], item}, ctx); err != nil {
							return nil, err
						}
					}

					return args[0].EvalFunction(stack.Parent(), []values.Value{item, item}, ctx)
				}),
			"reverse": NewMethodLikeNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					if this != nil {
						props := values.AssertArrayProperties(this.Properties())
						props.ChangeOrder()
					}
					return this, nil
				}),
			"shift": NewNormalFunction(&None{}, popShiftFn),
			"slice": NewNormalFunction(&And{&Opt{Int}, &Opt{Int}},
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
			"some": NewNormalFunction(&Function{}, filterFn),
			"sort": NewMethodLikeNormalFunction(&Opt{&Function{}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {

					props := values.AssertArrayProperties(this.Properties())

					if len(args) == 1 {
						dummyIdx := NewInt(ctx)
						item, err := this.GetIndex(stack, dummyIdx, ctx)
						if err != nil {
							panic("unexpected")
						}

						if sortRes, err := args[0].EvalFunction(stack.Parent(),
							[]values.Value{
                item.Copy(values.NewCopyCache()), 
                item.Copy(values.NewCopyCache()),
              }, ctx); err != nil {
							return nil, err
						} else if !sortRes.IsInstanceOf(Number) {
							return nil, ctx.NewError("Error: sort function doesn't return a number")
						}
					}

					props.ChangeOrder()

					return this, nil
				}),
			"splice": NewMethodLikeNormalFunction(&And{Int, &And{&Opt{Int}, &Opt{&Rest{&Any{}}}}},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					props := values.AssertArrayProperties(this.Properties())

					for i := 2; i < len(args); i++ {
						props.AppendItem(args[i])
					}

					return values.NewInstance(Array, props.RemoveLiteralness(), ctx), nil
				}),
			"unshift": NewMethodLikeNormalFunction(&Any{}, pushUnshiftFn),
		},
		NewConstructorGeneratorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if len(args) < 2 {
				if err := CheckInputs(&Opt{Int}, args, ctx); err != nil {
					return nil, err
				}

				return NewArray([]values.Value{}, ctx), nil
			} else {
				// can contain anything
				// XXX: should all array entries have same base prototype?
				return NewLiteralArray(args, ctx), nil
			}
		}, func(stack values.Stack, keys []string, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if keys != nil {
				return nil, ctx.NewError("Error: unexpected keyed content type for Array")
			}

			if args == nil {
				//return nil, ctx.NewError("Error: can't generate Array without content type")
				return NewArray([]values.Value{}, ctx), nil
			} else if len(args) == 1 {
				return NewArray(args, ctx), nil
			} else {
				return nil, ctx.NewError("Error: expected single content typed entry for Array")
			}
		}),
	},
	}

	return true
}

var _ArrayOk = generateArrayPrototype()
