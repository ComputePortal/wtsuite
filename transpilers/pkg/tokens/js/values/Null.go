package values

import (
	"fmt"
	"os"
	//"reflect"

	"../../context"
)

// Null is the most permissive value, if proto == nil it should pass every test
type Null struct {
	outer Interface
	proto Prototype // can be nil, in which case this Value can be All
	ValueData
}

func NewNull(p Prototype, ctx context.Context) Value {
	return &Null{p, p, ValueData{ctx}}
}

func NewAllNull(ctx context.Context) Value {
	return &Null{nil, nil, ValueData{ctx}}
}

func (v *Null) TypeName() string {
	if v.proto == nil {
		return "All"
	} else {
		return v.proto.Name()
	}
}

func (v *Null) Copy(cache CopyCache) Value {
	return NewNull(v.proto, v.Context())
}

func (v *Null) Cast(ntype *NestedType, ctx context.Context) (Value, error) {
	if _, ok := ntype.Interface().(*AllPrototype); ok || ntype.Name() == "any" {
		return NewNull(v.proto, ctx), nil
	} else if ntype.Name() == "function" {
		return v, nil
	}

	if ntype.Name() == "void" {
		errCtx := ctx
		return nil, errCtx.NewError("Error: can't cast void to " + ntype.Name())
	} else if ntype.Name() == "class" {
		errCtx := ctx
		return nil, errCtx.NewError("Error: can't cast class to " + ntype.Name())
	}

	interf := ntype.Interface()

	if v.proto != nil {
		if !v.outer.HasAncestor(interf) {
			return nil, ctx.NewError("Error: " + interf.Name() + " is not an ancestor of " + v.proto.Name())
		}

		return &Null{interf, v.proto, ValueData{ctx}}, nil
	} else if proto, ok := interf.(Prototype); ok {
		return &Null{interf, proto, ValueData{ctx}}, nil
	} else {
		return NewAllNull(ctx), nil
	}
}

func (v *Null) IsInstanceOf(ps ...Prototype) bool {
	if v.proto == nil {
		return true
	}

	for _, p := range ps {
		if v.proto.HasAncestor(p) {
			return true
		}
	}

	return false
}

func (v *Null) MaybeInstanceOf(p Prototype) bool {
	return v.IsInstanceOf(p)
}

func (v *Null) Merge(other_ Value) Value {
	other_ = UnpackContextValue(other_)

	other, ok := other_.(*Null)
	if !ok {
		return nil
	}

	if v.proto != other.proto {
		return nil
	}

	//return v

	// TODO: also use outer
	if v.outer == other.outer {
		return v
	} else if v.outer == nil && other.outer != nil {
		return &Null{other.outer, v.proto, ValueData{v.Context()}}
	} else if v.outer != nil && other.outer == nil {
		return &Null{v.outer, v.proto, ValueData{v.Context()}}
	} else {
		// XXX: or return v anyway?
		//return v
		return nil
	}
}

func (v *Null) LoopNestedPrototypes(fn func(Prototype)) {
	if v.proto != nil {
		fn(v.proto)
	}
}

func (v *Null) RemoveLiteralness(all bool) Value {
	return v
}

// null::Type() -> null::All
func (v *Null) EvalFunction(stack Stack, args []Value, ctx context.Context) (Value, error) {
	if v.proto == nil || v.proto.Name() == "Function" {
		return NewAllNull(ctx), nil
	} else {
		return nil, ctx.NewError("Error: null of type " + v.proto.Name() + " not callable")
	}
}

func (v *Null) EvalMethod(stack Stack, args []Value, ctx context.Context) error {
	if v.proto == nil || v.proto.Name() == "Function" {
		return nil
	} else {
		return ctx.NewError("Error: null of type " + v.proto.Name() + " not callable")
	}
}

func (v *Null) EvalConstructor(stack Stack, args []Value,
	ctx context.Context) (Value, error) {
	return NewAllNull(ctx), nil
}

func (v *Null) EvalAsEntryPoint(stack Stack, ctx context.Context) error {
	return ctx.NewError("Error: null can never be used as a placeholder for a Class or Function")
}

// TODO: try a little harder than that!
func (v *Null) GetMember(stack Stack, key string, includePrivate bool,
	ctx context.Context) (Value, error) {
	if v.proto == nil {
		return NewAllNull(ctx), nil
	} else {
		// can we get keys/args?
		this_, err := v.proto.GenerateInstance(stack, nil, nil, ctx)
		if err != nil {
			if VERBOSITY < 2 {
				fmt.Fprintf(os.Stderr,
					"Warning: unable to generate instance of %s, evaluating as all null (while evaluating %s)\n",
					v.proto.Name(), key)
				return NewAllNull(ctx), nil
			} else {
				return nil, err
			}
		}

		this := AssertInstance(this_)

		return v.proto.GetMember(stack, this, key, includePrivate, ctx)
	}
}

func (v *Null) SetMember(stack Stack, key string, arg Value, includePrivate bool,
	ctx context.Context) error {
	if v.proto == nil {
		if arg.IsFunction() {
			return ctx.NewError("Error: cannot set function member of null")
		}
		return nil
	} else {
		this_, err := v.proto.GenerateInstance(stack, nil, nil, ctx)
		if err != nil {
			if arg.IsFunction() {
				return ctx.NewError("Error: cannot set function member of null")
			}

			return nil
		}

		this := AssertInstance(this_)

		return v.proto.SetMember(stack, this, key, arg, includePrivate, ctx)
	}
}

func (v *Null) GetIndex(stack Stack, index Value, ctx context.Context) (Value, error) {
	if v.proto == nil {
		return NewAllNull(ctx), nil
	} else {
		return v.proto.GetIndex(stack, nil, index, ctx)
	}
}

func (v *Null) SetIndex(stack Stack, index Value, arg Value, ctx context.Context) error {
	if v.proto == nil {
		return nil
	} else {
		return v.proto.SetIndex(stack, nil, index, arg, ctx)
	}
}

// assume the best
func (v *Null) LoopForOf(fn func(Value) error, ctx context.Context) error {
	return nil
}

// assume the best
func (v *Null) LoopForIn(fn func(Value) error, ctx context.Context) error {
	return nil
}

func (v *Null) IsInstance() bool {
	return true
}

func (v *Null) IsNull() bool {
	return true
}

func (v *Null) IsVoid() bool {
	if v.proto == nil { // TODO: use actual void
		return true // can act as Void in case of unknown recursive call return value
	} else {
		return false
	}
}

func (v *Null) GetInstancePrototype() (Prototype, bool) {
	if v.proto == nil {
		// special All Prototype, false?
		return &AllPrototype{}, true
	} else {
		return v.proto, true
	}
}

func (v *Null) GetNullPrototype() (Prototype, bool) {
	// returned proto can be nil
	return v.proto, true
}

func (v *Null) ChangeInstancePrototype(p Prototype, inPlace bool) (Value, bool) {
	if inPlace {
		v.proto = p
		return v, true
	} else {
		return NewNull(p, v.Context()), true
	}
}

func (v *Null) ChangeInstanceInterface(interf Interface, inPlace bool, checkOuter bool) (Value, bool) {
	if inPlace {
		v.outer = interf
		return v, true
	} else {
		return &Null{interf, v.proto, ValueData{v.Context()}}, true
	}
}

func IsAllNull(v_ Value) bool {
	v_ = UnpackContextValue(v_)
	v, ok := v_.(*Null)
	if ok {
		return v.proto == nil
	}

	return false
}

func IsVoid(v_ Value) bool {
	return v_ == nil || (v_.IsVoid() && !IsAllNull(v_))
}
