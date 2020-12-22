package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
	"../patterns"
)

type LiteralString struct {
	value string
	LiteralData
}

func NewLiteralString(value string, ctx context.Context) *LiteralString {
	return &LiteralString{value, newLiteralData(ctx)}
}

func (t *LiteralString) Value() string {
	return t.value
}

func (t *LiteralString) Dump(indent string) string {
	return indent + "LiteralString('" + t.value + "')\n"
}

func (t *LiteralString) WriteExpression() string {
	m1 := patterns.SQ_REGEXP.MatchString(t.value)
	m2 := patterns.DQ_REGEXP.MatchString(t.value)

	if m1 && m2 { // prefer single quotes
		v := strings.Replace(t.value, `"`, `\"`, -1)

		return `'` + v + `'`
	} else if !m1 {
		return `'` + t.value + `'`
	} else {
		return `"` + t.value + `"`
	}
}

func (t *LiteralString) EvalExpression() (values.Value, error) {
	return prototypes.NewLiteralString(t.value, t.Context()), nil
}

// for refactoring
func (t *LiteralString) InnerContext() context.Context {
	n := len(t.value)
	if n == t.ctx.Len()-2 {
		return t.ctx.NewContext(1, n+1)
	} else {
		return t.TokenData.Context()
	}
}

func (t *LiteralString) Walk(fn WalkFunc) error {
  return fn(t)
}
