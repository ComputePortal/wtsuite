package values

import (
	"reflect"
	"strings"

	"../../context"
)

type Instance struct {
	orig  Prototype // actual instance child Prototype, used for abstract member implementations
	inner Prototype // used for member calling
	outer Interface // used for type checking (outer must be ancestor of inner)

	props Properties

	//props  map[string]Value // props are private if they start with "_"
	//locked bool             // are just the properties locked?
	removingLiteralness bool // to avoid recursive calls of same cached object

	ValueData
}

func newInstance(orig, inner Prototype, outer Interface, props Properties, ctx context.Context) *Instance {
	return &Instance{orig, inner, outer, props, false, ValueData{ctx}}
}

func NewDummyInstance(proto Prototype, ctx context.Context) *Instance {
	return newInstance(proto, proto, proto, NewProperties(ctx), ctx)
}

func NewInstance(p Prototype, props Properties, ctx context.Context) *Instance {
	if p == nil {
		panic("shouldn't be nil")
	}

	return newInstance(p, p, p, props, ctx)
}

func (v *Instance) TypeName() string {
	return v.outer.Name()
}

func (v *Instance) copy(cache CopyCache) *Instance{
	return newInstance(v.orig, v.inner, v.outer, v.props.Copy(cache), v.ctx)
}

// the only place where copy is really used
func (v *Instance) Copy(cache CopyCache) Value {
  if prev, ok := cache[v]; ok {
    return prev
  } else {
    new := newInstance(v.orig, v.inner, v.outer, nil, v.ctx)
    cache[v] = new

    newProps := v.props.Copy(cache)
    new.props = newProps

    return new
  }
}

func (v *Instance) Cast(ntype *NestedType, ctx context.Context) (Value, error) {
	if _, ok := ntype.Interface().(*AllPrototype); ok || ntype.Name() == "any" {
		// XXX: do we really need to copy?
		return v.Copy(NewCopyCache()), nil
	} else if ntype.Name() == "class" {
		errCtx := ntype.Context()
		return nil, errCtx.NewError("Error: can't cast instance to class")
	} else if ntype.Name() == "function" {
		errCtx := ntype.Context()
		return nil, errCtx.NewError("Error: can't cast instance to function")
	} else if ntype.Name() == "void" {
		errCtx := ntype.Context()
		return nil, errCtx.NewError("Error: can't use void here")
	}

	// Casting logic needs to be in interface, due to differing behaviour depending in ntype.Children
	return ntype.Interface().CastInstance(v, ntype.Children(), ctx)
}

func (v *Instance) IsInstanceOf(ps ...Prototype) bool {
	for _, p := range ps {
		if v.outer == nil || v.inner == nil {
			panic("inner and outer should be set")
		}
		if v.outer.HasAncestor(p) {
			return true
		}
	}

	return false
}

func (v *Instance) MaybeInstanceOf(p Prototype) bool {
	return v.IsInstanceOf(p)
}

func (v *Instance) Merge(other_ Value) Value {
	other_ = UnpackContextValue(other_)

	other, ok := other_.(*Instance)
	if !ok {
		return nil
	}

	if v.inner == nil {
		panic("shouldn't be nil")
	}
	if v.outer == nil {
		panic("shouldn't be nil")
	}
	if other.inner == nil {
		panic("shouldn't be nil")
	}
	if other.outer == nil {
		panic("shouldn't be nil")
	}

	if !(v.inner == other.inner && v.outer == other.outer) {
		return nil
	}

	mergedProps := v.props.Merge(other.props)

	if mergedProps == nil {
		return nil
	}

	return newInstance(v.inner, v.inner, v.outer, mergedProps, v.Context())
}

