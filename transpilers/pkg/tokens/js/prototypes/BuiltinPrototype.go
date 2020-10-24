package prototypes

import (
	"../values"

	"../../context"
)

var VERBOSITY = 0

// so we can inherit and still detect the right methods
type BuiltinPrototypeInterface interface {
	HasMember(this *values.Instance, key string, includePrivate bool) bool
	FindMember(key string) BuiltinFunction

	GetParent() values.Prototype
}

// very different from ../Prototype.go
type BuiltinPrototype struct {
	name   string
	parent values.Prototype // builtin types only have one parent

	// use a map because here we can combine getters/setters into a single function,
	//  in ../Prototype.go a list must be used because getters and setters might have
	//  the same name and cant be treated together
	members map[string]BuiltinFunction

	constructor BuiltinConstructor // the builtin constructors must return a value!
}

var isClassMacro func(string, string) bool = nil

func RegisterIsClassMacro(fn func(string, string) bool) bool {
	isClassMacro = fn
	return true
}

func allocBuiltinPrototype() *BuiltinPrototype {
	return &BuiltinPrototype{"", nil, map[string]BuiltinFunction{}, nil}
}

// implement ArgCheck, so we can use it directly for input argument checking
func (p *BuiltinPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *BuiltinPrototype) Name() string {
	return p.name
}

func (p *BuiltinPrototype) IsImplementedBy(other_ values.Prototype) (string, bool) {
	return "", false
}

func (p *BuiltinPrototype) GetParent() values.Prototype {
	return p.parent
}

func (p *BuiltinPrototype) HasAncestor(other_ values.Interface) bool {
	if _, ok := other_.IsImplementedBy(p); ok {
		return true
	} else if other, ok := other_.(BuiltinPrototypeInterface); ok {
		if other == p {
			return true
		} else if p.GetParent() != nil {
			return p.GetParent().HasAncestor(other_)
		} else {
			return false
		}
	} else if _, ok := other_.(*values.AllPrototype); ok {
		return true
	} else {
		return false
	}
}

func (p *BuiltinPrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	if typeChildren != nil {
		return nil, ctx.NewError("Error: " + p.Name() + " can't have content types")
	}

	newV, ok := v.ChangeInstanceInterface(p, false, true)
	if !ok {
		return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + p.Name())
	}

	return newV, nil
}

func (p *BuiltinPrototype) EvalConstructor(stack values.Stack, args []values.Value,
	childProto values.Prototype, ctx context.Context) (values.Value, error) {
	if p.constructor == nil || !p.constructor.hasConstructorFn() {
		return nil, ctx.NewError("Error: " + p.Name() + " doesn't have a constructor")
	}

	// TODO: should we pass on childProto?

	return p.constructor.EvalConstructor(stack, args, ctx)
}

func (p *BuiltinPrototype) GenerateInstance(stack values.Stack, keys []string, args []values.Value,
	ctx context.Context) (values.Value, error) {
	if p.constructor == nil || !p.constructor.hasInstanceGenerator() {
		return nil, ctx.NewError("Error: " + p.Name() + " doesn't have a generator")
	}

	return p.constructor.GenerateInstance(stack, keys, args, ctx)
}

func (p *BuiltinPrototype) EvalAsEntryPoint(stack values.Stack, ctx context.Context) error {
	return nil
}

func (p *BuiltinPrototype) FindMember(key string) BuiltinFunction {
	if member, ok := p.members[key]; ok {
		return member
	} else if p.parent != nil {
		parent, ok := p.parent.(BuiltinPrototypeInterface)
		if !ok {
			panic("unexpected")
		}
		return parent.FindMember(key)
	} else {
		return nil
	}
}

func (p *BuiltinPrototype) FindMemberCheckStatic(key string, this *values.Instance,
	ctx context.Context) (BuiltinFunction, error) {
	member := p.FindMember(key)
	if member == nil {
		err := ctx.NewError("Error: " + p.Name() + " doesn't have a member named " + key)
		if isClassMacro(p.Name(), key) {
			err.AppendString("Hint: class macros always require parentheses (i.e. are not getters)")
		}
		return nil, err
	}

	if this != nil && IsStatic(member) {
		return nil, ctx.NewError("Error: cannot call static method on instance")
	} else if this == nil && !IsStatic(member) {
		return nil, ctx.NewError("Error: cannot call non-static method on class")
	}

	return member, nil
}

func (p *BuiltinPrototype) HasMember(this *values.Instance, key string, includePrivate bool) bool {
	if _, ok := p.members[key]; ok {
		return true
	} else if key == "constructor" && p.constructor != nil {
		return true
	} else if p.parent != nil {
		parent, ok := p.parent.(BuiltinPrototypeInterface)
		if !ok {
			panic("unexpected")
		}
		return parent.HasMember(this, key, includePrivate)
	} else {
		return false
	}
}

func (p *BuiltinPrototype) GetMember(stack values.Stack, this *values.Instance, key string, includePrivate bool,
	ctx context.Context) (values.Value, error) {
	member, err := p.FindMemberCheckStatic(key, this, ctx)
	if err != nil {
		return nil, err
	}

	switch {
	case IsStaticGetter(member):
		return member.EvalFunction(stack, nil, []values.Value{}, ctx)
	case IsStatic(member):
		return values.NewFunction(member, stack, this, ctx), nil
	case IsGetter(member): // must come before 'case IsSetter()' because a member might be both a setter and getter at the same time
		return member.EvalFunction(stack, this, []values.Value{}, ctx)
	case IsSetter(member):
		return nil, ctx.NewError("Error: " + p.Name() + "." + key + " is a setter")
	case IsNormal(member):
		return values.NewFunction(member, stack, this, ctx), nil
	default:
		panic("unhandled")
	}
}

func (p *BuiltinPrototype) SetMember(stack values.Stack, this *values.Instance, key string, arg values.Value,
	includePrivate bool, ctx context.Context) error {
	member, err := p.FindMemberCheckStatic(key, this, ctx)
	if err != nil {
		return err
	}

	switch {
	case IsSetter(member):
		return member.EvalMethod(stack, this, []values.Value{arg}, ctx)
	default:
		return ctx.NewError("Error: " + p.Name() + "." + key + " is not a setter")
	}
}

func (p *BuiltinPrototype) GetIndex(stack values.Stack, this *values.Instance,
	index values.Value, ctx context.Context) (values.Value, error) {
	if p.parent == nil {
		return nil, ctx.NewError("Error: not indexable")
	} else {
		return p.parent.GetIndex(stack, this, index, ctx)
	}
}

func (p *BuiltinPrototype) SetIndex(stack values.Stack, this *values.Instance,
	index values.Value, arg values.Value, ctx context.Context) error {
	if p.parent == nil {
		return ctx.NewError("Error: not indexable")
	} else {
		return p.parent.SetIndex(stack, this, index, arg, ctx)
	}
}

func (p *BuiltinPrototype) LoopForOf(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	if p.parent == nil {
		return ctx.NewError("Error: 'for of' not possible")
	} else {
		return p.parent.LoopForOf(this, fn, ctx)
	}
}

func (p *BuiltinPrototype) LoopForIn(this *values.Instance, fn func(values.Value) error,
	ctx context.Context) error {
	if p.parent == nil {
		return ctx.NewError("Error: 'for in' not possible")
	} else {
		return p.parent.LoopForIn(this, fn, ctx)
	}
}

func (p *BuiltinPrototype) IsUniversal() bool {
	return true
}
