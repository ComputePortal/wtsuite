package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type LiteralBool struct {
	value bool
	LiteralData
}

func NewLiteralBool(value bool, ctx context.Context) *LiteralBool {
	return &LiteralBool{value, newLiteralData(ctx)}
}

func (t *LiteralBool) Value() bool {
	return t.value
}

func (t *LiteralBool) Dump(indent string) string {
	return indent + "LiteralBool(" + t.WriteExpression() + ")\n"
}

func (t *LiteralBool) WriteExpression() string {
	s := "false"
	if t.value {
		s = "true"
	}

	return s
}

func (t *LiteralBool) EvalExpression() (values.Value, error) {
  return values.NewScalar("bool", t.Context()), nil
}

func IsLiteralBool(t Expression) bool {
	_, ok := t.(*LiteralBool)
	return ok
}
