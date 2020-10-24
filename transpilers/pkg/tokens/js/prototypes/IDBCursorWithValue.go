package prototypes

import (
	"../values"
	"fmt"

	"../../context"
)

type IDBCursorWithValuePrototype struct {
	BuiltinPrototype
}

var IDBCursorWithValue *IDBCursorWithValuePrototype = allocIDBCursorWithValuePrototype()

func allocIDBCursorWithValuePrototype() *IDBCursorWithValuePrototype {
	return &IDBCursorWithValuePrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

func (p *IDBCursorWithValuePrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *IDBCursorWithValuePrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*IDBCursorWithValuePrototype); ok {
		if other == p {
			return true
		} else {
			return false
		}
	} else {
		_, ok = other_.IsImplementedBy(p)
		return ok
	}
}

func (p *IDBCursorWithValuePrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	newV_, ok := v.ChangeInstanceInterface(p, false, true)
	if !ok {
		return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + p.Name())
	}

	newV, ok := newV_.(*values.Instance)
	if !ok {
		panic("unexpected")
	}

	if typeChildren == nil {
		return newV, nil
	} else {
		if len(typeChildren) != 1 {
			return nil, ctx.NewError(fmt.Sprintf("Error: IDBCursorWithValue expects 1 type child, got %d", len(typeChildren)))
		}

		typeChild := typeChildren[0]

		// now cast all the items
		props := values.AssertIDBCursorWithValueProperties(newV.Properties())

		cursorValue := props.Value()
		var err error
		cursorValue, err = cursorValue.Cast(typeChild, ctx)
		if err != nil {
			return nil, err
		}

		newV = NewIDBCursorWithValue(cursorValue, ctx)
		return newV, nil
	}
}

func NewIDBCursorWithValue(v values.Value, ctx context.Context) *values.Instance {
	return values.NewInstance(IDBCursorWithValue, values.NewIDBCursorWithValueProperties(v, ctx), ctx)
}

func generateIDBCursorWithValuePrototype() bool {
	*IDBCursorWithValue = IDBCursorWithValuePrototype{BuiltinPrototype{
		"IDBCursorWithValue", IDBCursor,
		map[string]BuiltinFunction{
			"value": NewGetterFunction(func(stack values.Stack, this *values.Instance,
				args []values.Value, ctx context.Context) (values.Value, error) {
				props := values.AssertIDBCursorWithValueProperties(this.Properties())
				return props.Value(), nil
			}),
		},
		NewNoContentGenerator(IDBCursorWithValue),
	}}

	return true
}

var _IDBCursorWithValueOk = generateIDBCursorWithValuePrototype()
