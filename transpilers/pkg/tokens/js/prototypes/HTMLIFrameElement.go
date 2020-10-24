package prototypes

import ()

var HTMLIFrameElement *BuiltinPrototype = allocBuiltinPrototype()

func generateHTMLIFrameElementPrototype() bool {
	*HTMLIFrameElement = BuiltinPrototype{
		"HTMLIFrameElement", HTMLElement,
		map[string]BuiltinFunction{
			"contentDocument": NewGetter(Document),
			"contentWindow":   NewGetter(Window),
			"src":             NewGetterSetter(String),
		},
		NewNoContentGenerator(HTMLIFrameElement),
	}

	return true
}

var _HTMLIFrameElementOk = generateHTMLIFrameElementPrototype()
