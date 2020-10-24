package values

import (
	"../../context"
)

type Properties interface {
	GetProperty(key string) (Value, bool)
	SetProperty(key string, v Value) error
  InstanceProperties() map[string]*Instance

	Lock()
	Unlock()
	UnlockForStack(Stack)
	IsLocked(Stack) bool

  // the copyCache isnt used by the Properties directly, only by the values
	Copy(CopyCache) Properties
	Merge(other Properties) Properties
	IsEmpty() bool
	Context() context.Context

	LoopNestedPrototypes(fn func(Prototype))
	NonAllNull() bool
}

type PropertiesData struct {
	props    map[string]Value
	locked   bool
	unlocked []Stack // certain stacks are still allowed to mutate the propertiesdata even if locked==true
	ctx      context.Context
}

func newPropertiesData(ctx context.Context) PropertiesData {
	return PropertiesData{make(map[string]Value), true, []Stack{}, ctx}
}

func newPropertiesDataWithContent(m map[string]Value, ctx context.Context) PropertiesData {
	return PropertiesData{m, true, []Stack{}, ctx}
}

// generic properties
func NewProperties(ctx context.Context) Properties {
	pd := newPropertiesData(ctx)
	return &pd
}

func NewPropertiesWithContent(m map[string]Value, ctx context.Context) Properties {
	pd := newPropertiesDataWithContent(m, ctx)
	return &pd
}

// used by grapher
func (p *PropertiesData) InstanceProperties() map[string]*Instance {
  res := make(map[string]*Instance)

  for key, val_ := range p.props {
    valUnpacked := UnpackContextValue(val_)

    if val, ok := valUnpacked.(*Instance); ok {
      res[key] = val
    }
  }

  return res
}

func (p *PropertiesData) copy(cache CopyCache) PropertiesData {
	props := make(map[string]Value)

	for k, v := range p.props {
		props[k] = v.Copy(cache)
	}

	return PropertiesData{props, p.locked, p.unlocked, p.ctx}
}

func (p *PropertiesData) Copy(cache CopyCache) Properties {
	pd := p.copy(cache)
	return &pd
}

func (p *PropertiesData) merge(otherProps map[string]Value) map[string]Value {
	for k, _ := range otherProps {
		if _, ok := p.props[k]; !ok {
			return nil
		}
	}

	//merged := make(map[string]Value) // this uses a lot of memory in parallel mode!

	for k, val := range p.props {
		otherVal, ok := otherProps[k]
		if !ok {
			return nil
		}

		if v := val.Merge(otherVal); v == nil {
			return nil
		} else {
			//otherProps[k] = v
			if len(otherProps) > len(p.props) {
				if IsAllNull(otherVal) {
					otherProps[k] = v
				}
			} else {
				if IsAllNull(val) {
					p.props[k] = v
				}
			}
			//p.props[k] = v
			//merged[k] = v // uses a lot of memory, so maybe merging on the original is good enough
		}
	}

	if len(otherProps) > len(p.props) {
		return otherProps
	} else {
		return p.props
	}
}

func (p *PropertiesData) Merge(other_ Properties) Properties {
	other, ok := other_.(*PropertiesData)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	return &PropertiesData{mergedProps, true, append(p.unlocked, other.unlocked...), p.Context()}
}

func (p *PropertiesData) Lock() {
	p.locked = true
}

func (p *PropertiesData) Unlock() {
	p.locked = false
}

func (p *PropertiesData) UnlockForStack(s Stack) {
	p.unlocked = append(p.unlocked, s)
}

func (p *PropertiesData) IsLocked(s Stack) bool {
	if p.locked {
		if s == nil {
			return true
		} else {
			for _, us := range p.unlocked {
				if us == s {
					return false
				}
			}
			return true
		}
	} else {
		return false
	}
	return p.locked
}

func (p *PropertiesData) Context() context.Context {
	return p.ctx
}

func (p *PropertiesData) IsEmpty() bool {
	return len(p.props) == 0
}

func (p *PropertiesData) GetProperty(key string) (Value, bool) {
	v, ok := p.props[key]
	return v, ok
}

func (p *PropertiesData) SetProperty(key string, v Value) error {
	if v.IsInterface() {
		panic("object property can't be an interface")
	}

	// we must remove the literalness, because properties are clearly intended for mutation
	v = v.RemoveLiteralness(false)

	// should we set to MultiValue?
	if prev, ok := p.props[key]; ok && prev != v {
		if !prev.IsNull() {
			if prev.IsClass() && !v.IsClass() {
				errCtx := v.Context()
				return errCtx.NewError("Error: variable previously contained a class value")
			} else if prev.IsFunction() && !v.IsFunction() {
				errCtx := v.Context()
				return errCtx.NewError("Error: variable previously contained a function value")
			} else if prev.IsInstance() && !v.IsInstance() {
				errCtx := v.Context()
				return errCtx.NewError("Error: variable previously contained an instance")
			}
		}

		p.props[key] = NewMulti([]Value{prev, v}, v.Context())
	} else {
		p.props[key] = v
	}

	return nil
}

func (p *PropertiesData) RemoveLiteralness(all bool) *PropertiesData {
	newProps := make(map[string]Value)

	for k, v := range p.props {
		newProps[k] = v.RemoveLiteralness(all)
	}

	return &PropertiesData{newProps, p.locked, p.unlocked, p.ctx}
}

func (p *PropertiesData) LoopNestedPrototypes(fn func(Prototype)) {
	for _, v := range p.props {
		v.LoopNestedPrototypes(fn)
	}
}

func (p *PropertiesData) NonAllNull() bool {
	for _, v := range p.props {
		if IsAllNull(v) {
			return false
		}
	}

	return true
}
