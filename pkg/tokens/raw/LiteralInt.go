package raw

import (
	"fmt"
	"strconv"
	"strings"

	"../context"
)

// distinct from Float because it is used for indexing in lists
type LiteralInt struct {
	value int
	TokenData
}

func NewValueLiteralInt(v int, ctx context.Context) *LiteralInt {
	return &LiteralInt{v, TokenData{ctx}}
}

func NewLiteralInt(x interface{}, ctx context.Context) (*LiteralInt, error) {
	switch v := x.(type) {
	case string:
		value, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, ctx.NewError("Syntax Error: invalid literal int")
		}
		return NewValueLiteralInt(int(value), ctx), nil
	case int:
		return NewValueLiteralInt(v, ctx), nil
	default:
		panic("expected string or int")
	}
}

func NewHexLiteralInt(x string, ctx context.Context) (*LiteralInt, error) {
	x = strings.Replace(x, "0x", "", 1)
	value, err := strconv.ParseInt(x, 16, 64)
	if err != nil {
		return nil, ctx.NewError("Syntax Error: invalid literal hex int")
	}

	return NewValueLiteralInt(int(value), ctx), nil
}

func (t *LiteralInt) Value() int {
	return t.value
}

func (t *LiteralInt) Dump(indent string) string {
	s := fmt.Sprintf("%d", t.value)
	return indent + "LiteralInt(" + s + ")\n"
}

func IsLiteralInt(t Token) bool {
	_, ok := t.(*LiteralInt)
	return ok
}

func AssertLiteralInt(t Token) (*LiteralInt, error) {
	if i, ok := t.(*LiteralInt); !ok {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected literal int")
	} else {
		return i, nil
	}
}
