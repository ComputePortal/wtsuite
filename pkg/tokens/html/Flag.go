package html

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Flag struct {
	TokenData
}

func NewFlag(ctx context.Context) *Flag {
	return &Flag{TokenData{ctx}}
}

func (t *Flag) Eval(scope Scope) (Token, error) {
	return t, nil
}

func (t *Flag) Write() string {
	return ""
}

func (t *Flag) Dump(indent string) string {
	return indent + "Flag\n"
}

func IsFlag(t Token) bool {
	switch s := t.(type) {
	case *Flag:
		return true
	case *String:
		return s.Value() == ""
	default:
		return false
	}
}

func AssertFlag(t Token) error {
	// can also be an empty string
	switch s := t.(type) {
	case *Flag:
		return nil
	case *String:
		if s.Value() == "" {
			return nil
		} else {
			errCtx := t.Context()
			return errCtx.NewError("Syntax Error: expected flag (empty string)")
		}
	default:
		errCtx := t.Context()
		return errCtx.NewError("Syntax Error: expected flag")
	}
}

func (a *Flag) IsSame(other Token) bool {
	return IsFlag(other)
}
