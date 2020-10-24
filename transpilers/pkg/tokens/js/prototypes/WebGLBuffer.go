package prototypes

var WebGLBuffer *BuiltinPrototype = allocBuiltinPrototype()

func generateWebGLBufferPrototype() bool {
	*WebGLBuffer = BuiltinPrototype{
		"WebGLBuffer", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _WebGLBufferOk = generateWebGLBufferPrototype()
