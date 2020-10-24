package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type ClassInterface struct {
	clType  *TypeExpression
	extends *TypeExpression // optional, can be nil, TODO: use it

	members            []*FunctionInterface
	cache              map[values.Prototype]string // only check implements once
	explicitImplements []values.Prototype          // all explicit, or all implicit (no mixing)

	TokenData
}

func NewClassInterface(clType *TypeExpression, extends *TypeExpression,
	ctx context.Context) (*ClassInterface, error) {
	ci := &ClassInterface{
		clType,
		extends,
		make([]*FunctionInterface, 0),
		make(map[values.Prototype]string),
		nil,
		TokenData{ctx},
	}

	// change variable so we can register implements classes during resolvenames stage
	ci.clType.ref = NewVariable(ci.Name(), true, ci.clType.Context())
	ci.clType.ref.SetObject(ci)

	return ci, nil
}

func (t *ClassInterface) AddMember(member *FunctionInterface) error {
	// members can have same names, but something must be different (excluding arg names)
	t.members = append(t.members, member)

	return nil
}

func (t *ClassInterface) RegisterExplicitImplementation(cl *Class) {
	if t.explicitImplements == nil {
		t.explicitImplements = make([]values.Prototype, 0)
	}

	t.explicitImplements = append(t.explicitImplements, cl)
}

func (t *ClassInterface) Name() string {
	return t.clType.Name()
}

func (t *ClassInterface) GetVariable() Variable {
	return t.clType.GetVariable()
}

func (t *ClassInterface) AddStatement(st Statement) {
	panic("not a block")
}

func (t *ClassInterface) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Interface(")
	b.WriteString(strings.Replace(t.clType.Dump(""), "\n", "", -1))
	b.WriteString(")\n")

	for _, member := range t.members {
		if prototypes.IsGetter(member) {
			b.WriteString("\n")
			b.WriteString(indent + "  ")
			b.WriteString("getter ")
		} else if prototypes.IsSetter(member) {
			b.WriteString("\n")
			b.WriteString(indent + "  ")
			b.WriteString("setter ")
		} else {
			b.WriteString(indent + "  ")
		}
		b.WriteString(strings.Replace(member.Dump(), "\n", "", -1))
	}

	return b.String()
}

func (t *ClassInterface) WriteStatement(indent string) string {
	return ""
}

func (t *ClassInterface) HasMember(this *values.Instance, key string,
	includePrivate bool) bool {
	// includPrivate is ignored because all interface members are public
	for _, member := range t.members {
		if member.Name() == key {
			return true
		}
	}

	return false
}

func (t *ClassInterface) HoistNames(scope Scope) error {
	return nil
}

func (t *ClassInterface) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	if typeChildren != nil {
		return nil, ctx.NewError("Error: user defined interface cannot have content types")
	}

	if vProto, ok := v.GetInstancePrototype(); ok {
		if msg, ok := t.IsImplementedBy(vProto); !ok {
			return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't implement " + t.Name() + "(" + msg + ")")
		}
	}

	newV, ok := v.ChangeInstanceInterface(t, false, true)
	if !ok {
		return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + t.Name())
	}

	return newV, nil
}

func (t *ClassInterface) ResolveStatementNames(scope Scope) error {
	if scope.HasVariable(t.Name()) {
		errCtx := t.Context()
		err := errCtx.NewError("Error: '" + t.Name() + "' already defined " +
			"(interface needs unique name)")
		other, _ := scope.GetVariable(t.Name())
		err.AppendContextString("Info: defined here ", other.Context())
		return err
	} else {
		if err := scope.SetVariable(t.Name(), t.GetVariable()); err != nil {
			return err
		}

		// interface members cant have default arguments, so inner scope is irrelevant
		for _, member := range t.members {
			if err := member.ResolveInterfaceNames(scope); err != nil {
				return err
			}
		}

		return nil
	}
}

func (t *ClassInterface) HoistValues(stack values.Stack) error {
	return nil
}

func (t *ClassInterface) EvalStatement(stack values.Stack) error {
	// add self immediately, because we might might refer to it in the members
	val := values.NewClassInterface(t, t.Context())
	if err := stack.SetValue(t.GetVariable(), val, false, t.Context()); err != nil {
		return err
	}

	for _, member := range t.members {
		if err := member.EvalInterface(stack); err != nil {
			return err
		}
	}

	return nil
}

func (t *ClassInterface) ResolveStatementActivity(usage Usage) error {
	return nil
}

func (t *ClassInterface) UniversalStatementNames(ns Namespace) error {
	return nil
}

func (t *ClassInterface) UniqueStatementNames(ns Namespace) error {
	return nil
}

func (t *ClassInterface) Walk(fn WalkFunc) error {
  if err := t.clType.Walk(fn); err != nil {
    return err
  }

  if t.extends != nil {
    if err := t.extends.Walk(fn); err != nil {
      return err
    }
  }

  for _, member := range t.members {
    if err := member.Walk(fn); err != nil {
      return err
    }
  }

  return fn(t)
}

func (t *ClassInterface) HasAncestor(interf values.Interface) bool {
	// TODO: should interfaces be able to extends other interfaces?
	if interf == t {
		return true
	} else {
		return false
	}
}

// TODO: reimplement cache in a concurrently safe if it adds a speed advantage
func (t *ClassInterface) IsImplementedBy(proto values.Prototype) (string, bool) {
	//if implements, ok := t.cache[proto]; ok {
	//return implements, implements == ""
	//	} else
	if _, ok := proto.(*values.AllPrototype); ok {
		return "", true
	} else if t.explicitImplements != nil {
		found := false
		for _, testProto := range t.explicitImplements {
			if testProto == proto {
				found = true
				break
			}
		}

		if !found {
			return "all explicit or all implicit", false
		}
	}

	msg := ""

	// check that at least all the members exist
	for _, member := range t.members {
		var ok bool
		msg, ok = member.IsImplementedBy(proto)
		if !ok {
			if msg == "" {
				panic("empty implementedBy hint not allowed")
			}

			break
		}
	}

	//t.cache[proto] = msg

	return msg, msg == ""
}

func (t *ClassInterface) GenerateInstance(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
	if t.explicitImplements == nil {
		return nil, ctx.NewError("Error: interface can only be generated if it is explicitely implemented")
	}

	if keys != nil || args != nil {
		return nil, ctx.NewError("Error: parametric interface not supported")
	}

	var err error
	vs := make([]values.Value, len(t.explicitImplements))

	for i, proto := range t.explicitImplements {
		vs[i], err = proto.GenerateInstance(stack, nil, nil, ctx)
		if err != nil {
			return nil, err
		}
	}

	return values.NewMulti(vs, ctx), nil
}
