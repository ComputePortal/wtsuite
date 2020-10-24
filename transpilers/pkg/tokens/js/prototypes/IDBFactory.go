package prototypes

var IDBFactory *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBFactoryPrototype() bool {
	*IDBFactory = BuiltinPrototype{
		"IDBFactory", nil,
		map[string]BuiltinFunction{
			"open": NewNormal(&And{String, &Opt{Int}}, IDBOpenDBRequest),
		},
		nil,
	}

	return true
}

var _IDBFactoryOk = generateIDBFactoryPrototype()
