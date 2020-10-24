package prototypes

import (
	"../values"

	"../../context"
)

type EvalFunctionType func(values.Stack, *values.Instance, []values.Value, context.Context) (values.Value, error)

// builtin functions are always part of a builtin class
type BuiltinFunction interface {
	Name() string // for debugging
	Length() int

	Role() FunctionRole
	CheckInterface(argsAndRet []*values.NestedType, ctx context.Context) error
	EvalFunction(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) (values.Value, error)
	EvalFunctionNoReturn(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) (values.Value, error)
	EvalMethod(stack values.Stack, this *values.Instance, args []values.Value,
		ctx context.Context) error
	EvalAsEntryPoint(stack values.Stack, this *values.Instance, ctx context.Context) error

	CheckArgs(args []values.Value) bool
}

type BuiltinFunctionData struct {
	methodLike bool // return values can be ignored safely
	role       FunctionRole
	c          ArgCheck
	fn         EvalFunctionType
}

func NewInstance(proto values.Prototype, ctx context.Context) *values.Instance {
	if proto == nil {
		panic("should've been caught before")
	}

	var props values.Properties = nil
	switch {
	case proto.HasAncestor(String):
		props = values.NewStringProperties(false, "", ctx)
	case proto.HasAncestor(Int):
		props = values.NewIntProperties(false, 0, ctx)
	case proto.HasAncestor(Boolean):
		props = values.NewBooleanProperties(false, false, ctx)
	case proto.HasAncestor(Array):
		props = values.NewArrayProperties(false, []values.Value{}, ctx)
	case proto.HasAncestor(Event):
		props = values.NewEventProperties(nil, ctx)
	case proto.HasAncestor(IDBCursorWithValue):
		props = values.NewIDBCursorWithValueProperties(nil, ctx)
	case proto.HasAncestor(IDBRequest):
		props = values.NewIDBRequestProperties(nil, ctx)
	case proto.HasAncestor(Map):
		props = values.NewMapProperties(ctx)
	case proto.HasAncestor(Set):
		props = values.NewSetProperties(ctx)
	case proto.HasAncestor(Number):
		props = values.NewNumberProperties(false, 0.0, ctx)
	case proto.HasAncestor(Object):
		props = values.NewObjectProperties(nil, ctx)
	case proto.HasAncestor(Promise):
		props = values.NewPromiseProperties(ctx)
	default:
		props = values.NewProperties(ctx)
	}

	res := values.NewInstance(proto, props, ctx)
	if res.IsClass() {
		panic("shouldn't be possible")
	}
	return res
}

func (f *BuiltinFunctionData) Role() FunctionRole {
	return f.role
}

func (f *BuiltinFunctionData) Name() string {
	return ""
}

func (f *BuiltinFunctionData) Length() int {
	return -1
}

func (f *BuiltinFunctionData) CheckArgs(args []values.Value) bool {
	if err := CheckInputs(f.c, args, context.NewDummyContext()); err != nil {
		return false
	}

	return true
}

func (f *BuiltinFunctionData) CheckInterface(argsAndRef []*values.NestedType, ctx context.Context) error {
	// because the builtin functions are highly variable, this is very hard to do
	return nil
}

func (f *BuiltinFunctionData) EvalFunction(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	if err := CheckInputs(f.c, args, ctx); err != nil {
		return nil, err
	}

	res, err := f.fn(stack, this, args, ctx)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, ctx.NewError("Error: function doesnt have a return value")
	}

	return res, nil
}

func (f *BuiltinFunctionData) EvalFunctionNoReturn(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	return f.fn(stack, this, args, ctx)
}

func (f *BuiltinFunctionData) EvalMethod(stack values.Stack, this *values.Instance, args []values.Value,
	ctx context.Context) error {
	if f.fn == nil {
		panic("shouldn't be nil (" + f.Name() + ")")
	}

	if err := CheckInputs(f.c, args, ctx); err != nil {
		return err
	}

	res, err := f.fn(stack, this, args, ctx)
	if err != nil {
		return err
	}

	if res == nil || f.methodLike {
		return nil
	} else {
		return ctx.NewError("Error: function has return value (hint: use void)")
	}
}

