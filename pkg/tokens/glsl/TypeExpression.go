package glsl

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type TypeExpression struct {
  VarExpression
}

func NewTypeExpression(name string, ctx context.Context) *TypeExpression {
  return &TypeExpression{newVarExpression(name, ctx)}
}

func (t *TypeExpression) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Type(")
  b.WriteString(t.VarExpression.Name())
  b.WriteString(")")

  return b.String()
}

func (t *TypeExpression) WriteExpression() string {
  return t.VarExpression.WriteExpression()
}

func (t *TypeExpression) EvalExpression() (values.Value, error) {
  errCtx := t.Context()
  panic(errCtx.NewError("use instantiate instead").Error())
}

func (t *TypeExpression) InstantiateUniform(ctx context.Context) (values.Value, error) {
  valType, err := t.VarExpression.EvalExpression()
  if err != nil {
    return nil, err
  }

  return valType.Instantiate(ctx)
}

func (t *TypeExpression) Instantiate(ctx context.Context) (values.Value, error) {
  val, err := t.InstantiateUniform(ctx)
  if err != nil {
    return nil, err
  }

  if values.IsSampler2D(val) {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: sampler2D only available as uniform")
  } else if values.IsSamplerCube(val) {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: samplerCube only available as uniform")
  }

  return val, nil
}
