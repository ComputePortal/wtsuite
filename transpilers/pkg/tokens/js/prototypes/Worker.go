package prototypes

var Worker *BuiltinPrototype = allocBuiltinPrototype()

func generateWorkerPrototype() bool {
	*Worker = BuiltinPrototype{
		"Worker", nil,
		map[string]BuiltinFunction{
			"onmessage":      NewSetterFunction(&Function{}, messagePortCallback),
			"onmessageerror": NewSetterFunction(&Function{}, messagePortCallback),
			"postMessage":    NewNormal(&Any{}, nil),
		},
		NewConstructor(&And{String, &Opt{String}}, Worker),
	}

	return true
}

var _WorkerOk = generateWorkerPrototype()
