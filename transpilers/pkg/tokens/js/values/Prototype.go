package values

import (
	"../../context"
)

type Prototype interface {
	Interface

	GetParent() Prototype // can be nil

	// childProto can be used to inject implementations of abstract members (usually childProto is nil though)
	EvalConstructor(stack Stack, args []Value, childProto Prototype, ctx context.Context) (Value, error)

	// so that globally exported classes can be evaluated with generated arguments
	// builtin prototypes ignore this
	EvalAsEntryPoint(stack Stack, ctx context.Context) error
	IsUniversal() bool // if true: can be exported to databases etc.

	// evaluates Getters, but returns Normal or Static functions
	GetMember(stack Stack, this *Instance, key string, includePrivate bool,
		ctx context.Context) (Value, error)

	// evaluates Setters
	SetMember(stack Stack, this *Instance, key string, arg Value, includePrivate bool,
		ctx context.Context) error

	// special for String, Array, Object
	GetIndex(stack Stack, this *Instance, index Value,
		ctx context.Context) (Value, error)
	SetIndex(stack Stack, this *Instance, index Value, arg Value,
		ctx context.Context) error

	LoopForIn(this *Instance, fn func(Value) error, ctx context.Context) error
	LoopForOf(this *Instance, fn func(Value) error, ctx context.Context) error
}