func (v *Instance) RemoveLiteralness(all bool) Value {
	// the same instance can be nested, so mark that literalness has been removed
	if v.removingLiteralness {
		return v
	}

	v.removingLiteralness = true
	// changes are inplace
	if all {
		if p, ok := v.props.(*StringProperties); ok {
			if _, ok := p.LiteralValue(); ok {
				v.props = p.RemoveLiteralness()
				//return newInstance(v.inner, v.inner, v.outer, p.RemoveLiteralness(), v.Context())
			}
		}

		if p, ok := v.props.(*IntProperties); ok {
			if _, ok := p.LiteralValue(); ok {
				v.props = p.RemoveLiteralness()
				//return newInstance(v.inner, v.inner, v.outer, p.RemoveLiteralness(), v.Context())
			}
		}

		if p, ok := v.props.(*ArrayProperties); ok {
			if _, ok := p.LiteralValues(); ok {
				v.props = p.RemoveLiteralness()
				//return newInstance(v.inner, v.inner, v.outer, p.RemoveLiteralness(), v.Context())
			}
		}

		if p, ok := v.props.(*PropertiesData); ok {
			v.props = p.RemoveLiteralness(all)
			//return newInstance(v.inner, v.inner, v.outer, , v.Context())
		}
	}

	if p, ok := v.props.(*NumberProperties); ok {
		if _, ok := p.LiteralValue(); ok {
			v.props = p.RemoveLiteralness()
			//return newInstance(v.inner, v.inner, v.outer,  v.Context())
		}
	}

	if p, ok := v.props.(*BooleanProperties); ok {
		if _, ok := p.LiteralValue(); ok {
			v.props = p.RemoveLiteralness()
			//return newInstance(v.inner, v.inner, v.outer, v.Context())
		}
	}

	v.removingLiteralness = false

	return v
}

func (v *Instance) hasMember(key string, includePrivate bool) bool {
	if v.outer.HasMember(v, key, includePrivate) {
		return true
	} else if _, ok := v.props.GetProperty(key); ok {
		return true
	} else {
		return false
	}
}

func (v *Instance) assertHasMember(key string, includePrivate bool,
	ctx context.Context) error {
	if _, ok := v.props.GetProperty(key); ok {
		if !(includePrivate || !strings.HasPrefix(key, "_")) {
			err := ctx.NewError("Error: " + v.outer.Name() + "." +
				key + " is private")
			return err
		}
	} else {
		if includePrivate {
			if !v.inner.HasMember(v, key, includePrivate) {
				err := ctx.NewError("Error: " + v.inner.Name() + "." + key + " undefined (have private access)")
				return err
			}
		} else {
			if !v.outer.HasMember(v, key, includePrivate) {
				err := ctx.NewError("Error: " + v.outer.Name() + "." + key + " undefined (don't have private access)")
				return err
			} // if v.outer == v.inner than error will be caught later
		}
	}

	return nil
}

func (v *Instance) GetMember(stack Stack, key string, includePrivate bool,
	ctx context.Context) (Value, error) {
	if prop, ok := v.props.GetProperty(key); ok && (includePrivate || !strings.HasPrefix(key, "_")) {
		return prop, nil
	}

	if err := v.assertHasMember(key, includePrivate, ctx); err != nil {
		return nil, err
	}

	return v.inner.GetMember(stack, v, key, includePrivate, ctx)
}

func (v *Instance) SetMember(stack Stack, key string, arg Value, includePrivate bool,
	ctx context.Context) error {

	if !v.hasMember(key, true) && !v.props.IsLocked(stack) {
		return v.props.SetProperty(key, arg)
	} else if _, ok := v.props.GetProperty(key); ok && (includePrivate || !strings.HasPrefix(key, "_")) {
		return v.props.SetProperty(key, arg)
	} else {
		if err := v.assertHasMember(key, includePrivate, ctx); err != nil {
			return err
		}

		return v.inner.SetMember(stack, v, key, arg, includePrivate, ctx)
	}
}

func (v *Instance) GetIndex(stack Stack, index Value, ctx context.Context) (Value, error) {
	return v.inner.GetIndex(stack, v, index, ctx)
}

