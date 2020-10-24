package prototypes

var GLEnum *BuiltinPrototype = allocBuiltinPrototype()

func generateGLEnumPrototype() bool {
	*GLEnum = BuiltinPrototype{
		"GLEnum", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _GLEnumOk = generateGLEnumPrototype()
