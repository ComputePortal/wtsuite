package math

import (
	"fmt"
	"math"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	INT_X0 = 0.00001
	INT_Y0 = 0.00001

	INT_X01A = 0.0001
	INT_Y01A = -0.0613

	INT_X01B = 0.0162
	INT_Y01B = -0.1146

	INT_X1 = 0.04901
	INT_Y1 = -0.15999

	INT_X12A = 0.0825
	INT_Y12A = -0.2061

	INT_X12B = 0.1260
	INT_Y12B = -0.22899

	INT_X2 = 0.18001 // top most point
	INT_Y2 = -0.22899

	INT_X23A = 0.2095
	INT_Y23A = -0.22899

	INT_X23B = 0.2350
	INT_Y23B = -0.2216

	INT_X3 = 0.25701
	INT_Y3 = -0.20599

	INT_X34A = 0.2794
	INT_Y34A = -0.1903

	INT_X34B = 0.2920
	INT_Y34B = -0.1698

	INT_X4 = 0.2920 // right-most point
	INT_Y4 = -0.1442

	INT_X45A = 0.2920
	INT_Y45A = -0.1306

	INT_X45B = 0.2868
	INT_Y45B = -0.1204

	INT_X5 = 0.27801
	INT_Y5 = -0.11299

	INT_X56A = 0.2704
	INT_Y56A = -0.1062

	INT_X56B = 0.2602
	INT_Y56B = -0.10299

	INT_X6 = 0.24901
	INT_Y6 = -0.10299

	INT_X67A = 0.2381
	INT_Y67A = -0.10299

	INT_X67B = 0.2277
	INT_Y67B = -0.1061

	INT_X7 = 0.21901
	INT_Y7 = -0.11299

	INT_X78A = 0.2114
	INT_Y78A = -0.12

	INT_X78B = 0.20701
	INT_Y78B = -0.1304

	INT_X8 = 0.20701
	INT_Y8 = -0.14399

	INT_X89A = 0.20701
	INT_Y89A = -0.1558

	INT_X89B = 0.2112
	INT_Y89B = -0.1658

	INT_X9 = 0.21801
	INT_Y9 = -0.17299

	INT_X910A = 0.2254
	INT_Y910A = -0.18

	INT_X910B = 0.2344
	INT_Y910B = -0.1844

	INT_X10 = 0.24501 // inflection point
	INT_Y10 = -0.18599

	INT_X1011A = 0.2278
	INT_Y1011A = -0.1997

	INT_X1011B = 0.2060
	INT_Y1011B = -0.20699

	INT_X11 = 0.18001
	INT_Y11 = -0.20699

	INT_X1112A = 0.1154
	INT_Y1112A = -0.20699

	INT_X1112B = 0.08301
	INT_Y1112B = -0.138

	INT_X12 = 0.08301
	INT_Y12 = 0.00001

	INT_Y13 = 0.64201 // for the height of the straight part

	INT_X0_NEXT = 0.35001 // for the space between two integrals

)

// contour related
const (
	INT_CONTOUR_INNER_R = 0.249
	INT_CONTOUR_OUTER_R = 0.274
)

// metrics
const (
	INT_HOR_D_SPACING = 0.1
	INT_HOR_S_SPACING = 0.1
	INT_BOUND_SPACING = 0.05

	INT_VER_CENTER = plusMinusFracYOffset - 0.5*lineThickness
)

type Integral struct {
	integrand Token
	boundA    Token // TODO: actually use these
	boundB    Token
	d         Token // TODO: should the be array?
	contour   bool
	TokenData
}

