package prototypes

var CanvasGradient *BuiltinPrototype = allocBuiltinPrototype()

func generateCanvasGradientPrototype() bool {
	*CanvasGradient = BuiltinPrototype{
		"CanvasGradient", nil,
		map[string]BuiltinFunction{
			"addColorStop": NewNormal(&And{Number, String}, nil),
		},
		nil,
	}

	return true
}

var _CanvasGradientOk = generateCanvasGradientPrototype()
