package html

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Function struct {
	name string
	args []Token
	TokenData
}

func NewValueFunction(name string, args []Token, ctx context.Context) *Function {
	return NewFunction(name, args, ctx)
}

func NewFunction(name string, args []Token, ctx context.Context) *Function {
	return &Function{name, args, TokenData{ctx}}
}

func (t *Function) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent + "Function " + t.name + "\n")

	for _, arg := range t.args {
		b.WriteString(arg.Dump(indent + "  "))
	}

	return b.String()
}

func IsAnyFunction(t Token) bool {
	_, ok := t.(*Function)
	return ok
}

func IsFunction(t Token, name string) bool {
	if fn, ok := t.(*Function); ok {
		return fn.name == name
	}

	return false
}

func AssertFunction(t Token) (*Function, error) {
	if fn, ok := t.(*Function); ok {
		return fn, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected a function")
	}
}

func (t *Function) Eval(scope Scope) (Token, error) {
	res, err := scope.Eval(t.name, t.args, t.Context())
	if err != nil {
		return nil, err
	}

	if _, ok := res.(*Function); ok {
		panic("result of an eval can't be another function")
	}

	return res, nil
}

func (a *Function) IsSame(other Token) bool {
	if b, ok := other.(*Function); ok {
		if a.name == b.name {
			if len(a.args) == len(b.args) {
				for i, _ := range a.args {
					if !a.args[i].IsSame(b.args[i]) {
						return false
					}
				}

				return true
			}
		}
	}

	return false
}

func (t *Function) Args() []Token {
	return t.args[:]
}
