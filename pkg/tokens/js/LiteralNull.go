package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type LiteralNull struct {
  LiteralData
}

func NewLiteralNull(ctx context.Context) *LiteralNull {
	return &LiteralNull{newLiteralData(ctx)}
}

func (t *LiteralNull) Dump(indent string) string {
	return indent + t.WriteExpression() + "\n"
}

func (t *LiteralNull) WriteExpression() string {
	return "null"
}

func (t *LiteralNull) EvalExpression() (values.Value, error) {
	return values.NewAll(t.Context()), nil
}

func (t *LiteralNull) Walk(fn WalkFunc) error {
  return fn(t)
}
