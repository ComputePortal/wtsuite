package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
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
