package js

import (
	"strconv"
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type EnumMember struct {
	key   *Word
	val Expression
}

// enum is a statement (and also a Class!)
type Enum struct {
	clType        *TypeExpression
	extends       *TypeExpression // Int, String or something else, never nil
	cachedExtends values.Prototype
  members       []*EnumMember
	TokenData
}

func NewEnum(clType *TypeExpression, extends *TypeExpression, keys []*Word,
	vs []Expression, ctx context.Context) (*Enum, error) {
  members := make([]*EnumMember, len(keys))
  for i, key := range keys {
    members[i] = &EnumMember{key, vs[i]}
  }

	return &Enum{
		clType,
		extends,
		nil, // evaluated later
    members,
		TokenData{ctx},
	}, nil
}

func (t *Enum) Name() string {
	return t.clType.Name()
}

func (t *Enum) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return prototypes.CheckPrototype(t, args, pos, ctx)
}

func (t *Enum) GetParent() values.Prototype {
	return t.cachedExtends
}

func (t *Enum) GetVariable() Variable {
	return t.clType.GetVariable()
}

func (t *Enum) AddStatement(st Statement) {
	panic("not available")
}

func (t *Enum) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Enum(")
	b.WriteString(t.clType.Dump(""))
	b.WriteString(") extends ")
	b.WriteString(t.extends.Dump(""))
	b.WriteString("\n")

	for _, member := range t.members {
		key := member.key
		b.WriteString(member.val.Dump(indent + key.Value() + ":"))
	}

	return b.String()
}

