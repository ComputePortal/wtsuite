package js

import (
	"strconv"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type EnumMember struct {
	key   *Word
	val Expression
}

// enum is a statement (and also a Class!)
type Enum struct {
	nameExpr   *TypeExpression
	parentExpr *TypeExpression // Int, String or something else, never nil
  members    []*EnumMember
	TokenData
}

func NewEnum(nameExpr *TypeExpression, parentExpr *TypeExpression, keys []*Word,
	vs []Expression, ctx context.Context) (*Enum, error) {
  members := make([]*EnumMember, len(keys))
  for i, key := range keys {
    if key.Value() == "values" || key.Value() == "value" || key.Value() == "keys" || key.Value() == "key" {
      errCtx := key.Context()
      return nil, errCtx.NewError("Error: forbidden name for enum member")
    }

    members[i] = &EnumMember{key, vs[i]}
  }

  en := &Enum{
		nameExpr,
		parentExpr,
    members,
		TokenData{ctx},
	}

  en.nameExpr.GetVariable().SetObject(en)

  return en, nil
}

func (t *Enum) Name() string {
	return t.nameExpr.Name()
}

func (t *Enum) GetPrototypes() ([]values.Prototype, error) {
  // doesn't need to return self
  return []values.Prototype{}, nil
}

func (t *Enum) GetParent() (values.Prototype, error) {
  if t.parentExpr == nil {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: enum needs to extend something")
  }

  proto := t.parentExpr.GetPrototype()

  if proto != nil {
    return proto, nil
  } else {
    errCtx := t.parentExpr.Context()
    return nil, errCtx.NewError("Error: not a prototype")
  }
}

func (t *Enum) GetInterfaces() ([]values.Interface, error) {
  return []values.Interface{}, nil
}

func (t *Enum) GetVariable() Variable {
	return t.nameExpr.GetVariable()
}

// can never be directly constructed
func (t *Enum) GetClassValue() (*values.Class, error) {
  return nil, nil
}

func (t *Enum) GetEnumValue() (*values.Enum, error) {
  return values.NewEnum(t, t.Context()), nil
}

func (t *Enum) AddStatement(st Statement) {
	panic("not available")
}

func (t *Enum) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Enum(")
	b.WriteString(t.nameExpr.Dump(""))
	b.WriteString(") extends ")
	b.WriteString(t.parentExpr.Dump(""))
	b.WriteString("\n")

	for _, member := range t.members {
		key := member.key
		b.WriteString(member.val.Dump(indent + key.Value() + ":"))
	}

	return b.String()
}

func (t *Enum) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	var b strings.Builder

	name := t.nameExpr.WriteExpression()
	b.WriteString(indent)
	b.WriteString("class ")
	b.WriteString(name)
	b.WriteString(" extends ")
	b.WriteString(t.parentExpr.Name())
	b.WriteString("{")

	b.WriteString(nl)
	b.WriteString(indent + tab)
	b.WriteString("static get values(){return Object.freeze([")
	for i, member := range t.members {
		b.WriteString(member.val.WriteExpression())

		if i < len(t.members)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("])}")

	b.WriteString(nl)
	b.WriteString(indent + tab)
	b.WriteString("static get keys(){return Object.freeze([")
	for i, member := range t.members {
		b.WriteString("'")
		b.WriteString(member.key.Value())
		b.WriteString("'")
		if i < len(t.members)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("])}")

	for i, member := range t.members {
		b.WriteString(nl)
		b.WriteString(indent + tab)
		b.WriteString("static get ")
		b.WriteString(member.key.Value())
		b.WriteString("(){return ")
		b.WriteString(name)
		b.WriteString(".values[")
		b.WriteString(strconv.Itoa(i))
		b.WriteString("]}")
	}

	// TODO: is the value runtime getter really necessary?
	b.WriteString(nl)
	b.WriteString(indent + tab)
	b.WriteString("static value(key){return ")
	b.WriteString(name)
	b.WriteString(".values[{")
	for i, member := range t.members {
		b.WriteString(member.key.Value())
		b.WriteString(":")
		b.WriteString(strconv.Itoa(i))

		if i < len(t.members)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("}[key]]}")

	b.WriteString(nl)
	b.WriteString(indent + tab)
	b.WriteString("constructor(){throw new Error(\"Cannot call constructor on enum\")}")

	// the name getter is probably not used in high performance code, so we can use the simplistic approach
	b.WriteString(nl)
	b.WriteString(indent + tab)
	b.WriteString("static key(v){for(var i=0;i<")
	b.WriteString(name)
	b.WriteString(".keys.length;i++){")
	b.WriteString("if(v==")
	b.WriteString(name)
	b.WriteString(".values[i]){return ")
	b.WriteString(name)
	b.WriteString(".keys[i]}}}")

	b.WriteString(nl)
	b.WriteString(indent)
	b.WriteString("}")

	return b.String()
}

func (t *Enum) HoistNames(scope Scope) error {
	return nil
}

func (t *Enum) ResolveStatementNames(scope Scope) error {
	if scope.HasVariable(t.Name()) {
		errCtx := t.Context()
		err := errCtx.NewError("Error: '" + t.Name() + "' already defined " +
			"(enum needs unique name)")
		other, _ := scope.GetVariable(t.Name())
		err.AppendContextString("Info: defined here ", other.Context())
		return err
	} else {
		if err := t.parentExpr.ResolveExpressionNames(scope); err != nil {
			return err
		}

		if err := scope.SetVariable(t.Name(), t.GetVariable()); err != nil {
			return err
		}

		// all member expressions end up in the same list, so they can share a scope
		subScope := NewSubScope(scope)
		for _, member := range t.members {
			if err := member.val.ResolveExpressionNames(subScope); err != nil {
				return err
			}
		}

		return nil
	}

}

func (t *Enum) evalInternal() error {
  parent, err := t.GetParent()
  if err != nil {
    return err
  }

  parentClassVal, err := parent.GetClassValue()
  if err != nil {
    return err
  }

  parentVal, err := parentClassVal.EvalConstructor(nil, t.Context())
  if err != nil {
    return err
  }
  // now evaluate each member, and check that they respect the parent

  for _, member := range t.members {
    mVal, err := member.Eval()
    if err != nil {
      return err
    }

    if err := parentVal.Check(mVal, member.Context()); err != nil {
      return err
    }
  }

  return nil
}

func (m *EnumMember) Context() context.Context {
  return m.key.Context()
}

func (m *EnumMember) Eval() (values.Value, error) {
  return m.val.EvalExpression()
}

// XXX: should enums be available as expressions
func (t *Enum) EvalStatement() error {
  if err := t.evalInternal(); err != nil {
    return err
  }

  variable := t.GetVariable()

  val, err := t.GetEnumValue()
  if err != nil {
    return err
  }

  variable.SetValue(val)

  return nil
}

func (t *Enum) ResolveStatementActivity(usage Usage) error {
  if err := t.parentExpr.ResolveExpressionActivity(usage); err != nil {
    return err
  }

	if usage.InFunction() {
		nameVar := t.nameExpr.GetVariable()

		if err := usage.Rereference(nameVar, t.Context()); err != nil {
			return err
		}
	}

	tmp := usage.InFunction()
	usage.SetInFunction(false)

	// in reverse order
	for i := len(t.members) - 1; i >= 0; i-- {
    member := t.members[i]
		if err := member.val.ResolveExpressionActivity(usage); err != nil {
			usage.SetInFunction(tmp)
			return err
		}
	}

	usage.SetInFunction(tmp)

	return nil
}

func (t *Enum) UniversalStatementNames(ns Namespace) error {
	if err := t.parentExpr.UniversalExpressionNames(ns); err != nil {
		return err
	}

	for _, member := range t.members {
		if err := member.val.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Enum) UniqueStatementNames(ns Namespace) error {
	// enums aren't actually instances of enum classes (they remain instances of String or Int), so universal classname isn't necessary
	if err := ns.ClassName(t.nameExpr.GetVariable()); err != nil {
		return err
	}

	if err := t.parentExpr.UniqueExpressionNames(ns); err != nil {
		return err
	}

	for _, member := range t.members {
		if err := member.val.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Enum) Walk(fn WalkFunc) error {
  if err := t.nameExpr.Walk(fn); err != nil {
    return err
  }

  if err := t.parentExpr.Walk(fn); err != nil {
    return err
  }

  for _, member := range t.members {
    if err := member.Walk(fn); err != nil {
      return err
    }
  }

  return fn(t)
}

func (m *EnumMember) Walk(fn WalkFunc) error {
  if err := m.key.Walk(fn); err != nil {
    return err
  }

  if err := m.val.Walk(fn); err != nil {
    return err
  }

  return fn(m)
}

func (t *Enum) Check(other_ values.Interface, ctx context.Context) error {
  // only exact match is possible
  if other, ok := other_.(*Enum); ok {
    if other == t {
      return nil
    } else {
      return ctx.NewError("Error: expected enum " + t.Name() + ", got enum " + other.Name())
    }
  } else {
    return ctx.NewError("Error: not an enum")
  }
}

func (t *Enum) IsUniversal() bool {
  parent, err := t.GetParent()
  if err != nil {
    panic("should've been caught before")
  }

	return parent.IsUniversal()
}

func (t *Enum) IsRPC() bool {
  return false
}

func (t *Enum) IsAbstract() bool {
  return false
}

func (t *Enum) IsFinal() bool {
  return true
}

func (t *Enum) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  parent, err := t.GetParent()
  if err != nil {
    return nil, err
  }

  return parent.GetInstanceMember(key, includePrivate, ctx)
}

func (t *Enum) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  parent, err := t.GetParent()
  if err != nil {
    return err
  }

  return parent.SetInstanceMember(key, includePrivate, arg, ctx)
}

func (t *Enum) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  switch key {
  case "values":
    return prototypes.NewArray(values.NewInstance(t, ctx), ctx), nil
  case "keys":
    return prototypes.NewArray(prototypes.NewString(ctx), ctx), nil
  case "value":
    return values.NewFunction([]values.Value{prototypes.NewString(ctx), values.NewInstance(t, ctx)}, ctx), nil
  case "key":
    return values.NewFunction([]values.Value{values.NewInstance(t, ctx), prototypes.NewString(ctx)}, ctx), nil
  default:
    for _, member := range t.members {
      if member.key.Value() == key {
        return values.NewInstance(t, ctx), nil
      }
    }

    parent, err := t.GetParent()
    if err != nil {
      return nil, err
    }

    return parent.GetClassMember(key, includePrivate, ctx)
  }
}
