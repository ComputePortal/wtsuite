package values

import (
	"../../context"
)

type Function struct {
	fn    Callable
	stack Stack     // original function stack, can be nil for eg super() in which case incoming stack is used
	this  *Instance // can be nil
	ValueData
}

func NewFunction(fn Callable, stack Stack, this *Instance, ctx context.Context) *Function {
	if this != nil && this.IsClass() && this.IsFunction() {
		panic("impossible1")
	}
	return &Function{fn, stack, this, ValueData{ctx}}
}

func NewFunctionFunction(fn func(Stack, *Instance, []Value, context.Context) (Value, error), stack Stack, this *Instance, ctx context.Context) *Function {
	if this != nil && this.IsClass() && this.IsFunction() {
		panic("impossible2")
	}
	return &Function{&CallableData{fn}, stack, this, ValueData{ctx}}
}

func (v *Function) TypeName() string {
	return "function"
}

func (v *Function) Length() int {
	return v.fn.Length()
}

func (v *Function) Context() context.Context {
	return v.ctx
}

func (v *Function) Copy(cache CopyCache) Value {
  // PAST: dont  copy this (might've lead to circular copying)
  // TODO: can avoid circular dependency of copied function using cache?
	return &Function{v.fn, v.stack, v.this, ValueData{v.ctx}}
}

func (v *Function) Cast(ntype *NestedType, ctx context.Context) (Value, error) {
	if ntype.Name() == "function" || ntype.Name() == "any" {
		return v, nil
	} else {
		if ntype.Name() == "class" {
			return nil, ctx.NewError("Error: have function, want " + ntype.Name())
		} else {
			return nil, ctx.NewError("Error: have function, want instance of " + ntype.Name())
		}
	}
}

func (v *Function) Merge(other_ Value) Value {
	other_ = UnpackContextValue(other_)

	other, ok := other_.(*Function)
	if !ok {
		return nil
	}

	if v.fn != other.fn {
		return nil
	}

	// everything must be equal
	if v.stack != other.stack {
		return nil
	}

	if v.this != other.this {
		return nil
	}

	return v
}

func (v *Function) LoopNestedPrototypes(fn func(Prototype)) {
}

func (v *Function) RemoveLiteralness(all bool) Value {
	return v
}

func (v *Function) EvalFunctionNoReturn(stack Stack, args []Value, ctx context.Context) (Value, error) {
	if v.stack != nil {
		// use DualStack because IsRecursive method relies on different stack
		return v.fn.EvalFunctionNoReturn(NewDualStack(v.stack, stack), v.this, args, ctx)
	} else {
		return v.fn.EvalFunctionNoReturn(stack, v.this, args, ctx)
	}
}

func (v *Function) EvalFunction(stack Stack, args []Value, ctx context.Context) (Value, error) {
	if v.stack != nil {
		return v.fn.EvalFunction(NewDualStack(v.stack, stack), v.this, args, ctx)
	} else {
		return v.fn.EvalFunction(stack, v.this, args, ctx)
	}
}

func (v *Function) EvalAsEntryPoint(stack Stack, ctx context.Context) error {
	if v.this != nil && v.this.IsClass() && v.this.IsFunction() {
		panic("impossible3")
	}

	if v.stack != nil {
		return v.fn.EvalAsEntryPoint(NewDualStack(v.stack, stack), v.this, ctx)
	} else {
		return v.fn.EvalAsEntryPoint(stack, v.this, ctx)
	}
}

func (v *Function) EvalMethod(stack Stack, args []Value, ctx context.Context) error {
	if v.stack != nil {
		// why is this so slow
		return v.fn.EvalMethod(NewDualStack(v.stack, stack), v.this, args, ctx)
	} else {
		return v.fn.EvalMethod(stack, v.this, args, ctx)
	}
}

func (v *Function) IsFunction() bool {
	return true
}