func (t *Enum) WriteStatement(indent string) string {
	var b strings.Builder

	name := t.clType.WriteExpression()
	b.WriteString(indent)
	b.WriteString("class ")
	b.WriteString(name)
	b.WriteString(" extends ")
	b.WriteString(t.extends.Name())
	b.WriteString("{")

	b.WriteString(NL)
	b.WriteString(indent + TAB)
	b.WriteString("static get values(){return Object.freeze([")
	for i, member := range t.members {
		b.WriteString(member.val.WriteExpression())

		if i < len(t.members)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("])}")

	b.WriteString(NL)
	b.WriteString(indent + TAB)
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
		b.WriteString(NL)
		b.WriteString(indent + TAB)
		b.WriteString("static get ")
		b.WriteString(member.key.Value())
		b.WriteString("(){return ")
		b.WriteString(name)
		b.WriteString(".values[")
		b.WriteString(strconv.Itoa(i))
		b.WriteString("]}")
	}

	// TODO: is the value runtime getter really necessary?
	b.WriteString(NL)
	b.WriteString(indent + TAB)
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

	b.WriteString(NL)
	b.WriteString(indent + TAB)
	b.WriteString("constructor(){throw new Error(\"Cannot call constructor on enum\")}")

	// the name getter is probably not used in high performance code, so we can use the simplistic approach
	b.WriteString(NL)
	b.WriteString(indent + TAB)
	b.WriteString("static name(v){for(var i=0;i<")
	b.WriteString(name)
	b.WriteString(".keys.length;i++){")
	b.WriteString("if(v==")
	b.WriteString(name)
	b.WriteString(".values[i]){return ")
	b.WriteString(name)
	b.WriteString(".keys[i]}}}")

	b.WriteString(NL)
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
		if err := t.extends.ResolveExpressionNames(scope); err != nil {
			return err
		}

		var err error
		t.cachedExtends, err = getExtendsClass(t.extends.GetVariable(), t.extends.Context())
		if err != nil {
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

/*func (t *Enum) cacheExtends(stack values.Stack) error {
	var err error
	t.cachedExtends, err = cacheClassExtends(stack, t.extends, t.Context())
	if err != nil {
		return err
	}

	return nil
}*/

func (t *Enum) HoistValues(stack values.Stack) error {
	return nil
}

func (t *Enum) EvalStatement(stack values.Stack) error {
	//if err := t.cacheExtends(stack); err != nil {
	//return err
	//}

	if _, err := t.GenerateInstance(stack, nil, nil, t.Context()); err != nil {
		return err
	}

	// class cannot be used for inheriting
	val := values.NewClass(t, t.Context())

	return stack.SetValue(t.GetVariable(), val, false, t.Context())
}

func (t *Enum) ResolveStatementActivity(usage Usage) error {
	if parent, ok := t.cachedExtends.(*Class); ok {
		if err := parent.ResolveStatementActivity(usage); err != nil {
			return err
		}
	}

	if usage.InFunction() {
		clVar := t.clType.GetVariable()

		if err := usage.Rereference(clVar, t.Context()); err != nil {
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
	if err := t.extends.UniversalExpressionNames(ns); err != nil {
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
	if err := ns.ClassName(t.clType.GetVariable()); err != nil {
		return err
	}

	if err := t.extends.UniqueExpressionNames(ns); err != nil {
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
  if err := t.clType.Walk(fn); err != nil {
    return err
  }

  if err := t.extends.Walk(fn); err != nil {
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

func (t *Enum) IsImplementedBy(proto values.Prototype) (string, bool) {
	return "", false
}

func (t *Enum) HasAncestor(interf values.Interface) bool {
	if _, ok := interf.IsImplementedBy(t); ok {
		return true
	} else if _, ok := interf.(*values.AllPrototype); ok {
		return true
	} else if t == interf {
		return true
	} else {
		return t.cachedExtends.HasAncestor(interf)
	}
}

func (t *Enum) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	if typeChildren != nil {
		return nil, ctx.NewError("Error: enum can't have content types")
	}

	newV, ok := v.ChangeInstanceInterface(t, false, true)
	if !ok {
		return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + t.Name())
	}

	return newV, nil
}

func (t *Enum) EvalConstructor(stack values.Stack, args []values.Value,
	childProto values.Prototype, ctx context.Context) (values.Value, error) {
	return nil, ctx.NewError("Error: enum cannot be constructed")
}

func (t *Enum) generateValue(stack values.Stack, i int) (values.Value, error) {
	expr := t.members[i].val
	v, err := expr.EvalExpression(stack)
	if err != nil {
		return nil, err
	}

	if _, ok := v.ChangeInstancePrototype(t, true); !ok {
		errCtx := expr.Context()
		return nil, errCtx.NewError("Error: not an instance")
	}

	return v, nil
}

func (t *Enum) generateValues(stack values.Stack, ctx context.Context) ([]values.Value, error) {
	vs := make([]values.Value, 0)

	for i, _ := range t.members {
		v, err := t.generateValue(stack, i)
		if err != nil {
			return nil, err
		}

		vs = append(vs, v)
	}

	return vs, nil
}

func (t *Enum) generateKeys(stack values.Stack, ctx context.Context) ([]values.Value, error) {
	ks := make([]values.Value, 0)
	for _, member := range t.members {
    key := member.key
		kStr := key.Value()

		k := prototypes.NewLiteralString(kStr, ctx)

		ks = append(ks, k)
	}

	return ks, nil
}

func (t *Enum) GenerateInstance(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
	if keys != nil || args != nil {
		return nil, ctx.NewError("Error: parametric enums not yet supported")
	}

	vs, err := t.generateValues(stack, ctx)
	if err != nil {
		return nil, err
	}

	return values.NewMulti(vs, ctx), nil
}

func (t *Enum) EvalAsEntryPoint(stack values.Stack, ctx context.Context) error {
	_, err := t.GenerateInstance(stack, nil, nil, ctx)

	return err
}

func (t *Enum) IsUniversal() bool {
	//if !patterns.JS_UNIVERSAL_CLASS_NAME_REGEXP.MatchString(t.Name()) {
	//return ctx.NewError("Error: " + t.Name() + " is an invalid name for a universal enum")
	//}

	//t.isUniversal = true
	return t.cachedExtends.IsUniversal()
}

func (t *Enum) HasMember(this *values.Instance, key string, includePrivate bool) bool {
	// key can be any of the members, values, value, keys, name
	switch key {
	case "values", "value", "keys", "name":
		return true
	default:
		for _, member := range t.members {
			if member.key.Value() == key {
				return true
			}
		}

		return t.cachedExtends.HasMember(this, key, includePrivate)
	}
}

func (t *Enum) GetMember(stack values.Stack, this *values.Instance, key string, includePrivate bool, ctx context.Context) (values.Value, error) {
	if this == nil {
		switch key {
		case "values":
			values, err := t.generateValues(stack, ctx)
			if err != nil {
				return nil, err
			}
			return prototypes.NewLiteralArray(values, ctx), nil
		case "keys":
			keys, err := t.generateKeys(stack, ctx)
			if err != nil {
				return nil, err
			}
			return prototypes.NewLiteralArray(keys, ctx), nil
		case "value":
			return values.NewFunctionFunction(func(stack_ values.Stack,
				this_ *values.Instance, args_ []values.Value,
				ctx_ context.Context) (values.Value, error) {
				if err := prototypes.CheckInputs(prototypes.String, args_, ctx_); err != nil {
					return nil, err
				}
				// TODO: wrap in own prototype
				return t.GenerateInstance(stack, nil, nil, ctx_)
			}, stack, this, ctx), nil
		case "name":
			return values.NewFunctionFunction(func(stack_ values.Stack,
				this_ *values.Instance, args_ []values.Value,
				ctx_ context.Context) (values.Value, error) {
				if err := prototypes.CheckInputs(t, args_, ctx_); err != nil {
					return nil, err
				}

				return prototypes.NewString(ctx_), nil
			}, stack, this, ctx), nil
		default:
			for i, member := range t.members {
				if member.key.Value() == key {
					return t.generateValue(stack, i)
				}
			}

			//return nil, ctx.NewError("Error: enum doesnt have static method " + key)
			return t.cachedExtends.GetMember(stack, nil, key, includePrivate, ctx)
		}
	} else {
		return t.cachedExtends.GetMember(stack, this, key, includePrivate, ctx)
	}
}

func (t *Enum) SetMember(stack values.Stack, this *values.Instance, key string, arg values.Value, includePrivate bool, ctx context.Context) error {
	return t.cachedExtends.SetMember(stack, this, key, arg, includePrivate, ctx)
}

func (t *Enum) GetIndex(stack values.Stack, this *values.Instance, index values.Value,
	ctx context.Context) (values.Value, error) {
	return t.cachedExtends.GetIndex(stack, this, index, ctx)
}

func (t *Enum) SetIndex(stack values.Stack, this *values.Instance, index values.Value,
	arg values.Value, ctx context.Context) error {
	return t.cachedExtends.SetIndex(stack, this, index, arg, ctx)
}

func (t *Enum) LoopForIn(this *values.Instance, fn func(values.Value) error, ctx context.Context) error {
	return t.cachedExtends.LoopForIn(this, fn, ctx)
}

func (t *Enum) LoopForOf(this *values.Instance, fn func(values.Value) error, ctx context.Context) error {
	return t.cachedExtends.LoopForOf(this, fn, ctx)
}
