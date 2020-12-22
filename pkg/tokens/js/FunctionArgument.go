package js

import (
	"strings"

	"./values"

	"../context"
	"../patterns"
)

type FunctionArgument struct {
	nameExpr   *VarExpression
	typeExpr   *TypeExpression // can't be nil (must at least be any)
	def        Expression // can be nil
	TokenData
}

func NewFunctionArgument(name string, typeExpr *TypeExpression, def Expression,
	ctx context.Context) (*FunctionArgument, error) {
  if typeExpr == nil {
    panic("must at least be any")
  }

	return &FunctionArgument{NewVarExpression(name, ctx), typeExpr, def,
		TokenData{ctx}}, nil
}

func (fa *FunctionArgument) Name() string {
	return fa.nameExpr.Name()
}

func (fa *FunctionArgument) TypeName() string {
  return fa.typeExpr.Name()
}

func (fa *FunctionArgument) GetVariable() Variable {
  return fa.nameExpr.GetVariable()
}

func (fa *FunctionArgument) HasDefault() bool {
  return fa.def != nil
}

func (fa *FunctionArgument) AssertNoDefault() error {
  if fa.HasDefault() {
    errCtx := fa.def.Context()
    return errCtx.NewError("Error: unexpected default value")
  }

  return nil
}

func (fa *FunctionArgument) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Arg(")

	b.WriteString(fa.Name())

  b.WriteString(patterns.DCOLON)
  b.WriteString(fa.typeExpr.Dump(""))

	if fa.HasDefault() {
		b.WriteString(patterns.EQUAL)
		b.WriteString(fa.def.Dump(""))
	}

	b.WriteString(")\n")

	return b.String()
}

func (fa *FunctionArgument) Write() string {
	var b strings.Builder

	b.WriteString(fa.Name())

	if fa.HasDefault() {
		b.WriteString("=")
		b.WriteString(fa.def.WriteExpression())
	}

	return b.String()
}

func (fa *FunctionArgument) ResolveInterfaceNames(scope Scope) error {
	if fa.HasDefault() {
		errCtx := fa.Context()
		return errCtx.NewError("Error: interface member cant have default")
	}

  if err := fa.typeExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

	return nil
}

func (fa *FunctionArgument) ResolveNames(scope Scope) error {
	if fa.HasDefault() {
		if err := fa.def.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

  if err := fa.typeExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }


  name := fa.nameExpr.Name()
	if name != "_" {
    variable := fa.nameExpr.GetVariable()


		if err := scope.SetVariable(name, variable); err != nil {
			return err
		}
	}

  if fa.HasDefault() {
    if err := fa.def.ResolveExpressionNames(scope); err != nil {
      return err
    }
  }

	return nil
}

func (fa *FunctionArgument) GetValue() (values.Value, error) {
	val, err := fa.typeExpr.EvalExpression()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (fa *FunctionArgument) Eval() error {
  argVal, err := fa.GetValue()
  if err != nil {
    return err
  }

  variable := fa.nameExpr.GetVariable()
  variable.SetValue(argVal)

  // also check that the default respects this type
  if fa.HasDefault() {
    defVal, err := fa.def.EvalExpression()
    if err != nil {
      return err
    }

    if err := argVal.Check(defVal, fa.Context()); err != nil {
      return err
    }
  }

  return nil
}

func (fa *FunctionArgument) UniversalNames(ns Namespace) error {
	if fa.typeExpr != nil {
		if err := fa.typeExpr.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	if fa.HasDefault() {
		if err := fa.def.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (fa *FunctionArgument) UniqueNames(ns Namespace) error {
	ns.ArgName(fa.nameExpr.GetVariable())

	if fa.typeExpr != nil {
		if err := fa.typeExpr.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	if fa.HasDefault() {
		if err := fa.def.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (fa* FunctionArgument) Walk(fn WalkFunc) error {
  if err := fa.nameExpr.Walk(fn); err != nil {
    return err
  }

  if fa.typeExpr != nil {
    if err := fa.typeExpr.Walk(fn); err != nil {
      return err
    }
  }

  if fa.HasDefault() {
    if err := fa.def.Walk(fn); err != nil {
      return err
    }
  }

  return fn(fa)
}
