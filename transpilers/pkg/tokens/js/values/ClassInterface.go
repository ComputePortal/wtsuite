package values

import (
	"../../context"
)

type ClassInterface struct {
	interf Interface
	ValueData
}

func NewClassInterface(interf Interface, ctx context.Context) Value {
	return &ClassInterface{interf, ValueData{ctx}}
}

func (v *ClassInterface) TypeName() string {
	return v.interf.Name()
}

func (v *ClassInterface) Copy(cache CopyCache) Value {
	return NewClassInterface(v.interf, v.Context())
}

func (v *ClassInterface) Cast(ntype *NestedType, ctx context.Context) (Value, error) {
	return nil, ctx.NewError("Error: can't cast an Interface (can't use an interface as a value)")
}

func (v *ClassInterface) Merge(other_ Value) Value {
	other_ = UnpackContextValue(other_)

	other, ok := other_.(*ClassInterface)
	if !ok {
		return nil
	}

	if v.interf != other.interf {
		return nil
	}

	return v
}

func (v *ClassInterface) LoopNestedPrototypes(fn func(Prototype)) {
}

func (v *ClassInterface) RemoveLiteralness(all bool) Value {
	return v
}

func (v *ClassInterface) IsInterface() bool {
	return true
}

func (v *ClassInterface) GetClassInterface() (Interface, bool) {
	return v.interf, true
}
