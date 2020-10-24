package prototypes

var WebGLShader *BuiltinPrototype = allocBuiltinPrototype()

func generateWebGLShaderPrototype() bool {
	*WebGLShader = BuiltinPrototype{
		"WebGLShader", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _WebGLShaderOk = generateWebGLShaderPrototype()
