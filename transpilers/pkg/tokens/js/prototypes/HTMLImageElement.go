package prototypes

import (
	"../values"

	"../../context"
)

var HTMLImageElement *BuiltinPrototype = allocBuiltinPrototype()

func imageOnloadCallback(stack values.Stack, this *values.Instance, args []values.Value,
	ctx context.Context) (values.Value, error) {
	arg := args[0]

	event := NewEvent(this, ctx)
	if err := arg.EvalMethod(stack.Parent(), []values.Value{event}, ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func generateHTMLImageElementPrototype() bool {
	*HTMLImageElement = BuiltinPrototype{
		"HTMLImageElement", HTMLElement,
		map[string]BuiltinFunction{
			"height": NewGetterSetter(Int),
			"width":  NewGetterSetter(Int),
			"onload": NewSetterFunction(&Function{}, imageOnloadCallback),
			"src":    NewSetter(String),
		},
		NewNoContentGenerator(HTMLImageElement),
	}

	return true
}

var _HTMLImageElementOk = generateHTMLImageElementPrototype()
