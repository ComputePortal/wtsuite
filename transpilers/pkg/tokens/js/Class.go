package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// only support single inheritance (easier to maintain code, and similar to java)
type Class struct {
	nameExpr         *TypeExpression
	parentExpr       *TypeExpression  // can be nil
	interfExprs      []*VarExpression // can't be nil, can be zero length, can't contain nil
  constructor      *Function
	members          []ClassMember    // list because getter and setter can have same name
	universalName    string
	TokenData
}

func NewClass(nameExpr *TypeExpression, parentExpr *TypeExpression, ctx context.Context) (*Class, error) {
	cl := &Class{
		nameExpr,
		parentExpr,
    make([]*VarExpression, 0),
		nil, // set later
		make([]ClassMember, 0),
		"",
		TokenData{ctx},
	}

	cl.nameExpr.GetVariable().SetObject(cl)

	return cl, nil
}

func NewUniversalClass(nameExpr *TypeExpression, parentExpr *TypeExpression, interfExprs []*VarExpression, universalName string, ctx context.Context) (*Class, error) {
	cl, err := NewClass(nameExpr, parentExpr, ctx)
	if err != nil {
		panic(err)
	}

  for _, interfExpr := range interfExprs {
    if interfExpr == nil {
      panic("interfExpr can't be nil")
    }
  }

	cl.interfExprs = interfExprs
	cl.universalName = universalName

	return cl, nil
}

func (t *Class) newScope(parent Scope) Scope {
	return NewClassScope(t, parent)
}

func (t *Class) Name() string {
	return t.nameExpr.Name()
}

func (t *Class) IsUniversal() bool {
	return t.universalName != ""
}

func (t *Class) GetVariable() Variable {
	return t.nameExpr.GetVariable()
}

func (t *Class) GetPrototypes() ([]values.Prototype, error) {
  // doesn't need to include itself
  return []values.Prototype{}, nil
}

// available after resolve names stage
// returns nil if not found
func (t *Class) GetParent() (values.Prototype, error) {
  if t.parentExpr == nil {
    return nil, nil
  }

  val, err := t.parentExpr.EvalExpression()
  if err != nil {
    return nil, err
  }

  proto := values.GetPrototype(val)

  if proto != nil {
    return proto, nil
  } else {
    errCtx := t.parentExpr.Context()
    return nil, errCtx.NewError("Error: not a prototype")
  }
}

func (t *Class) GetInterfaces() ([]values.Interface, error) {
  interfs := make([]values.Interface, 0)

  for _, interfExpr := range t.interfExprs {
    interf := interfExpr.GetInterface()
    if interf != nil {
      interfs = append(interfs, interf)
    } else {
      errCtx := interfExpr.Context()
      return nil, errCtx.NewError("Error: not an interface")
    }
  }

  return interfs, nil
}

func (t *Class) getMember(name string, preferSetter bool) ClassMember {
  if preferSetter {
    for _, member := range t.members {
      if member.Name() == name {
        if prototypes.IsSetter(member) {
          return member
        }
      }
    }
  } else {
    for _, member := range t.members {
      if member.Name() == name {
        if !prototypes.IsSetter(member) {
          return member
        }
      }
    }
  }

  for _, member := range t.members {
    if member.Name() == name {
      return member
    }
  }

  return nil
}

func (t *Class) AddProperty(name *Word, expr *TypeExpression) error {
  if prev := t.getMember(name.Value(), false); prev != nil {
    errCtx := name.Context()
    err := errCtx.NewError("Error: already have a member named " + name.Value())
    err.AppendContextString("Info: previously declared here", prev.Context())
    return err
  }

  t.members = append(t.members, NewClassProperty(name, expr))

  return nil
}

func (t *Class) Properties() (map[string]values.Value, error) {
  ctx := t.Context()
  props := make(map[string]values.Value)
  for _, member := range t.members {
    if prop, ok := member.(*ClassProperty); ok {
      v, err := prop.GetValue(ctx)
      if err != nil {
        return nil, err
      }

      props[member.Name()] = v
    }
  }

  return props, nil
}

