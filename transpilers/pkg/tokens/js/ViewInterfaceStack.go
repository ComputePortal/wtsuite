package js

import (
	"./values"
)

type ViewInterfaceStack struct {
	viewInterf    *ViewInterface
	allViewInterf map[string]*ViewInterface
	values.StackData
}

func NewViewInterfaceStack(vif *ViewInterface, all map[string]*ViewInterface, parent values.Stack) *ViewInterfaceStack {
	return &ViewInterfaceStack{
		vif,
		all,
		values.NewStackData(parent),
	}
}

func (s *ViewInterfaceStack) GetViewInterface(args ...string) values.ViewInterface {

	if len(args) == 0 {
		return s.viewInterf
	} else if len(args) == 1 {
		arg := args[0]

		if vif, ok := s.allViewInterf[arg]; ok {
			return vif
		} else {
			return nil
		}
	} else {
		panic("bad number of arguments")
	}
}
