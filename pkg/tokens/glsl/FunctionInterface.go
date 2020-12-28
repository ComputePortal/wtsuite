package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type FunctionInterface struct {
  retType *TypeExpression
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

func (fi *FunctionInterface) Name() string {
  return fi.nameExpr.Name()
}

func (fi *FunctionInterface) Dump(indent string) string {
	var b strings.Builder

	// dumping of name can be done here, but writing can't be done below because we need exact control on Function
  b.WriteString(fi.retType.Dump(indent))

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

  b.WriteString(fi.retType.WriteExpression())

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
