package prototypes

import ()

var HTMLLinkElement *BuiltinPrototype = allocBuiltinPrototype()

func generateHTMLLinkElementPrototype() bool {
	*HTMLLinkElement = BuiltinPrototype{
		"HTMLLinkElement", HTMLElement,
		map[string]BuiltinFunction{
			"download": NewGetterSetter(String),
			"href":     NewGetterSetter(String),
			"rel":      NewGetterSetter(String),
		},
		NewNoContentGenerator(HTMLLinkElement),
	}

	return true
}

var _HTMLLinkElementOk = generateHTMLLinkElementPrototype()
