package prototypes

var URL *BuiltinPrototype = allocBuiltinPrototype()

func generateURLPrototype() bool {
	*URL = BuiltinPrototype{
		"URL", nil,
		map[string]BuiltinFunction{
			"searchParams":    NewGetter(URLSearchParams),
			"createObjectURL": NewStatic(Blob, String),
			"revokeObjectURL": NewStatic(String, nil),
		},
		NewConstructor(&And{String, &Opt{String}}, URL),
	}

	return true
}

var _URLOk = generateURLPrototype()
