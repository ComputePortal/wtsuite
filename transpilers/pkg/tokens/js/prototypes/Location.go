package prototypes

var Location *BuiltinPrototype = allocBuiltinPrototype()

func generateLocationPrototype() bool {
	*Location = BuiltinPrototype{
		"Location", nil,
		map[string]BuiltinFunction{
			"hash":     NewGetterSetter(String),
			"href":     NewGetterSetter(String),
			"origin":   NewGetter(String),
			"pathname": NewGetter(String),
			"search":   NewGetter(String),
		},
		nil,
	}

	return true
}

var _LocationOk = generateLocationPrototype()
