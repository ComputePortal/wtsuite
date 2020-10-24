package prototypes

var Text *BuiltinPrototype = allocBuiltinPrototype()

func generateTextPrototype() bool {
	*Text = BuiltinPrototype{
		"Text", Node,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _TextOk = generateTextPrototype()
