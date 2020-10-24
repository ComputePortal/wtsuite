package prototypes

import (
	"../values"

	"../../context"
)

var MessagePort *BuiltinPrototype = allocBuiltinPrototype()

func messagePortCallback(stack values.Stack, this *values.Instance, args []values.Value,
	ctx context.Context) (values.Value, error) {
	arg := args[0]

	event := NewAltEvent(MessageEvent, this, ctx)
	if err := arg.EvalMethod(stack.Parent(), []values.Value{event}, ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func generateMessagePortPrototype() bool {
	*MessagePort = BuiltinPrototype{
		"MessagePort", EventTarget,
		map[string]BuiltinFunction{
			"close":          NewNormal(&None{}, nil),
			"onmessage":      NewSetterFunction(&Function{}, messagePortCallback),
			"onmessageerror": NewSetterFunction(&Function{}, messagePortCallback),
			"postMessage":    NewNormal(&Any{}, nil),
			"start":          NewNormal(&None{}, nil),
		},
		nil,
	}

	return true
}

var _MessagePortOk = generateMessagePortPrototype()