func (t *Class) AddConstructor(fn *Function) error {
  if t.constructor != nil {
    errCtx := fn.Context()
    return errCtx.NewError("Error: constructor already defined")
  }

  if fn.Role() != prototypes.NORMAL {
    errCtx := fn.Context()
    return errCtx.NewError("Error: constructor can't have any modifiers")
  }

  t.constructor = fn

  return nil
}

func (t *Class) AddFunction(fn *Function) error {
	if fn.Name() == "" {
		errCtx := fn.Context()
		err := errCtx.NewError("Error: member function doesn't have a name")
		return err
	} else if fn.Name() == "constructor" {
    return t.AddConstructor(fn)
  }

  if prototypes.IsGetter(fn) {
    if prev := t.getMember(fn.Name(), false); prev != nil && !prototypes.IsSetter(prev) {
      errCtx := fn.Context()
      err := errCtx.NewError("Error: already have a member named " + fn.Name())
      err.AppendContextString("Info: previously declared here", prev.Context())
      return err
    }
  } else if prototypes.IsSetter(fn) {
    if prev := t.getMember(fn.Name(), false); prev != nil && !prototypes.IsGetter(prev) {
      errCtx := fn.Context()
      err := errCtx.NewError("Error: already have a member named " + fn.Name())
      err.AppendContextString("Info: previously declared here", prev.Context())
      return err
    }
  } else {
    if prev := t.getMember(fn.Name(), false); prev != nil {
      errCtx := fn.Context()
      err := errCtx.NewError("Error: already have a member named " + fn.Name())
      err.AppendContextString("Info: previously declared here", prev.Context())
      return err
    }
  }

	// check that there is only one constructor
	if fn.Name() == "constructor" {
    panic("use AddConstructor for constructor")
	}

	member := NewClassFunction(fn)
	errCtx := fn.Context()

	switch {
	case prototypes.IsStatic(member) && prototypes.IsConst(member):
		return errCtx.NewError("Error: a static class function can't be const " +
			"(hint: const applies to instances)")
	case prototypes.IsPublic(member) && prototypes.IsPrivate(member):
		return errCtx.NewError("Error: a class member can't be public and " +
			"private at the same time")
	case prototypes.IsOverride(member) && prototypes.IsAbstract(member):
		return errCtx.NewError("Error: a class member can't be abstract and " +
			"override at the same time")
	case prototypes.IsAsync(member) && (prototypes.IsGetter(member) || prototypes.IsSetter(member)):
		return errCtx.NewError("Error: a class member can't be async and " +
			"getter/setter at the same time")
	}

	t.members = append(t.members, member)

	return nil
}

func (t *Class) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Class(")
	b.WriteString(t.nameExpr.Dump(""))
	b.WriteString(")")

	if t.parentExpr != nil {
		b.WriteString(" extends ")

		b.WriteString(t.parentExpr.Dump(""))
	}

	b.WriteString("\n")

  if t.constructor != nil {
		b.WriteString(t.constructor.Dump(indent + "  "))
  }

	b.WriteString(indent)
	for _, member := range t.members {
		b.WriteString(member.Dump(indent + "  "))
	}

	return b.String()
}

func (t *Class) WriteExpression() string {
	return t.WriteStatement("")
}

func (t *Class) WriteStatement(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("class ")

	// class expressions can be anonymous
	if t.nameExpr != nil {
		b.WriteString(t.nameExpr.WriteExpression())
	}

	if t.parentExpr != nil {
		b.WriteString(" extends ")
		b.WriteString(t.parentExpr.WriteExpression())
	}

	b.WriteString("{")

	hasContent := false
  if t.constructor != nil {
    b.WriteString(NL)
    b.WriteString(indent + TAB)
    b.WriteString("constructor")
    b.WriteString(t.constructor.writeBody(indent + TAB, NL, TAB))
    hasContent = true
  }

	for _, member := range t.members {
		s := member.WriteStatement(indent + TAB)

		if s != "" {
			b.WriteString(NL)
			b.WriteString(s)
			hasContent = true
		}
	}

	//b.WriteString(t.writeNonOverriddenGettersAndSetters(indent + TAB))

	if hasContent {
		b.WriteString(NL)
		b.WriteString(indent)
	}
	b.WriteString("}")

	return b.String()
}

