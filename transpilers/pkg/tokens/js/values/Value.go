package values

import (
	"reflect"

	"../../context"
)

var VERBOSITY = 0

type Value interface {
	TypeName() string
	Context() context.Context

	Copy(CopyCache) Value
	Cast(ntype *NestedType, ctx context.Context) (Value, error)

	// cannot be instance of Interface
	IsInstanceOf(ps ...Prototype) bool // ps[0] or ps[1] or ps[2] ... etc.
	MaybeInstanceOf(p Prototype) bool  // only relevant for Multi
	Merge(other Value) Value           // nil if not mergeable (can be used for IsSame test
	RemoveLiteralness(all bool) Value  // if all==false then only boolean and number are reset

	EvalFunction(stack Stack, args []Value, ctx context.Context) (Value, error)
	EvalFunctionNoReturn(stack Stack, args []Value, ctx context.Context) (Value, error)
	EvalMethod(stack Stack, args []Value, ctx context.Context) error
	EvalConstructor(stack Stack, args []Value, ctx context.Context) (Value, error)
	EvalAsEntryPoint(stack Stack, ctx context.Context) error

	GetMember(stack Stack, key string, includePrivate bool,
		ctx context.Context) (Value, error)
	SetMember(stack Stack, key string, value Value, includePrivate bool,
		ctx context.Context) error
	GetIndex(stack Stack, index Value, ctx context.Context) (Value, error)
	SetIndex(stack Stack, index Value, value Value, ctx context.Context) error

	LoopForIn(fn func(Value) error, ctx context.Context) error
	LoopForOf(fn func(Value) error, ctx context.Context) error
	LoopNestedPrototypes(fn func(Prototype))

	// value type checking
	IsClass() bool
	IsFunction() bool
	IsInstance() bool
	IsNull() bool
	IsVoid() bool
	IsInterface() bool

	// for literal operators
	LiteralIntValue() (int, bool)
	LiteralBooleanValue() (bool, bool)
	LiteralNumberValue() (float64, bool)
	LiteralStringValue() (string, bool)
	LiteralArrayValues() ([]Value, bool)

	GetClassPrototype() (Prototype, bool)
	GetClassInterface() (Interface, bool)
	GetInstancePrototype() (Prototype, bool)
	GetNullPrototype() (Prototype, bool)
	ChangeInstancePrototype(p Prototype, inPlace bool) (Value, bool)
	ChangeInstanceInterface(interf Interface, inPlace bool, checkOuter bool) (Value, bool)
}

type ValueData struct {
	ctx context.Context
}

func (v *ValueData) Context() context.Context {
	return v.ctx
}

func (v *ValueData) IsInstanceOf(ps ...Prototype) bool {
	return false
}

func (v *ValueData) MaybeInstanceOf(p Prototype) bool {
	return false
}

func (v *ValueData) EvalFunction(stack Stack, args []Value, ctx context.Context) (Value, error) {
	return nil, ctx.NewError("Error: not a function (" + reflect.TypeOf(v).String() + ")")
}

func (v *ValueData) EvalFunctionNoReturn(stack Stack, args []Value,
	ctx context.Context) (Value, error) {
	return nil, ctx.NewError("Error: not a function (" + reflect.TypeOf(v).String() + ")")
}

func (v *ValueData) EvalMethod(stack Stack, args []Value, ctx context.Context) error {
	return ctx.NewError("Error: not a function (" + reflect.TypeOf(v).String() + ")")
}

func (v *ValueData) EvalConstructor(stack Stack, args []Value, ctx context.Context) (Value, error) {
	return nil, ctx.NewError("Error: not a class")
}

func (v *ValueData) EvalAsEntryPoint(stack Stack, ctx context.Context) error {
	// valid entrypoints are classes and functions
	return ctx.NewError("Error: not an entry point")
}

func (v *ValueData) GetMember(stack Stack, key string, includePrivate bool,
	ctx context.Context) (Value, error) {
	return nil, ctx.NewError("Error: not a class or instance")
}

func (v *ValueData) SetMember(stack Stack, key string, value Value, includePrivate bool,
	ctx context.Context) error {
	return ctx.NewError("Error: not a class or instance")
}

