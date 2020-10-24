package prototypes

import (
	"../values"

	"../../context"
)

var Buffer *BuiltinPrototype = allocBuiltinPrototype()

func generateBufferPrototype() bool {
	*Buffer = BuiltinPrototype{
		"Buffer", nil,
		map[string]BuiltinFunction{
      "concat": NewStaticFunction(Array, 
      func(stack values.Stack, this *values.Instance,
      args []values.Value, ctx context.Context) (values.Value, error) {
        arg0_ := values.UnpackContextValue(args[0])
        arg0, ok := arg0_.(*values.Instance)
        if !ok {
          errCtx := arg0_.Context()
          return nil, errCtx.NewError("Error: not an Array instance")
        }

        props := values.AssertArrayProperties(arg0.Properties())
        item := props.GetItem()

        if !item.IsInstanceOf(Buffer) {
          return nil, ctx.NewError("Error: argument is not an Array of Buffer's (" + item.TypeName() + ")")
        }

        return NewInstance(Buffer, ctx), nil
      }),
      "from": NewNormal(&Or{Buffer, &Or{Array, &And{String, &Opt{String}}}}, Buffer),
      "toString": NewNormal(&And{&Opt{String}, &And{&Opt{Int}, &Opt{Int}}}, String),
    },
		nil,
	}

	return true
}

var _BufferOk = generateBufferPrototype()
