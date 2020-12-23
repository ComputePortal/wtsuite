package math

import (
	"fmt"
	"math"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

const (
	SUM_X0 = 0.00001
	SUM_Y0 = 0.00003

	SUM_X1 = 0.24201
	SUM_Y1 = -0.40797

	SUM_X2 = 0.00901
	SUM_Y2 = -0.87497

	SUM_X3 = 0.60801
	SUM_Y3 = -0.87497

	SUM_X4 = 0.61801
	SUM_Y4 = -0.70997

	SUM_X5 = 0.59201
	SUM_Y5 = -0.70997

	SUM_X56A = 0.5873
	SUM_Y56A = -0.7346

	SUM_X56B = 0.5827
	SUM_Y56B = -0.7533

	SUM_X6 = 0.57801
	SUM_Y6 = -0.7635

	SUM_X67A = 0.5734
	SUM_Y67A = -0.7793

	SUM_X67B = 0.5654
	SUM_Y67B = -0.7916

	SUM_X7 = 0.55401
	SUM_Y7 = -0.80297

	SUM_X78A = 0.5434
	SUM_Y78A = -0.8150

	SUM_X78B = 0.5280
	SUM_Y78B = -0.8230

	SUM_X8 = 0.5094
	SUM_Y8 = -0.8237

	SUM_X89A = 0.4887
	SUM_Y89A = -0.8316

	SUM_X89B = 0.4627
	SUM_Y89B = -0.8340

	SUM_X9 = 0.43001
	SUM_Y9 = -0.83397

	SUM_X10 = 0.16201
	SUM_Y10 = -0.83397

	SUM_X11 = 0.35001
	SUM_Y11 = -0.46397

	SUM_X12 = 0.14901
	SUM_Y12 = -0.11197

	SUM_X13 = 0.47701
	SUM_Y13 = -0.11197

	SUM_X1314A = 0.5343
	SUM_Y1314A = -0.1120

	SUM_X1314B = 0.5740
	SUM_Y1314B = -0.1180

	SUM_X14 = 0.59601
	SUM_Y14 = -0.12997

	SUM_X1415A = 0.6187
	SUM_Y1415A = -0.1426

	SUM_X1415B = 0.6383
	SUM_Y1415B = -0.1693

	SUM_X15 = 0.65501
	SUM_Y15 = -0.20997

	SUM_X16 = 0.68101
	SUM_Y16 = -0.20997

	SUM_X17 = 0.62701
	SUM_Y17 = 0.00003
)

const (
	SUM_HOR_SPACING       = 0.1
	SUM_VER_SPACING       = 0.1
	SUM_LOWER_UPPER_SCALE = 0.6
)

type Sum struct {
	rhs   Token
	lower Token
	upper Token
	TokenData
}

// simply vertical scaling is ugly, we must keep the angle
func GenSumPath(scope Scope, x float64, y float64, w float64, h float64, ctx context.Context) (boundingbox.BB, error) {
	n01x := SUM_X1 - SUM_X0
	n01y := SUM_Y1 - SUM_Y0

	n12x := SUM_X2 - SUM_X1
	n12y := SUM_Y2 - SUM_Y1

	hMin := SUM_Y0 - SUM_Y2
	hExtra := 0.0
	if h > hMin {
		hExtra = h - hMin
	}

	wExtra1 := hExtra / (math.Abs(n01y/n01x) + math.Abs(n12y/n12x))
	hExtra1 := wExtra1 * math.Abs(n01y/n01x)

	wMin := SUM_X16 - SUM_X0
	wExtra := wExtra1
	if w > wMin+wExtra1 {
		wExtra += w - wMin - wExtra1
	}

	x0 := x
	y0 := y

	xRef := SUM_X0
	yRef := SUM_Y0

	fn := func(xIn, yIn float64) (float64, float64) {
		xOut := x0 + (xIn - xRef)
		yOut := y0 + (yIn - yRef)

		return xOut, yOut
	}

	x1, y1 := fn(SUM_X1+wExtra1, SUM_Y1-hExtra1)
	x2, y2 := fn(SUM_X2, SUM_Y2-hExtra)
	x3, y3 := fn(SUM_X3+wExtra, SUM_Y3-hExtra)
	x4, y4 := fn(SUM_X4+wExtra, SUM_Y4-hExtra)
	x5, y5 := fn(SUM_X5+wExtra, SUM_Y5-hExtra)
	x56a, y56a := fn(SUM_X56A+wExtra, SUM_Y56A-hExtra)
	x56b, y56b := fn(SUM_X56B+wExtra, SUM_Y56B-hExtra)
	x6, y6 := fn(SUM_X6+wExtra, SUM_Y6-hExtra)
	x67b, y67b := fn(SUM_X67B+wExtra, SUM_Y67B-hExtra)
	x7, y7 := fn(SUM_X7+wExtra, SUM_Y7-hExtra)
	x78b, y78b := fn(SUM_X78B+wExtra, SUM_Y78B-hExtra)
	x8, y8 := fn(SUM_X8+wExtra, SUM_Y8-hExtra)
	x89b, y89b := fn(SUM_X89B+wExtra, SUM_Y89B-hExtra)
	x9, y9 := fn(SUM_X9+wExtra, SUM_Y9-hExtra)

	x10, y10 := fn(SUM_X10, SUM_Y10-hExtra)
	x11, y11 := fn(SUM_X11+wExtra1, SUM_Y11-hExtra1)
	x12, y12 := fn(SUM_X12, SUM_Y12)
	x13, y13 := fn(SUM_X13+wExtra, SUM_Y13)
	x1314a, y1314a := fn(SUM_X1314A+wExtra, SUM_Y1314A)
	x1314b, y1314b := fn(SUM_X1314B+wExtra, SUM_Y1314B)
	x14, y14 := fn(SUM_X14+wExtra, SUM_Y14)
	x1415b, y1415b := fn(SUM_X1415B+wExtra, SUM_Y1415B)
	x15, y15 := fn(SUM_X15+wExtra, SUM_Y15)
	x16, y16 := fn(SUM_X16+wExtra, SUM_Y16)
	x17, y17 := fn(SUM_X17+wExtra, SUM_Y17)

	var b strings.Builder

	b.WriteString(fmt.Sprintf("M%g %gL%g %gL%g %g", x0, y0, x1, y1, x2, y2))
	b.WriteString(fmt.Sprintf("L%g %gL%g %gL%g %g", x3, y3, x4, y4, x5, y5))
	b.WriteString(fmt.Sprintf("C%g %g %g %g %g %g", x56a, y56a, x56b, y56b, x6, y6))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x67b, y67b, x7, y7, x78b, y78b, x8, y8))
	b.WriteString(fmt.Sprintf("S%g %g %g %gL%g %gL%g %g", x89b, y89b, x9, y9, x10, y10, x11, y11))
	b.WriteString(fmt.Sprintf("L%g %gL%g %g", x12, y12, x13, y13))
	b.WriteString(fmt.Sprintf("C%g %g %g %g %g %g", x1314a, y1314a, x1314b, y1314b, x14, y14))
	b.WriteString(fmt.Sprintf("S%g %g %g %gL%g %gL%g %gZ", x1415b, y1415b, x15, y15, x16, y16, x17, y17))

	if err := scope.BuildMathPath(b.String(), ctx); err != nil {
		return nil, err
	}

	return boundingbox.NewBB(x0, y2, x16, y0), nil
}

