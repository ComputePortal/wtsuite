package prototypes

var DedicatedWorkerGlobalScope *BuiltinPrototype = allocBuiltinPrototype()

func generateDedicatedWorkerGlobalScopePrototype() bool {
	*DedicatedWorkerGlobalScope = BuiltinPrototype{
		"DedicatedWorkerGlobalScope", WorkerGlobalScope,
		map[string]BuiltinFunction{
			"postMessage": NewStatic(&Any{}, nil),
		},
		nil,
	}

	return true
}

var _DedicatedWorkerGlobalScopeOk = generateDedicatedWorkerGlobalScopePrototype()
