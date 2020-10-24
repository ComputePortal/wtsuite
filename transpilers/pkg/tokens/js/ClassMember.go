package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// implements the values.Callable interface
type ClassMember struct {
	function *Function
	proto    values.Prototype // prototype to which this ClassMember belongs (needed to set the correct super in the stack)
	super    Variable
}

func NewClassMember(fn *Function, proto values.Prototype) *ClassMember {
	return &ClassMember{fn, proto, nil}
}

func (m *ClassMember) Name() string {
	return m.function.Name()
}

func (m *ClassMember) Length() int {
	return m.function.Length()
}

func (m *ClassMember) Role() prototypes.FunctionRole {
	return m.function.Role()
}

func (m *ClassMember) GetThisVariable() Variable {
	return m.function.GetThisVariable()
}

func (m *ClassMember) getModifierString() string {
	s := ""

	if prototypes.IsGetter(m) {
		s += "get "
	}

	if prototypes.IsSetter(m) {
		s += "set "
	}

	if prototypes.IsStatic(m) {
		s += "static "
	}

	if prototypes.IsAsync(m) {
		s += "async "
	}

	return s
}

func (m *ClassMember) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString(m.getModifierString())

	b.WriteString(strings.TrimLeft(m.function.Dump(indent+"  "), " "))

	return b.String()
}

func (m *ClassMember) WriteStatement(indent string) string {
	if prototypes.IsAbstract(m) {
		return ""
	}

	var b strings.Builder

	b.WriteString(indent)

	b.WriteString(m.getModifierString())

	fn := m.function
	b.WriteString(fn.Name())

	b.WriteString(fn.writeBody(indent, NL, TAB))

	return b.String()
}

func (m *ClassMember) ResolveNames(scope Scope) error {
	// get the super variable, for later!
	if scope.HasVariable("super") {
		var err error = nil
		m.super, err = scope.GetVariable("super")
		if err != nil {
			panic(err)
		}
	}

	return m.function.ResolveExpressionNames(scope)
}

func (m *ClassMember) GenerateArgInstances(stack values.Stack, ctx context.Context) ([]values.Value, error) {
	return m.function.GenerateArgInstances(stack, ctx)
}

func (m *ClassMember) setSuper(stack values.Stack, this *values.Instance, ctx context.Context) error {
	if m.super != nil && !prototypes.IsStatic(m) {
		superObj, ok := this.ChangeInstancePrototype(m.proto.GetParent(), false)
		if !ok {
			panic("unexpected")
		}

		if err := stack.SetValue(m.super, superObj, false, ctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *ClassMember) EvalFunction(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	if prototypes.IsAbstract(m) {
		err := ctx.NewError("Error: can't call abstract function on " + this.TypeName() + " (hint: wrap with if (t !== null){...})")
		return nil, err
	}

	if err := m.setSuper(stack, this, ctx); err != nil {
		return nil, err
	}

	if values.ALLOW_CACHING {

		if v, ok := stack.GetCacheValue(m.function); ok {
			if v != nil {
				return values.NewContextValue(v, ctx), nil
			} else {
				return m.function.EvalFunction(stack, this, args, ctx)
			}
		} else {
			var thisGen *values.Instance = nil
			thisGenOk := true
			if !prototypes.IsStatic(m) {
				thisGen, thisGenOk = stack.GetGeneratedInstance(m.proto.(*Class))
			}

			if thisGenOk {
				genRes, err := m.function.evalFunctionAsEntryPoint(stack, thisGen, ctx)
				if err != nil {
					stack.SetCacheValue(m.function, nil)
					return m.function.EvalFunction(stack, this, args, ctx)
				} else {
					stack.SetCacheValue(m.function, genRes)
					return genRes, nil
				}
			} else {
				stack.SetCacheValue(m.function, nil)
				return m.function.EvalFunction(stack, this, args, ctx)
			}
		}
	} else {
		return m.function.EvalFunction(stack, this, args, ctx)
	}
}

func (m *ClassMember) EvalFunctionNoReturn(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	if prototypes.IsAbstract(m) {
		return nil, ctx.NewError("Error: can't call abstract function on " + this.TypeName() + " (hint: wrap with if (t !== null){...})")
	}

	if err := m.setSuper(stack, this, ctx); err != nil {
		return nil, err
	}

	return m.function.EvalFunctionNoReturn(stack, this, args, ctx)
}

func (m *ClassMember) EvalMethod(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) error {
	if prototypes.IsAbstract(m) {
		return ctx.NewError("Error: can't call abstract method")
	}

	if err := m.setSuper(stack, this, ctx); err != nil {
		return err
	}

	if values.ALLOW_CACHING {
		// try using a cached value first
		if v, ok := stack.GetCacheValue(m.function); ok {
			if v != nil {
				if !v.IsVoid() {
					errCtx := ctx
					return errCtx.NewError("Error: should be void for method")
				}
				return nil
			} else {
				return m.function.EvalMethod(stack, this, args, ctx)
			}
		} else {
			var thisGen *values.Instance = nil
			thisGenOk := true
			if !prototypes.IsStatic(m) {
				thisGen, thisGenOk = stack.GetGeneratedInstance(m.proto.(*Class))
			}
			if thisGenOk {
				err := m.function.evalMethodAsEntryPoint(stack, thisGen, ctx)
				if err != nil {
					stack.SetCacheValue(m.function, nil)
					return m.function.EvalMethod(stack, this, args, ctx)
				} else {
					if m.function.IsAsync() {
						promise, err := prototypes.NewVoidPromise(ctx)
						if err != nil {
							panic(err)
						}
						stack.SetCacheValue(m.function, promise)
					} else {
						stack.SetCacheValue(m.function, values.NewVoid(ctx))
					}
					return nil
				}
			} else {
				stack.SetCacheValue(m.function, nil)
				return m.function.EvalMethod(stack, this, args, ctx)
			}
		}
	} else {
		return m.function.EvalMethod(stack, this, args, ctx)
	}
}

func (m *ClassMember) EvalAsEntryPoint(stack values.Stack, this *values.Instance, ctx context.Context) error {
	if prototypes.IsAbstract(m) {
		return ctx.NewError("Error: can't call abstract method")
	}

	if err := m.setSuper(stack, this, ctx); err != nil {
		return err
	}

	return m.function.EvalAsEntryPoint(stack, this, ctx)
}

func (m *ClassMember) ResolveActivity(usage Usage) error {
	return m.function.ResolveStatementActivity(usage)
}

func (m *ClassMember) UniversalNames(ns Namespace) error {
	return m.function.UniversalExpressionNames(ns)
}

func (m *ClassMember) UniqueNames(ns Namespace) error {
	return m.function.UniqueExpressionNames(ns)
}

func (m *ClassMember) Walk(fn WalkFunc) error {
  if err := m.function.Walk(fn); err != nil {
    return err
  }

  return fn(m)
}
