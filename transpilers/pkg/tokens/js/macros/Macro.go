package macros

import (
	"../"

	"../values"

	"../../context"
)

var COMPACT = false

type Macro struct {
	args []js.Expression
	ctx  context.Context
}

type MacroConstructor func([]js.Expression, context.Context) (js.Expression, error)

type MacroGroup struct {
	macros map[string]MacroConstructor
}

func newMacro(args []js.Expression, ctx context.Context) Macro {
	return Macro{args, ctx}
}

func (m *Macro) Context() context.Context {
	return m.ctx
}

func (m *Macro) ResolveExpressionNames(scope js.Scope) error {
	for _, arg := range m.args {
		if err := arg.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (m *Macro) ResolveExpressionActivity(usage js.Usage) error {
	for _, arg := range m.args {
		if err := arg.ResolveExpressionActivity(usage); err != nil {
			return err
		}
	}

	return nil
}

func (m *Macro) UniversalExpressionNames(ns js.Namespace) error {
	for _, arg := range m.args {
		if err := arg.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (m *Macro) UniqueExpressionNames(ns js.Namespace) error {
	for _, arg := range m.args {
		if err := arg.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (m *Macro) evalArgs(stack values.Stack) ([]values.Value, error) {
	res := make([]values.Value, len(m.args))

	for i, expr := range m.args {
		arg, err := expr.EvalExpression(stack)
		if err != nil {
			return nil, err
		}

		res[i] = arg
	}

	return res, nil
}

// TODO: implement for each specific macro
func (m *Macro) Walk(fn js.WalkFunc) error {
  for _, arg := range m.args {
    if err := arg.Walk(fn); err != nil {
      return err
    }
  }

  return fn(m)
}
