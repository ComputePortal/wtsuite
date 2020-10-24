package prototypes

var DOMRect *BuiltinPrototype = allocBuiltinPrototype()

func generateDOMRectPrototype() bool {
	*DOMRect = BuiltinPrototype{
		"DOMRect", nil,
		map[string]BuiltinFunction{
			"bottom": NewGetter(Number),
			"height": NewGetter(Number),
			"left":   NewGetter(Number),
			"right":  NewGetter(Number),
			"top":    NewGetter(Number),
			"width":  NewGetter(Number),
			"x":      NewGetter(Number),
			"y":      NewGetter(Number),
		},
		nil,
	}

	return true
}

var _DOMRectOk = generateDOMRectPrototype()
