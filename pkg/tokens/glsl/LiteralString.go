package glsl

import (
  "github.com/computeportal/wtsuite/pkg/tokens/context"
)

type LiteralString struct {
	value string
	LiteralData
}

func NewLiteralString(value string, ctx context.Context) *LiteralString {
	return &LiteralString{value, newLiteralData(ctx)}
}

func (t *LiteralString) Value() string {
	return t.value
}
