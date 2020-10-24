package prototypes

var Console *BuiltinPrototype = allocBuiltinPrototype()

func generateConsolePrototype() bool {
	*Console = BuiltinPrototype{
		"Console", nil,
		map[string]BuiltinFunction{
			"log": NewNormal(&Rest{&Any{}}, nil),
		},
		nil,
	}

	return true
}

var _ConsoleOk = generateConsolePrototype()
