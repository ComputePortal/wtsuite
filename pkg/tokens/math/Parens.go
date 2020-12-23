package math

import (
	"fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	PARENS_X0 = 0.00047
	PARENS_Y0 = -0.00205

	PARENS_X01A = 0.00047
	PARENS_Y01A = -0.0496

	PARENS_X01B = 0.0071
	PARENS_Y01B = -0.0955

	PARENS_X1 = 0.01847
	PARENS_Y1 = -0.13805

	PARENS_X12A = 0.0311
	PARENS_Y12A = -0.1816

	PARENS_X12B = 0.0456
	PARENS_Y12B = -0.2173

	PARENS_X2 = 0.06247
	PARENS_Y2 = -0.24505

	PARENS_X23A = 0.0795
	PARENS_Y23A = -0.2738

	PARENS_X23B = 0.1003
	PARENS_Y23B = -0.3011

	PARENS_X3 = 0.12547
	PARENS_Y3 = -0.32705

	PARENS_X34A = 0.1513
	PARENS_Y34A = -0.3542

	PARENS_X34B = 0.1725
	PARENS_Y34B = -0.3734

	PARENS_X4 = 0.18947
	PARENS_Y4 = -0.38605

	PARENS_X45A = 0.2062
	PARENS_Y45A = -0.3991

	PARENS_X45B = 0.2257
	PARENS_Y45B = -0.4122

	PARENS_X5 = 0.24747
	PARENS_Y5 = -0.42605

	PARENS_X6 = 0.25647
	PARENS_Y6 = -0.41005

	PARENS_X67A = 0.2269
	PARENS_Y67A = -0.3863

	PARENS_X67B = 0.2034
	PARENS_Y67B = -0.3641

	PARENS_X7 = 0.1870
	PARENS_Y7 = -0.34505

	PARENS_X78A = 0.1675
	PARENS_Y78A = -0.3256

	PARENS_X78B = 0.1503
	PARENS_Y78B = -0.3005

	PARENS_X8 = 0.13347
	PARENS_Y8 = -0.26905

	PARENS_X89A = 0.1167
	PARENS_Y89A = -0.2384

	PARENS_X89B = 0.1047
	PARENS_Y89B = -0.2015

	PARENS_X9 = 0.09747
	PARENS_Y9 = -0.15805

	PARENS_X910A = 0.0902
	PARENS_Y910A = -0.1155

	PARENS_X910B = 0.08647
	PARENS_Y910B = -0.0640

	PARENS_X10 = 0.08647
	PARENS_Y10 = -0.00205 // change to 0

	parensYCenter      = 0.220
	extraParensSpacing = 0.1 // horizontal, inner, based on content level
)

type Parens struct {
	content Token
	left    *Symbol
	right   *Symbol
	TokenData
}

func NewParens(content Token, ctx context.Context) (Token, error) {
	left := newSymbol("(", ctx)

	right := newSymbol(")", ctx)

	return &Parens{content, left, right, newTokenData(ctx)}, nil
}

func (t *Parens) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Parens()\n")
	b.WriteString(t.content.Dump(indent + "  "))

	return b.String()
}

