package values

import (
	"../../context"
)

// returned by AllNull
type AllPrototype struct {
}

func (p *AllPrototype) Name() string {
	return "All"
}

func (p *AllPrototype) IsImplementedBy(proto Prototype) (string, bool) {
	return "", false
}

func (p *AllPrototype) HasAncestor(interf Interface) bool {
	return true
}

func (p *AllPrototype) GetParent() Prototype {
	return nil
}

func (p *AllPrototype) IsUniversal() bool {
	return true // true produces less errors than false
}

func (p *AllPrototype) CastInstance(v *Instance, typeChildren []*NestedType, ctx context.Context) (Value, error) {
	// impossible to cast
	return v, nil
}

func (p *AllPrototype) EvalConstructor(stack Stack, args []Value,
	childProto Prototype, ctx context.Context) (Value, error) {
	return nil, ctx.NewError("Error: cant construct null")
}

func (p *AllPrototype) GenerateInstance(stack Stack, keys []string, args []Value,
	ctx context.Context) (Value, error) {
	return NewAllNull(ctx), nil
}

func (p *AllPrototype) EvalAsEntryPoint(stack Stack, ctx context.Context) error {
	return nil
}

func (p *AllPrototype) HasMember(this *Instance, key string, includePrivate bool) bool {
	return true
}

func (p *AllPrototype) GetMember(stack Stack, this *Instance, key string, includePrivate bool,
	ctx context.Context) (Value, error) {
	return NewAllNull(ctx), nil
}

func (p *AllPrototype) SetMember(stack Stack, this *Instance, key string, arg Value,
	includePrivate bool, ctx context.Context) error {
	return nil
}

func (p *AllPrototype) GetIndex(stack Stack, this *Instance, index Value,
	ctx context.Context) (Value, error) {

	return NewAllNull(ctx), nil
}

func (p *AllPrototype) SetIndex(stack Stack, this *Instance, index Value, arg Value,
	ctx context.Context) error {
	return nil
}

func (p *AllPrototype) LoopForIn(this *Instance, fn func(Value) error, ctx context.Context) error {
	return nil
}

func (p *AllPrototype) LoopForOf(this *Instance, fn func(Value) error, ctx context.Context) error {
	return nil
}
