package prototypes

import (
	"fmt"

	"../values"

	"../../context"
)

type PromisePrototype struct {
	BuiltinPrototype
}

var Promise *PromisePrototype = allocPromisePrototype()

func allocPromisePrototype() *PromisePrototype {
	return &PromisePrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

func NewResolvedPromise(v values.Value, ctx context.Context) (values.Value, error) {
	props := values.NewPromiseProperties(ctx)
	if err := props.SetResolveArgs([]values.Value{v}, ctx); err != nil {
		return nil, err
	}

	props.SetRejectArgs([]values.Value{NewInstance(Error, ctx)})
	promise := values.NewInstance(Promise, props, ctx)
	return promise, nil
}

func NewVoidPromise(ctx context.Context) (values.Value, error) {
	props := values.NewPromiseProperties(ctx)
	if err := props.SetResolveArgs([]values.Value{}, ctx); err != nil {
		return nil, err
	}

	props.SetRejectArgs([]values.Value{NewInstance(Error, ctx)})
	promise := values.NewInstance(Promise, props, ctx)
	return promise, nil
}

func (p *PromisePrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *PromisePrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*PromisePrototype); ok {
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

func (p *PromisePrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
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
			return nil, ctx.NewError(fmt.Sprintf("Error: Promise expects 1 type child, got %d", len(typeChildren)))
		}

		typeChild := typeChildren[0]

		props := values.AssertPromiseProperties(newV.Properties())
		resArgs := props.GetResolveArgs()
		if len(resArgs) == 0 {
			props.SetPostponedCheck(func(args []values.Value, ctx_ context.Context) ([]values.Value, error) {
				if len(args) == 0 && typeChild.Name() == "void" {
					return args, nil
				} else if len(args) != 1 {
					return nil, ctx_.NewError("Error: Promise resolve expects 1 argument")
				}

				arg := args[0]

				newArg, err := arg.Cast(typeChild, ctx)
				if err != nil {
					return nil, err
				}

				return []values.Value{newArg}, nil
			})

			return newV, nil
		} else {
			newProps_ := props.Copy(values.NewCopyCache())
			newProps := newProps_.(*values.PromiseProperties)

			newProps.ClearResolveArgs()
			for _, resArg_ := range resArgs {
				if len(resArg_) == 0 && typeChild.Name() == "void" {
					if err := newProps.SetResolveArgs([]values.Value{}, ctx); err != nil {
						return nil, err
					}
					continue
				} else if len(resArg_) != 1 {
					errCtx := props.Context()
					return nil, errCtx.NewError("Error: promise resolve takes 1 argument")
				}

				resArg := resArg_[0]

				newResArg, err := resArg.Cast(typeChild, ctx)
				if err != nil {
					return nil, err
				}

				if err := newProps.SetResolveArgs([]values.Value{newResArg}, ctx); err != nil {
					return nil, err
				}
			}

			return values.NewInstance(Promise, newProps, ctx), nil
		}
	}
}

