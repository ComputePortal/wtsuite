package values

import (
	"../../context"
)

type Void struct {
	ValueData
}

func NewVoid(ctx context.Context) Value {
	return &Void{ValueData{ctx}}
}

func (v *Void) TypeName() string {
	return "Void"
}

func (v *Void) Copy(cache CopyCache) Value {
	return NewVoid(v.Context())
}

func (v *Void) Cast(ntype *NestedType, ctx context.Context) (Value, error) {
	// return error, instead of panicking, so we can use this as a 'NeverNull' when checking interfaces implemented by BuiltinPrototypes
	return nil, ctx.NewError("Error: void cant be cast")
}

func (v *Void) Merge(other Value) Value {
	other = UnpackContextValue(other)

	if _, ok := other.(*Void); !ok {
		return nil
	}

	return v
}

func (v *Void) LoopNestedPrototypes(fn func(Prototype)) {
}

func (v *Void) RemoveLiteralness(all bool) Value {
	return v
}

func (v *Void) IsVoid() bool {
	return true
}
