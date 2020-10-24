package prototypes

var NodeJS_http_ServerResponse *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_http_ServerResponsePrototype() bool {
	*NodeJS_http_ServerResponse = BuiltinPrototype{
		"http.ServerResponse", NodeJS_EventEmitter,
		map[string]BuiltinFunction{
			"end": NewMethodLikeNormal(&And{&Opt{&Or{String, Buffer}},
				&Opt{String}}, NodeJS_http_ServerResponse),
			"statusCode":    NewSetter(Int),
			"statusMessage": NewSetter(String),
			"write":         NewMethodLikeNormal(&And{&Or{String, Buffer}, &And{&Opt{String}, &Opt{&Function{}}}}, Boolean),
			"writeHead": NewMethodLikeNormal(&Or{&And{Int, &And{String, Object}},
				&And{Int, Object}},
				NodeJS_http_ServerResponse),
		},
		nil,
	}

	return true
}

var _NodeJS_http_ServerResponseOk = generateNodeJS_http_ServerResponsePrototype()
