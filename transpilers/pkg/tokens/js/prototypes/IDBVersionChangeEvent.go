package prototypes

var IDBVersionChangeEvent *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBVersionChangeEventPrototype() bool {
	*IDBVersionChangeEvent = BuiltinPrototype{
		"IDBVersionChangeEvent", Event,
		map[string]BuiltinFunction{
			"oldVersion": NewGetter(Int),
			"newVersion": NewGetter(Int),
		},
		nil,
	}

	return true
}

var _IDBVersionChangeEventOk = generateIDBVersionChangeEventPrototype()
