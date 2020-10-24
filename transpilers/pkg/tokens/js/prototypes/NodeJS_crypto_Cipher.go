package prototypes

var NodeJS_crypto_Cipher *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_crypto_CipherPrototype() bool {
	*NodeJS_crypto_Cipher = BuiltinPrototype{
		"crypto.Cipher", NodeJS_EventEmitter,
		map[string]BuiltinFunction{
      "final":  NewNormal(String, String),
      "update": NewNormal(&And{String, &And{String, String}}, String),
		},
		nil,
	}

	return true
}

var _NodeJS_crypto_CipherOk = generateNodeJS_crypto_CipherPrototype()
