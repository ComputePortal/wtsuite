package prototypes

import (
	"../values"

	"../../context"
)

var IDBIndex *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBIndex() bool {
	*IDBIndex = BuiltinPrototype{
		"IDBIndex", nil,
		map[string]BuiltinFunction{
			"getAll": NewNormalFunction(&And{&Opt{IDBKeyRange}, &Opt{Int}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewIDBRequest(NewArray([]values.Value{NewObject(nil, ctx)}, ctx), ctx), nil
				}),
			"getAllKeys": NewNormalFunction(&And{&Opt{IDBKeyRange}, &Opt{Int}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					// why would we be interested in the Int keys?
					return NewIDBRequest(NewArray([]values.Value{NewInt(ctx)}, ctx), ctx), nil
				}),
			"openCursor": NewNormalFunction(&And{&Opt{&Or{Int, IDBKeyRange}}, &Opt{String}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewIDBRequest(NewIDBCursorWithValue(NewInstance(Object, ctx), ctx), ctx), nil
				}),
		},
		nil,
	}

	return true
}

var _IDBIndexOk = generateIDBIndex()
