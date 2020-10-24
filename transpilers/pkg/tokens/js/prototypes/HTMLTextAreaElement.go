package prototypes

import ()

var HTMLTextAreaElement *BuiltinPrototype = allocBuiltinPrototype()

func generateHTMLTextAreaElementPrototype() bool {
	*HTMLTextAreaElement = BuiltinPrototype{
		"HTMLTextAreaElement", HTMLElement,
		map[string]BuiltinFunction{
			"checkValidity":     NewNormal(&None{}, Boolean),
			"setCustomValidity": NewNormal(String, nil),
			"value":             NewGetterSetter(String),
			"selectionStart":    NewGetterSetter(Int),
			"selectionEnd":      NewGetterSetter(Int),
		},
		NewNoContentGenerator(HTMLTextAreaElement),
	}

	return true
}

var _HTMLTextAreaElementOk = generateHTMLTextAreaElementPrototype()
