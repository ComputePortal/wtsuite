package glsl

import (
	"fmt"
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type LiteralFloat struct {
	value float64
	LiteralData
}

func NewLiteralFloat(value float64, ctx context.Context) *LiteralFloat {
	return &LiteralFloat{value, newLiteralData(ctx)}
}

func (t *LiteralFloat) Value() float64 {
	return t.value
}

func (t *LiteralFloat) Dump(indent string) string {
	return indent + "LiteralFloat(" + t.WriteExpression() + ")\n"
}

func (t *LiteralFloat) WriteExpression() string {
  res := fmt.Sprintf("%g", t.value)
  if !strings.ContainsAny(res, ".e") {
    res += ".0"
  }

  return res
}

func (t *LiteralFloat) EvalExpression() (values.Value, error) {
	return values.NewScalar("float", t.Context()), nil
}

func IsLiteralFloat(t Expression) bool {
	_, ok := t.(*LiteralFloat)
	return ok
}