func NewSum(rhs Token, lower Token, upper Token, ctx context.Context) (*Sum, error) {
	return &Sum{
		rhs,
		lower,
		upper,
		newTokenData(ctx),
	}, nil
}

func (t *Sum) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Sum(")

	if t.lower != nil {
		b.WriteString(t.lower.Dump(""))
	}
	b.WriteString(" ")
	if t.upper != nil {
		b.WriteString(t.upper.Dump(""))
	}
	b.WriteString(")\n")
	b.WriteString(t.rhs.Dump(indent + "  "))

	return b.String()
}

func (t *Sum) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	subScopeRhs := scope.NewSubScope()
	bbRhs, err := t.rhs.GenerateTags(subScopeRhs, 0, y)
	if err != nil {
		return nil, err
	}

	symbolW := 0.0

	var subScopeLower SubScope = nil
	var bbLower boundingbox.BB
	if t.lower != nil {
		subScopeLower = scope.NewSubScope()
		bbLower, err = t.lower.GenerateTags(subScopeLower, 0, 0)
		if err != nil {
			return nil, err
		}

		symbolW = bbLower.Width() * SUM_LOWER_UPPER_SCALE
	}

	var subScopeUpper SubScope = nil
	var bbUpper boundingbox.BB
	if t.upper != nil {
		subScopeUpper = scope.NewSubScope()
		bbUpper, err = t.upper.GenerateTags(subScopeUpper, 0, 0)
		if err != nil {
			return nil, err
		}

		symbolW = math.Max(symbolW, bbUpper.Width()*SUM_LOWER_UPPER_SCALE)
	}

	symbolH := bbRhs.Height()
	symbolY := bbRhs.Bottom()

	bbSymbol, err := GenSumPath(scope, x, symbolY, symbolW, symbolH, t.Context())
	if err != nil {
		return nil, err
	}

	dxRhs := bbSymbol.Right() + SUM_HOR_SPACING
	subScopeRhs.Transform(dxRhs, 0.0, 1.0, 1.0)
	bbRhs = bbRhs.Translate(dxRhs, 0)

	bbAll := boundingbox.Merge(bbSymbol, bbRhs)

	if t.lower != nil {
		dxLower := bbSymbol.Left() + bbSymbol.Width()*0.5 - bbLower.Width()*SUM_LOWER_UPPER_SCALE*0.5
		dyLower := bbSymbol.Bottom() + bbLower.Height()*SUM_LOWER_UPPER_SCALE + SUM_VER_SPACING
		subScopeLower.Transform(dxLower, dyLower, SUM_LOWER_UPPER_SCALE, SUM_LOWER_UPPER_SCALE)
		bbLower = bbLower.Scale(SUM_LOWER_UPPER_SCALE, SUM_LOWER_UPPER_SCALE)
		bbLower = bbLower.Translate(dxLower, dyLower)

		bbAll = boundingbox.Merge(bbAll, bbLower)
	}

	if t.upper != nil {
		dxUpper := bbSymbol.Left() + bbSymbol.Width()*0.5 - bbUpper.Width()*SUM_LOWER_UPPER_SCALE*0.5
		dyUpper := bbSymbol.Top() - SUM_VER_SPACING
		subScopeUpper.Transform(dxUpper, dyUpper, SUM_LOWER_UPPER_SCALE, SUM_LOWER_UPPER_SCALE)
		bbUpper = bbUpper.Scale(SUM_LOWER_UPPER_SCALE, SUM_LOWER_UPPER_SCALE)
		bbUpper = bbUpper.Translate(dxUpper, dyUpper)

		bbAll = boundingbox.Merge(bbAll, bbUpper)
	}

	return bbAll, nil
}
