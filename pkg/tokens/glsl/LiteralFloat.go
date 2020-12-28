package glsl

import (
	"fmt"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
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
	return fmt.Sprintf("%g", t.value)
}

func IsLiteralFloat(t Expression) bool {
	_, ok := t.(*LiteralFloat)
	return ok
}
