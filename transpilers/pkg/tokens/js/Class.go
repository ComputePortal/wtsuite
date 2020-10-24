package js

import (
	"fmt"
	"reflect"
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// only support single inheritance (easier to maintain code)
// implements values.Prototype, Statement and Expression interfaces
type Class struct {
	clType           *TypeExpression
	extends          *TypeExpression // can be nil
	implements       *VarExpression  // can be nil
	super            Variable        // will have constructor value of parent
	cachedExtends    values.Prototype // resolved during names stage
	cachedImplements values.Interface // might be useful in future
	members          []*ClassMember
	childProto       values.Prototype // used as temporary value during construction (should be nil otherwise!)
	universalName    string
	TokenData
}

func NewClass(clType *TypeExpression, extends *TypeExpression, ctx context.Context) (*Class, error) {
	cl := &Class{
		clType,
		extends,
		nil,
		NewVariable("super", true, ctx),
		nil,
		nil,
		make([]*ClassMember, 0),
		nil,
		"",
		TokenData{ctx},
	}

	cl.clType.ref.SetObject(cl)

	return cl, nil
}

func NewUniversalClass(clType *TypeExpression, extends *TypeExpression, implements *VarExpression, universalName string, ctx context.Context) (*Class, error) {
	cl, err := NewClass(clType, extends, ctx)
	if err != nil {
		panic(err)
	}

	cl.implements = implements
	cl.universalName = universalName
	return cl, nil
}

func (t *Class) newScope(parent Scope) Scope {
	return NewClassScope(t, parent)
}

func (t *Class) Name() string {
	return t.clType.Name()
}

func (t *Class) GetParent() values.Prototype {
	return t.cachedExtends
}

func (t *Class) GetInterface() values.Interface {
  return t.cachedImplements
}

func (t *Class) IsUniversal() bool {
	return t.universalName != ""
}

/*func (t *Class) ToUniverse(ctx context.Context) error {
	if !patterns.JS_UNIVERSAL_CLASS_NAME_REGEXP.MatchString(t.Name()) {
		return ctx.NewError("Error: " + t.Name() + " is an invalid name for a universal class")
	}

	t.isUniversal = true
	return nil
}*/

func (t *Class) GetVariable() Variable {
	return t.clType.GetVariable()
}

func (t *Class) AddMember(fn *Function) error {
	if t.cachedExtends != nil {
		panic("cannot add member functions to class after ancestors have been collected")
	}

	if fn.Name() == "" {
		errCtx := fn.Context()
		err := errCtx.NewError("Error: member function doesn't have a name")
		panic(err)
		return err
	}

	// check that there is only one constructor
	if fn.Name() == "constructor" {
		for _, m := range t.members {
			if m.Name() == "constructor" {
				errCtx := fn.Context()
				err := errCtx.NewError("Error: there can only be one constructor per class")
				err.AppendContextString("Info: prev constructor defined here", m.function.Context())
				return err
			}
		}

		if fn.Role() != prototypes.NORMAL {
			errCtx := fn.Context()
			return errCtx.NewError("Error: constructor can't have any modifiers")
		}
	}

	member := NewClassMember(fn, t)
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
	b.WriteString(t.clType.Dump(""))
	b.WriteString(")")

	if t.extends != nil {
		b.WriteString(" extends ")

		b.WriteString(t.extends.Dump(""))
	}

	b.WriteString("\n")

	for _, member := range t.members {
		b.WriteString(member.Dump(indent))
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
	if t.clType != nil {
		b.WriteString(t.clType.WriteExpression())
	}

	if t.extends != nil {
		b.WriteString(" extends ")
		b.WriteString(t.extends.WriteExpression())
	}

	b.WriteString("{")

	hasContent := false
	for _, member := range t.members {
		s := member.WriteStatement(indent + TAB)

		if s != "" {
			b.WriteString(NL)
			b.WriteString(s)
			hasContent = true
		}
	}

	b.WriteString(t.writeNonOverriddenGettersAndSetters(indent + TAB))

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
func getExtendsClass(v Variable, ctx context.Context) (values.Prototype, error) {
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
}

func (t *Class) ResolveExpressionNames(scope Scope) error {
	if t.extends != nil {
		if err := t.extends.ResolveExpressionNames(scope); err != nil {
			return err
		}

		var err error
		t.cachedExtends, err = getExtendsClass(t.extends.GetVariable(), t.extends.Context())
		if err != nil {
			return err
		}
	}

	if t.implements != nil {
		if err := t.implements.ResolveExpressionNames(scope); err != nil {
			return err
		}

		// add self to classInterface
		implementsInterface_ := t.implements.GetVariable().GetObject()

		if implementsInterface_ == nil {
			errCtx := t.implements.Context()
			return errCtx.NewError("Error: bad implements")
		}

		switch implementsInterface := implementsInterface_.(type) {
		case *ClassInterface:
			t.cachedImplements = implementsInterface
			implementsInterface.RegisterExplicitImplementation(t)
		default:
			errCtx := t.implements.Context()
			return errCtx.NewError("Error: not an interface")
		}
	}

	for _, member := range t.members {
		subScope := t.newScope(scope)

		// super is available in every function, not just the constructor
		if t.extends != nil {
			if err := subScope.SetVariable(t.super.Name(), t.super); err != nil {
				return err
			}
		}

		if err := member.ResolveNames(subScope); err != nil {
			return err
		}
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

// reused by Enum
/*func cacheClassExtends(stack values.Stack, extends *TypeExpression, ctx context.Context) (values.Prototype, error) {
	if extends != nil {
		val, err := extends.EvalExpression(stack) // TODO: how should a TypeExpression be 'evaluated
		if err != nil {
			return nil, err
		}

		if !val.IsClass() {
			errCtx := extends.Context()

			if val.IsInterface() {
				return nil, errCtx.NewError("Error: can't inherit from an interface")
			} else {
				return nil, errCtx.NewError("Error: parent is not a class")
			}
		}

		cachedExtends, ok := val.GetClassPrototype()
		if !ok {
			panic("unexpected")
		}

		if _, isNok := cachedExtends.(*Enum); isNok {
			errCtx := extends.Context()
			return nil, errCtx.NewError("Error: cannot inherit from enum")
		}

		if _, _, isNok := IsNodeJSModule(cachedExtends); isNok {
			errCtx := extends.Context()
			return nil, errCtx.NewError("Error: cannot inherit from a backend module")
		}

		return cachedExtends, nil
	}

	return nil, nil
}*/

func (t *Class) HoistValues(stack values.Stack) error {
	return nil
}

func (t *Class) EvalExpression(stack values.Stack) (values.Value, error) {
	if t.implements != nil {
		// must happen during eval stage, because only then are all implementations registered
		if msg, ok := t.cachedImplements.IsImplementedBy(t); !ok {
			errCtx := t.implements.Context()
			return nil, errCtx.NewError("Error: " + t.Name() + " doesn't implement " + t.cachedImplements.Name() + " (" + msg + ")")
		}
	}

	return values.NewClass(t, t.Context()), nil
}

func (t *Class) EvalStatement(stack values.Stack) error {
	// add self immediately, because we might refer to it in the members
	val := values.NewClass(t, t.Context())
	if err := stack.SetValue(t.GetVariable(), val, false, t.Context()); err != nil {
		return err
	}

	if _, err := t.EvalExpression(stack); err != nil {
		return err
	}

	return nil
}

func (t *Class) EvalAsEntryPoint(stack values.Stack, ctx context.Context) error {
	// loop the public members, and evaluate the as entry points
	//  (special treatment of constructor though)
	this_, err := t.GenerateInstance(stack, nil, nil, ctx)
	if err != nil {
		return err
	}

	this := values.AssertInstance(this_)

	for _, member := range t.members {
		if member.Name() == "constructor" {
			// already checked via GenerateInstance, skip
			continue
		} else if prototypes.IsPublic(member) {
			if err := member.EvalAsEntryPoint(stack, this, ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Class) IsImplementedBy(p values.Prototype) (string, bool) {
	return "", false
}

func (t *Class) HasAncestor(interf values.Interface) bool {
	if _, ok := interf.IsImplementedBy(t); ok {
		return true
	} else if _, ok := interf.(*values.AllPrototype); ok {
		return true
	} else if t == interf {
		return true
	} else if t.extends != nil {
		return t.cachedExtends.HasAncestor(interf)
	} else {
		return false
	}
}

// potentially return Setter/Getter as second ClassMember
// return non-nil prototypes.BuiltinPrototype for further specific searching
func (t *Class) FindMember(key string, includePrivate bool,
	preferSetter bool) *ClassMember {
	for _, member := range t.members {
		if prototypes.IsPublic(member) || includePrivate {
			if member.Name() == key {
				if preferSetter && prototypes.IsGetter(member) {
					continue
				} else if (!preferSetter) && prototypes.IsSetter(member) {
					continue
				} else {
					return member
				}
			}
		}
	}

	return nil
}

// doesn't yet look in extends, but might look in child
func (t *Class) FindMemberCheckStatic(key string, this *values.Instance, includePrivate bool,
	preferSetter bool, ctx context.Context) (*ClassMember, error) {
	member := t.FindMember(key, includePrivate, preferSetter)
	if member == nil {
		// try looking in child, before allowing caller to look in extends
		if this != nil {
			if instanceClass, ok := this.GetOriginalPrototype().(*Class); ok && instanceClass != t {
				member = instanceClass.FindMember(key, includePrivate, preferSetter)
			}
		}

		if member == nil {
			return nil, nil
		}
	}

	if this != nil && prototypes.IsStatic(member) {
		return nil, ctx.NewError("Error: cannot call static method on instance")
	} else if this == nil && !prototypes.IsStatic(member) {
		return nil, ctx.NewError("Error: cannot call non-static method on class")
	}

	// look in child in case of abstract
	if prototypes.IsAbstract(member) {
		/*if t.childProto != nil {
			childClass, ok := t.childProto.(*Class)
			if !ok {
				panic("unexpected")
			}

			var err error
			member, err = childClass.FindMemberCheckStatic(key, this, includePrivate, false, ctx)
			if err != nil {
				return nil, ctx.NewError("Error: " + t.Name() + "." + key + " is abstract")
				//return nil, err
			}
		} */
		if this != nil {
			// can the following completely replace childProto?
			instanceProto := this.GetOriginalPrototype()

			if instanceClass, ok := instanceProto.(*Class); ok && instanceClass != t {
				var err error
				member, err = instanceClass.FindMemberCheckStatic(key, this, includePrivate, false, ctx)
				if err != nil {
					return nil, ctx.NewError("Error: " + t.Name() + "." + key + " is abstract")
					//return nil, err
				}
			}
		}

		// TODO: we should check that the functioninterfaces corresponds (except the abstract part)
	}

	return member, nil
}

func (t *Class) HasMember(this *values.Instance, key string, includePrivate bool) bool {
	for _, member := range t.members {
		if prototypes.IsPublic(member) || includePrivate {
			if member.Name() == key {
				return true
			}
		}
	}

	// try parent
	if t.extends == nil {
		return false
	}

	return t.cachedExtends.HasMember(this, key, includePrivate)
}

func (t *Class) isAbstract(implemented map[string]*ClassMember) bool {
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
}

// XXX: can return value always be instance?
func (t *Class) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	if typeChildren != nil {
		return nil, ctx.NewError("Error: user defined class cannot have content types")
	}

	newV, ok := v.ChangeInstanceInterface(t, false, true)
	if !ok {
		err := ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + t.Name())
		return nil, err
	}

	return newV, nil
}

func (t *Class) EvalConstructor(stack values.Stack, args []values.Value,
	childProto values.Prototype, ctx context.Context) (values.Value, error) {
	if childProto == nil {
		childProto = t
	}

	constructor := t.FindMember("constructor", false, false)
	if constructor == nil {
		if t.extends == nil {
			return nil, ctx.NewError("Error: " + t.Name() + ".constructor undefined")
		} else {
			// convert res to this
			// cachedExtends should have access to exact implementations of abstract functions
			res, err := t.cachedExtends.EvalConstructor(stack, args, childProto, ctx)
			if err != nil {
				return nil, err
			}

			if _, ok := res.ChangeInstancePrototype(t, true); !ok {
				errCtx := res.Context()
				return nil, errCtx.NewError("Error: not an instance")
			}

			return res, nil
		}
	} else {
		props := values.NewProperties(ctx)
		this := values.NewInstance(t, props, ctx)
		if childProto != t {
			this.SetOriginalPrototype(childProto)
		}
		thisRef := constructor.GetThisVariable()

		// mutation must be checked during ResolveNames stage?
		//props.Unlock() // TODO: this is not good enough, only allowed to add props in constructor scope (even more restrictive: only constructor top-level scope)

		superCalled := false
		if t.extends != nil {
			// set super
			// TODO: once super() has been called, "super" should refer to a parent casted version of this
			superFn := values.NewFunctionFunction(func(stack_ values.Stack,
				this_ *values.Instance, args_ []values.Value,
				ctx_ context.Context) (values.Value, error) {
				superCalled = true

				// ignore 'this_', and overwrite 'this'
				// cachedExtends doesnt necessarily need a constructor for super() to be called
				if !t.cachedExtends.HasMember(nil, "constructor", false) {
					if len(args_) != 0 {
						errCtx := args_[0].Context()
						return nil, errCtx.NewError("Error: unexpected args for implicit constructor")
					}

					if err := constructor.setSuper(stack_, this_, ctx_); err != nil {
						return nil, err
					}

					return nil, nil
				}

				// returns an instance, the extract the properties and merge with the 	l
				newThis_, err := t.cachedExtends.EvalConstructor(stack_, args_, childProto, ctx_)
				if err != nil {
					return nil, err
				}

				if !this.Properties().IsEmpty() {
					return nil, ctx_.NewError("Error: super called after this was mutated")
				}

				newThis := values.AssertInstance(newThis_)

				newThis.Properties().Unlock()
				if _, ok := newThis.ChangeInstancePrototype(t, true); !ok {
					panic("not an instance")
				}
				//newThis.SetOriginalPrototype(childProto)

				// should never be reset
				if err := stack_.SetValue(thisRef, newThis, false,
					ctx_); err != nil {
					return nil, err
				}

				if err := constructor.setSuper(stack_, this_, ctx_); err != nil {
					return nil, err
				}

				return nil, nil
			}, nil, this, ctx)

			// should never be reset
			if err := stack.SetValue(t.super, superFn, false,
				ctx); err != nil {
				return nil, err
			}
		}

		// might be overwritten by super (the few lines above are actually executed after the next line)
		// should never be reset
		err := stack.SetValue(thisRef, this, false, ctx)
		if err != nil {
			return nil, err
		}

		// avoid internal setting of alternate super
		if err := constructor.function.EvalAsConstructor(stack, this, args, ctx); err != nil {
			return nil, err
		}

		result_, err := stack.GetValue(thisRef, ctx)
		if err != nil {
			return nil, err
		}

		if result_ == nil {
			errCtx := ctx
			return nil, errCtx.NewError("Error: variable declared, but doesn't have a value")
		}

		result := values.AssertInstance(result_)

		result.Properties().Lock()

		if t.extends != nil && !superCalled {
			errCtx := constructor.function.Context()

			return nil, errCtx.NewError("Error: super not called")
		}

		return result, nil
	}
}

func (t *Class) EvalCachedConstructor(stack values.Stack, args []values.Value,
	childProto values.Prototype, ctx context.Context) (values.Value, error) {
	if cachedInstance, ok := stack.GetGeneratedInstance(t); ok && cachedInstance != nil {
		return values.NewContextValue(cachedInstance, ctx), nil
	} else {
		res, err := t.EvalConstructor(stack, args, childProto, ctx)
		if err != nil {
			return nil, err
		}

		// cache an instance
		if !ok {
			if genRes_, genErr := t.GenerateInstance(stack, nil, nil, ctx); genErr == nil {
				genRes, instanceOk := genRes_.(*values.Instance)
				if !instanceOk {
					panic("unexpected")
				}

				// none of the properties can be all null
				if genRes.NoPropertiesAllNull() {
					stack.SetGeneratedInstance(t, genRes)
				} else {
					stack.SetGeneratedInstance(t, nil)
				}
			} else {
				if VERBOSITY >= 3 {
					fmt.Println("unable to generate instance of ", t.Name(), " due to: ", genErr.Error())
				}
				stack.SetGeneratedInstance(t, nil)
			}
		}

		// return the evalConstructor result though
		return res, nil
	}
}

func (t *Class) GenerateInstance(stack values.Stack, keys []string, args []values.Value, ctx context.Context) (values.Value, error) {
	if keys != nil || args != nil {
		return nil, ctx.NewError("Error: parametric classes not yet supported")
	}

	// constructor simply needs to be of full type
	constructor := t.FindMember("constructor", false, false)
	if constructor == nil {
		if t.extends == nil {
			return nil, ctx.NewError("Error: " + t.Name() + ".constructor undefined")
		} else {
			// convert res to this
			res, err := t.cachedExtends.GenerateInstance(stack, keys, args, ctx)
			if err != nil {
				return nil, err
			}

			if _, ok := res.ChangeInstancePrototype(t, true); !ok {
				errCtx := res.Context()
				return nil, errCtx.NewError("Error: not an instance")
			}

			return res.RemoveLiteralness(true), nil
		}
	} else {
		if cachedInstance, ok := stack.GetGeneratedInstance(t); ok && cachedInstance != nil {
			return cachedInstance, nil
		} else {
			instanceThatIsChangedLater := values.NewDummyInstance(t, ctx)
			stack.SetGeneratedInstance(t, instanceThatIsChangedLater)

			// args might recursively depend on this stack
			args, err := constructor.GenerateArgInstances(stack, ctx)
			if err != nil {
				return nil, err
			}

			res_, err := t.EvalConstructor(stack, args, nil, ctx)
			if err != nil {
				return nil, err
			}
			res_ = res_.RemoveLiteralness(true)
			res := values.AssertInstance(res_)

			instanceThatIsChangedLater.CopyInPlace(res)

			return instanceThatIsChangedLater, nil // cache in stack, so we can use it recursively
		}
	}
}

func (t *Class) GetMember(stack values.Stack, this *values.Instance, key string,
	includePrivate bool, ctx context.Context) (values.Value, error) {
	member, err := t.FindMemberCheckStatic(key, this, includePrivate, false, ctx)
	if err != nil {
		return nil, err
	}

	if member == nil {
		if t.extends != nil {
			// XXX: how to do abstract implementation injection?
			return t.cachedExtends.GetMember(stack, this, key, includePrivate, ctx)
		} else {
			if t.HasMember(this, key, includePrivate) {
				return nil, ctx.NewError("Error: " + t.Name() + "." + key + " is a setter")
			} else {

				havePrivateAccess := " (don't have private access)"
				if includePrivate {
					havePrivateAccess = " (have private access)"
				}

				err := ctx.NewError("Error: " + t.Name() + "." + key + " undefined" + havePrivateAccess)
				return nil, err
			}
		}
	} else {
		switch {
		case prototypes.IsStaticGetter(member):
			return member.EvalFunction(stack, nil, []values.Value{}, ctx)
		case prototypes.IsStatic(member):
			return values.NewFunction(member, stack, this, ctx), nil
		case prototypes.IsGetter(member):
			return member.EvalFunction(stack, this, []values.Value{}, ctx)
		case prototypes.IsSetter(member):
			return nil, ctx.NewError("Error: " + t.Name() + "." + key + " is a setter")
		case prototypes.IsNormal(member):
			return values.NewFunction(member, stack, this, ctx), nil
		default:
			panic("unhandled")
		}
	}
}

func (t *Class) SetMember(stack values.Stack, this *values.Instance, key string,
	arg values.Value, includePrivate bool, ctx context.Context) error {
	member, err := t.FindMemberCheckStatic(key, this, includePrivate, true, ctx)
	if err != nil {
		return err
	}

	if member == nil {
		if t.extends != nil {
			return t.cachedExtends.SetMember(stack, this, key, arg, includePrivate, ctx)
		} else {
			if t.HasMember(this, key, includePrivate) {
				return ctx.NewError("Error: " + t.Name() + "." + key + " is a getter")
			} else {
				return ctx.NewError("Error: " + t.Name() +
					" doesn't have a member named " + key)
			}
		}
	} else {
		switch {
		case prototypes.IsSetter(member):
			err := member.EvalMethod(stack, this, []values.Value{arg}, ctx)
			return err
		default:
			return ctx.NewError("Error: " + t.Name() + "." + key + " is not a setter")
		}
	}
}

func (t *Class) GetIndex(stack values.Stack, this *values.Instance, index values.Value,
	ctx context.Context) (values.Value, error) {
  // it is indexable if the parent is
  if t.cachedExtends != nil {
    return t.cachedExtends.GetIndex(stack, this, index, ctx)
  } else {
    return nil, ctx.NewError("Error: not indexable")
  }
}

func (t *Class) SetIndex(stack values.Stack, this *values.Instance, index values.Value,
	arg values.Value, ctx context.Context) error {
  if t.cachedExtends != nil {
    return t.cachedExtends.SetIndex(stack, this, index, arg, ctx)
  } else {
    return ctx.NewError("Error: not indexable")
  }
}

func (t *Class) ResolveExpressionActivity(usage Usage) error {
	tmp := usage.InFunction()
	usage.SetInFunction(false)

	for _, mf := range t.members {
		if err := mf.ResolveActivity(usage); err != nil {
			usage.SetInFunction(tmp)
			return err
		}
	}

	usage.SetInFunction(tmp)

	return nil
}

func (t *Class) ResolveStatementActivity(usage Usage) error {
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

	return t.ResolveExpressionActivity(usage)
}

func (t *Class) UniversalExpressionNames(ns Namespace) error {
	if t.extends != nil {
		if err := t.extends.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	if t.IsUniversal() {
		if err := ns.UniversalName(t.clType.GetVariable(), t.universalName); err != nil {
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
	if t.extends != nil {
		if err := t.extends.UniqueExpressionNames(ns); err != nil {
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
		if err := ns.ClassName(t.clType.GetVariable()); err != nil {
			return err
		}
	}

	return t.UniqueExpressionNames(ns)
}

func (t *Class) Walk(fn WalkFunc) error {
  if err := t.clType.Walk(fn); err != nil {
    return err
  }

  if t.extends != nil {
    if err := t.extends.Walk(fn); err != nil {
      return err
    }
  }

  if t.implements != nil {
    if err := t.implements.Walk(fn); err != nil {
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

func (t *Class) LoopForIn(this *values.Instance, fn func(values.Value) error, ctx context.Context) error {
	return ctx.NewError("Error: inheriting from containers not yet implemented")
}

func (t *Class) LoopForOf(this *values.Instance, fn func(values.Value) error, ctx context.Context) error {
	return ctx.NewError("Error: inheriting from containers not yet implemented")
}

func (t *Class) writeNonOverriddenGettersAndSetters(indent string) string {
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
}
