package html

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Null struct {
	TokenData
}

func NewNull(ctx context.Context) *Null {
	return &Null{TokenData{ctx}}
}

func (t *Null) Eval(scope Scope) (Token, error) {
	return t, nil
}

func (t *Null) Dump(indent string) string {
	return indent + "Null\n"
}

func IsNull(t Token) bool {
	switch t.(type) {
	case *Null:
		return true
	default:
		return false
	}
}

func (a *Null) IsSame(other Token) bool {
	return IsNull(other)
}
