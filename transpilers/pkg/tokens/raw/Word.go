package raw

import (
	"../context"
	"../patterns"
)

// use Symbol for KeyWords!
type Word struct {
	value string
	TokenData
}

func NewValueWord(value string, ctx context.Context) *Word {
	return &Word{value, TokenData{ctx}}
}

func NewWord(value string, ctx context.Context) (*Word, error) {
	return NewValueWord(value, ctx), nil
}

func (t *Word) Value() string { // this function is not named 'String()' to avoid confusion with the tokens.String type
	return t.value
}

func IsAnyWord(t Token) bool {
	if _, ok := t.(*Word); ok {
		return true
	} else if s, ok := t.(*Symbol); ok {
		return patterns.IsSimpleWord(s.Value())
	} else {
		return false
	}
}

func IsWord(t Token, s string) bool {
	if w, ok := t.(*Word); ok {
		return w.value == s
	} else if sym, ok := t.(*Symbol); ok {
		if patterns.IsSimpleWord(sym.Value()) {
			return sym.value == s
		}
	}

	return false
}

func AssertWord(t Token) (*Word, error) {
	if w, ok := t.(*Word); ok {
		return w, nil
	} else if s, ok := t.(*Symbol); ok {
		if patterns.IsSimpleWord(s.Value()) {
			return NewWord(s.Value(), s.Context())
		}
	}

	errCtx := t.Context()
	err := errCtx.NewError("Error: expected a word")
	return nil, err
}

func AssertNotWord(t Token) error {
	if _, ok := t.(*Word); ok {
		errCtx := t.Context()
		return errCtx.NewError("Error: didn't expect a word")
	} else if s, ok := t.(*Symbol); ok && patterns.IsSimpleWord(s.Value()) {
		errCtx := t.Context()
		return errCtx.NewError("Error: didn't expect a word (wordlike-symbol actually)")
	}

	return nil
}

func (t *Word) Dump(indent string) string {
	return indent + "Word(" + t.Value() + ")\n"
}