func GenIntegralPath(scope Scope, x float64, y float64, h float64, ctx context.Context) (boundingbox.BB, error) {
	hMiddle := INT_Y13 - INT_Y0
	hCurve := INT_Y0 - INT_Y2

	if h > 2*hCurve+hMiddle {
		hMiddle = h - 2*hCurve
	}

	x0 := x + (INT_X4 - INT_X12)
	y0 := y - hCurve - hMiddle

	xRef := INT_X0
	yRef := INT_Y0

	fn := func(xIn, yIn float64) (float64, float64) {
		xOut := x0 + (xIn - xRef)
		yOut := y0 + (yIn - yRef)

		return xOut, yOut
	}

	x01a, y01a := fn(INT_X01A, INT_Y01A)
	x01b, y01b := fn(INT_X01B, INT_Y01B)
	x1, y1 := fn(INT_X1, INT_Y1)

	//x12a, y12a := fn(INT_X12A, INT_Y12A)
	x12b, y12b := fn(INT_X12B, INT_Y12B)
	x2, y2 := fn(INT_X2, INT_Y2)

	x23a, y23a := fn(INT_X23A, INT_Y23A)
	x23b, y23b := fn(INT_X23B, INT_Y23B)
	x3, y3 := fn(INT_X3, INT_Y3)

	//x34a, y34a := fn(INT_X34A, INT_Y34A)
	x34b, y34b := fn(INT_X34B, INT_Y34B)
	x4, y4 := fn(INT_X4, INT_Y4)

	x45a, y45a := fn(INT_X45A, INT_Y45A)
	x45b, y45b := fn(INT_X45B, INT_Y45B)
	x5, y5 := fn(INT_X5, INT_Y5)

	//x56a, y56a := fn(INT_X56A, INT_Y56A)
	x56b, y56b := fn(INT_X56B, INT_Y56B)
	x6, y6 := fn(INT_X6, INT_Y6)

	//x67a, y67a := fn(INT_X67A, INT_Y67A)
	x67b, y67b := fn(INT_X67B, INT_Y67B)
	x7, y7 := fn(INT_X7, INT_Y7)

	//x78a, y78a := fn(INT_X78A, INT_Y78A)
	x78b, y78b := fn(INT_X78B, INT_Y78B)
	x8, y8 := fn(INT_X8, INT_Y8)

	//x89a, y89a := fn(INT_X89A, INT_Y89A)
	x89b, y89b := fn(INT_X89B, INT_Y89B)
	x9, y9 := fn(INT_X9, INT_Y9)

	//x910a, y910a := fn(INT_X910A, INT_Y910A)
	x910b, y910b := fn(INT_X910B, INT_Y910B)
	x10, y10 := fn(INT_X10, INT_Y10)

	x1011a, y1011a := fn(INT_X1011A, INT_Y1011A)
	x1011b, y1011b := fn(INT_X1011B, INT_Y1011B)
	x11, y11 := fn(INT_X11, INT_Y11)

	x1112a, y1112a := fn(INT_X1112A, INT_Y1112A)
	x1112b, y1112b := fn(INT_X1112B, INT_Y1112B)
	x12, y12 := fn(INT_X12, INT_Y12)

	y0_ := y0 + hMiddle
	x0_ := x12

	xRef = x0
	yRef = y0

	fn_ := func(xIn, yIn float64) (float64, float64) {
		xOut := x0_ - (xIn - xRef)
		yOut := y0_ - (yIn - yRef)

		return xOut, yOut
	}

	x01a_, y01a_ := fn_(x01a, y01a)
	x01b_, y01b_ := fn_(x01b, y01b)
	x1_, y1_ := fn_(x1, y1)
	x12b_, y12b_ := fn_(x12b, y12b)
	x2_, y2_ := fn_(x2, y2)
	x23a_, y23a_ := fn_(x23a, y23a)
	x23b_, y23b_ := fn_(x23b, y23b)
	x3_, y3_ := fn_(x3, y3)
	x34b_, y34b_ := fn_(x34b, y34b)
	x4_, y4_ := fn_(x4, y4)
	x45a_, y45a_ := fn_(x45a, y45a)
	x45b_, y45b_ := fn_(x45b, y45b)
	x5_, y5_ := fn_(x5, y5)
	x56b_, y56b_ := fn_(x56b, y56b)
	x6_, y6_ := fn_(x6, y6)
	//x67a_, y67a_ := fn_(x67a, y67a)
	x67b_, y67b_ := fn_(x67b, y67b)
	x7_, y7_ := fn_(x7, y7)
	//x78a_, y78a_ := fn_(x78a, y78a)
	x78b_, y78b_ := fn_(x78b, y78b)
	x8_, y8_ := fn_(x8, y8)
	//x89a_, y89a_ := fn_(x89a, y89a)
	x89b_, y89b_ := fn_(x89b, y89b)
	x9_, y9_ := fn_(x9, y9)
	//x910a_, y910a_ := fn_(x910a, y910a)
	x910b_, y910b_ := fn_(x910b, y910b)
	x10_, y10_ := fn_(x10, y10)
	x1011a_, y1011a_ := fn_(x1011a, y1011a)
	x1011b_, y1011b_ := fn_(x1011b, y1011b)
	x11_, y11_ := fn_(x11, y11)
	x1112a_, y1112a_ := fn_(x1112a, y1112a)
	x1112b_, y1112b_ := fn_(x1112b, y1112b)
	x12_, y12_ := fn_(x12, y12)

	var b strings.Builder

	b.WriteString(fmt.Sprintf("M%g %gC%g %g %g %g %g %g", x0, y0, x01a, y01a, x01b, y01b, x1, y1))
	b.WriteString(fmt.Sprintf("S%g %g %g %gC%g %g %g %g %g %g", x12b, y12b, x2, y2, x23a, y23a, x23b, y23b, x3, y3))
	b.WriteString(fmt.Sprintf("S%g %g %g %gC%g %g %g %g %g %g", x34b, y34b, x4, y4, x45a, y45a, x45b, y45b, x5, y5))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x56b, y56b, x6, y6, x67b, y67b, x7, y7))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x78b, y78b, x8, y8, x89b, y89b, x9, y9))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x910b, y910b, x10, y10))
	b.WriteString(fmt.Sprintf("C%g %g %g %g %g %g", x1011a, y1011a, x1011b, y1011b, x11, y11))
	b.WriteString(fmt.Sprintf("C%g %g %g %g %g %g", x1112a, y1112a, x1112b, y1112b, x12, y12))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x0_, y0_, x01a_, y01a_, x01b_, y01b_, x1_, y1_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gC%g %g %g %g %g %g", x12b_, y12b_, x2_, y2_, x23a_, y23a_, x23b_, y23b_, x3_, y3_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gC%g %g %g %g %g %g", x34b_, y34b_, x4_, y4_, x45a_, y45a_, x45b_, y45b_, x5_, y5_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x56b_, y56b_, x6_, y6_, x67b_, y67b_, x7_, y7_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x78b_, y78b_, x8_, y8_, x89b_, y89b_, x9_, y9_))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x910b_, y910b_, x10_, y10_))
	b.WriteString(fmt.Sprintf("C%g %g %g %g %g %g", x1011a_, y1011a_, x1011b_, y1011b_, x11_, y11_))
	b.WriteString(fmt.Sprintf("C%g %g %g %g %g %gZ", x1112a_, y1112a_, x1112b_, y1112b_, x12_, y12_))

	if err := scope.BuildMathPath(b.String(), ctx); err != nil {
		return nil, err
	}

	return boundingbox.NewBB(x4_, y2, x4, y2_), nil
}

