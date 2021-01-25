package html

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type Bool struct {
	value bool
	TokenData
}

func NewValueBool(b bool, ctx context.Context) *Bool {
	return &Bool{b, TokenData{ctx}}
}

func NewBool(x interface{}, ctx context.Context) (*Bool, error) {
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

	return NewValueBool(b, ctx), nil
}

func (t *Bool) Value() bool {
	return t.value
}

func (t *Bool) IsPrimitive() bool {
	return false
}

func (t *Bool) Eval(scope Scope) (Token, error) {
	return t, nil
}

func (t *Bool) EvalLazy(tag FinalTag) (Token, error) {
  return t, nil
}

func (t *Bool) Write() string {
	if t.value {
		return "true"
	} else {
		return "false"
	}
}

func (t *Bool) Dump(indent string) string {
	return indent + "Bool(" + t.Write() + ")\n"
}

func IsBool(t Token) bool {
	_, ok := t.(*Bool)
	return ok
}

func IsFalseBool(t_ Token) bool {
	t, ok := t_.(*Bool)
	if ok {
		return !t.Value()
	} else {
		return false
	}
}

func IsTrueBool(t_ Token) bool {
	t, ok := t_.(*Bool)
	if ok {
		return t.Value()
	} else {
		return false
	}
}

func AssertBool(t Token) (*Bool, error) {
	if b, ok := t.(*Bool); ok {
		return b, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected bool")
	}
}

func (a *Bool) IsSame(other Token) bool {
	if b, ok := other.(*Bool); ok {
		return a.value == b.value
	} else {
		return false
	}
}
