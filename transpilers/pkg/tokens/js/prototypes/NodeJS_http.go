package prototypes

var NodeJS_http *BuiltinPrototype = allocBuiltinPrototype()

// is actually a builtin nodejs module
func generateNodeJS_httpPrototype() bool {
	*NodeJS_http = BuiltinPrototype{
		"http", nil,
		map[string]BuiltinFunction{
			"createServer":    NewStatic(&None{}, NodeJS_http_Server),
			"IncomingMessage": NewStaticClassGetter(NodeJS_http_IncomingMessage),
			"Server":          NewStaticClassGetter(NodeJS_http_Server),
			"ServerResponse":  NewStaticClassGetter(NodeJS_http_ServerResponse),
		},
		nil,
	}

	return true
}

var _NodeJS_httpOk = generateNodeJS_httpPrototype()