func (t *Class) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Class) HoistNames(scope Scope) error {
	return nil
}

// reused by enum
/*func getExtendsClass(v Variable, ctx context.Context) (values.Prototype, error) {
	extendsClass_ := v.GetObject()
	if extendsClass_ == nil {
		return nil, ctx.NewError("Error: bad extends (object not set)")
	}

	switch extendsClass := extendsClass_.(type) {
	case *Class:
		return extendsClass, nil
	case values.Prototype:
		return extendsClass, nil
	case *ClassInterface:
		return nil, ctx.NewError("Error: can't extend interface")
	case *Enum:
		return nil, ctx.NewError("Error: cannot inherit from enum")

	default:
		return nil, ctx.NewError("Error: bad extends (" + reflect.TypeOf(extendsClass_).String() + ")")
	}
}*/

func (t *Class) resolveConstructorNames(scope Scope) error {
  // resolve the constructor
  if t.constructor != nil {
    subScope := t.newScope(scope)
    if t.parentExpr != nil {
      parent, err := t.GetParent()
      if err != nil {
        return err
      }

      super := NewVariable("super", true, t.parentExpr.Context())
      parentClassVal, err := parent.GetClassValue()
      if err != nil {
        return err
      }

      if parentClassVal == nil {
        errCtx := t.parentExpr.Context()
        return errCtx.NewError("Error: parent constructor not found")
      }

      superVal, err := values.NewSuper(parentClassVal, t.parentExpr.Context())
      if err != nil {
        return err
      }

      super.SetValue(superVal)

      subScope.SetVariable(super.Name(), super)
    }

    thisVar := t.constructor.GetThisVariable()
    thisVal := values.NewInstance(t, t.Context())
    thisValWrapper := values.NewThis(thisVal, t.Context())
    thisVar.SetValue(thisValWrapper)

    if err := t.constructor.ResolveExpressionNames(subScope); err != nil {
      return err
    }

  }

  return nil
}

func (t *Class) resolveMemberNames(scope Scope) error {
	for _, member_ := range t.members {
    // type assertion is easiest here
    switch member := member_.(type) {
    case *ClassFunction:
      subScope := t.newScope(scope)

      // super is available in every function, not just the constructor
      if t.parentExpr != nil {
        parent, err := t.GetParent()
        if err != nil {
          return err
        }

        if parent != nil {
          superVal := values.NewInstance(parent, member.Context())
          superVar := NewVariable("super", true, t.parentExpr.Context())
          superVar.SetValue(superVal)

          if err := subScope.SetVariable("super", superVar); err != nil {
            return err
          }
        }
      }

      if !prototypes.IsStatic(member) {
        thisVar := member.GetThisVariable()
        thisVal := values.NewInstance(t, member.Context())
        thisVar.SetValue(thisVal)
      }

      if err := member.ResolveNames(subScope); err != nil {
        return err
      }
    case *ClassProperty:
      if err := member.ResolveNames(scope); err != nil {
        return err
      }
    default:
      panic("not implemented")
    }
	}

	return nil
}

func (t *Class) checkUniversalness() error {
  ctx := t.Context()

  if t.universalName == "" {
    return nil
  }

	for _, member := range t.members {
    if !member.IsUniversal() {
      errCtx := ctx
      err := errCtx.NewError("Error: not universal")

      err.AppendContextString("Info: member not universal", member.Context())
    }
  }

  return nil
}

