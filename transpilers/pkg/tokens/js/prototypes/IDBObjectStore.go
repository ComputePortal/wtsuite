package prototypes

import (
	"../values"

	"../../context"
)

var IDBObjectStore *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBObjectStorePrototype() bool {
	*IDBObjectStore = BuiltinPrototype{
		"IDBObjectStore", nil,
		map[string]BuiltinFunction{
			"add": NewNormalFunction(&Any{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					// an autoincrement key is always an int
					return NewIDBRequest(NewInstance(Int, ctx), ctx), nil
				}),
			"clear": NewNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewIDBRequest(nil, ctx), nil
				}),
			"count": NewNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewIDBRequest(NewInstance(Int, ctx), ctx), nil
				}),
			"createIndex": NewMethodLikeNormalFunction(&And{String, &And{String, &Opt{Object}}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					// TODO check content of object
					return NewInstance(IDBIndex, ctx), nil
				}),
			"delete": NewNormalFunction(&Or{Int, String},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewIDBRequest(nil, ctx), nil
				}),
			"get": NewNormalFunction(&Or{Int, String},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewIDBRequest(NewInstance(Object, ctx), ctx), nil
				}),
			"index": NewNormal(String, IDBIndex),
			"openCursor": NewNormalFunction(&And{&Opt{&Or{Int, IDBKeyRange}}, &Opt{String}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewIDBRequest(NewIDBCursorWithValue(NewInstance(Object, ctx), ctx), ctx), nil
				}),
			"put": NewNormalFunction(&And{&Any{}, &Opt{&Or{Int, String}}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewIDBRequest(nil, ctx), nil
				}),
		},
		nil,
	}

	return true
}

var _IDBObjectStoreOk = generateIDBObjectStorePrototype()
