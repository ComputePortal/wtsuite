package prototypes

var ImageData *BuiltinPrototype = allocBuiltinPrototype()

func generateImageDataPrototype() bool {
	*ImageData = BuiltinPrototype{
		"ImageData", nil,
		map[string]BuiltinFunction{
			"data":   NewGetter(Uint8ClampedArray),
			"height": NewGetter(Int),
			"width":  NewGetter(Int),
		},
		nil,
	}

	return true
}

var _ImageDataOk = generateImageDataPrototype()
