package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type FunctionInterface struct {
  retType *TypeExpression // can be nil for void return
  nameExpr *VarExpression
  args []*FunctionArgument
}

func NewFunctionInterface(retTypeExpr *TypeExpression, name string, args []*FunctionArgument, ctx context.Context) *FunctionInterface {
  return &FunctionInterface{
    retTypeExpr,
    NewVarExpression(name, ctx),
    args,
  }
}

func (fi *FunctionInterface) Context() context.Context {
  return fi.nameExpr.Context()
}

func (fi *FunctionInterface) Name() string {
  return fi.nameExpr.Name()
}

func (fi *FunctionInterface) GetVariable() Variable {
	return fi.nameExpr.GetVariable()
}

func (fi *FunctionInterface) Dump(indent string) string {
	var b strings.Builder

	// dumping of name can be done here, but writing can't be done below because we need exact control on Function
  if fi.retType != nil {
    b.WriteString(fi.retType.Dump(indent))
  }

	if fi.Name() != "" {
		b.WriteString(fi.Name())
	}

	b.WriteString("(")

	for i, arg := range fi.args {
		b.WriteString(arg.Dump(indent + "  "))

		if i < len(fi.args)-1 {
			b.WriteString(patterns.COMMA)
		}
	}

	b.WriteString(")")

	b.WriteString("\n")

	return b.String()
}

func (fi *FunctionInterface) WriteInterface() string {
  var b strings.Builder

  if fi.retType == nil {
    b.WriteString("void")
  } else {
    b.WriteString(fi.retType.WriteExpression())
  }

  b.WriteString(" ")
  b.WriteString(fi.nameExpr.WriteExpression())
  b.WriteString("(")

  for i, arg := range fi.args {
    b.WriteString(arg.WriteArgument())
    if i < len(fi.args) - 1 {
      b.WriteString(",")
    }
  }

  b.WriteString(")")

  return b.String()
}

func (fi *FunctionInterface) ResolveNames(scope Scope) error {
	if fi.retType != nil {
		if err := fi.retType.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	for _, arg := range fi.args {
		if err := arg.ResolveNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (fi *FunctionInterface) EvalCall(args []values.Value, ctx context.Context) (values.Value, error) {
  if args != nil {
    if len(args) != len(fi.args) {
      errCtx := ctx
      return nil, errCtx.NewError("Error: expected " + strconv.Itoa(len(fi.args)) + " args, got " + strconv.Itoa(len(args)) + " args")
    }

    for i, argVal := range args {
      argType, err := fi.args[i].GetTypeValue()
      if err != nil {
        return nil, err
      }

      if err := argType.Check(argVal, argVal.Context()); err != nil {
        return nil, err
      }
    }
  }

  if fi.retType == nil {
    return nil, nil
  } else {
    return fi.retType.Instantiate(ctx)
  }
}

func (fi *FunctionInterface) UniqueNames(ns Namespace) error {
	for _, arg := range fi.args {
		if err := arg.UniqueNames(ns); err != nil {
			return err
		}
	}

	return nil
}
