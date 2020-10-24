package prototypes

var Path *BuiltinPrototype = allocBuiltinPrototype()

func generatePathPrototype() bool {
	*Path = BuiltinPrototype{
		"Path", nil,
		map[string]BuiltinFunction{
			"join": NewStatic(&And{String, &And{String, &Rest{String}}}, String),
		},
		nil,
	}

	return true
}

var _PathOk = generatePathPrototype()
