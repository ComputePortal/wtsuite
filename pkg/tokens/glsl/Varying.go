package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Varying struct {
  precType PrecisionType
  Pointer
}

func NewVarying(precType PrecisionType, typeExpr *TypeExpression, name string, ctx context.Context) *Varying {
  return &Varying{precType, newPointer(typeExpr, NewVarExpression(name, ctx), -1, ctx)}
}

func (t *Varying) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Varying(")
  b.WriteString(PrecisionTypeToString(t.precType))
  b.WriteString(" ")
  b.WriteString(t.typeExpr.Dump(""))
  b.WriteString(" ")
  b.WriteString(t.nameExpr.Dump(""))
  b.WriteString("\n")

  return b.String()
}

func (t *Varying) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  // TODO: check if actually used
  b.WriteString(indent)
  b.WriteString("varying ")
  if t.precType != DEFAULTP {
    b.WriteString(PrecisionTypeToString(t.precType))
    b.WriteString(" ")
  }
  b.WriteString(t.typeExpr.WriteExpression())
  b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())
  b.WriteString(";")

  return b.String()
}

func (t *Varying) EvalStatement() error {
  if err := t.Pointer.EvalStatement(); err != nil {
    return err
  }

  variable := t.GetVariable()

  val := variable.GetValue()

  if !values.IsSimple(val) {
    errCtx := val.Context()
    return errCtx.NewError("Error: expected simple type, got " +val.TypeName())
  }

  return nil
}

func (t *Varying) ResolveStatementActivity(usage Usage) error {
  // can be set without usage in case of vertex shader
  if TARGET == "fragment" {
    return t.Pointer.ResolveStatementActivity(usage)
  } else {
    return nil
  }
}

func (t *Varying) Collect(varyings map[string]string) error {
  // expecting only simple types
  varyings[t.Name()] = t.typeExpr.WriteExpression()

  return nil
}
