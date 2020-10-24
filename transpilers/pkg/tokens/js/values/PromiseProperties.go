package values

import (
	"../../context"
)

type PostponedResolveCheck func(args []Value, ctx context.Context) ([]Value, error)

type PromiseProperties struct {
	resolveFn      Value
	rejectFn       Value
	resolveArgs    [][]Value
	rejectArgs     [][]Value
	postponedCheck PostponedResolveCheck
	PropertiesData
}

func NewPromiseProperties(ctx context.Context) *PromiseProperties {
	return &PromiseProperties{nil, nil, make([][]Value, 0), make([][]Value, 0), nil, newPropertiesData(ctx)}
}

// literals cannot be returned through async promise calls
//  (if we knew the values in advance, why would we do an async call?)
func ClearLiterals(vals []Value) []Value {
	res := make([]Value, len(vals))

	for i, v := range vals {
		ctx := v.Context()
		if _, ok := v.LiteralBooleanValue(); ok {
			proto, ok := v.GetInstancePrototype()
			if !ok {
				panic("unexpected")
			}
			res[i] = NewInstance(proto, NewBooleanProperties(false, false, ctx), ctx)
		} else if _, ok := v.LiteralIntValue(); ok {
			proto, ok := v.GetInstancePrototype()
			if !ok {
				panic("unexpected")
			}
			res[i] = NewInstance(proto, NewIntProperties(false, 0, ctx), ctx)
		} else if _, ok := v.LiteralNumberValue(); ok {
			proto, ok := v.GetInstancePrototype()
			if !ok {
				panic("unexpected")
			}
			res[i] = NewInstance(proto, NewNumberProperties(false, 0.0, ctx), ctx)
		} else if _, ok := v.LiteralStringValue(); ok {
			proto, ok := v.GetInstancePrototype()
			if !ok {
				panic("unexpected")
			}
			res[i] = NewInstance(proto, NewStringProperties(false, "", ctx), ctx)
		} else if items, ok := v.LiteralArrayValues(); ok {
			proto, ok := v.GetInstancePrototype()
			if !ok {
				panic("unexpected")
			}
			res[i] = NewInstance(proto, NewArrayProperties(false, items, ctx), ctx)
		} else {
			res[i] = v
		}
	}

	return res
}

func (p *PromiseProperties) Copy(cache CopyCache) Properties {
	var resolveFnCpy Value = nil
	if p.resolveFn != nil {
		resolveFnCpy = p.resolveFn.Copy(cache)
	}

	var rejectFnCpy Value = nil
	if p.rejectFn != nil {
		rejectFnCpy = p.rejectFn.Copy(cache)
	}

	resolveArgsCpy := make([][]Value, len(p.resolveArgs))
	for i, grp := range p.resolveArgs {
		resolveArgsCpy[i] = copyValueList(grp, cache)
	}

	rejectArgsCpy := make([][]Value, len(p.rejectArgs))
	for i, grp := range p.rejectArgs {
		rejectArgsCpy[i] = copyValueList(grp, cache)
	}

	innerCpy := p.PropertiesData.copy(cache)

	return &PromiseProperties{resolveFnCpy, rejectFnCpy, resolveArgsCpy, rejectArgsCpy,
		p.postponedCheck, innerCpy}
}

