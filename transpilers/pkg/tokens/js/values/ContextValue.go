package values

import (
	"../../context"
)

// give a different context to a value
type ContextValue struct {
	val Value
	ctx context.Context
}

func NewContextValue(val Value, ctx context.Context) Value {
	if ctxVal, ok := val.(*ContextValue); ok {
		return NewContextValue(ctxVal.val, ctx)
	} else {
		return &ContextValue{val, ctx}
	}
}

func (v *ContextValue) TypeName() string {
	return v.val.TypeName()
}

func (v *ContextValue) Context() context.Context {
	return v.ctx
}

func (v *ContextValue) Copy(cache CopyCache) Value {
	// after a Copy
	return NewContextValue(v.val.Copy(cache), v.Context())
}

func (v *ContextValue) Cast(ntype *NestedType, ctx context.Context) (Value, error) {
	return v.val.Cast(ntype, ctx)
}

func (v *ContextValue) IsInstanceOf(ps ...Prototype) bool {
	return v.val.IsInstanceOf(ps...)
}

func (v *ContextValue) MaybeInstanceOf(p Prototype) bool {
	return v.val.MaybeInstanceOf(p)
}

func UnpackContextValue(val_ Value) Value {
	switch val := val_.(type) {
	case *ContextValue:
		return val.val
	default:
		return val_
	}
}

func (v *ContextValue) Merge(other Value) Value {
	res := v.val.Merge(other)
	if res == nil {
		return nil
	}

	return &ContextValue{res, v.Context()}
}

func (v *ContextValue) RemoveLiteralness(all bool) Value {
	return &ContextValue{v.val.RemoveLiteralness(all), v.Context()}
}

func (v *ContextValue) EvalFunction(stack Stack, args []Value,
	ctx context.Context) (Value, error) {
	return v.val.EvalFunction(stack, args, ctx)
}

func (v *ContextValue) EvalFunctionNoReturn(stack Stack, args []Value,
	ctx context.Context) (Value, error) {
	return v.val.EvalFunctionNoReturn(stack, args, ctx)
}

func (v *ContextValue) EvalMethod(stack Stack, args []Value, ctx context.Context) error {
	return v.val.EvalMethod(stack, args, ctx)
}

func (v *ContextValue) EvalConstructor(stack Stack, args []Value,
	ctx context.Context) (Value, error) {
	return v.val.EvalConstructor(stack, args, ctx)
}

func (v *ContextValue) EvalAsEntryPoint(stack Stack, ctx context.Context) error {
	return v.val.EvalAsEntryPoint(stack, ctx)
}

func (v *ContextValue) GetMember(stack Stack, key string, includePrivate bool,
	ctx context.Context) (Value, error) {
	return v.val.GetMember(stack, key, includePrivate, ctx)
}

func (v *ContextValue) SetMember(stack Stack, key string, arg Value, includePrivate bool,
	ctx context.Context) error {
	return v.val.SetMember(stack, key, arg, includePrivate, ctx)
}

func (v *ContextValue) GetIndex(stack Stack, index Value,
	ctx context.Context) (Value, error) {
	return v.val.GetIndex(stack, index, ctx)
}

func (v *ContextValue) SetIndex(stack Stack, index Value, arg Value,
	ctx context.Context) error {
	return v.val.SetIndex(stack, index, arg, ctx)
}

func (v *ContextValue) LoopForOf(fn func(Value) error, ctx context.Context) error {
	return v.val.LoopForOf(fn, ctx)
}

func (v *ContextValue) LoopForIn(fn func(Value) error, ctx context.Context) error {
	return v.val.LoopForIn(fn, ctx)
}

func (v *ContextValue) LoopNestedPrototypes(fn func(Prototype)) {
	v.val.LoopNestedPrototypes(fn)
}

func (v *ContextValue) IsClass() bool {
	return v.val.IsClass()
}

func (v *ContextValue) IsFunction() bool {
	return v.val.IsFunction()
}

func (v *ContextValue) IsInstance() bool {
	return v.val.IsInstance()
}

func (v *ContextValue) IsNull() bool {
	return v.val.IsNull()
}

func (v *ContextValue) IsVoid() bool {
	return v.val.IsVoid()
}

func (v *ContextValue) IsInterface() bool {
	return v.val.IsInterface()
}

func (v *ContextValue) LiteralBooleanValue() (bool, bool) {
	return v.val.LiteralBooleanValue()
}

func (v *ContextValue) LiteralArrayValues() ([]Value, bool) {
	return v.val.LiteralArrayValues()
}

func (v *ContextValue) LiteralIntValue() (int, bool) {
	return v.val.LiteralIntValue()
}

func (v *ContextValue) LiteralNumberValue() (float64, bool) {
	return v.val.LiteralNumberValue()
}

func (v *ContextValue) LiteralStringValue() (string, bool) {
	return v.val.LiteralStringValue()
}

func (v *ContextValue) GetClassPrototype() (Prototype, bool) {
	return v.val.GetClassPrototype()
}

func (v *ContextValue) GetClassInterface() (Interface, bool) {
	return v.val.GetClassInterface()
}

func (v *ContextValue) GetInstancePrototype() (Prototype, bool) {
	return v.val.GetInstancePrototype()
}

func (v *ContextValue) GetNullPrototype() (Prototype, bool) {
	return v.val.GetNullPrototype()
}

func (v *ContextValue) ChangeInstancePrototype(p Prototype, inPlace bool) (Value, bool) {
	return v.val.ChangeInstancePrototype(p, inPlace)
}

func (v *ContextValue) ChangeInstanceInterface(interf Interface, inPlace bool, checkOuter bool) (Value, bool) {
	return v.val.ChangeInstanceInterface(interf, inPlace, checkOuter)
}
