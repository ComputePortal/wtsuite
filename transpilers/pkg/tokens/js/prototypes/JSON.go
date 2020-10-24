package prototypes

var JSON *BuiltinPrototype = allocBuiltinPrototype()

func generateJSONPrototype() bool {
	*JSON = BuiltinPrototype{
		"JSON", nil,
		map[string]BuiltinFunction{
			"stringify": NewStatic(Object, String),
			"parse":     NewStatic(String, Object),
		},
		nil,
	}

	return true
}

var _JSONOk = generateJSONPrototype()
