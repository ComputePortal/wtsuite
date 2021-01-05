package glsl

import (
	"fmt"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type LiteralInt struct {
	value int
	LiteralData
}

func NewLiteralInt(value int, ctx context.Context) *LiteralInt {
	return &LiteralInt{value, newLiteralData(ctx)}
}

func (t *LiteralInt) Value() int {
	return t.value
}

func (t *LiteralInt) Dump(indent string) string {
	return indent + "LiteralInt(" + t.WriteExpression() + ")\n"
}

func (t *LiteralInt) WriteExpression() string {
	return fmt.Sprintf("%d", t.value)
}

func (t *LiteralInt) EvalExpression() (values.Value, error) {
  return values.NewLiteralInt(t.value, t.Context()), nil
}

func IsLiteralInt(t Expression) bool {
	_, ok := t.(*LiteralInt)
	return ok
}
