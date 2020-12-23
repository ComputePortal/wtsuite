package raw

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type SpecialNumber struct {
	value string
	TokenData
}

func NewSpecialNumber(value string, ctx context.Context) *SpecialNumber {
	return &SpecialNumber{value, TokenData{ctx}}
}

func (t *SpecialNumber) Dump(indent string) string {
	return indent + "SpecialNumber(" + t.value + ")\n"
}

func (t *SpecialNumber) Value() string {
	return t.value
}

func IsSpecialNumber(t Token) bool {
	_, ok := t.(*SpecialNumber)
	return ok
}

func AssertSpecialNumber(t Token) (*SpecialNumber, error) {
	if sn, ok := t.(*SpecialNumber); ok {
		return sn, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected special number")
	}
}
