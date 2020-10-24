package prototypes

import (
	"../values"

	"../../context"
)

type BuiltinInterfaceMember struct {
	args ArgCheck
	role FunctionRole
	ret  values.Prototype // can be nil for no return value
}

type BuiltinInterface struct {
	name string

	Members map[string]*BuiltinInterfaceMember

	cache map[values.Prototype]string // only check implements once
}

var PrototypeImplements func(proto values.Prototype, interf values.Interface) (string, bool) = nil

func RegisterPrototypeImplements(fn func(proto values.Prototype, interf values.Interface) (string, bool)) bool {
	PrototypeImplements = fn

	return true
}

func allocBuiltinInterface() *BuiltinInterface {
	return &BuiltinInterface{"", map[string]*BuiltinInterfaceMember{}, map[values.Prototype]string{}}
}

func (m *BuiltinInterfaceMember) Role() FunctionRole {
	return m.role
}

func (m *BuiltinInterfaceMember) Check(args []interface{}) bool {
	pos := 0

	var err error
	pos, err = m.args.Check(args, pos, context.NewDummyContext())
	if err != nil {
		return false
	}

	if pos != len(args) {
		return false
	}

	return true
}

func (p *BuiltinInterfaceMember) CheckRetType(retName string) bool {
	// builtin interfaces can't be part of packages this way!
	if p.ret == nil {
		return retName == ""
	} else {
		return retName == p.ret.Name()
	}
}

func (p *BuiltinInterface) Name() string {
	return p.name
}

func (p *BuiltinInterface) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckInterface(p, args, pos, ctx)
}

func (p *BuiltinInterface) HasMember(this *values.Instance, key string, includePrivate bool) bool {
	_, ok := p.Members[key]
	return ok
}

func (p *BuiltinInterface) HasAncestor(interf values.Interface) bool {
	return false
}

func (p *BuiltinInterface) IsImplementedBy(proto values.Prototype) (string, bool) {
	if p.cache == nil {
		p.cache = make(map[values.Prototype]string)
	}

	if implements, ok := p.cache[proto]; ok {
		return implements, implements == ""
	} else if _, ok := proto.(*values.AllPrototype); ok {
		return "", true
	}

	msg, _ := PrototypeImplements(proto, p)

	p.cache[proto] = msg

	return msg, msg == ""
}

func (p *BuiltinInterface) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	if typeChildren != nil {
		return nil, ctx.NewError("Error: " + p.Name() + " can't have content types")
	}

	newV, ok := v.ChangeInstanceInterface(p, false, true)
	if !ok {
		return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + p.Name())
	}

	return newV, nil

}

func (p *BuiltinInterface) GenerateInstance(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
	return nil, ctx.NewError("Error: builtin interface can't yet be generated")
}