// same inputs as integral path!
func GenContourPath(scope Scope, x float64, y float64, h float64, ctx context.Context) (boundingbox.BB, error) {
	hMiddle := INT_Y13 - INT_Y0
	hCurve := INT_Y0 - INT_Y2

	if h > 2*hCurve+hMiddle {
		hMiddle = h - 2*hCurve
	}

	thickness := INT_CONTOUR_OUTER_R - INT_CONTOUR_INNER_R
	// top point of circle
	x0 := x + (INT_X4 - INT_X12) + 0.5*(INT_X12-INT_X0)
	y0 := y - hCurve - hMiddle*0.5 - INT_CONTOUR_OUTER_R

	x1 := x0
	y1 := y0 + 2*INT_CONTOUR_OUTER_R

	// inner points
	x2 := x0
	y2 := y0 + thickness

	x3 := x1
	y3 := y1 - thickness

	var b strings.Builder

	r01 := INT_CONTOUR_OUTER_R
	b.WriteString(fmt.Sprintf("M%g %gA%g %g 0 0 0 %g %g", x0, y0, r01, r01, x1, y1))
	b.WriteString(fmt.Sprintf("A%g %g 0 0 0 %g %gZ", r01, r01, x0, y0))

	r23 := INT_CONTOUR_INNER_R
	b.WriteString(fmt.Sprintf("M%g %gA%g %g 0 0 1 %g %g", x2, y2, r23, r23, x3, y3))
	b.WriteString(fmt.Sprintf("A%g %g 0 0 1 %g %g", r23, r23, x2, y2))

	if err := scope.BuildMathPath(b.String(), ctx); err != nil {
		return nil, err
	}

	return boundingbox.NewBB(x0-r01, y0, x0+r01, y1), nil
}

