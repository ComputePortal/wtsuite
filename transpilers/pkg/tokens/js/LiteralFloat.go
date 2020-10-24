package js

import (
	"fmt"

	"./prototypes"
	"./values"

	"../context"
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

func (t *LiteralFloat) EvalExpression(stack values.Stack) (values.Value, error) {
	return values.NewInstance(prototypes.Number, values.NewNumberProperties(true, t.value, t.Context()), t.Context()), nil
}

func (t *LiteralFloat) Walk(fn WalkFunc) error {
  return fn(t)
}
