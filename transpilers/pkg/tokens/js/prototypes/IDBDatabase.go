package prototypes

var IDBDatabase *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBDatabasePrototype() bool {
	*IDBDatabase = BuiltinPrototype{
		"IDBDatabase", nil,
		map[string]BuiltinFunction{
			"createObjectStore": NewMethodLikeNormal(&And{String, &Opt{Object}}, IDBObjectStore),
			"transaction":       NewNormal(&And{Array, &Opt{String}}, IDBTransaction),
		},
		NewNoContentGenerator(IDBDatabase),
	}

	return true
}

var _IDBDatabaseOk = generateIDBDatabasePrototype()
