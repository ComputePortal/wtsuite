package prototypes

var WorkerGlobalScope *BuiltinPrototype = allocBuiltinPrototype()

func generateWorkerGlobalScopePrototype() bool {
	*WorkerGlobalScope = BuiltinPrototype{
		"WorkerGlobalScope", EventTarget,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _WorkerGlobalScopeOk = generateWorkerGlobalScopePrototype()
