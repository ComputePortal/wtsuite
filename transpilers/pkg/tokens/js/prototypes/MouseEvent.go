package prototypes

var MouseEvent *EventPrototype = allocEventPrototype()

func generateMouseEventPrototype() bool {
	*MouseEvent = EventPrototype{BuiltinPrototype{
		"MouseEvent", Event,
		map[string]BuiltinFunction{
			"altKey":   NewGetter(Boolean),
			"clientX":  NewGetter(Number),
			"clientY":  NewGetter(Number),
			"ctrlKey":  NewGetter(Boolean),
			"metaKey":  NewGetter(Boolean),
			"shiftKey": NewGetter(Boolean),
		},
		nil,
	}}

	return true
}

var _MouseEventOk = generateMouseEventPrototype()
