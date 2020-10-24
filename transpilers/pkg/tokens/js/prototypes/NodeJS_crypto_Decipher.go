package prototypes

var NodeJS_crypto_Decipher *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_crypto_DecipherPrototype() bool {
	*NodeJS_crypto_Decipher = BuiltinPrototype{
		"crypto.Decipher", NodeJS_EventEmitter,
		map[string]BuiltinFunction{
      "final":  NewNormal(String, String),
      "update": NewNormal(&And{String, &And{String, String}}, String),
		},
		nil,
	}

	return true
}

var _NodeJS_crypto_DecipherOk = generateNodeJS_crypto_DecipherPrototype()
