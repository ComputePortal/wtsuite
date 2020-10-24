package prototypes

var IDBTransaction *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBTransactionPrototype() bool {
	*IDBTransaction = BuiltinPrototype{
		"IDBTransaction", EventTarget,
		map[string]BuiltinFunction{
			"commit":      NewNormal(&None{}, nil),
			"objectStore": NewNormal(String, IDBObjectStore),
			"onerror":     NewSetterFunction(&Function{}, idbRequestCallback),
			"oncomplete":  NewSetterFunction(&Function{}, idbRequestCallback),
		},
		NewNoContentGenerator(IDBTransaction),
	}

	return true
}

var _IDBTransactionOk = generateIDBTransactionPrototype()
