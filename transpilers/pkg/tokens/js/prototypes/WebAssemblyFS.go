package prototypes

var WebAssemblyFS *BuiltinInterface = allocBuiltinInterface()

func generateWebAssemblyFSInterface() bool {
	*WebAssemblyFS = BuiltinInterface{
		"WebAssemblyFS", map[string]*BuiltinInterfaceMember{
			"exists": &BuiltinInterfaceMember{String, NORMAL, Boolean},
			"open":   &BuiltinInterfaceMember{String, NORMAL, Int},
			"create": &BuiltinInterfaceMember{String, NORMAL, Int},
			"close":  &BuiltinInterfaceMember{Int, NORMAL, nil},
			"read":   &BuiltinInterfaceMember{&And{Int, Int}, NORMAL, Uint8Array},
			"write":  &BuiltinInterfaceMember{&And{Int, Uint8Array}, NORMAL, nil},
			"seek":   &BuiltinInterfaceMember{&And{Int, Int}, NORMAL, nil},
			"tell":   &BuiltinInterfaceMember{Int, NORMAL, Int},
			"size":   &BuiltinInterfaceMember{Int, NORMAL, Int},
		},
		nil,
	}

	return true
}

var _WebAssemblyFSOk = generateWebAssemblyFSInterface()
