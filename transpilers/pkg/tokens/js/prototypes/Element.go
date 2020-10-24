package prototypes

import (
//"../values"

//"../../context"
)

var Element *BuiltinPrototype = allocBuiltinPrototype()

func generateElementPrototype() bool {
	*Element = BuiltinPrototype{
		"Element", Node,
		map[string]BuiltinFunction{
			"className":             NewGetterSetter(String),
			"getAttribute":          NewNormal(String, String),
			"getBoundingClientRect": NewNormal(&None{}, DOMRect),
			"hasAttribute":          NewNormal(String, Boolean),
			"id":                    NewGetterSetter(String),
			"innerHTML":             NewGetterSetter(String),
			"removeAttribute":       NewNormal(String, nil),
			"scrollLeft":            NewGetterSetter(Number),
			"scrollWidth":           NewGetter(Number),
			"scrollTo":              NewNormal(&And{Number, Number}, nil),
			"scrollIntoView":        NewNormal(&Or{&None{}, &Or{Boolean, Object}}, nil),
			"scrollTop":             NewGetterSetter(Number),
			"scrollHeight":          NewGetter(Number),
			"clientLeft":            NewGetter(Number),
			"clientWidth":           NewGetter(Number),
			"clientTop":             NewGetter(Number),
			"clientHeight":          NewGetter(Number),
			"setAttribute":          NewNormal(&And{String, String}, nil),
			"tagName":               NewGetter(String),
		},
		nil,
	}

	return true
}

var _ElementOk = generateElementPrototype()