func (t *Parens) genParensPath(scope Scope, x float64, y float64, h float64, flipX bool) (boundingbox.BB, error) {
	hCurve := (PARENS_Y0 - PARENS_Y5)

	hMiddle := 0.0
	if h > 2*hCurve {
		hMiddle = h - 2*hCurve
	}

	x0 := x
	fX := 1.0
	if flipX {
		x0 = x + (PARENS_X6 - PARENS_X0)
		fX = -1.0
	}

	y0 := y - hCurve - hMiddle

	xRef := PARENS_X0
	yRef := PARENS_Y0

	x01a := x0 + fX*(PARENS_X01A-xRef)
	y01a := y0 + (PARENS_Y01A - yRef)

	x01b := x0 + fX*(PARENS_X01B-xRef)
	y01b := y0 + (PARENS_Y01B - yRef)

	x1 := x0 + fX*(PARENS_X1-xRef)
	y1 := y0 + (PARENS_Y1 - yRef)

	x12a := x0 + fX*(PARENS_X12A-xRef)
	y12a := y0 + (PARENS_Y12A - yRef)

	x12b := x0 + fX*(PARENS_X12B-xRef)
	y12b := y0 + (PARENS_Y12B - yRef)

	x2 := x0 + fX*(PARENS_X2-xRef)
	y2 := y0 + (PARENS_Y2 - yRef)

	x23a := x0 + fX*(PARENS_X23A-xRef)
	y23a := y0 + (PARENS_Y23A - yRef)

	x23b := x0 + fX*(PARENS_X23B-xRef)
	y23b := y0 + (PARENS_Y23B - yRef)

	x3 := x0 + fX*(PARENS_X3-xRef)
	y3 := y0 + (PARENS_Y3 - yRef)

	x34a := x0 + fX*(PARENS_X34A-xRef)
	y34a := y0 + (PARENS_Y34A - yRef)

	x34b := x0 + fX*(PARENS_X34B-xRef)
	y34b := y0 + (PARENS_Y34B - yRef)

	x4 := x0 + fX*(PARENS_X4-xRef)
	y4 := y0 + (PARENS_Y4 - yRef)

	x45a := x0 + fX*(PARENS_X45A-xRef)
	y45a := y0 + (PARENS_Y45A - yRef)

	x45b := x0 + fX*(PARENS_X45B-xRef)
	y45b := y0 + (PARENS_Y45B - yRef)

	x5 := x0 + fX*(PARENS_X5-xRef)
	y5 := y0 + (PARENS_Y5 - yRef)

	x6 := x0 + fX*(PARENS_X6-xRef)
	y6 := y0 + (PARENS_Y6 - yRef)

	x67a := x0 + fX*(PARENS_X67A-xRef)
	y67a := y0 + (PARENS_Y67A - yRef)

	x67b := x0 + fX*(PARENS_X67B-xRef)
	y67b := y0 + (PARENS_Y67B - yRef)

	x7 := x0 + fX*(PARENS_X7-xRef)
	y7 := y0 + (PARENS_Y7 - yRef)

	x78a := x0 + fX*(PARENS_X78A-xRef)
	y78a := y0 + (PARENS_Y78A - yRef)

	x78b := x0 + fX*(PARENS_X78B-xRef)
	y78b := y0 + (PARENS_Y78B - yRef)

	x8 := x0 + fX*(PARENS_X8-xRef)
	y8 := y0 + (PARENS_Y8 - yRef)

	x89a := x0 + fX*(PARENS_X89A-xRef)
	y89a := y0 + (PARENS_Y89A - yRef)

	x89b := x0 + fX*(PARENS_X89B-xRef)
	y89b := y0 + (PARENS_Y89B - yRef)

	x9 := x0 + fX*(PARENS_X9-xRef)
	y9 := y0 + (PARENS_Y9 - yRef)

	x910a := x0 + fX*(PARENS_X910A-xRef)
	y910a := y0 + (PARENS_Y910A - yRef)

	x910b := x0 + fX*(PARENS_X910B-xRef)
	y910b := y0 + (PARENS_Y910B - yRef)

	x10 := x0 + fX*(PARENS_X10-xRef)
	y10 := y0 + (PARENS_Y10 - yRef)

	// the other side is has the same
	yRef = y0
	y0_ := y - hCurve

	y01a_ := y0_ - (y01a - yRef)
	y1_ := y0_ - (y1 - yRef)
	y12a_ := y0_ - (y12a - yRef)
	y2_ := y0_ - (y2 - yRef)
	y23a_ := y0_ - (y23a - yRef)
	y3_ := y0_ - (y3 - yRef)
	y34a_ := y0_ - (y34a - yRef)
	y4_ := y0_ - (y4 - yRef)
	y45a_ := y0_ - (y45a - yRef)
	y45b_ := y0_ - (y45b - yRef)
	y5_ := y0_ - (y5 - yRef)
	y6_ := y0_ - (y6 - yRef)
	y67a_ := y0_ - (y67a - yRef)
	//y67b_ := y0_ - (y67b - yRef)
	y7_ := y0_ - (y7 - yRef)
	y78a_ := y0_ - (y78a - yRef)
	y8_ := y0_ - (y8 - yRef)
	y89a_ := y0_ - (y89a - yRef)
	y9_ := y0_ - (y9 - yRef)
	y910a_ := y0_ - (y910a - yRef)
	y910b_ := y0_ - (y910b - yRef)
	y10_ := y0_ - (y10 - yRef)

	var b strings.Builder

	b.WriteString(fmt.Sprintf("M%g %gC%g %g %g %g %g %g", x0, y0, x01a, y01a, x01b, y01b, x1, y1))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x12b, y12b, x2, y2, x23b, y23b, x3, y3))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x34b, y34b, x4, y4, x45b, y45b, x5, y5))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x6, y6, x67a, y67a, x67b, y67b, x7, y7))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x78b, y78b, x8, y8, x89b, y89b, x9, y9))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x910b, y910b, x10, y10))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x10, y10_, x910b, y910b_, x910a, y910a_, x9, y9_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x89a, y89a_, x8, y8_, x78a, y78a_, x7, y7_))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x67a, y67a_, x6, y6_))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x5, y5_, x45b, y45b_, x45a, y45a_, x4, y4_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x34a, y34a_, x3, y3_, x23a, y23a_, x2, y2_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %gZ", x12a, y12a_, x1, y1_, x01a, y01a_, x0, y0_))

	if err := scope.BuildMathPath(b.String(), t.Context()); err != nil {
		return nil, err
	}

	if flipX {
		return boundingbox.NewBB(x6, y5, x0, y5_), nil
	} else {
		return boundingbox.NewBB(x0, y5, x6, y5_), nil
	}
}

func (t *Parens) genLeftPath(scope Scope, x float64, y float64, h float64) (boundingbox.BB, error) {
	return t.genParensPath(scope, x, y, h, false)
}

func (t *Parens) genRightPath(scope Scope, x float64, y float64, h float64) (boundingbox.BB, error) {
	return t.genParensPath(scope, x, y, h, true)
}

func (t *Parens) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	// if content is frac return that directly
	if _, ok := t.content.(*Frac); ok {
		return t.content.GenerateTags(scope, x, y)
	}

	subScope := scope.NewSubScope()

	bbContent, err := t.content.GenerateTags(subScope, 0, y)
	if err != nil {
		return nil, err
	}

	bbLeft, err := t.genLeftPath(scope, x, bbContent.Bottom(), bbContent.Height())
	if err != nil {
		return nil, err
	}

	dxContent := bbLeft.Right() + extraParensSpacing
	subScope.Transform(dxContent, 0.0, 1.0, 1.0)
	bbContent = bbContent.Translate(dxContent, y) // XXX: do we really the y translation here?

	bbRight, err := t.genRightPath(scope, bbContent.Right()+extraParensSpacing, bbContent.Bottom(), bbContent.Height())
	if err != nil {
		return nil, err
	}

	return boundingbox.Merge(bbLeft, bbContent, bbRight), nil
}