func (v *ValueData) GetIndex(stack Stack, index Value, ctx context.Context) (Value, error) {
	return nil, ctx.NewError("Error: not indexable")
}

func (v *ValueData) SetIndex(stack Stack, index Value, value Value, ctx context.Context) error {
	return ctx.NewError("Error: not indexable")
}

func (v *ValueData) LoopForIn(fn func(Value) error, ctx context.Context) error {
	return ctx.NewError("Error: for in doesnt apply")
}

func (v *ValueData) LoopForOf(fn func(Value) error, ctx context.Context) error {
	return ctx.NewError("Error: not iterable")
}

func (v *ValueData) IsClass() bool {
	return false
}

func (v *ValueData) IsFunction() bool {
	return false
}

func (v *ValueData) IsInstance() bool {
	return false
}

func (v *ValueData) IsNull() bool {
	return false
}

func (v *ValueData) IsVoid() bool {
	return false
}

func (v *ValueData) IsInterface() bool {
	return false
}

func (v *ValueData) LiteralBooleanValue() (bool, bool) {
	return false, false
}

func (v *ValueData) LiteralIntValue() (int, bool) {
	return 0, false
}

func (v *ValueData) LiteralNumberValue() (float64, bool) {
	return 0.0, false
}

func (v *ValueData) LiteralStringValue() (string, bool) {
	return "", false
}

func (v *ValueData) LiteralArrayValues() ([]Value, bool) {
	return nil, false
}

func (v *ValueData) GetClassPrototype() (Prototype, bool) {
	return nil, false
}

func (v *ValueData) GetClassInterface() (Interface, bool) {
	return nil, false
}

func (v *ValueData) GetInstancePrototype() (Prototype, bool) {
	return nil, false
}

func (v *ValueData) GetNullPrototype() (Prototype, bool) {
	return nil, false
}

func (v *ValueData) ChangeInstancePrototype(p Prototype, inPlace bool) (Value, bool) {
	return nil, false
}

func (v *ValueData) ChangeInstanceInterface(interf Interface, inPlace bool, checkOuter bool) (Value, bool) {
	return nil, false
}

// l1 and l2 might have different order
func mergeValueLists(l1, l2 []Value) []Value {
	if len(l1) != len(l2) {
		return nil
	}

	// values might be in different order
	res := make([]Value, len(l1)) // XXX: hopefully filled with null by default

Outer:
	for _, v1 := range l1 {
		for i, v2 := range l2 {
			if res[i] != nil {
				continue
			} else if vres := v2.Merge(v1); vres != nil {
				res[i] = vres
				continue Outer
			}
		}

		return nil
	}

	return res
}

func stringListsSameContent(l1, l2 []string) bool {
	if len(l1) != len(l2) {
		return false
	}

	done := make([]bool, len(l1))

Outer:
	for _, v1 := range l1 {
		for i, v2 := range l2 {
			if done[i] {
				continue
			} else if v2 == v1 {
				done[i] = true
				continue Outer
			}
		}

		return false
	}

	return true
}

func mergeValueListsSameOrder(l1, l2 []Value) []Value {
	if len(l1) != len(l2) {
		return nil
	}

	res := make([]Value, len(l1))

	for i, v1 := range l1 {
		v2 := l2[i]

		v := v1.Merge(v2)
		if v == nil {
			return nil
		}

		res[i] = v
	}

	return res
}

func copyValueList(l []Value, cache CopyCache) []Value {
	res := make([]Value, len(l))

	for i, v := range l {
		res[i] = v.Copy(cache)
	}

	return res
}

func InstancePrototypeIsAncestorOf(v Value, ps ...Prototype) bool {
	if proto, ok := v.GetInstancePrototype(); ok {
		for _, pr := range ps {
			if pr.HasAncestor(proto) {
				return true
			}
		}
	}

	return false
}

func IsLiteral(v Value) bool {
	if _, ok := v.LiteralStringValue(); ok {
		return true
	} else if _, ok := v.LiteralIntValue(); ok {
		return true
	} else if _, ok := v.LiteralBooleanValue(); ok {
		return true
	} else if _, ok := v.LiteralNumberValue(); ok {
		return true
	} else if _, ok := v.LiteralArrayValues(); ok {
		return true
	} else {
		return false
	}
}
