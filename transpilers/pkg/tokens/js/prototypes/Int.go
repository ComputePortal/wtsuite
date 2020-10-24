package prototypes

import (
	"math"
	"strconv"

	"../values"

	"../../context"
)

var Int *BuiltinPrototype = allocBuiltinPrototype()

func NewInt(ctx context.Context) values.Value {
	return values.NewInstance(Int, values.NewIntProperties(false, 0, ctx), ctx)
}

func NewLiteralInt(value int, ctx context.Context) values.Value {
	return values.NewInstance(Int, values.NewIntProperties(true, value, ctx), ctx)
}

func generateIntPrototype() bool {
	*Int = BuiltinPrototype{
		"Int", Number,
		map[string]BuiltinFunction{
			"toString": NewNormalFunction(&Opt{Int},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					if len(args) == 0 {
						// TODO: handle different radixes
						if litInt, ok := this.LiteralIntValue(); ok {
							return NewLiteralString(strconv.Itoa(litInt), ctx), nil
						}
					}

					return NewString(ctx), nil
				}),
		},
		NewConstructorGeneratorFunction(
			func(stack values.Stack, args []values.Value,
				ctx context.Context) (values.Value, error) {
				// XXX: conversion of literals
				if err := CheckInputs(&Or{String, &Or{Int, Number}}, args, ctx); err != nil {
					return nil, err
				}

				if s, ok := args[0].LiteralStringValue(); ok {
					i, err := strconv.Atoi(s)
					if err != nil {
						errCtx := args[0].Context()
						return nil, errCtx.NewError("Error: can't convert " + s + " to int")
					}

					return NewLiteralInt(i, ctx), nil
				} else if f, ok := args[0].LiteralNumberValue(); ok {
					// TODO: check consistency with javascript
					i := int(math.Round(f))
					return NewLiteralInt(i, ctx), nil
				} else {
					return NewInt(ctx), nil
				}
			}, func(stack values.Stack, keys []string, args []values.Value,
				ctx context.Context) (values.Value, error) {
				if keys != nil || args != nil {
					return nil, ctx.NewError("Error: unexpected content types")
				}
				return NewInt(ctx), nil
			}),
	}

	return true
}

var _IntOk = generateIntPrototype()
