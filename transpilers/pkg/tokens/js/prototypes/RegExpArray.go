package prototypes

var RegExpArray *BuiltinPrototype = allocBuiltinPrototype()

func generateRegExpArray() bool {
	*RegExpArray = BuiltinPrototype{
		"RegExpArray", Array,
		map[string]BuiltinFunction{
			"index": NewGetter(Int),
			"input": NewGetter(String),
		},
		nil,
	}

	return true
}

var _RegExpArrayOk = generateRegExpArray()
