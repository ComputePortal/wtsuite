package prototypes

var WebGLExtension *BuiltinPrototype = allocBuiltinPrototype()

func generateWebGLExtension() bool {
	*WebGLExtension = BuiltinPrototype{
		"WebGLExtension", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _WebGLExtensionOk = generateWebGLExtension()
