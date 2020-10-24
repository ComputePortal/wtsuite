package prototypes

import (
	"../values"

	"../../context"
)

var MessageEvent *BuiltinPrototype = allocBuiltinPrototype()

func generateMessageEventPrototype() bool {
	*MessageEvent = BuiltinPrototype{
		"MessageEvent", Event,
		map[string]BuiltinFunction{
			"data": NewGetter(&values.AllPrototype{}),
			"ports": NewGetterFunction(
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					content := NewInstance(MessagePort, ctx)
					return NewArray([]values.Value{content}, ctx), nil
				}),
		},
		NewConstructorGeneratorFunction(nil,
			func(stack values.Stack, keys []string, args []values.Value,
				ctx context.Context) (values.Value, error) {
				if keys != nil || args != nil {
					return nil, ctx.NewError("Error: unexpected content types")
				}

				target := NewInstance(MessagePort, ctx)
				return NewAltEvent(MessageEvent, target, ctx), nil
			}),
	}

	return true
}

var _MessageEventOk = generateMessageEventPrototype()
