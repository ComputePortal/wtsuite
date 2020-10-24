package prototypes

import (
	"../values"

	"../../context"
)

var Node *BuiltinPrototype = allocBuiltinPrototype()

func generateNodePrototype() bool {
	*Node = BuiltinPrototype{
		"Node", EventTarget,
		map[string]BuiltinFunction{
			"appendChild": NewMethodLikeNormalFunction(Node,
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return args[0], nil
				}),
			"firstChild": NewGetter(Node),
			"insertBefore": NewMethodLikeNormalFunction(&And{Node, Node},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return args[0], nil
				}),
			"lastChild":  NewGetter(Node),
			"parentNode": NewGetter(Node),
			"removeChild": NewMethodLikeNormalFunction(Node,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return args[0], nil
				}),
			"replaceChild": NewMethodLikeNormalFunction(&And{Node, Node},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {

					// args[0] is new
					// args[1] is old
					// return old
					return args[1], nil
				}),
			"contains": NewNormal(Node, Boolean),
		},
		nil,
	}

	return true
}

var _NodeOK = generateNodePrototype()
