package prototypes

var SharedWorker *BuiltinPrototype = allocBuiltinPrototype()

func generateSharedWorkerPrototype() bool {
	*SharedWorker = BuiltinPrototype{
		"SharedWorker", nil,
		map[string]BuiltinFunction{
			"port": NewGetter(MessagePort),
		},
		NewConstructor(&And{String, &Opt{String}}, SharedWorker),
	}

	return true
}

var _SharedWorkerOk = generateSharedWorkerPrototype()
