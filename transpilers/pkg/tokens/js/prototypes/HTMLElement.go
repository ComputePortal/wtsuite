package prototypes

import (
	"../values"

	"../../context"
)

var HTMLElement *BuiltinPrototype = allocBuiltinPrototype()

func NewHTMLElement(ctx context.Context) *values.Instance {
	return NewInstance(HTMLElement, ctx)
}

func generateHTMLElementPrototype() bool {
	*HTMLElement = BuiltinPrototype{
		"HTMLElement", Element,
		map[string]BuiltinFunction{
			"blur":          NewNormal(&None{}, nil),
			"cellIndex":     NewGetter(Int),            // only available for td
			"children":      NewGetter(HTMLCollection), // valid for all nodes, but this would be weird
			"click":         NewNormal(&None{}, nil),
			"focus":         NewNormal(&None{}, nil),
			"rowIndex":      NewGetter(Int),
			"style":         NewGetter(CSSStyleDeclaration),
			"parentElement": NewGetter(HTMLElement), // actually an element, but then we would be able to get the rowIndex and the cellIndex
			"querySelector": NewNormal(String, HTMLElement),
			"offsetWidth":   NewGetter(Int),
			"offsetHeight":  NewGetter(Int),
		},
		NewNoContentGenerator(HTMLElement),
	}

	return true
}

var _HTMLElementOk = generateHTMLElementPrototype()
