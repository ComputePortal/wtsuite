package html

import (
	"strings"

	"../context"
)

type Parens struct {
	values []Token
	alts   []Token // rhs of argDefaults for function and class
	TokenData
}

func NewParens(values []Token, alts []Token, ctx context.Context) *Parens {
	if len(values) != len(alts) {
		panic("expected same lenghts")
	}

	return &Parens{values, alts, TokenData{ctx}}
}

func (t *Parens) Values() []Token {
	return t.values
}

func (t *Parens) Alts() []Token {
	return t.alts
}

func (t *Parens) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Parens(\n")

	for i, v := range t.values {
		b.WriteString(v.Dump(indent + "  "))
		b.WriteString("\n")
		if t.alts[i] != nil {
			b.WriteString(t.alts[i].Dump(indent + "= "))
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (t *Parens) Eval(scope Scope) (Token, error) {
	if len(t.values) != 1 || t.alts[0] != nil {
		errCtx := t.Context()
		err := errCtx.NewError("Error: bad parens (not a function or class declaration)")
		panic(err)
		return nil, err
	}

	return t.values[0].Eval(scope)
}

func (t *Parens) Len() int {
	return len(t.values)
}

func (t *Parens) Loop(fn func(i int, value Token, alt Token) error) error {
	for i, v := range t.values {
		a := t.alts[i]

		if err := fn(i, v, a); err != nil {
			return err
		}
	}

	return nil
}

// only relevant for first token
func (t *Parens) IsSame(other Token) bool {
	if len(t.values) == 0 {
		return false
	}

	return t.values[0].IsSame(other)
}

func IsParens(t Token) bool {
	_, ok := t.(*Parens)
	return ok
}

func AssertParens(t Token) (*Parens, error) {
	p, ok := t.(*Parens)
	if !ok {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected Parens")
	}

	return p, nil
}
