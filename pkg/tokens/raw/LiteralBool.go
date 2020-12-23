package raw

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type LiteralBool struct {
	value bool
	TokenData
}

func NewValueLiteralBool(b bool, ctx context.Context) *LiteralBool {
	return &LiteralBool{b, TokenData{ctx}}
}

func NewLiteralBool(x interface{}, ctx context.Context) (*LiteralBool, error) {
	b := false
	switch s := x.(type) {
	case bool:
		b = s
	case string:
		switch s {
		case patterns.BOOL_TRUE:
			b = true
		case patterns.BOOL_FALSE:
			b = false
		default:
			return nil, ctx.NewError("Syntax Error: invalid bool")
		}
	default:
		panic("expected string or bool")
	}

	return NewValueLiteralBool(b, ctx), nil
}

func (t *LiteralBool) Value() bool {
	return t.value
}

func (t *LiteralBool) Dump(indent string) string {
	s := "false"
	if t.value {
		s = "true"
	}

	return indent + "LiteralBool(" + s + ")\n"
}

func IsLiteralBool(t Token) bool {
	_, ok := t.(*LiteralBool)
	return ok
}

func AssertLiteralBool(t Token) (*LiteralBool, error) {
	if b, ok := t.(*LiteralBool); ok {
		return b, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected literal bool")
	}
}
