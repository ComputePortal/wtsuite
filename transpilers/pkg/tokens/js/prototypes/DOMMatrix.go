package prototypes

var DOMMatrix *BuiltinPrototype = allocBuiltinPrototype()

func generateDOMMatrixPrototype() bool {
	*DOMMatrix = BuiltinPrototype{
		"DOMMatrix", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _DOMMatrixOk = generateDOMMatrixPrototype()
