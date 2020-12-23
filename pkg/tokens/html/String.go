package html

import (
	"fmt"
	"os"
	"reflect"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type String struct {
	wasWord bool // used by ui for special checking
	value   string
	TokenData
}

func NewValueString(value string, ctx context.Context) *String {
	return &String{false, value, TokenData{ctx}}
}

func NewWordString(value string, ctx context.Context) *String {
	return &String{true && VERBOSITY >= 3, value, TokenData{ctx}}
}

func NewString(value string, ctx context.Context) (*String, error) {
	return NewValueString(value, ctx), nil
}

func NewDummyContextString(value string) *String {
	return &String{false, value, TokenData{context.NewDummyContext()}}
}

func (t *String) Value() string {
	return t.value
}

func (t *String) Len() int {
	return len(t.value)
}

func (t *String) Eval(scope Scope) (Token, error) {
	if t.wasWord {
		ctx := t.Context()
		nonWord := &String{false, t.value, TokenData{ctx}}
		res, err := scope.Eval("get", []Token{nonWord, NewNull(ctx)}, ctx)
		if err != nil {
			panic("should never return error")
		}

		if !IsNull(res) {
			errCtx := ctx
			err := errCtx.NewError("Warning: word is also a variable")
			fmt.Fprintf(os.Stderr, err.Error())

			t.wasWord = false // only warn once
		}

		return nonWord, nil
	}

	return t, nil
}

func (t *String) Write() string {
	// without the quotes
	return t.value
}

func (t *String) Dump(indent string) string {
	return indent + "String(" + t.Write() + ")\n"
}

func IsString(t Token) bool {
	_, ok := t.(*String)
	return ok
}

func AssertString(t Token) (*String, error) {
	if s, ok := t.(*String); !ok {
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected string (got " + reflect.TypeOf(t).String() + ")")
		return nil, err
	} else {
		return s, nil
	}
}

func (t *String) InnerContext() context.Context {
	n := len(t.value)
	if n == t.ctx.Len()-2 {
		return t.ctx.NewContext(1, n+1)
	} else {
		return t.TokenData.Context()
	}
}

func (a *String) IsSame(other Token) bool {
	if b, ok := other.(*String); ok {
		return a.value == b.value
	} else {
		return false
	}
}

func (t *String) WasWord() bool {
	return t.wasWord
}
