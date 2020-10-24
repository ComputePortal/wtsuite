package prototypes

import (
	"../values"

	"../../context"
)

var Blob *BuiltinPrototype = allocBuiltinPrototype()

func generateBlob() bool {
	*Blob = BuiltinPrototype{
		"Blob", nil,
		map[string]BuiltinFunction{
			"arrayBuffer": NewNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewResolvedPromise(NewInstance(ArrayBuffer, ctx), ctx)
				}),
			"size":  NewGetter(Int),
			"slice": NewNormal(&And{&Opt{Int}, &And{&Opt{Int}, &Opt{String}}}, Blob),
			"type":  NewGetter(String),
		},
		NewConstructorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if err := CheckInputs(&And{&Or{Array, String}, &Opt{Object}}, args, ctx); err != nil {
				return nil, err
			}

			item, err := args[0].GetIndex(stack, NewInt(ctx), ctx)
			if err != nil {
				return nil, err
			}

			if !item.IsInstanceOf(ArrayBuffer) && !item.IsInstanceOf(String) {
				errCtx := ctx
				return nil, errCtx.NewError("Error: expected an array of ArrayBuffers or Strings, got an array of " + item.TypeName())
			}

			return NewInstance(Blob, ctx), nil
		}),
	}

	return true
}

var _BlobOk = generateBlob()
