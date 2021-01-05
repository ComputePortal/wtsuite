package glsl

import (
  "github.com/computeportal/wtsuite/pkg/tokens/context"
  "github.com/computeportal/wtsuite/pkg/tokens/glsl/values"
)

type LiteralString struct {
	value string
	LiteralData
}

func NewLiteralString(value string, ctx context.Context) *LiteralString {
	return &LiteralString{value, newLiteralData(ctx)}
}

func (t *LiteralString) EvalExpression() (values.Value, error) {
  panic("should'nt appear anywhere where it might be evaluated")
}

func (t *LiteralString) Value() string {
	return t.value
}
