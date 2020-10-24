package prototypes

var CSSStyleDeclaration *BuiltinPrototype = allocBuiltinPrototype()

func generateCSSStyleDeclarationPrototype() bool {
	*CSSStyleDeclaration = BuiltinPrototype{
		"CSSStyleDeclaration", nil,
		map[string]BuiltinFunction{
			"display":  NewGetterSetter(String),
			"fontSize": NewGetter(String),
			// style is not the same as getBoundingClientRect!
			"height":           NewGetterSetter(String),
			"width":            NewGetterSetter(String),
			"top":              NewGetterSetter(String),
			"bottom":           NewGetterSetter(String),
			"left":             NewGetterSetter(String),
			"right":            NewGetterSetter(String),
			"position":         NewGetterSetter(String),
			"getPropertyValue": NewNormal(String, String),
			"removeProperty":   NewMethodLikeNormal(String, String),
			"setProperty":      NewNormal(&And{String, &And{&Opt{String}, &Opt{String}}}, nil),
		},
		nil,
	}

	return true
}

var _CSSStyleDeclarationOk = generateCSSStyleDeclarationPrototype()
