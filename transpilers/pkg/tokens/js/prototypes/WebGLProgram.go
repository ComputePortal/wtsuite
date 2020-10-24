package prototypes

var WebGLProgram *BuiltinPrototype = allocBuiltinPrototype()

func generateWebGLProgramPrototype() bool {
	*WebGLProgram = BuiltinPrototype{
		"WebGLProgram", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _WebGLProgramOk = generateWebGLProgramPrototype()
