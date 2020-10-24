package prototypes

var NodeJS_crypto *BuiltinPrototype = allocBuiltinPrototype()

// is actually a builtin nodejs module
func generateNodeJS_cryptoPrototype() bool {
	*NodeJS_crypto = BuiltinPrototype{
		"crypto", nil,
		map[string]BuiltinFunction{
      "createCipheriv":   NewStatic(&And{String, &And{Buffer, Buffer}}, 
        NodeJS_crypto_Cipher),
      "createDecipheriv": NewStatic(&And{String, &And{Buffer, Buffer}}, 
        NodeJS_crypto_Decipher),
			"randomBytes":      NewStatic(Int, Buffer),
			"Cipher":           NewStaticClassGetter(NodeJS_crypto_Cipher),
			"Decipher":         NewStaticClassGetter(NodeJS_crypto_Decipher),
		},
		nil,
	}

	return true
}

var _NodeJS_cryptoOk = generateNodeJS_cryptoPrototype()
