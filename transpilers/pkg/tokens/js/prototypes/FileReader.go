package prototypes

import (
	"../values"

	"../../context"
)

var FileReader *BuiltinPrototype = allocBuiltinPrototype()

func fileReaderCallback(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	arg := args[0]

	event := NewEvent(this, ctx)
	if err := arg.EvalMethod(stack.Parent(), []values.Value{event}, ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func generateFileReader() bool {
	*FileReader = BuiltinPrototype{
		"FileReader", nil,
		map[string]BuiltinFunction{
			"onload":  NewSetterFunction(&Function{}, fileReaderCallback),
			"onerror": NewSetterFunction(&Function{}, fileReaderCallback),
			"result": NewGetterFunction(func(stack values.Stack, this *values.Instance,
				args []values.Value, ctx context.Context) (values.Value, error) {
				return NewInstance(ArrayBuffer, ctx), nil
			}),
			"readAsArrayBuffer": NewNormal(Blob, nil),
		},
		NewConstructor(&None{}, FileReader),
	}

	return true
}

var _FileReaderOk = generateFileReader()
