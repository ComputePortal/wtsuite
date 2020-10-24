package prototypes

import (
	"../values"

	"../../context"
)

var NodeJS_fs *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_fsPrototype() bool {
	*NodeJS_fs = BuiltinPrototype{
		"fs", nil,
		map[string]BuiltinFunction{
			"existsSync": NewStatic(String, Boolean),
			"readFileSync": NewStaticFunction(&And{String, &Opt{String}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					var returnProto values.Prototype = Buffer
					if len(args) == 2 {
						// always string if encoding is specified
						if str, ok := args[1].LiteralStringValue(); ok {
							switch str {
							case "buffer":
								returnProto = Buffer
							default:
								returnProto = String
							}
						} else {
							returnProto = String
						}
					}

					return NewInstance(returnProto, ctx), nil
				}),
			"unlinkSync": NewStatic(&Or{String, Buffer}, nil),
		},
		nil,
	}

	return true
}

var _NodeJS_fsOk = generateNodeJS_fsPrototype()
