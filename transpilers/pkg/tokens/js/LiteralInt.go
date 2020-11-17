package js

import (
	"fmt"

	"./prototypes"
	"./values"

	"../context"
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
	return prototypes.NewLiteralInt(t.value, t.Context()), nil
}

func (t *LiteralInt) Walk(fn WalkFunc) error {
  return fn(t)
}

func IsLiteralInt(t Expression) bool {
	_, ok := t.(*LiteralInt)
	return ok
}
