package js

import (
	"./values"

	"../context"
)

type FunctionStack struct {
	function    *Function
	returnValue values.Value
	values.StackData
	RootStack
}

func NewFunctionStack(fn *Function, parent values.Stack) *FunctionStack {
	return &FunctionStack{
		fn,
		nil,
		values.NewStackData(parent),
		newRootStack(),
	}
}

func (s *FunctionStack) GetReturn(ctx context.Context) (values.Value, error) {
	return s.returnValue, nil
}

func (s *FunctionStack) SetReturn(v values.Value, ctx context.Context) error {
	// type checking is done in return statement itself
	if s.returnValue != nil {
		v = values.NewMulti([]values.Value{v, s.returnValue}, ctx)
	}

	s.returnValue = v
	return nil
}

func (s *FunctionStack) HasValue(ptr interface{}) bool {
	return s.RootStack.HasValue(ptr) || s.StackData.HasValue(ptr)
}

func (s *FunctionStack) GetValue(ptr interface{}, ctx context.Context) (values.Value, error) {
	if s.StackData.HasValue(ptr) {
		return s.StackData.GetValue(ptr, ctx)
	} else {
		return s.RootStack.GetValue(ptr, ctx)
	}
}

func (s *FunctionStack) SetValue(ptr interface{}, v values.Value,
	allowBranching bool, ctx context.Context) error {
	if s.StackData.HasValue(ptr) {
		return s.StackData.SetValue(ptr, v, allowBranching, ctx)
	} else {
		return s.RootStack.SetValue(ptr, v, false, ctx)
	}
}

// also returns the current return value
func (s *FunctionStack) IsRecursive(fn_ interface{}) (bool, values.Value) {
	if fn, ok := fn_.(*Function); ok {
		if fn == s.function {
			return true, s.returnValue
		} else {
			return s.StackData.IsRecursive(fn)
		}
	} else {
		panic("unexpected")
	}
}

func (s *FunctionStack) ResolveAwait(t interface{}) (values.Value, bool, error) {
	return s.Parent().ResolveAwait(t)
}

func (s *FunctionStack) GetGeneratedInstance(ptr interface{}) (*values.Instance, bool) {
	return s.Parent().GetGeneratedInstance(ptr)
}

func (s *FunctionStack) SetGeneratedInstance(ptr interface{}, inst *values.Instance) {
	s.Parent().SetGeneratedInstance(ptr, inst)
}

func (s *FunctionStack) GetCacheValue(ptr interface{}) (values.Value, bool) {
	return s.Parent().GetCacheValue(ptr)
}

func (s *FunctionStack) SetCacheValue(ptr interface{}, val values.Value) {
	s.Parent().SetCacheValue(ptr, val)
}
