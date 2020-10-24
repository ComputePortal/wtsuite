package prototypes

import (
	"../values"

	"../../context"
)

var NodeJS_http_Server *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_http_ServerPrototype() bool {
	*NodeJS_http_Server = BuiltinPrototype{
		"http.Server", NodeJS_EventEmitter,
		map[string]BuiltinFunction{
			"addListener": NewNormalFunction(&And{String, &Function{}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					var callbackArgs []values.Value = nil

					if str, ok := args[0].LiteralStringValue(); ok {
						switch str {
						case "request":
							callbackArgs = []values.Value{NewInstance(NodeJS_http_IncomingMessage, ctx),
								NewInstance(NodeJS_http_ServerResponse, ctx)}
						}
					}

					if callbackArgs != nil {
						callbackCtx := args[1].Context()
						if err := args[1].EvalMethod(stack.Parent(), callbackArgs, callbackCtx); err != nil {
							return nil, err
						}

						return nil, nil
					} else {
						eventEmitterFn, err := NodeJS_EventEmitter.GetMember(stack, this, "addListener", false,
							ctx)
						if err != nil {
							panic(err)
						}

						if err := eventEmitterFn.EvalMethod(stack.Parent(), args, ctx); err != nil {
							return nil, err
						}

						return nil, nil
					}
				}),
			"listen": NewMethodLikeNormal(&Or{String, &And{&Opt{Int}, &Opt{String}}}, nil),
		},
		nil,
	}

	return true
}

var _NodeJS_http_ServerOk = generateNodeJS_http_ServerPrototype()
