package math

import (
	"strings"

	"./boundingbox"

	"../context"
)

const (
	braceToExprSpacing = 0.5                                      // horizontal, inner
	caseVerSpacing     = 0.1                                      // vertical
	braceYCenter       = plusMinusFracYOffset - 0.5*lineThickness // vertical, in order to center the inflection of the brace
	innerYCenter       = 0.1
	exprToCondSpacing  = 1.5 // horizontal
)

type IfElse struct {
	align *Align // used during CalcBB and GenerateTags
	TokenData
}

func NewIfElse(conds []Token, exprs []Token, ctx context.Context) (*IfElse, error) {
	if len(exprs) != len(conds) {
		panic("should've been caught before")
	}

	eqs := make([][]Token, 0)
	for i, expr := range exprs {
		eqs = append(eqs, []Token{expr, conds[i]})
	}

	return &IfElse{
		newAlign(exprToCondSpacing, caseVerSpacing, false, true, eqs, ctx),
		newTokenData(ctx),
	}, nil
}

func newIfElse(exprs []Token, conds []Token, ctx context.Context) *IfElse {
	ie, err := NewIfElse(exprs, conds, ctx)
	if err != nil {
		panic(err)
	}

	return ie
}

func (t *IfElse) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("IfElse\n")

	b.WriteString(t.align.Dump("| "))

	return b.String()
}

func (t *IfElse) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	subScope := scope.NewSubScope()

	bbAlign, err := t.align.GenerateTags(subScope, 0.0, 0.0)
	if err != nil {
		return nil, err
	}

	// TODO: generate brace

	h := bbAlign.Height()
	bbBrace, err := GenLeftBrace(scope, x, y-braceYCenter+0.5*h, h, t.Context())

	dx := bbBrace.Right() + braceToExprSpacing
	dy := y - braceYCenter - bbAlign.Bottom() + h*0.5

	subScope.Transform(dx, dy, 1.0, 1.0)
	bbAlign = bbAlign.Translate(dx, dy)

	return boundingbox.Merge(bbBrace, bbAlign), nil
}
