package prototypes

var WebGLTexture *BuiltinPrototype = allocBuiltinPrototype()

func generateWebGLTexturePrototype() bool {
	*WebGLTexture = BuiltinPrototype{
		"WebGLTexture", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _WebGLTextureOk = generateWebGLTexturePrototype()
