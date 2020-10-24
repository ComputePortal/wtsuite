package prototypes

// not yet supported by most browsers
var Path2D *BuiltinPrototype = allocBuiltinPrototype()

func generatePath2DPrototype() bool {
	*Path2D = BuiltinPrototype{
		"Path2D", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _Path2DOk = generatePath2DPrototype()
