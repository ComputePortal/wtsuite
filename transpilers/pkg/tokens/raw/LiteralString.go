package raw

import (
	"../context"
)

type LiteralString struct {
	wasWord bool
	value   string
	TokenData
}

func NewValueLiteralString(value string, ctx context.Context) *LiteralString {
	return &LiteralString{false, value, TokenData{ctx}}
}

// used by ui
func NewWordLiteralString(value string, ctx context.Context) *LiteralString {
	return &LiteralString{true, value, TokenData{ctx}}
}

func NewLiteralString(value string, ctx context.Context) (*LiteralString, error) {
	return NewValueLiteralString(value, ctx), nil
}

func (t *LiteralString) Value() string {
	return t.value
}

func (t *LiteralString) Dump(indent string) string {
	return indent + "LiteralString(" + t.value + ")\n"
}

func IsLiteralString(t Token) bool {
	_, ok := t.(*LiteralString)
	return ok
}

func AssertLiteralString(t Token) (*LiteralString, error) {
	if s, ok := t.(*LiteralString); !ok {
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected literal string")
		return nil, err
	} else {
		return s, nil
	}
}

func (t *LiteralString) InnerContext() context.Context {
	n := len(t.value)
	if n == t.ctx.Len()-2 {
		return t.ctx.NewContext(1, n+1)
	} else {
		return t.TokenData.Context()
	}
}

func (t *LiteralString) WasWord() bool {
	return t.wasWord
}
