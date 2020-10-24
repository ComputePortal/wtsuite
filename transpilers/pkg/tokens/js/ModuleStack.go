package js

import (
	"./values"

	"../context"
)

type ModuleStack struct {
	module *ControlModule
	values.StackData
}

func (s *ModuleStack) GetReturn(ctx context.Context) (values.Value, error) {
	return nil, ctx.NewError("Error: not in a function")
}

func (s *ModuleStack) SetReturn(v values.Value, ctx context.Context) error {
	return ctx.NewError("Error: not in a function")
}

func (s *ModuleStack) IsRecursive(fn interface{}) (bool, values.Value) {
	return false, nil
}
