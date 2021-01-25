package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type FunctionArgument struct {
  role FunctionArgumentRole
  typeExpr *TypeExpression
  nameExpr *VarExpression
  length int
  TokenData
}

func NewFunctionArgument(role FunctionArgumentRole, typeExpr *TypeExpression, name string, length int, ctx context.Context) *FunctionArgument {
  return &FunctionArgument{
    role,
    typeExpr,
    NewVarExpression(name, ctx),
    length,
    newTokenData(ctx),
  }
}

func (fa *FunctionArgument) Name() string {
  return fa.nameExpr.Name()
}

func (fa *FunctionArgument) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Arg(")

  b.WriteString(fa.typeExpr.Dump(""))
	b.WriteString(fa.Name())

	b.WriteString(")\n")

	return b.String()
}

func (fa *FunctionArgument) WriteArgument() string {
  var b strings.Builder

  b.WriteString(RoleToString(fa.role))

  b.WriteString(fa.typeExpr.WriteExpression())
  b.WriteString(" ")
  b.WriteString(fa.Name())

  if fa.length > 0 {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(fa.length))
    b.WriteString("]")
  }

  return b.String()
}

func (fa *FunctionArgument) GetTypeValue() (values.Value, error) {
  return fa.typeExpr.Instantiate(fa.Context())
}

func (fa *FunctionArgument) ResolveNames(scope Scope) error {
  if err := fa.typeExpr.ResolveExpressionNames(scope); err != nil {
    return err
  }

  name := fa.nameExpr.Name()
  variable := fa.nameExpr.GetVariable()

  val, err := fa.typeExpr.Instantiate(fa.Context())
  if err != nil {
    return err
  }

  if fa.length > 0 {
    val = values.NewArray(val, fa.length, fa.Context())
  }

  variable.SetValue(val)

  if err := scope.SetVariable(name, variable); err != nil {
    return err
  }

  return nil
}

func (fa *FunctionArgument) UniqueNames(ns Namespace) error {
	ns.ArgName(fa.nameExpr.GetVariable())
  return nil
}
