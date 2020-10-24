package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
	"../patterns"
)

type FunctionArgument struct {
	name       *VarExpression
	constraint *TypeExpression // can be nil
	def        Expression // can be nil
	rest       bool
	TokenData
}

func NewFunctionArgument(name string, constraint *TypeExpression, def Expression,
	ctx context.Context) (*FunctionArgument, error) {
	return &FunctionArgument{NewConstantVarExpression(name, ctx), constraint, def, false,
		TokenData{ctx}}, nil
}

func NewRestFunctionArgument(name string, ctx context.Context) (*FunctionArgument, error) {
	return &FunctionArgument{NewConstantVarExpression(name, ctx), nil, nil, true,
		TokenData{ctx}}, nil
}

func (fa *FunctionArgument) Name() string {
	return fa.name.Name()
}

func (fa *FunctionArgument) ConstraintName() string {
	if fa.constraint == nil {
		return ""
	} else {
		return fa.constraint.Name()
	}
}

func (fa *FunctionArgument) Rest() bool {
	return fa.rest
}

func (fa *FunctionArgument) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Arg(")

	b.WriteString(fa.Name())

	if fa.constraint != nil {
		b.WriteString(patterns.DCOLON)
		b.WriteString(fa.constraint.Dump(""))
	}

	if fa.def != nil {
		b.WriteString(patterns.EQUAL)
		b.WriteString(fa.def.Dump(""))
	}

	b.WriteString(")\n")

	return b.String()
}

func (fa *FunctionArgument) Write(rest bool) string {
	var b strings.Builder

	if rest {
		b.WriteString("...")
	}

	b.WriteString(fa.Name())

	if fa.def != nil {
		b.WriteString("=")
		b.WriteString(fa.def.WriteExpression())
	}

	return b.String()
}

func (fa *FunctionArgument) ResolveInterfaceNames(scope Scope) error {
	if fa.def != nil {
		errCtx := fa.Context()
		return errCtx.NewError("Error: interface member cant have default")
	}

	if fa.constraint != nil {
		if err := fa.constraint.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (fa *FunctionArgument) ResolveNames(outer Scope, inner Scope) error {
	if fa.def != nil {
		if err := fa.def.ResolveExpressionNames(inner); err != nil {
			return err
		}
	}

	if fa.constraint != nil {
		if err := fa.constraint.ResolveExpressionNames(outer); err != nil {
			return err
		}
	}

	if fa.name.Name() != "_" {
		if err := inner.SetVariable(fa.name.Name(), fa.name.GetVariable()); err != nil {
			return err
		}
	}

	return nil
}

func (fa *FunctionArgument) EvalInterface(stack values.Stack) error {
	if fa.constraint != nil {
		constraintClassVal, err := fa.constraint.EvalExpression(stack)
		if err != nil {
			return err
		}

		_, ok := constraintClassVal.GetClassInterface()
		if !ok {
			errCtx := fa.constraint.Context()
			return errCtx.NewError("Error: not a class or interface")
		}
	}

	return nil
}

// TODO: should this be nested, or the nesting elsewhere?
func (fa *FunctionArgument) constrainArg(stack values.Stack, arg values.Value,
	ctx context.Context) (values.Value, error) {

	if fa.constraint != nil {
		return fa.constraint.Constrain(stack, arg)
	}

	return arg, nil
}

func (fa *FunctionArgument) GenerateArgInstance(stack values.Stack, ctx context.Context) (values.Value, error) {
	if fa.constraint == nil {
		return nil, ctx.NewError("Error: arg type not specified")
	}

	val, err := fa.constraint.GenerateInstance(stack, ctx)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (fa *FunctionArgument) EvalArg(stack values.Stack, arg values.Value,
	ctx context.Context) error {
	isNull := arg.IsNull()

	var err error
	arg, err = fa.constrainArg(stack, arg, ctx)
	if err != nil {
		return err
	}

	if isNull {
		// try to instantiate
		altArg, err := fa.GenerateArgInstance(stack, ctx)
		if err == nil {
			arg = altArg
		}
	}

	return stack.SetValue(fa.name.GetVariable(), arg, false, ctx)
}

func (fa *FunctionArgument) EvalDef(stack values.Stack, ctx context.Context) error {
	if fa.def == nil {
		errCtx := fa.Context()
		err := errCtx.NewError("Error: argument doesn't have a default")
		if VERBOSITY >= 1 {
			err.AppendContextString("Info: called here", ctx)
		}
		return err
	}

	val, err := fa.def.EvalExpression(stack)
	if err != nil {
		return err
	}

	val, err = fa.constrainArg(stack, val, ctx)
	if err != nil {
		return err
	}

	return stack.SetValue(fa.name.GetVariable(), val, false, ctx)
}

func (fa *FunctionArgument) EvalRest(stack values.Stack, args []values.Value,
	ctx context.Context) error {
	val := prototypes.NewLiteralArray(args, ctx)

	return stack.SetValue(fa.name.GetVariable(), val, false, ctx)
}

func (fa *FunctionArgument) UniversalNames(ns Namespace) error {
	if fa.constraint != nil {
		if err := fa.constraint.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	if fa.def != nil {
		if err := fa.def.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (fa *FunctionArgument) UniqueNames(ns Namespace) error {
	ns.ArgName(fa.name.GetVariable())

	if fa.constraint != nil {
		if err := fa.constraint.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	if fa.def != nil {
		if err := fa.def.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (fa* FunctionArgument) Walk(fn WalkFunc) error {
  if err := fa.name.Walk(fn); err != nil {
    return err
  }

  if fa.constraint != nil {
    if err := fa.constraint.Walk(fn); err != nil {
      return err
    }
  }

  if fa.def != nil {
    if err := fa.def.Walk(fn); err != nil {
      return err
    }
  }

  return fn(fa)
}

func (fa *FunctionArgument) IsImplementedByOtherArg(other *FunctionArgument) bool {
	// TODO: --exact-interface to assert matching arg names as well
	if fa.constraint == nil {
		return other.constraint == nil
	} else if other.constraint == nil {
		return false
	}

	return fa.constraint.Dump("") == other.constraint.Dump("")
}

