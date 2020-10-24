package values

import (
	"../../context"
)

type Stack interface {
	Parent() Stack

	GetReturn(ctx context.Context) (Value, error)
	SetReturn(v Value, ctx context.Context) error

	HasValue(ptr interface{}) bool

	// the stored value can be nil (in case of a var declaration without rhs)
	GetValue(ptr interface{}, ctx context.Context) (Value, error)

	SetValue(ptr interface{}, v Value, allowBranching bool, ctx context.Context) error

	// eg. for click event listeners
	ResolveAwait(t interface{}) (Value, bool, error)

	// also returns return-value from parent stack where recursion is detected
	IsRecursive(fn interface{}) (bool, Value)

	GetGeneratedInstance(ptr interface{}) (*Instance, bool)
	SetGeneratedInstance(ptr interface{}, inst *Instance)

	GetCacheValue(ptr interface{}) (Value, bool)
	SetCacheValue(ptr interface{}, val Value)

	// optional path
	GetViewInterface(args ...string) ViewInterface
}

type StackData struct {
	parent Stack
}

func NewStackData(parent Stack) StackData {
	return StackData{parent}
}

func (s *StackData) Parent() Stack {
	return s.parent
}

func (s *StackData) GetReturn(ctx context.Context) (Value, error) {
	return s.parent.GetReturn(ctx)
}

func (s *StackData) SetReturn(v Value, ctx context.Context) error {
	return s.parent.SetReturn(v, ctx)
}

func (s *StackData) HasValue(ptr interface{}) bool {
	return s.parent.HasValue(ptr)
}

func (s *StackData) GetValue(ptr interface{}, ctx context.Context) (Value, error) {
	return s.parent.GetValue(ptr, ctx)
}

func (s *StackData) SetValue(ptr interface{}, v Value,
	allowBranching bool, ctx context.Context) error {
	return s.parent.SetValue(ptr, v, allowBranching, ctx)
}

func (s *StackData) IsRecursive(fn interface{}) (bool, Value) {
	return s.parent.IsRecursive(fn)
}

func (s *StackData) GetViewInterface(args ...string) ViewInterface {
	return s.parent.GetViewInterface(args...)
}

func (s *StackData) ResolveAwait(t interface{}) (Value, bool, error) {
	return s.parent.ResolveAwait(t)
}

func (s *StackData) GetGeneratedInstance(ptr interface{}) (*Instance, bool) {
	return s.parent.GetGeneratedInstance(ptr)
}

func (s *StackData) SetGeneratedInstance(ptr interface{}, inst *Instance) {
	s.parent.SetGeneratedInstance(ptr, inst)
}

func (s *StackData) GetCacheValue(ptr interface{}) (Value, bool) {
	return s.parent.GetCacheValue(ptr)
}

func (s *StackData) SetCacheValue(ptr interface{}, val Value) {
	s.parent.SetCacheValue(ptr, val)
}
