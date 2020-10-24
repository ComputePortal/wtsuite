package js

import (
	"./prototypes"
	"./values"

	"../context"
)

type LiteralBoolean struct {
	value bool
	LiteralData
}

func NewLiteralBoolean(value bool, ctx context.Context) *LiteralBoolean {
	return &LiteralBoolean{value, newLiteralData(ctx)}
}

func (t *LiteralBoolean) Value() bool {
	return t.value
}

func (t *LiteralBoolean) Dump(indent string) string {
	return indent + "LiteralBoolean(" + t.WriteExpression() + ")\n"
}

func (t *LiteralBoolean) WriteExpression() string {
	s := "false"
	if t.value {
		s = "true"
	}

	return s
}

func (t *LiteralBoolean) EvalExpression(stack values.Stack) (values.Value, error) {
	return values.NewInstance(prototypes.Boolean, values.NewBooleanProperties(true, t.value, t.Context()), t.Context()), nil
}

func (t *LiteralBoolean) Walk(fn WalkFunc) error {
  return fn(t)
}

func IsLiteralTrue(t Expression) bool {
	if lit, ok := t.(*LiteralBoolean); ok {
		return lit.value
	} else {
		return false
	}
}

