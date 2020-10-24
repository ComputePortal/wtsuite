package prototypes

var WheelEvent *EventPrototype = allocEventPrototype()

func generateWheelEventPrototype() bool {
	*WheelEvent = EventPrototype{BuiltinPrototype{
		"WheelEvent", MouseEvent,
		map[string]BuiltinFunction{
			"deltaX": NewGetter(Number),
			"deltaY": NewGetter(Number),
			"deltaZ": NewGetter(Number),
		},
		nil,
	}}

	return true
}

var _WheelEventOk = generateWheelEventPrototype()
