package values

import (
	"../../context"
)

type DualStack struct {
	valStack Stack // Stack used for values
	recStack Stack // Stack used for recursive checking
}

func NewDualStack(valStack Stack, recStack Stack) *DualStack {
	return &DualStack{valStack, recStack}
}

func (s *DualStack) Parent() Stack {
	return s.valStack.Parent()
}

func (s *DualStack) GetReturn(ctx context.Context) (Value, error) {
	return s.valStack.GetReturn(ctx)
}

func (s *DualStack) SetReturn(v Value, ctx context.Context) error {
	return s.valStack.SetReturn(v, ctx)
}

func (s *DualStack) HasValue(ptr interface{}) bool {
	return s.valStack.HasValue(ptr)
}

func (s *DualStack) GetValue(ptr interface{}, ctx context.Context) (Value, error) {
	return s.valStack.GetValue(ptr, ctx)
}

func (s *DualStack) SetValue(ptr interface{}, v Value,
	allowBranching bool, ctx context.Context) error {
	return s.valStack.SetValue(ptr, v, allowBranching, ctx)
}

func (s *DualStack) IsRecursive(fn interface{}) (bool, Value) {
	return s.recStack.IsRecursive(fn)
}

func (s *DualStack) GetViewInterface(args ...string) ViewInterface {
	return s.valStack.GetViewInterface(args...)
}

func (s *DualStack) ResolveAwait(t interface{}) (Value, bool, error) {
	return s.valStack.ResolveAwait(t)
}

func (s *DualStack) GetGeneratedInstance(ptr interface{}) (*Instance, bool) {
	return s.valStack.GetGeneratedInstance(ptr)
}

func (s *DualStack) SetGeneratedInstance(ptr interface{}, inst *Instance) {
	s.valStack.SetGeneratedInstance(ptr, inst)
}

func (s *DualStack) GetCacheValue(ptr interface{}) (Value, bool) {
	return s.valStack.GetCacheValue(ptr)
}

func (s *DualStack) SetCacheValue(ptr interface{}, val Value) {
	s.valStack.SetCacheValue(ptr, val)
}
