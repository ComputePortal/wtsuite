package prototypes

var CanvasPattern *BuiltinPrototype = allocBuiltinPrototype()

func generateCanvasPatternPrototype() bool {
	*CanvasPattern = BuiltinPrototype{
		"CanvasPattern", nil,
		map[string]BuiltinFunction{},
		nil,
	}

	return true
}

var _CanvasPatternOk = generateCanvasPatternPrototype()
