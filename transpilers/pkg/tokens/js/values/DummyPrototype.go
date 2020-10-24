package values

import (
	"../../context"
)

// a placeholder for "any", "void" and "class" in TypeExpressions
type DummyPrototype struct {
	name string
}

func NewDummyPrototype(name string) *DummyPrototype {
	return &DummyPrototype{name}
}

func (p *DummyPrototype) Name() string {
	return p.name
}

func (p *DummyPrototype) CastInstance(v *Instance, typeChildren []*NestedType, ctx context.Context) (Value, error) {
	panic("shouldn't be called")
}

func (p *DummyPrototype) HasMember(this *Instance, key string, includePrivate bool) bool {
	panic("shouldn't be called")
}

func (p *DummyPrototype) HasAncestor(interf Interface) bool {
	panic("shouldn't be called")
}

func (p *DummyPrototype) IsImplementedBy(proto Prototype) (string, bool) {
	panic("shouldn't be called")
}

func (p *DummyPrototype) GetParent() Prototype {
	panic("shouldn't be called")
}

func (p *DummyPrototype) EvalConstructor(stack Stack, args []Value, childProto Prototype, ctx context.Context) (Value, error) {
	panic("shouldn't be called")
}

func (p *DummyPrototype) GenerateInstance(stack Stack, keys []string, args []Value, ctx context.Context) (Value, error) {
	panic("shouldn't be called")
}

func (p *DummyPrototype) EvalAsEntryPoint(stack Stack, ctx context.Context) error {
	panic("shouldn't be called")
}

func (p *DummyPrototype) IsUniversal() bool {
	panic("shouldn't be called")
}

func (p *DummyPrototype) GetMember(stack Stack, this *Instance, key string, includePrivate bool, ctx context.Context) (Value, error) {
	panic("shouldn't be called")
}

func (p *DummyPrototype) SetMember(stack Stack, this *Instance, key string, arg Value, includePrivate bool, ctx context.Context) error {
	panic("shouldn't be called")
}

func (p *DummyPrototype) GetIndex(stack Stack, this *Instance, index Value, ctx context.Context) (Value, error) {
	panic("shouldn't be called")
}

func (p *DummyPrototype) SetIndex(stack Stack, this *Instance, index Value, arg Value, ctx context.Context) error {
	panic("shouldn't be called")
}

func (p *DummyPrototype) LoopForIn(this *Instance, fn func(Value) error, ctx context.Context) error {
	panic("shouldn't be called")
}

func (p *DummyPrototype) LoopForOf(this *Instance, fn func(Value) error, ctx context.Context) error {
	panic("shouldn't be called")
}
