package values

import (
	"../../context"
)

var ALLOW_CACHING = true

type CachedConstructorPrototype interface {
	Prototype
	EvalCachedConstructor(stack Stack, args []Value, childProto Prototype, ctx context.Context) (Value, error)
}

type Class struct {
	proto Prototype
	ValueData
}

func NewClass(proto Prototype, ctx context.Context) Value {
	return &Class{proto, ValueData{ctx}}
}

func (v *Class) TypeName() string {
	return v.proto.Name() + ".prototype"
}

func (v *Class) Copy(cache CopyCache) Value {
	return NewClass(v.proto, v.Context())
}

func (v *Class) Cast(ntype *NestedType, ctx context.Context) (Value, error) {
	if ntype.Name() == "class" || ntype.Name() == "any" {
		return v, nil
	} else {
		if ntype.Name() == "function" {
			return nil, ctx.NewError("Error: have class, want " + ntype.Name())
		} else {
			return nil, ctx.NewError("Error: have class, want instance of " + ntype.Name())
		}
	}
}

func (v *Class) Merge(other_ Value) Value {
	other_ = UnpackContextValue(other_)

	other, ok := other_.(*Class)
	if !ok {
		return nil
	}

	if v.proto != other.proto {
		return nil
	}

	return v
}

func (v *Class) RemoveLiteralness(all bool) Value {
	return v
}

func (v *Class) EvalConstructor(stack Stack, args []Value,
	ctx context.Context) (Value, error) {
	if ccProto, ok := v.proto.(CachedConstructorPrototype); ok && ALLOW_CACHING {
		return ccProto.EvalCachedConstructor(stack, args, nil, ctx)
	} else {
		return v.proto.EvalConstructor(stack, args, nil, ctx)
	}
}

func (v *Class) EvalAsEntryPoint(stack Stack, ctx context.Context) error {
	return v.proto.EvalAsEntryPoint(stack, ctx)
}

func (v *Class) GetMember(stack Stack, key string, includePrivate bool,
	ctx context.Context) (Value, error) {
	return v.proto.GetMember(stack, nil, key, includePrivate, ctx)
}

func (v *Class) SetMember(stack Stack, key string, arg Value, includePrivate bool,
	ctx context.Context) error {
	return ctx.NewError("Error: cannot set member of static class")
}

func (v *Class) LoopNestedPrototypes(fn func(Prototype)) {
	fn(v.proto)
}

func (v *Class) IsClass() bool {
	return true
}

func (v *Class) GetClassPrototype() (Prototype, bool) {
	return v.proto, true
}

func (v *Class) GetClassInterface() (Interface, bool) {
	return v.proto, true
}
