package prototypes

var HashChangeEvent *EventPrototype = allocEventPrototype()

func generateHashChangeEventPrototype() bool {
	*HashChangeEvent = EventPrototype{BuiltinPrototype{
		"HashChangeEvent", Event,
		map[string]BuiltinFunction{
			"newURL": NewGetter(String),
			"oldURL": NewGetter(String),
		},
		nil,
	}}

	return true
}

var _HashChangeEventOk = generateHashChangeEventPrototype()
