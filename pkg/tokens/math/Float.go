package math

import (
	"fmt"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Float struct {
	v float64
	Word
}

func NewFloat(value float64, ctx context.Context) (*Float, error) {
	strValue := fmt.Sprintf("%g", value)
	return &Float{value, Word{strValue, newTokenData(ctx)}}, nil
}

func (t *Float) Dump(indent string) string {
	return indent + "Float(" + t.value + ")\n"
}

func IsFloat(t Token) bool {
	_, ok := t.(*Float)
	return ok
}
