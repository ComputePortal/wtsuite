package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// eg. for key of member
// essentially a string-context pair
type Word struct {
	value string
	TokenData
}

func NewWord(value string, ctx context.Context) *Word {
	return &Word{value, TokenData{ctx}}
}

func (t *Word) Value() string {
	return t.value
}

func (t *Word) Dump(indent string) string {
	return indent + "Word(" + t.value + ")\n"
}

func AssertWord(t Token) (*Word, error) {
	if w, ok := t.(*Word); ok {
		return w, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected word")
	}
}
