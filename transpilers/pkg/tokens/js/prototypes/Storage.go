package prototypes

var Storage *BuiltinPrototype = allocBuiltinPrototype()

func generateStoragePrototype() bool {
	*Storage = BuiltinPrototype{
		"Storage", nil,
		map[string]BuiltinFunction{
			"clear":      NewNormal(&None{}, nil),
			"getItem":    NewNormal(String, String),
			"key":        NewNormal(Int, String),
			"length":     NewGetter(Int),
			"removeItem": NewNormal(String, nil),
			"setItem":    NewNormal(&And{String, String}, nil),
		},
		nil,
	}

	return true
}

var _StorageOk = generateStoragePrototype()
