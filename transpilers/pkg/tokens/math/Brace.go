package math

import (
	"fmt"
	"strings"

	"./boundingbox"

	"../context"
)

const (
	BRACE_X0 = 0.00002
	BRACE_Y0 = 0.0

	BRACE_X01A = 0.0335
	BRACE_Y01A = -0.0105

	BRACE_X01B = 0.056
	BRACE_Y01B = -0.0229

	BRACE_X1 = 0.06702
	BRACE_Y1 = -0.038

	BRACE_X12A = 0.0795
	BRACE_Y12A = -0.0531

	BRACE_X12B = 0.08502
	BRACE_Y12B = -0.0784

	BRACE_X2 = 0.08502
	BRACE_Y2 = -0.1275

	BRACE_X3 = 0.08502
	BRACE_Y3 = -0.291

	BRACE_X34A = 0.08502
	BRACE_Y34A = -0.3206

	BRACE_X34B = 0.0887
	BRACE_Y34B = -0.3434

	BRACE_X4 = 0.09502
	BRACE_Y4 = -0.362

	BRACE_X45A = 0.1016
	BRACE_Y45A = -0.3806

	BRACE_X45B = 0.1121
	BRACE_Y45B = -0.3949

	BRACE_X5 = 0.12802
	BRACE_Y5 = -0.404

	BRACE_X56A = 0.1434
	BRACE_Y56A = -0.4137

	BRACE_X56B = 0.1602
	BRACE_Y56B = -0.4208

	BRACE_X6 = 0.17802
	BRACE_Y6 = -0.424

	BRACE_X67A = 0.1966
	BRACE_Y67A = -0.4269

	BRACE_X67B = 0.2201
	BRACE_Y67B = -0.4290

	BRACE_X7 = 0.25002
	BRACE_Y7 = -0.43

	BRACE_X8 = 0.25002
	BRACE_Y8 = -0.419

	BRACE_X89A = 0.2153
	BRACE_Y89A = -0.419

	BRACE_X89B = 0.1917
	BRACE_Y89B = -0.3965

	BRACE_X9 = 0.17802
	BRACE_Y9 = -0.379

	BRACE_X910A = 0.1654
	BRACE_Y910A = -0.3610

	BRACE_X910B = 0.15902
	BRACE_Y910B = -0.3333

	BRACE_X10 = 0.15902
	BRACE_Y10 = -0.291

	BRACE_X11 = 0.15902
	BRACE_Y11 = -0.1275

	BRACE_X1112A = 0.15902
	BRACE_Y1112A = -0.0866

	BRACE_X1112B = 0.1527
	BRACE_Y1112B = -0.0577

	BRACE_X12 = 0.1415
	BRACE_Y12 = -0.0388

	BRACE_X1213A = 0.1273
	BRACE_Y1213A = -0.0226

	BRACE_X1213B = 0.1024
	BRACE_Y1213B = -0.0087

	BRACE_X13 = 0.06502
	BRACE_Y13 = 0.0
)