func (t *Class) ResolveExpressionNames(scope Scope) error {
	if t.parentExpr != nil {
		if err := t.parentExpr.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

  for _, interfExpr := range t.interfExprs {
    if err := interfExpr.ResolveExpressionNames(scope); err != nil {
      return err
    }
  }

  if err := t.resolveConstructorNames(scope); err != nil {
    return err
  }

  if err := t.resolveMemberNames(scope); err != nil {
    return err
  }

	return nil
}

func (t *Class) ResolveStatementNames(scope Scope) error {
	if scope.HasVariable(t.Name()) {
		errCtx := t.Context()
		err := errCtx.NewError("Error: '" + t.Name() + "' already defined " +
			"(class needs unique name)")
		other, _ := scope.GetVariable(t.Name())
		err.AppendContextString("Info: defined here ", other.Context())
		return err
	} else {
		if err := scope.SetVariable(t.Name(), t.GetVariable()); err != nil {
			return err
		}

		return t.ResolveExpressionNames(scope)
	}
}

func (t *Class) GetClassValue() (*values.Class, error) {
  ctx := t.Context()

  if t.constructor == nil {
    // use parent constructor
    if t.parentExpr == nil {
      return nil, nil
    } else {
      parent, err := t.GetParent()
      if err != nil {
        return nil, err
      }

      cv, err := parent.GetClassValue()
      if err != nil {
        return nil, err
      }

      cvArgs := cv.GetConstructorArgs()

      return values.NewClass(cvArgs, t, ctx), nil
    }
  } else {
    fnVal, err := t.constructor.GetFunctionValue()
    if err != nil {
      return nil, err
    }

    if !fnVal.IsVoid() {
      errCtx := t.constructor.Context()
      return nil, errCtx.NewError("Error: constructor must return void")
    }

    args := fnVal.GetArgs()
    return values.NewClass(args, t, ctx), nil
  }
}

func (t *Class) evalInternal() error {
  if t.constructor != nil {
    // the returned value is the function itself
    if _, err := t.constructor.EvalExpression(); err != nil {
      return err
    } 

    thisVar := t.constructor.GetThisVariable()
    thisVal, _ := thisVar.GetValue().(*values.This)

    // check that all properties are "touched"
    for _, member := range t.members {
      if prototypes.IsProperty(member) {
        if err := thisVal.AssertTouched(member.Name(), member.Context()); err != nil {
          return err
        }
      }
    }
  }

  for _, member := range t.members {
    if err := member.Eval(); err != nil {
      return err
    }
  }

  for _, interfExpr := range t.interfExprs {
    // interfExpr.EvalExpression would give an error, because no value is set
    interf := interfExpr.GetInterface()
    isPrototype := interfExpr.GetPrototype() != nil
    if isPrototype {
      errCtx := interfExpr.Context()
      return errCtx.NewError("Error: can't implement other class")
    }

    if interf == nil {
      errCtx := interfExpr.Context()
      return errCtx.NewError("Error: not an interface")
    }

    // implementation registration is done inside the check
    if err := interf.Check(t, interfExpr.Context()); err != nil {
      return err
    }

  }

  return nil
}

func (t *Class) EvalExpression() (values.Value, error) {
  if err := t.evalInternal(); err != nil {
    return nil, err
  }

	return t.GetClassValue()
}

func (t *Class) EvalStatement() error {
  variable := t.GetVariable()
  classVal, err := t.GetClassValue()
  if err != nil {
    return err
  }

  variable.SetValue(classVal)

  if err := t.evalInternal(); err != nil {
    return err
  }

	return nil
}

// potentially return Setter/Getter as second ClassMember
// return non-nil prototypes.BuiltinPrototype for further specific searching
func (t *Class) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  setterFound := false

	for _, member := range t.members {
    if member.Name() == key {
      if prototypes.IsSetter(member) {
        setterFound = true
      } else if prototypes.IsPublic(member) || includePrivate {
        return member.GetValue(ctx)
			} else {
        return nil, ctx.NewError("Error: member is private")
      }
		}
	}

  // try parent
  if setterFound {
    return nil, ctx.NewError("Error: is only a setter")
  } else if t.parentExpr != nil {
    parent, err := t.GetParent()
    if err != nil {
      return nil, err
    }

    return parent.GetInstanceMember(key, includePrivate, ctx)
  } else {
    return nil, nil
  }
}

