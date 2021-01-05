package glsl

import (
  "strconv"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type Uniform struct {
  Pointer
}

func NewUniform(typeExpr *TypeExpression, name string, n int, ctx context.Context) *Uniform {
  return &Uniform{newPointer(typeExpr, NewVarExpression(name, ctx), n, ctx)}
}

func (t *Uniform) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Uniform(")

  b.WriteString(t.typeExpr.Dump(""))
  b.WriteString(" ")
  b.WriteString(t.nameExpr.Dump(""))

  if (t.length > 0) {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

  b.WriteString(")\n")

  return b.String()
}

func (t *Uniform) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  var b strings.Builder

  // TODO: check if actually used
  b.WriteString(indent)
  b.WriteString("uniform ")
  b.WriteString(t.typeExpr.WriteExpression())
  b.WriteString(" ")
  b.WriteString(t.nameExpr.WriteExpression())

  if (t.length > 0) {
    b.WriteString("[")
    b.WriteString(strconv.Itoa(t.length))
    b.WriteString("]")
  }

  return b.String()
}

// different from Pointer.EvalStatement(), because different types are allowed
func (t *Uniform) EvalStatement() error {
  variable := t.nameExpr.GetVariable()

  val, err := t.typeExpr.InstantiateUniform(t.Context())
  if err != nil {
    return err
  }

  if t.length > 0 {
    val = values.NewArray(val, t.length, t.Context())
  }

  variable.SetValue(val)

  return nil
}
