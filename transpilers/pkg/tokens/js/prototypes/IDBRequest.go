package prototypes

import (
	"fmt"

	"../values"

	"../../context"
)

type IDBRequestPrototype struct {
	BuiltinPrototype
}

var IDBRequest *IDBRequestPrototype = allocIDBRequestPrototype()

func allocIDBRequestPrototype() *IDBRequestPrototype {
	return &IDBRequestPrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

func (p *IDBRequestPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *IDBRequestPrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*IDBRequestPrototype); ok {
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

func (p *IDBRequestPrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
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
			return nil, ctx.NewError(fmt.Sprintf("Error: IDBRequest expects 1 type child, got %d", len(typeChildren)))
		}

		typeChild := typeChildren[0]

		// now cast all the items
		props := values.AssertIDBRequestProperties(newV.Properties())

		result := props.Result()
		var err error
		result, err = result.Cast(typeChild, ctx)
		if err != nil {
			return nil, err
		}

		newV = NewIDBRequest(result, ctx)
		return newV, nil
	}
}

func NewIDBRequest(v values.Value, ctx context.Context) *values.Instance {
	return values.NewInstance(IDBRequest, values.NewIDBRequestProperties(v, ctx), ctx)
}

func generateIDBRequestCallback(eventProto values.Prototype) func(values.Stack, *values.Instance, []values.Value, context.Context) (values.Value, error) {
	return func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
		arg := args[0]

		event := NewAltEvent(eventProto, this, ctx)
		if err := arg.EvalMethod(stack.Parent(), []values.Value{event}, ctx); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func idbRequestCallback(stack values.Stack, this *values.Instance, args []values.Value,
	ctx context.Context) (values.Value, error) {
	arg := args[0]

	event := NewEvent(this, ctx)
	if err := arg.EvalMethod(stack.Parent(), []values.Value{event}, ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func generateIDBRequestPrototype() bool {
	*IDBRequest = IDBRequestPrototype{BuiltinPrototype{
		"IDBRequest", EventTarget,
		map[string]BuiltinFunction{
			"onerror":   NewSetterFunction(&Function{}, idbRequestCallback),
			"onsuccess": NewSetterFunction(&Function{}, idbRequestCallback),
			"result": NewGetterFunction(func(stack values.Stack, this *values.Instance,
				args []values.Value, ctx context.Context) (values.Value, error) {
				props := values.AssertIDBRequestProperties(this.Properties())
				return props.Result(), nil
			}),
		},
		nil,
	}}

	return true
}

var _IDBRequestOk = generateIDBRequestPrototype()