func (t *Class) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  for _, member := range t.members {
    if member.Name() == key {
      if prototypes.IsSetter(member) || prototypes.IsProperty(member) {
        if prototypes.IsPublic(member) || includePrivate {
          return member.SetValue(arg, ctx)
        }
      }
    }
  }

  if t.parentExpr == nil {
    return ctx.NewError("Error: setable member not found")
  } else {
    parent, err := t.GetParent()
    if err != nil {
      return err
    }

    return parent.SetInstanceMember(key, includePrivate, arg, ctx)
  }
}

func (t *Class) GetClassMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  nonStaticFound := false
  for _, member := range t.members {
    if member.Name() == key {
      if prototypes.IsStatic(member) {
        if prototypes.IsPublic(member) || includePrivate {
          return member.GetValue(ctx)
        } else {
          return nil, ctx.NewError("Error: member is private")
        }
      } else {
        nonStaticFound = true
      }
    }
  }

  if nonStaticFound {
    return nil, ctx.NewError("Error: member isn't static")
  } else if t.parentExpr != nil {
    parent, err := t.GetParent()
    if err != nil {
      return nil, err
    }

    return parent.GetClassMember(key, includePrivate, ctx)
  } else {
    return nil, nil
  }
}

/*func (t *Class) isAbstract(implemented map[string]*ClassMember) bool {
	for _, member := range t.members {
		if _, ok := implemented[member.Name()]; !ok {
			if prototypes.IsAbstract(member) {
				return true
			}

			implemented[member.Name()] = member
		}
	}

	if t.cachedExtends != nil {
		if parent, ok := t.cachedExtends.(*Class); ok {
			return parent.isAbstract(implemented)
		}
	}

	return false
}

func (t *Class) IsAbstract() bool {
	return t.isAbstract(make(map[string]*ClassMember))
}*/

func (t *Class) Check(other_ values.Interface, ctx context.Context) error {
  other, ok := other_.(values.Prototype) 
  if !ok {
    return ctx.NewError("Error: not a class instance")
  }

  // keep getting the parent of the other until they match us
  for other != nil {
    if otherClass, ok := other.(*Class); ok {
      if otherClass == t {
        return nil
      }
    }

    var err error
    other, err = other.GetParent()
    if err != nil {
      return err
    }
  }

  return ctx.NewError("Error: " + other_.Name() + " doesn't inherit from " + t.Name())
}

func (t *Class) ResolveExpressionActivity(usage Usage) error {
	tmp := usage.InFunction()
	usage.SetInFunction(false)

	for _, member := range t.members {
		if err := member.ResolveActivity(usage); err != nil {
			usage.SetInFunction(tmp)
			return err
		}
	}

	usage.SetInFunction(tmp)

	return nil
}

func (t *Class) ResolveStatementActivity(usage Usage) error {
  parent_, err := t.GetParent()
  if err != nil {
    return err
  }

  if parent, ok := parent_.(*Class); ok {
		if err := parent.ResolveStatementActivity(usage); err != nil {
			return err
		}
	}

	if usage.InFunction() {
		variable := t.GetVariable()

		if err := usage.Rereference(variable, t.Context()); err != nil {
			return err
		}
	}

	return t.ResolveExpressionActivity(usage)
}

