package values

import (
	"../../context"
)

type Interface interface {
	Name() string

	// ignore this
	HasMember(this *Instance, key string, includePrivate bool) bool

	HasAncestor(interf Interface) bool

	// the return strings states whatever is missing/differs
	IsImplementedBy(proto Prototype) (string, bool)

	// leave type assertion Value->Instance to implementation
	CastInstance(v *Instance, ntypes []*NestedType, ctx context.Context) (Value, error)

	// generate instance using content typing only
	GenerateInstance(stack Stack, keys []string, args []Value, ctx context.Context) (Value, error)
}
