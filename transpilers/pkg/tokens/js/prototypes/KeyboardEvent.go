package prototypes

import (
	"../values"

	"../../context"
)

var KeyboardEvent *EventPrototype = allocEventPrototype()

func generateKeyboardEventPrototype() bool {
	*KeyboardEvent = EventPrototype{BuiltinPrototype{
		"KeyboardEvent", Event,
		map[string]BuiltinFunction{
			"altKey":   NewGetter(Boolean),
			"ctrlKey":  NewGetter(Boolean),
			"key":      NewGetter(String),
			"metaKey":  NewGetter(Boolean),
			"shiftKey": NewGetter(Boolean),
		},
		NewConstructorGenerator(&And{String, &Opt{Object}}, KeyboardEvent,
			func(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
				if keys != nil || args != nil {
					return nil, ctx.NewError("Error: unexpected content types")
				}

				return NewInstance(KeyboardEvent, ctx), nil
			}),
	}}

	return true
}

var _KeyboardEventOk = generateKeyboardEventPrototype()