func (t *Class) UniversalExpressionNames(ns Namespace) error {
	if t.parentExpr != nil {
		if err := t.parentExpr.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

  if err := t.checkUniversalness(); err != nil {
    return err
  }

	if t.IsUniversal() {
		if err := ns.UniversalName(t.nameExpr.GetVariable(), t.universalName); err != nil {
			return err
		}
	}

	for _, member := range t.members {
		if err := member.UniversalNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Class) UniqueExpressionNames(ns Namespace) error {
	if t.parentExpr != nil {
		if err := t.parentExpr.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	for _, member := range t.members {
		if err := member.UniqueNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *Class) UniversalStatementNames(ns Namespace) error {
	return t.UniversalExpressionNames(ns)
}

func (t *Class) UniqueStatementNames(ns Namespace) error {
	if !t.IsUniversal() {
		if err := ns.ClassName(t.nameExpr.GetVariable()); err != nil {
			return err
		}
	}

	return t.UniqueExpressionNames(ns)
}

func (t *Class) Walk(fn WalkFunc) error {
  if err := t.nameExpr.Walk(fn); err != nil {
    return err
  }

  if t.parentExpr != nil {
    if err := t.parentExpr.Walk(fn); err != nil {
      return err
    }
  }

  for _, interfExpr := range t.interfExprs {
    if err := interfExpr.Walk(fn); err != nil {
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

// TODO: overriding a setter must also override a getter, this just print a warning if this case is detected
/*func (t *Class) writeNonOverriddenGettersAndSetters(indent string) string {
	var b strings.Builder

	classGetters, classSetters := t.collectNonOverriddenGettersAndSetters()

	for _, classSetter := range classSetters {
		class := classSetter[0]
		setter := classSetter[1]

		b.WriteString(NL)
		b.WriteString(indent)
		b.WriteString("set ")
		b.WriteString(setter)
		b.WriteString("(x){Object.getOwnPropertyDescriptor(")
		b.WriteString(class)
		b.WriteString(".prototype,\"")
		b.WriteString(setter)
		b.WriteString("\").set.call(this,x);}")
	}

	for _, classGetter := range classGetters {
		class := classGetter[0]
		getter := classGetter[1]

		b.WriteString(NL)
		b.WriteString(indent)
		b.WriteString("get ")
		b.WriteString(getter)
		b.WriteString("(){Object.getOwnPropertyDescriptor(")
		b.WriteString(class)
		b.WriteString(".prototype,\"")
		b.WriteString(getter)
		b.WriteString("\").get.call(this);}")
	}

	return b.String()
}

func (t *Class) collectNonOverriddenGettersAndSetters() ([][2]string, [][2]string) {
	classGetters := make([][2]string, 0)
	classSetters := make([][2]string, 0)

	for _, member := range t.members {
		if prototypes.IsGetter(member) {
			isAlone := true
			for _, check := range t.members {
				if prototypes.IsSetter(check) && check.Name() == member.Name() {
					isAlone = false
					break
				}
			}

			if !isAlone {
				continue
			}

			// look for a setter
			extends := t
			for {
				if extends.cachedExtends == nil {
					break
				}

				ok := false
				extends, ok = extends.cachedExtends.(*Class)
				if !ok {
					break
				}

				if setter := extends.FindMember(member.Name(), true, true); setter != nil && prototypes.IsSetter(setter) {
					classSetters = append(classSetters, [2]string{extends.clType.WriteExpression(),
						setter.Name()})
					break
				}
			}
		} else if prototypes.IsSetter(member) {
			isAlone := true
			for _, check := range t.members {
				if prototypes.IsGetter(check) && check.Name() == member.Name() {
					isAlone = false
					break
				}
			}

			if !isAlone {
				continue
			}

			extends := t
			for {
				if extends.cachedExtends == nil {
					break
				}

				var ok bool
				extends, ok = extends.cachedExtends.(*Class)
				if !ok {
					break
				}

				if getter := extends.FindMember(member.Name(), true, false); getter != nil && prototypes.IsGetter(getter) {
					classGetters = append(classGetters, [2]string{extends.clType.WriteExpression(),
						getter.Name()})
					break
				}
			}
		}
	}

	return classGetters, classSetters
}*/
