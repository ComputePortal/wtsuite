package prototypes

var SharedWorkerGlobalScope *BuiltinPrototype = allocBuiltinPrototype()

func generateSharedWorkerGlobalScopePrototype() bool {
	*SharedWorkerGlobalScope = BuiltinPrototype{
		"SharedWorkerGlobalScope", WorkerGlobalScope,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _SharedWorkerGlobalScopeOk = generateSharedWorkerGlobalScopePrototype()
