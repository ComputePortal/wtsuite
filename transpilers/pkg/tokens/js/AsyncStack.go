package js

import (
	"./values"
)

type AsyncStack struct {
	await interface{}
	val   values.Value
	values.StackData
}

func NewAsyncStack(parent values.Stack, await interface{}, args []values.Value) (*AsyncStack, error) {
	if len(args) > 1 {
		errCtx := args[2].Context()
		return nil, errCtx.NewError("Error: expected 0 or 1 arguments for async await resolution")
	}

	var arg values.Value = nil
	if len(args) == 1 {
		arg = args[0]
	}

	return &AsyncStack{await, arg, values.NewStackData(parent)}, nil
}

func (s *AsyncStack) ResolveAwait(t interface{}) (values.Value, bool, error) {
	if t == s.await {
		return s.val, true, nil
	} else {
		return s.StackData.ResolveAwait(t)
	}
}