func (f *BuiltinFunctionData) EvalAsEntryPoint(stack values.Stack, this *values.Instance, ctx context.Context) error {
	return nil
}

func newNormal(c ArgCheck, output values.Prototype, methodLike bool) *BuiltinFunctionData {
	return &BuiltinFunctionData{methodLike, NORMAL, c,
		func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
			if output == nil {
				return nil, nil
			} else {
				return NewInstance(output, ctx), nil
			}
		},
	}
}

func NewNormal(c ArgCheck, output values.Prototype) *BuiltinFunctionData {
	return newNormal(c, output, false)
}

func NewMethodLikeNormal(c ArgCheck, output values.Prototype) *BuiltinFunctionData {
	return newNormal(c, output, true)
}

func newNormalFunction(c ArgCheck, fn EvalFunctionType, methodLike bool) *BuiltinFunctionData {
	return &BuiltinFunctionData{methodLike, NORMAL, c, fn}
}

func NewNormalFunction(c ArgCheck, fn EvalFunctionType) *BuiltinFunctionData {
	return newNormalFunction(c, fn, false)
}

func NewMethodLikeNormalFunction(c ArgCheck, fn EvalFunctionType) *BuiltinFunctionData {
	return newNormalFunction(c, fn, true)
}

func NewGetter(output values.Prototype) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, GETTER, &None{},
		func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
			return NewInstance(output, ctx), nil
		},
	}
}

func NewGetterFunction(fn EvalFunctionType) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, GETTER, &None{}, fn}
}

func NewStaticGetter(output values.Prototype) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, STATIC | GETTER, &None{},
		func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
			return NewInstance(output, ctx), nil
		},
	}
}

func NewStaticClassGetter(classProto values.Prototype) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, STATIC | GETTER, &None{},
		func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
			return values.NewClass(classProto, ctx), nil
		},
	}
}

func NewStatic(c ArgCheck, output values.Prototype) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, STATIC, c,
		func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
			if output == nil {
				return nil, nil
			} else {
				return NewInstance(output, ctx), nil
			}
		},
	}
}

func NewStaticFunction(c ArgCheck, fn EvalFunctionType) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, STATIC, c, fn}
}

func NewStaticGetterFunction(fn EvalFunctionType) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, STATIC | GETTER, &None{}, fn}
}

func NewSetter(c ArgCheck) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, SETTER, c,
		func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
			return nil, nil
		},
	}
}

func NewSetterFunction(c ArgCheck, fn EvalFunctionType) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, SETTER, c, fn}
}

func NewGetterSetter(proto values.Prototype) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, GETTER | SETTER, &Rest{&Any{}}, // check is done internally
		func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
			if len(args) == 1 { // setter
				return nil, CheckInputs(&PrototypeCheck{proto}, args, ctx)
			} else if len(args) == 0 { // getter
				return NewInstance(proto, ctx), nil
			} else {
				panic("unexpected number of arguments")
			}
		},
	}
}

func NewGetterSetterFunction(proto values.Prototype, fn EvalFunctionType) *BuiltinFunctionData {
	return &BuiltinFunctionData{false, GETTER | SETTER, &Rest{&Any{}},
		func(stack values.Stack, this *values.Instance, args []values.Value, ctx context.Context) (values.Value, error) {
			if len(args) == 1 {
				if err := CheckInputs(&PrototypeCheck{proto}, args, ctx); err != nil {
					return nil, err
				}
			} else if len(args) == 0 {
				if err := CheckInputs(&None{}, args, ctx); err != nil {
					return nil, err
				}
			}

			return fn(stack, this, args, ctx)
		}}
}
