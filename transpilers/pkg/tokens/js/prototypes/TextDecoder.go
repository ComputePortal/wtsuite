package prototypes

var TextDecoder *BuiltinPrototype = allocBuiltinPrototype()

func generateTextDecoder() bool {
	*TextDecoder = BuiltinPrototype{
		"TextDecoder", nil,
		map[string]BuiltinFunction{
			"decode": NewNormal(Uint8Array, String),
		},
		NewConstructor(&And{&Opt{String}, &Opt{Boolean}}, TextDecoder),
	}

	return true
}

var _TextDecoderOk = generateTextDecoder()
