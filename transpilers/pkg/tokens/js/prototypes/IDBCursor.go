package prototypes

import (
	"../values"

	"../../context"
)

var IDBCursor *BuiltinPrototype = allocBuiltinPrototype()

func generateIDBCursorPrototype() bool {
	*IDBCursor = BuiltinPrototype{
		"IDBCursor", nil,
		map[string]BuiltinFunction{
			"advance":            NewNormal(Int, nil),
			"continue":           NewNormal(&Opt{Int}, nil),
			"continuePrimaryKey": NewNormal(&And{Int, Int}, nil),
			"delete": NewNormalFunction(&None{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					return NewInstance(IDBRequest, ctx), nil
				}),
			"key": NewGetter(Int),
			"update": NewNormalFunction(&Any{},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					// should updated value be same type are current value?
					// XXX: should request.result = args[0] ?
					return NewInstance(IDBRequest, ctx), nil
				}),
		},
		nil,
	}

	return true
}

var _IDBCursorOk = generateIDBCursorPrototype()
