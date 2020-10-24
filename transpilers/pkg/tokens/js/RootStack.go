package js

import (
	"fmt"

	"./values"

	"../context"
)

// map of variables -> values
// used in both GlobalStack and FunctionStack
type RootStack struct {
	vals map[interface{}]values.Value
}

func newRootStack() RootStack {
	return RootStack{make(map[interface{}]values.Value)}
}

func (s *RootStack) HasValue(ptr interface{}) bool {
	_, ok := s.vals[ptr]
	return ok
}

func (s *RootStack) GetValue(ptr interface{}, ctx context.Context) (values.Value, error) {
	if v, ok := s.vals[ptr]; ok {
		return v, nil
	} else {
		err := ctx.NewError(fmt.Sprintf("Error: variable not found (ref %p in %p)\n", ptr, s.vals))
		panic(err)
		return nil, err
	}
}

func (s *RootStack) SetValue(ptr interface{}, v values.Value,
	allowBranching bool, ctx context.Context) error {

	if _, ok := ptr.(*Package); ok {
		return ctx.NewError("Error: package can't have a value")
	}

	// type checks must be performed in the Expressions and Statements themselves
	// eg. super changes from function to instance, which would give an error if we type checked here, but clearly is ok

	s.vals[ptr] = v

	return nil
}
