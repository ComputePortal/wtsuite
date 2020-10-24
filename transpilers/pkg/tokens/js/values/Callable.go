package values

import (
	"../../context"
)

type Callable interface {
	Name() string // for debugging
	Length() int  // number of arguments, eg. used by Array.map callback

	EvalFunction(stack Stack, this *Instance, args []Value, ctx context.Context) (Value, error)
	EvalFunctionNoReturn(Stack, *Instance, []Value, context.Context) (Value, error) // if method -> return value == nil (used in Promises)
	EvalMethod(Stack, *Instance, []Value, context.Context) error

	EvalAsEntryPoint(Stack, *Instance, context.Context) error
}

//eg. used in prototypes.Promise
type CallableData struct {
	fn func(stack Stack, this *Instance, args []Value, ctx context.Context) (Value, error)
}

func (c *CallableData) Name() string {
	return ""
}

func (c *CallableData) Length() int {
	return -1 // unknown number of arguments
}

func (c *CallableData) EvalFunction(stack Stack, this *Instance, args []Value, ctx context.Context) (Value, error) {
	res, err := c.fn(stack, this, args, ctx)
	if err != nil {
		return nil, err
	}

	if res == nil {
		panic("no return value")
	}

	return res, nil
}

func (c *CallableData) EvalFunctionNoReturn(stack Stack, this *Instance, args []Value,
	ctx context.Context) (Value, error) {
	return c.fn(stack, this, args, ctx)
}

func (c *CallableData) EvalMethod(stack Stack, this *Instance, args []Value, ctx context.Context) error {
	ret, err := c.fn(stack, this, args, ctx)
	if err != nil {
		return err
	}

	if ret != nil {
		return ctx.NewError("Error: unexpected return value")
	}

	return nil
}

// CallableData is only used for builtin functions, so EvalAsEntryPoint can be ignored
func (c *CallableData) EvalAsEntryPoint(stack Stack, this *Instance, ctx context.Context) error {
	return nil
}
