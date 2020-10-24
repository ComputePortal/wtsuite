package prototypes

import (
	"../values"

	"../../context"
)

type EvalConstructorType func(stack values.Stack, args []values.Value, ctx context.Context) (values.Value, error)
type GenerateInstanceType func(stack values.Stack, keys []string, args []values.Value,
	ctx context.Context) (values.Value, error)

type BuiltinConstructor interface {
	hasConstructorFn() bool
	EvalConstructor(stack values.Stack, args []values.Value, ctx context.Context) (values.Value, error)
	hasInstanceGenerator() bool
	GenerateInstance(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error)
}

type BuiltinConstructorData struct {
	fn    EvalConstructorType
	fnGen GenerateInstanceType
}

func NewConstructor(c ArgCheck, output values.Prototype) *BuiltinConstructorData {
	if output == nil {
		panic("nil output not allowed")
	}

	return &BuiltinConstructorData{
		func(stack values.Stack, args []values.Value, ctx context.Context) (values.Value, error) {
			if err := CheckInputs(c, args, ctx); err != nil {
				return nil, err
			}

			return NewInstance(output, ctx), nil
		},
		nil,
	}
}

func NewConstructorFunction(fn EvalConstructorType) *BuiltinConstructorData {
	return &BuiltinConstructorData{
		fn,
		nil,
	}
}

func NewConstructorGenerator(c ArgCheck, output values.Prototype, fnGen GenerateInstanceType) *BuiltinConstructorData {
	return &BuiltinConstructorData{
		func(stack values.Stack, args []values.Value, ctx context.Context) (values.Value, error) {
			if err := CheckInputs(c, args, ctx); err != nil {
				return nil, err
			}

			return NewInstance(output, ctx), nil
		},
		fnGen,
	}
}

func NewConstructorGeneratorFunction(fn EvalConstructorType, fnGen GenerateInstanceType) *BuiltinConstructorData {
	return &BuiltinConstructorData{fn, fnGen}
}

func NewConstructorNoContentGenerator(c ArgCheck, output values.Prototype) *BuiltinConstructorData {
	return &BuiltinConstructorData{
		func(stack values.Stack, args []values.Value, ctx context.Context) (values.Value, error) {
			if err := CheckInputs(c, args, ctx); err != nil {
				return nil, err
			}

			return NewInstance(output, ctx), nil
		},
		func(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
			if keys != nil || args != nil {
				return nil, ctx.NewError("Error: unexpected content")
			}

			return NewInstance(output, ctx), nil
		},
	}
}

func NewConstructorFunctionNoContentGenerator(fn EvalConstructorType, output values.Prototype) *BuiltinConstructorData {
	return &BuiltinConstructorData{
		fn,
		func(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
			if keys != nil || args != nil {
				return nil, ctx.NewError("Error: unexpected content")
			}

			return NewInstance(output, ctx), nil
		},
	}
}

func NewNoContentGenerator(output values.Prototype) *BuiltinConstructorData {
	return &BuiltinConstructorData{
		nil,
		func(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
			if keys != nil || args != nil {
				return nil, ctx.NewError("Error: unexpected content")
			}

			return NewInstance(output, ctx), nil
		},
	}
}

func (c *BuiltinConstructorData) hasConstructorFn() bool {
	return c.fn != nil
}

func (c *BuiltinConstructorData) EvalConstructor(stack values.Stack, args []values.Value, ctx context.Context) (values.Value, error) {
	return c.fn(stack, args, ctx)
}

func (c *BuiltinConstructorData) hasInstanceGenerator() bool {
	return c.fnGen != nil
}

func (c *BuiltinConstructorData) GenerateInstance(stack values.Stack, keys []string, args []values.Value,
	ctx context.Context) (values.Value, error) {

	// c.fnGen != nil should've been checked before (using the HasInstanceGenerator() function)
	return c.fnGen(stack, keys, args, ctx)
}
