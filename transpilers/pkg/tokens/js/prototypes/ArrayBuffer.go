package prototypes

var ArrayBuffer *BuiltinPrototype = allocBuiltinPrototype()

func generateArrayBuffer() bool {
	*ArrayBuffer = BuiltinPrototype{
		"ArrayBuffer", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _ArrayBufferOk = generateArrayBuffer()