func generatePromisePrototype() bool {
	*Promise = PromisePrototype{BuiltinPrototype{
		"Promise", nil,
		map[string]BuiltinFunction{
			"all": NewStaticFunction(Array,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					item, err := args[0].GetIndex(stack, NewInt(ctx), ctx)
					if err != nil {
						return nil, err
					}

					if !item.IsInstanceOf(Promise) {
						errCtx := args[0].Context()
						return nil, errCtx.NewError("Error: not an array of promises (" + item.TypeName() + ")")
					}

					return NewResolvedPromise(NewArray([]values.Value{values.NewAllNull(ctx)}, ctx), ctx)
				}),
			"catch": NewMethodLikeNormalFunction(&Function{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {

					props := values.AssertPromiseProperties(this.Properties())

					if rejectArgs := props.GetRejectArgs(); len(rejectArgs) == 0 {
						props.SetRejectFn(args[0])
					} else {
						retVals := make([]values.Value, 0)

						for _, rArgs := range rejectArgs {
							retVal, err := args[0].EvalFunctionNoReturn(stack.Parent(), rArgs, ctx)
							if err != nil {
								return nil, err
							}

							if retVal != nil {
								retVals = append(retVals, retVal)
							}
						}

						props.ClearRejectArgs()
						if len(retVals) != 0 {
							retVal := values.NewMulti(retVals, ctx)

							if retVal.IsInstanceOf(Promise) {
								panic("chaining promises not yet supported")
								return retVal, nil
							} else {
								props.SetRejectArgs([]values.Value{retVal})
							}
						}
					}

					return this, nil
				}),
			"then": NewMethodLikeNormalFunction(&Function{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {

					props := values.AssertPromiseProperties(this.Properties())

					if resolveArgs := props.GetResolveArgs(); len(resolveArgs) == 0 {
						props.SetResolveFn(args[0])
					} else {
						retVals := make([]values.Value, 0)

						for _, rArgs := range resolveArgs {
							retVal, err := args[0].EvalFunctionNoReturn(stack.Parent(), rArgs, ctx)
							if err != nil {
								return nil, err
							}

							if retVal != nil {
								retVals = append(retVals, retVal)
							}
						}

						props.ClearResolveArgs()
						if len(retVals) != 0 {
							retVal := values.NewMulti(retVals, ctx)
							if retVal.IsInstanceOf(Promise) {
								panic("chaining of promises not yet supported, because they wouldnt work in async way")
								return retVal, nil
							} else {
								if err := props.SetResolveArgs([]values.Value{retVal}, ctx); err != nil {
									return nil, err
								}
							}
						}
					}

					return this, nil
				}),
			".awaitMethod": NewNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					props := values.AssertPromiseProperties(this.Properties())

					resolveArgs := props.GetResolveArgs()
					if len(resolveArgs) == 0 {
						return nil, NewAsyncRequest(props)
					}

					if _, err := props.ResolveAwait(); err != nil {
						return nil, err
					}

					return nil, nil
				}),
			".awaitFunction": NewNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					props := values.AssertPromiseProperties(this.Properties())

					resolveArgs := props.GetResolveArgs()
					if len(resolveArgs) == 0 {
						return nil, NewAsyncRequest(props)
					}

					return props.ResolveAwait()
				}),
		},
		NewConstructorGeneratorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if err := CheckInputs(&Function{}, args, ctx); err != nil {
				return nil, err
			}

			props := values.NewPromiseProperties(ctx)
			this := values.NewInstance(Promise, props, ctx)

			resolve := values.NewFunctionFunction(func(stack_ values.Stack,
				this_ *values.Instance, args_ []values.Value,
				ctx_ context.Context) (values.Value, error) {

				if len(args_) > 1 {
					return nil, ctx_.NewError("Error: promise resolve expects 0 or 1 arguments")
				}

				if resolveFn := props.GetResolveFn(); resolveFn == nil {
					if err := props.SetResolveArgs(args_, ctx_); err != nil {
						return nil, err
					}
				} else {
					retVal, err := resolveFn.EvalFunctionNoReturn(stack_.Parent(), values.ClearLiterals(args_), ctx_)
					if err != nil {
						return nil, err
					}

					if retVal != nil {
						if retVal.IsInstanceOf(Promise) {
							err := ctx.NewError("chaining of promises not yet supported")
							err.AppendContextString("Info: resolveFn defined here", resolveFn.Context())
							err.AppendContextString("Info: value defined here", retVal.Context())
							return nil, err
						}

						if err := props.SetResolveArgs([]values.Value{retVal}, ctx); err != nil {
							return nil, err
						}
					}
				}

				return nil, nil
			}, stack, this, ctx)

			reject := values.NewFunctionFunction(func(stack_ values.Stack,
				this_ *values.Instance, args_ []values.Value,
				ctx_ context.Context) (values.Value, error) {

				if len(args_) != 1 && len(args_) != 0 {
					return nil, ctx_.NewError("Error: promise reject expects 0 or 1 argument")
				}

				if len(args_) == 1 && !args_[0].IsInstanceOf(Error) {
					return nil, ctx_.NewError("Error: promise reject expects Error argument")
				}

				if rejectFn := props.GetRejectFn(); rejectFn == nil {
					props.SetRejectArgs(args_)
				} else {
					retVal, err := rejectFn.EvalFunctionNoReturn(stack_.Parent(), values.ClearLiterals(args_), ctx_)
					if err != nil {
						return nil, err
					}

					if retVal != nil {
						if retVal.IsInstanceOf(Promise) {
							panic("chaining of promises not yet supported")
						}

						props.SetRejectArgs([]values.Value{retVal})
					}
				}

				return nil, nil
			}, stack, this, ctx)

			// immediately check types of input function
			if err := args[0].EvalMethod(stack.Parent(), []values.Value{resolve, reject},
				ctx); err != nil {

				return nil, err
			}

			return this, nil
		},
			// generator
			func(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
				if keys != nil {
					return nil, ctx.NewError("Error: unexpected keyed content type for Promise")
				}

				if args == nil {
					err := ctx.NewError("Error: content type not specified")
					return nil, err
				}

				if len(args) != 1 {
					// what about void?
					return nil, ctx.NewError("Error: expected single content typed entry for Promise")
				}

				props := values.NewPromiseProperties(ctx)
				promise := values.NewInstance(Promise, props, ctx)

				errInstance, err := Error.GenerateInstance(stack, nil, nil, ctx)
				if err != nil {
					return nil, err
				}

				if err := props.SetResolveArgs([]values.Value{args[0]}, ctx); err != nil {
					return nil, err
				}
				props.SetRejectArgs([]values.Value{errInstance})

				return promise, nil
			},
		),
	}}

	return true
}

var _PromiseOk = generatePromisePrototype()