func (v *Instance) SetIndex(stack Stack, index Value, arg Value,
	ctx context.Context) error {
	return v.inner.SetIndex(stack, v, index, arg, ctx)
}

func (v *Instance) LoopNestedPrototypes(fn func(Prototype)) {
	fn(v.inner)

	v.props.LoopNestedPrototypes(fn)
}

func (v *Instance) IsInstance() bool {
	return true
}

func (v *Instance) GetInstancePrototype() (Prototype, bool) {
	return v.inner, true
}

// Cast loses type info (acts on outer), whereas this low-level function can add type info (acts on inner and outer) (orig isn't changed
func (v *Instance) ChangeInstancePrototype(newProto Prototype, inPlace bool) (Value, bool) {
	if inPlace {
		v.outer = newProto
		v.inner = newProto

		return v, true
	} else {
		return newInstance(v.orig, newProto, newProto, v.props, v.Context()), true
	}
}

func (v *Instance) ChangeInstanceInterface(newInterf Interface, inPlace bool, checkOuter bool) (Value, bool) {
	// return false if inner doesnt have the newInterf Interface
	if !checkOuter {
		if !v.inner.HasAncestor(newInterf) {
			return nil, false
		}
	} else {
		if !v.outer.HasAncestor(newInterf) {
			return nil, false
		}
	}

	if inPlace {
		v.outer = newInterf

		return v, true
	} else {
		return newInstance(v.inner, v.inner, newInterf, v.props, v.Context()), true
	}
}

func (v *Instance) LoopForIn(fn func(Value) error, ctx context.Context) error {
	return v.inner.LoopForIn(v, fn, ctx)
}

func (v *Instance) LoopForOf(fn func(Value) error, ctx context.Context) error {
	return v.inner.LoopForOf(v, fn, ctx)
}

func (v *Instance) LiteralStringValue() (string, bool) {
	if p, ok := v.props.(*StringProperties); ok {
		return p.LiteralValue()
	} else {
		return v.ValueData.LiteralStringValue()
	}
}

func (v *Instance) LiteralNumberValue() (float64, bool) {
	if p, ok := v.props.(*NumberProperties); ok {
		return p.LiteralValue()
	} else {
		return v.ValueData.LiteralNumberValue()
	}
}

func (v *Instance) LiteralIntValue() (int, bool) {
	if p, ok := v.props.(*IntProperties); ok {
		return p.LiteralValue()
	} else {
		return v.ValueData.LiteralIntValue()
	}
}

func (v *Instance) LiteralBooleanValue() (bool, bool) {
	if p, ok := v.props.(*BooleanProperties); ok {
		return p.LiteralValue()
	} else {
		return v.ValueData.LiteralBooleanValue()
	}
}

func (v *Instance) LiteralArrayValues() ([]Value, bool) {
	if p, ok := v.props.(*ArrayProperties); ok {
		return p.LiteralValues()
	} else {
		return v.ValueData.LiteralArrayValues()
	}
}

func AssertInstance(v_ Value) *Instance {
	if v, ok := v_.(*Instance); ok {
		return v
	} else {
		panic("not an Instance (" + reflect.TypeOf(v_).String())
	}
}

func (v *Instance) Properties() Properties {
	return v.props
}

func (v *Instance) UnlockForStack(s Stack) {
	v.props.UnlockForStack(s)
}

// needed for super
func (v *Instance) SetOriginalPrototype(proto Prototype) {
	v.orig = proto
}

// needed for abstract implementations
func (v *Instance) GetOriginalPrototype() Prototype {
	return v.orig
}

func (v *Instance) CopyInPlace(other *Instance) {
	v.orig = other.orig
	v.inner = other.inner
	v.outer = other.outer
	v.props = other.props
	v.ctx = other.ctx
}

func (v *Instance) NoPropertiesAllNull() bool {
	return v.props.NonAllNull()
}
