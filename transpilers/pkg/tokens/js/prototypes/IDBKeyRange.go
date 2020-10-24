package prototypes

var IDBKeyRange *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBKeyRangePrototype() bool {
	*IDBKeyRange = BuiltinPrototype{
		"IDBKeyRange", nil,
		map[string]BuiltinFunction{
			"lowerBound": NewStatic(&And{&Or{Int, String}, &Opt{Boolean}}, IDBKeyRange),
			"upperBound": NewStatic(&And{&Or{Int, String}, &Opt{Boolean}}, IDBKeyRange),
			"bound": NewStatic(&And{&Or{&And{Int, Int}, &And{String, String}}, &And{&Opt{Boolean}, &Opt{Boolean}}},
				IDBKeyRange),
			"only":     NewStatic(&Or{Int, String}, IDBKeyRange),
			"includes": NewNormal(Int, Boolean),
		},
		nil,
	}

	return true
}

var _IDBKeyRangeOk = generateIDBKeyRangePrototype()