func (p *PromiseProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*PromiseProperties)
	if !ok {
		return nil
	}

	/*if p.postponedCheck != other.postponedCheck {
		return nil
	}*/

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	var mergedResolveFn Value = nil
	if p.resolveFn == nil {
		if other.resolveFn != nil {
			return nil
		}
	} else {
		if other.resolveFn == nil {
			return nil
		}

		mergedResolveFn = p.resolveFn.Merge(other.resolveFn)
		if mergedResolveFn == nil {
			return nil
		}
	}

	var mergedRejectFn Value = nil
	if p.rejectFn == nil {
		if other.rejectFn != nil {
			return nil
		}
	} else {
		if other.rejectFn == nil {
			return nil
		}

		mergedRejectFn = p.rejectFn.Merge(other.rejectFn)
		if mergedRejectFn == nil {
			return nil
		}
	}

	if len(p.resolveArgs) != len(other.resolveArgs) {
		return nil
	}

	mergedResolveArgs := make([][]Value, len(p.resolveArgs))
	for i, _ := range p.resolveArgs {
		mergedResolveArg := mergeValueListsSameOrder(p.resolveArgs[i], other.resolveArgs[i])
		if mergedResolveArg == nil {
			return nil
		}

		mergedResolveArgs[i] = mergedResolveArg
	}

	if len(p.rejectArgs) != len(other.rejectArgs) {
		return nil
	}

	mergedRejectArgs := make([][]Value, len(p.rejectArgs))
	for i, _ := range p.rejectArgs {
		mergedRejectArg := mergeValueListsSameOrder(p.rejectArgs[i], other.rejectArgs[i])
		if mergedRejectArg == nil {
			return nil
		}

		mergedRejectArgs[i] = mergedRejectArg
	}

	return &PromiseProperties{mergedResolveFn, mergedRejectFn, mergedResolveArgs, mergedRejectArgs,
		p.postponedCheck, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func (p *PromiseProperties) GetResolveArgs() [][]Value {
	return p.resolveArgs
}

func (p *PromiseProperties) GetRejectArgs() [][]Value {
	return p.rejectArgs
}

func (p *PromiseProperties) ResolveAwait() (Value, error) {
	// all 1 value or all 0 values (and as a consequence nil)
	anyOne := false
	anyNone := false

	vs := make([]Value, 0)

	for _, rArgs := range p.resolveArgs {
		if len(rArgs) == 0 {
			anyNone = true
		} else if len(rArgs) == 1 {
			vs = append(vs, rArgs[0])
			anyOne = true
		} else {
			errCtx := p.Context()
			return nil, errCtx.NewError("Error: expected 0 or 1 resolve values")
		}
	}

	if anyNone && anyOne {
		errCtx := p.Context()
		return nil, errCtx.NewError("Error: expected either 0 or 1, not both possibilities")
	}

  // XXX: do we really need to clear?
	p.ClearResolveArgs()

	if anyOne {
		return NewMulti(vs, p.Context()), nil
	} else {
		return nil, nil
	}
}

func (p *PromiseProperties) ClearResolveArgs() {
	p.resolveArgs = make([][]Value, 0)
}

func (p *PromiseProperties) ClearRejectArgs() {
	p.rejectArgs = make([][]Value, 0)
}

func (p *PromiseProperties) SetResolveArgs(a_ []Value, ctx context.Context) error {
	a := ClearLiterals(a_)

	if p.postponedCheck != nil {
		var err error
		a, err = p.postponedCheck(a, ctx)
		if err != nil {
			return err
		}
	}

	p.resolveArgs = append(p.resolveArgs, a)

	return nil
}

func (p *PromiseProperties) SetRejectArgs(a []Value) {
	p.rejectArgs = append(p.rejectArgs, ClearLiterals(a))
}

func (p *PromiseProperties) SetPostponedCheck(fn PostponedResolveCheck) {
	p.postponedCheck = fn
}

// XXX: it would be better to wrap functions here, in order to clear the literals
//  but for now it is easier the clear the literals directly in the Promise prototype
func (p *PromiseProperties) GetResolveFn() Value {
	return p.resolveFn
}

// XXX: it would be better to wrap functions here, in order to clear the literals
//  but for now it is easier the clear the literals directly in the Promise prototype
func (p *PromiseProperties) GetRejectFn() Value {
	return p.rejectFn
}

func (p *PromiseProperties) SetResolveFn(fn Value) {
	p.resolveFn = fn
}

func (p *PromiseProperties) SetRejectFn(fn Value) {
	p.rejectFn = fn
}

func AssertPromiseProperties(p_ Properties) *PromiseProperties {
	p, ok := p_.(*PromiseProperties)
	if ok {
		return p
	} else {
		panic("not PromiseProperties")
	}
}

func (p *PromiseProperties) LoopNestedPrototypes(fn func(Prototype)) {
	if p.rejectFn != nil {
		p.rejectFn.LoopNestedPrototypes(fn)
	}

	if p.resolveFn != nil {
		p.resolveFn.LoopNestedPrototypes(fn)
	}

	for _, args := range p.rejectArgs {
		for _, a := range args {
			a.LoopNestedPrototypes(fn)
		}
	}

	for _, args := range p.resolveArgs {
		for _, a := range args {
			a.LoopNestedPrototypes(fn)
		}
	}
}
