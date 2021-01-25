package html

import (
	"fmt"
	"strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// distinct from Float because it is used for indexing in lists
type Int struct {
	value int
	TokenData
}

func NewValueInt(v int, ctx context.Context) *Int {
	return &Int{v, TokenData{ctx}}
}

func NewInt(x interface{}, ctx context.Context) (*Int, error) {
	switch v := x.(type) {
	case string:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, ctx.NewError("Syntax Error: invalid int")
		}
		return NewValueInt(int(value), ctx), nil
	case int:
		return NewValueInt(v, ctx), nil
	default:
		panic("expected string or int")
	}
}

func (t *Int) Value() int {
	return t.value
}

func (t *Int) Eval(scope Scope) (Token, error) {
	return t, nil
}

func (t *Int) EvalLazy(tag FinalTag) (Token, error) {
	return t, nil
}

func (t *Int) Write() string {
	return fmt.Sprintf("%d", t.value)
}

func (t *Int) Dump(indent string) string {
	return indent + "Int(" + t.Write() + ")\n"
}

func IsInt(t Token) bool {
	_, ok := t.(*Int)
	return ok
}

func AssertInt(t Token) (*Int, error) {
	if i, ok := t.(*Int); !ok {
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected int")
		return nil, err
	} else {
		return i, nil
	}
}

func (a *Int) IsSame(other Token) bool {
	if b, ok := other.(*Int); ok {
		return a.value == b.value
	} else {
		return false
	}
}
