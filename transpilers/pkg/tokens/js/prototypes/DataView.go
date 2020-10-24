package prototypes

var DataView *BuiltinPrototype = allocBuiltinPrototype()

func generateDataView() bool {
	*DataView = BuiltinPrototype{
		"DataView", nil,
		map[string]BuiltinFunction{
			"getInt8":      NewNormal(&And{Int, &Opt{Boolean}}, Int),
			"getUint8":     NewNormal(&And{Int, &Opt{Boolean}}, Int),
			"getInt16":     NewNormal(&And{Int, &Opt{Boolean}}, Int),
			"getUint16":    NewNormal(&And{Int, &Opt{Boolean}}, Int),
			"getInt32":     NewNormal(&And{Int, &Opt{Boolean}}, Int),
			"getUint32":    NewNormal(&And{Int, &Opt{Boolean}}, Int),
			"getBigInt64":  NewNormal(&And{Int, &Opt{Boolean}}, Int),
			"getBigUint64": NewNormal(&And{Int, &Opt{Boolean}}, Int),
			"getFloat32":   NewNormal(&And{Int, &Opt{Boolean}}, Number),
			"getFloat64":   NewNormal(&And{Int, &Opt{Boolean}}, Number),
			"setInt8":      NewNormal(&And{Int, &And{Int, &Opt{Boolean}}}, nil),
			"setUint8":     NewNormal(&And{Int, &And{Int, &Opt{Boolean}}}, nil),
			"setInt16":     NewNormal(&And{Int, &And{Int, &Opt{Boolean}}}, nil),
			"setUint16":    NewNormal(&And{Int, &And{Int, &Opt{Boolean}}}, nil),
			"setInt32":     NewNormal(&And{Int, &And{Int, &Opt{Boolean}}}, nil),
			"setUint32":    NewNormal(&And{Int, &And{Int, &Opt{Boolean}}}, nil),
			"setBigInt64":  NewNormal(&And{Int, &And{Int, &Opt{Boolean}}}, nil),
			"setBigUint64": NewNormal(&And{Int, &And{Int, &Opt{Boolean}}}, nil),
			"setFloat32":   NewNormal(&And{Int, &And{Number, &Opt{Boolean}}}, nil),
			"setFloat64":   NewNormal(&And{Int, &And{Number, &Opt{Boolean}}}, nil),
		},
		NewConstructor(ArrayBuffer, DataView),
	}

	return true
}

var _DataViewOk = generateDataView()
