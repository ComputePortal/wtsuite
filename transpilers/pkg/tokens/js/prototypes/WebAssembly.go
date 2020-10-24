package prototypes

import ()

var WebAssembly *BuiltinPrototype = allocBuiltinPrototype()

func generateWebAssemblyPrototype() bool {
	*WebAssembly = BuiltinPrototype{
		"WebAssembly", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _WebAssemblyOk = generateWebAssemblyPrototype()
