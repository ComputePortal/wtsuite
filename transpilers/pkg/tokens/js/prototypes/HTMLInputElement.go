package prototypes

import ()

var HTMLInputElement *BuiltinPrototype = allocBuiltinPrototype()

func generateHTMLInputElementPrototype() bool {
	*HTMLInputElement = BuiltinPrototype{
		"HTMLInputElement", HTMLElement,
		map[string]BuiltinFunction{
			"checked":           NewGetterSetter(Boolean),
			"checkValidity":     NewNormal(&None{}, Boolean),
			"setCustomValidity": NewNormal(String, nil),
			"value":             NewGetterSetter(String),
			"selectionStart":    NewGetterSetter(Int),
			"selectionEnd":      NewGetterSetter(Int),
		},
		NewNoContentGenerator(HTMLInputElement),
	}

	return true
}

var _HTMLInputElementOk = generateHTMLInputElementPrototype()
