package prototypes

// not yet supported by most browsers
var IDBOpenDBRequest *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBOpenDBRequestPrototype() bool {
	*IDBOpenDBRequest = BuiltinPrototype{
		"IDBOpenDBRequest", IDBRequest,
		map[string]BuiltinFunction{
			"result":          NewGetter(IDBDatabase),
			"onupgradeneeded": NewSetterFunction(&Function{}, generateIDBRequestCallback(IDBVersionChangeEvent)),
		},
		nil,
	}

	return true
}

var _IDBOpenDBRequestOk = generateIDBOpenDBRequestPrototype()
