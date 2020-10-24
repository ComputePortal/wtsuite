package prototypes

import (
	"../values"

	"../../context"
)

var NodeJS_stream_Readable *BuiltinPrototype = allocBuiltinPrototype()

func generateNodeJS_stream_ReadablePrototype() bool {
	*NodeJS_stream_Readable = BuiltinPrototype{
		"stream.Readable", NodeJS_EventEmitter,
		map[string]BuiltinFunction{
      "addListener": NewNormalFunction(&And{String, &Function{}}, 
        func(stack values.Stack, this *values.Instance,
        args []values.Value, ctx context.Context) (values.Value, error) {
          if str, ok := args[0].LiteralStringValue(); ok {
            switch str {
            case "data":
              // assume default encoding in properties
              // TODO: detect encoding and use Buffer if necessary
              if err := args[1].EvalMethod(stack.Parent(), []values.Value{NewString(ctx)}, ctx); err != nil {
                return nil, err
              }
              return nil, nil
            }
          }

          // default doesnt have arguments
					if err := args[1].EvalMethod(stack.Parent(), []values.Value{}, ctx); err != nil {
						return nil, err
					}

					return nil, nil
        }),
      "read": NewNormal(&Opt{Int}, Buffer),
		},
		nil,
	}

	return true
}

var _NodeJS_stream_ReadableOk = generateNodeJS_stream_ReadablePrototype()
