package prototypes

import ()

var HTMLSelectElement *BuiltinPrototype = allocBuiltinPrototype()

func generateHTMLSelectElementPrototype() bool {
	*HTMLSelectElement = BuiltinPrototype{
		"HTMLSelectElement", HTMLElement,
		map[string]BuiltinFunction{
			"checkValidity":     NewNormal(&None{}, Boolean),
			"selectedIndex":     NewGetterSetter(Int),
			"setCustomValidity": NewNormal(String, nil),
			"value":             NewGetterSetter(String),
		},
		NewNoContentGenerator(HTMLSelectElement),
	}

	return true
}

var _HTMLSelectElementOk = generateHTMLSelectElementPrototype()
