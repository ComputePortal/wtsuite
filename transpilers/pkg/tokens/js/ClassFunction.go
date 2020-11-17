package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

// implements the values.Callable interface
type ClassFunction struct {
	function *Function
}

func NewClassFunction(fn *Function) *ClassFunction {
	return &ClassFunction{fn}
}

func (m *ClassFunction) Context() context.Context {
  return m.function.Context()
}

func (m *ClassFunction) Name() string {
	return m.function.Name()
}

func (m *ClassFunction) Role() prototypes.FunctionRole {
	return m.function.Role()
}

func (m *ClassFunction) IsUniversal() bool {
  return true
}

func (m *ClassFunction) GetThisVariable() Variable {
	return m.function.GetThisVariable()
}

func (m *ClassFunction) getModifierString() string {
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

func (m *ClassFunction) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString(m.getModifierString())

	b.WriteString(strings.TrimLeft(m.function.Dump(indent+"  "), " "))

	return b.String()
}

func (m *ClassFunction) WriteStatement(indent string) string {
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

func (m *ClassFunction) ResolveNames(scope Scope) error {
	return m.function.ResolveExpressionNames(scope)
}

func (m *ClassFunction) GetValue(ctx context.Context) (values.Value, error) {
  fn, err := m.function.GetFunctionValue()
  if err != nil {
    return nil, err
  }

  if prototypes.IsGetter(m) {
    return fn.EvalFunction([]values.Value{}, false, ctx)
  } else if prototypes.IsSetter(m) {
    return nil, ctx.NewError("Error: is a setter")
  } else {
    return fn, nil
  }
}

func (m *ClassFunction) SetValue(v values.Value, ctx context.Context) error {
  fn, err := m.function.GetFunctionValue()
  if err != nil {
    return err
  }

  if prototypes.IsGetter(m) && !prototypes.IsSetter(m) {
    return ctx.NewError("Error: is a getter")
  } else {
    _, err = fn.EvalFunction([]values.Value{v}, true, ctx)
    return err
  }
}

func (m *ClassFunction) EvalExpression() (values.Value, error) {
  return m.function.EvalExpression()
}

func (m *ClassFunction) ResolveActivity(usage Usage) error {
	return m.function.ResolveStatementActivity(usage)
}

func (m *ClassFunction) UniversalNames(ns Namespace) error {
	return m.function.UniversalExpressionNames(ns)
}

func (m *ClassFunction) UniqueNames(ns Namespace) error {
	return m.function.UniqueExpressionNames(ns)
}

func (m *ClassFunction) Walk(fn WalkFunc) error {
  if err := m.function.Walk(fn); err != nil {
    return err
  }

  return fn(m)
}