func NewIntegral(integrand Token, boundA Token, boundB Token, d Token, ctx context.Context) (*Integral, error) {
	return &Integral{
		integrand,
		boundA,
		boundB,
		d,
		false,
		newTokenData(ctx),
	}, nil
}

func NewContourIntegral(integrand Token, d Token, ctx context.Context) (*Integral, error) {
	return &Integral{
		integrand,
		nil,
		nil,
		d,
		true,
		newTokenData(ctx),
	}, nil
}

func (t *Integral) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Integral(")
	if t.boundA != nil {
		b.WriteString(t.boundA.Dump(""))
	}
	b.WriteString(" ")
	if t.boundB != nil {
		b.WriteString(t.boundB.Dump(""))
	}
	b.WriteString(")\n")
	b.WriteString(t.integrand.Dump(indent + "  "))
	b.WriteString("\n")
	b.WriteString(t.d.Dump("d:"))

	return b.String()
}

func (t *Integral) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {

	subScope := scope.NewSubScope()
	bbIntegrand, err := t.integrand.GenerateTags(subScope, 0, y)
	if err != nil {
		return nil, err
	}

	bbD, err := t.d.GenerateTags(subScope, bbIntegrand.Right()+INT_HOR_D_SPACING, y)
	if err != nil {
		return nil, err
	}

	bbContent := boundingbox.Merge(bbIntegrand, bbD)

	symbolH := 2.0 * max(bbContent.Bottom()-y, y-bbContent.Top())
	symbolY := y + 0.5*symbolH - INT_VER_CENTER

	bbSymbol, err := GenIntegralPath(scope, x, symbolY, symbolH, t.Context())
	if err != nil {
		return nil, err
	}

	if t.contour {
		if _, err := GenContourPath(scope, x, symbolY, symbolH, t.Context()); err != nil {
			return nil, err
		}
	}

	boundWidth := 0.0
	boundStartX := bbSymbol.Right() + INT_BOUND_SPACING
	if t.boundA != nil {
		boundASubScope := scope.NewSubScope()
		bbBoundA, err := t.boundA.GenerateTags(boundASubScope, 0.0, 0.0)
		if err != nil {
			return nil, err
		}

		boundASubScope.Transform(boundStartX,
			bbSymbol.Bottom()-indexSubscriptOffset,
			indexLevelDecrease, indexLevelDecrease)

		boundWidth = bbBoundA.Width() + INT_BOUND_SPACING
	}

	if t.boundB != nil {
		boundBSubScope := scope.NewSubScope()
		bbBoundB, err := t.boundB.GenerateTags(boundBSubScope, 0.0, 0.0)
		if err != nil {
			return nil, err
		}

		boundBSubScope.Transform(boundStartX, bbSymbol.Top()+indexLevelDecrease-powExponentOffset,
			indexLevelDecrease, indexLevelDecrease)

		boundWidth = math.Max(boundWidth, bbBoundB.Width()+INT_BOUND_SPACING)
	}

	contentStartX := boundStartX + boundWidth

	subScope.Transform(contentStartX, 0.0, 1.0, 1.0)
	bbContent = bbContent.Translate(contentStartX, 0)

	return boundingbox.Merge(bbContent, bbSymbol), nil
}
