package prototypes

var NodeJS_stream *BuiltinPrototype = allocBuiltinPrototype()

// is actually a builtin nodejs module
func generateNodeJS_streamPrototype() bool {
	*NodeJS_stream = BuiltinPrototype{
		"stream", nil,
		map[string]BuiltinFunction{
			"Readable": NewStaticClassGetter(NodeJS_stream_Readable),
		},
		nil,
	}

	return true
}

var _NodeJS_streamOk = generateNodeJS_streamPrototype()