func GenBrace(scope Scope, x float64, y float64, h float64, flipX bool, ctx context.Context) (boundingbox.BB, error) {
	hTopCurve := BRACE_Y3 - BRACE_Y7
	hMiddleCurve := BRACE_Y2 - BRACE_Y3
	hBottomCurve := BRACE_Y0 - BRACE_Y2

	hExtra := 0.0
	if h > 2*(hTopCurve+hMiddleCurve+hBottomCurve) {
		hExtra = 0.5*h - hTopCurve - hBottomCurve - hMiddleCurve
	}

	x0 := x
	fX := 1.0
	if flipX {
		x0 = x + (BRACE_X7 - BRACE_X0)
		fX = -1.0
	}

	y0 := y - (hTopCurve + hMiddleCurve + hBottomCurve + hExtra)

	xRef := BRACE_X0
	yRef := BRACE_Y0

	x01a := x0 + fX*(BRACE_X01A-xRef)
	y01a := y0 + (BRACE_Y01A - yRef)

	x01b := x0 + fX*(BRACE_X01B-xRef)
	y01b := y0 + (BRACE_Y01B - yRef)

	x1 := x0 + fX*(BRACE_X1-xRef)
	y1 := y0 + (BRACE_Y1 - yRef)

	x12a := x0 + fX*(BRACE_X12A-xRef)
	y12a := y0 + (BRACE_Y12A - yRef)

	x12b := x0 + fX*(BRACE_X12B-xRef)
	y12b := y0 + (BRACE_Y12B - yRef)

	x2 := x0 + fX*(BRACE_X2-xRef)
	y2 := y0 + (BRACE_Y2 - yRef)

	x3 := x0 + fX*(BRACE_X3-xRef)
	y3 := y0 + (BRACE_Y3 - yRef) - hExtra

	x34a := x0 + fX*(BRACE_X34A-xRef)
	y34a := y0 + (BRACE_Y34A - yRef) - hExtra

	x34b := x0 + fX*(BRACE_X34B-xRef)
	y34b := y0 + (BRACE_Y34B - yRef) - hExtra

	x4 := x0 + fX*(BRACE_X4-xRef)
	y4 := y0 + (BRACE_Y4 - yRef) - hExtra

	x45a := x0 + fX*(BRACE_X45A-xRef)
	y45a := y0 + (BRACE_Y45A - yRef) - hExtra

	x45b := x0 + fX*(BRACE_X45B-xRef)
	y45b := y0 + (BRACE_Y45B - yRef) - hExtra

	x5 := x0 + fX*(BRACE_X5-xRef)
	y5 := y0 + (BRACE_Y5 - yRef) - hExtra

	x56a := x0 + fX*(BRACE_X56A-xRef)
	y56a := y0 + (BRACE_Y56A - yRef) - hExtra

	x56b := x0 + fX*(BRACE_X56B-xRef)
	y56b := y0 + (BRACE_Y56B - yRef) - hExtra

	x6 := x0 + fX*(BRACE_X6-xRef)
	y6 := y0 + (BRACE_Y6 - yRef) - hExtra

	x67a := x0 + fX*(BRACE_X67A-xRef)
	y67a := y0 + (BRACE_Y67A - yRef) - hExtra

	x67b := x0 + fX*(BRACE_X67B-xRef)
	y67b := y0 + (BRACE_Y67B - yRef) - hExtra

	x7 := x0 + fX*(BRACE_X7-xRef)
	y7 := y0 + (BRACE_Y7 - yRef) - hExtra

	x8 := x0 + fX*(BRACE_X8-xRef)
	y8 := y0 + (BRACE_Y8 - yRef) - hExtra

	x89a := x0 + fX*(BRACE_X89A-xRef)
	y89a := y0 + (BRACE_Y89A - yRef) - hExtra

	x89b := x0 + fX*(BRACE_X89B-xRef)
	y89b := y0 + (BRACE_Y89B - yRef) - hExtra

	x9 := x0 + fX*(BRACE_X9-xRef)
	y9 := y0 + (BRACE_Y9 - yRef) - hExtra

	x910a := x0 + fX*(BRACE_X910A-xRef)
	y910a := y0 + (BRACE_Y910A - yRef) - hExtra

	x910b := x0 + fX*(BRACE_X910B-xRef)
	y910b := y0 + (BRACE_Y910B - yRef) - hExtra

	x10 := x0 + fX*(BRACE_X10-xRef)
	y10 := y0 + (BRACE_Y10 - yRef) - hExtra

	x11 := x0 + fX*(BRACE_X11-xRef)
	y11 := y0 + (BRACE_Y11 - yRef)

	x1112a := x0 + fX*(BRACE_X1112A-xRef)
	y1112a := y0 + (BRACE_Y1112A - yRef)

	x1112b := x0 + fX*(BRACE_X1112B-xRef)
	y1112b := y0 + (BRACE_Y1112B - yRef)

	x12 := x0 + fX*(BRACE_X12-xRef)
	y12 := y0 + (BRACE_Y12 - yRef)

	x1213a := x0 + fX*(BRACE_X1213A-xRef)
	y1213a := y0 + (BRACE_Y1213A - yRef)

	x1213b := x0 + fX*(BRACE_X1213B-xRef)
	y1213b := y0 + (BRACE_Y1213B - yRef)

	x13 := x0 + fX*(BRACE_X13-xRef)
	y13 := y0 + (BRACE_Y13 - yRef)

	yRef = y0
	y0_ := y0

	y01a_ := y0_ - (y01a - yRef)
	y1_ := y0_ - (y1 - yRef)
	y12a_ := y0_ - (y12a - yRef)
	y12b_ := y0_ - (y12b - yRef)
	y2_ := y0_ - (y2 - yRef)
	y3_ := y0_ - (y3 - yRef)
	y34a_ := y0_ - (y34a - yRef)
	y4_ := y0_ - (y4 - yRef)
	y45a_ := y0_ - (y45a - yRef)
	y5_ := y0_ - (y5 - yRef)
	y56a_ := y0_ - (y56a - yRef)
	y6_ := y0_ - (y6 - yRef)
	y67a_ := y0_ - (y67a - yRef)
	y67b_ := y0_ - (y67b - yRef)
	y7_ := y0_ - (y7 - yRef)
	y8_ := y0_ - (y8 - yRef)
	y89a_ := y0_ - (y89a - yRef)
	y9_ := y0_ - (y9 - yRef)
	y910a_ := y0_ - (y910a - yRef)
	y910b_ := y0_ - (y910b - yRef)
	y10_ := y0_ - (y10 - yRef)
	y11_ := y0_ - (y11 - yRef)
	y1112a_ := y0_ - (y1112a - yRef)
	y12_ := y0_ - (y12 - yRef)
	y1213a_ := y0_ - (y1213a - yRef)
	y1213b_ := y0_ - (y1213b - yRef)

	var b strings.Builder

	b.WriteString(fmt.Sprintf("M%g %gC%g %g %g %g %g %g", x0, y0, x01a, y01a, x01b, y01b, x1, y1))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x12b, y12b, x2, y2))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x3, y3, x34a, y34a, x34b, y34b, x4, y4))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x45b, y45b, x5, y5, x56b, y56b, x6, y6))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x67b, y67b, x7, y7))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x8, y8, x89a, y89a, x89b, y89b, x9, y9))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x910b, y910b, x10, y10))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x11, y11, x1112a, y1112a, x1112b, y1112b, x12, y12))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x1213b, y1213b, x13, y13))
	b.WriteString(fmt.Sprintf("C%g %g %g %g %g %g", x1213b, y1213b_, x1213a, y1213a_, x12, y12_))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x1112a, y1112a_, x11, y11_))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x10, y10_, x910b, y910b_, x910a, y910a_, x9, y9_))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x89a, y89a_, x8, y8_))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x7, y7_, x67b, y67b_, x67a, y67a_, x6, y6_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gS%g %g %g %g", x56a, y56a_, x5, y5_, x45a, y45a_, x4, y4_))
	b.WriteString(fmt.Sprintf("S%g %g %g %g", x34a, y34a_, x3, y3_))
	b.WriteString(fmt.Sprintf("L%g %gC%g %g %g %g %g %g", x2, y2_, x12b, y12b_, x12a, y12a_, x1, y1_))
	b.WriteString(fmt.Sprintf("S%g %g %g %gZ", x01a, y01a_, x0, y0_))

	if err := scope.BuildMathPath(b.String(), ctx); err != nil {
		return nil, err
	}

	if flipX {
		return boundingbox.NewBB(x7, y7, x0, y7_), nil
	} else {
		return boundingbox.NewBB(x0, y7, x7, y7_), nil
	}
}

func GenLeftBrace(scope Scope, x float64, y float64, h float64, ctx context.Context) (boundingbox.BB, error) {
	return GenBrace(scope, x, y, h, false, ctx)
}

func GenRightBrace(scope Scope, x float64, y float64, h float64, ctx context.Context) (boundingbox.BB, error) {
	return GenBrace(scope, x, y, h, true, ctx)
}
