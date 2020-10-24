package js

import (
	"./values"

	"../context"
)

var ALLOW_DUMMY_VIEW_INTERFACE = false

type GlobalStack struct {
	// all values can be saved here
	cache     *CacheStack
	RootStack
}

func NewGlobalStack(cs *CacheStack) *GlobalStack {
	return &GlobalStack{
		cs,
		newRootStack(),
	}
}

func NewDummyStack() values.Stack {
	return NewGlobalStack(nil)
}

func (s *GlobalStack) Parent() values.Stack {
	return nil
}

func (s *GlobalStack) GetReturn(ctx context.Context) (values.Value, error) {
	return nil, ctx.NewError("Error: not inside a function")
}

func (s *GlobalStack) SetReturn(v values.Value, ctx context.Context) error {
	return ctx.NewError("Error: not inside a function")
}

func (s *GlobalStack) IsRecursive(fn interface{}) (bool, values.Value) {
	return false, nil
}

func (s *GlobalStack) GetViewInterface(args ...string) values.ViewInterface {
  if ALLOW_DUMMY_VIEW_INTERFACE {
    return values.NewDummyViewInterface()
  } else {
    panic("shouldn't get this, ViewInterfaceStack should be below")
  }
}

// TODO: remove this too
func (s *GlobalStack) ResolveAwait(t interface{}) (values.Value, bool, error) {
	return nil, false, nil
}

func (s *GlobalStack) SetGeneratedInstance(ptr interface{}, inst *values.Instance) {
	if s.cache != nil { // otherwise this GlobalStack is just used as DummyStack
		s.cache.SetGeneratedInstance(ptr, inst)
	}
}

func (s *GlobalStack) GetGeneratedInstance(ptr interface{}) (*values.Instance, bool) {
	if s.cache != nil {
		return s.cache.GetGeneratedInstance(ptr)
	} else {
		return nil, false
	}
}

func (s *GlobalStack) SetCacheValue(ptr interface{}, val values.Value) {
	if s.cache != nil {
		s.cache.SetCacheValue(ptr, val)
	}
}

func (s *GlobalStack) GetCacheValue(ptr interface{}) (values.Value, bool) {
	if s.cache != nil {
		return s.cache.GetCacheValue(ptr)
	} else {
		return nil, false
	}
}
