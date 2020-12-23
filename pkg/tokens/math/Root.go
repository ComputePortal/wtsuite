package math

import (
	"fmt"
	"math"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

//                                 ________________________
//                                9                        8
//           __5                 /                         |
//      ____/   \               /                          |
//    _/         \             /    6______________________7
//   0            \           /    /
//    \         2  \         /    /
//     \   ____/ \  \       /    /
//      1_/       \  \     /    /
//                 \  \   /    /
//                  \  \ /    /
//                   \  4    /
//                    \     /
//                     \   /
//                      \ /
//                       3
//
//
//
//

// the first 6 points of the root, relative to the left bottom
const (
	ROOT_X0 = 0.0
	ROOT_Y0 = -0.47921

	ROOT_X1 = 0.02718
	ROOT_Y1 = -0.42137

	ROOT_X2 = 0.12247
	ROOT_Y2 = -0.46849

	ROOT_X3 = 0.37698
	ROOT_Y3 = 0.0

	ROOT_X4 = 0.35201
	ROOT_Y4 = -0.24255

	ROOT_X5 = 0.16953
	ROOT_Y5 = -0.56133

	// in radians
	ROOT_ANGLE     = 0.16071
	ROOT_THICKNESS = 0.0642

	ROOT_VER_SPACING  = 0.1
	ROOT_HOR_SPACING  = 0.1
	ROOT_LEFT_PADDING = 0.1
)

type Root struct {
	BinaryOp // XXX: second argument is ignored for now
}

func NewRoot(a Token, b Token, ctx context.Context) (Token, error) {
	return &Root{BinaryOp{"root", a, b, newTokenData(ctx)}}, nil
}

func (t *Root) genRootPath(scope Scope, x float64, y float64, innerWidth float64, innerHeight float64) (boundingbox.BB, error) {
	y67 := y - innerHeight - ROOT_VER_SPACING
	y89 := y67 - ROOT_THICKNESS
	y6 := y67
	y7 := y67
	y8 := y89
	y9 := y89

	x6 := x - ROOT_HOR_SPACING
	x7 := x + innerWidth
	x8 := x7

	x3 := x6 - innerHeight*math.Tan(ROOT_ANGLE)
	x2 := x3 - (ROOT_X3 - ROOT_X2)
	x1 := x3 - (ROOT_X3 - ROOT_X1)
	x0 := x3 - (ROOT_X3 - ROOT_X0)

	x4 := x3 - (ROOT_X3 - ROOT_X4)
	x5 := x3 - (ROOT_X3 - ROOT_X5)

	y3 := y
	y2 := y3 - (ROOT_Y3 - ROOT_Y2)
	y1 := y3 - (ROOT_Y3 - ROOT_Y1)
	y0 := y3 - (ROOT_Y3 - ROOT_Y0)
	y5 := y3 - (ROOT_Y3 - ROOT_Y5)
	y4 := y3 - (ROOT_Y3 - ROOT_Y4)

	// x9 remains
	x9 := x4 + math.Tan(ROOT_ANGLE)*(y4-y89)

	var b strings.Builder

	b.WriteString(fmt.Sprintf("M%g %gL%g %gL%g %g", x0, y0, x1, y1, x2, y2))
	b.WriteString(fmt.Sprintf("L%g %gL%g %gL%g %g", x3, y3, x6, y6, x7, y7))
	b.WriteString(fmt.Sprintf("L%g %gL%g %gL%g %gL%g %gZ", x8, y8, x9, y9, x4, y4, x5, y5))

	if err := scope.BuildMathPath(b.String(), t.Context()); err != nil {
		return nil, err
	}

	return boundingbox.NewBB(x0, y89, x7, y3), nil
}

func (t *Root) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	subScope := scope.NewSubScope()

	// we will keep using the t.a baseline
	bba, err := t.a.GenerateTags(subScope, 0.0, y)
	if err != nil {
		return nil, err
	}

	bbRoot, err := t.genRootPath(subScope, bba.Left(), bba.Bottom(), bba.Width(), bba.Height())
	if err != nil {
		return nil, err
	}

	bbRoot = boundingbox.NewBB(bbRoot.Left()-ROOT_LEFT_PADDING, bbRoot.Top(), bbRoot.Right(), bbRoot.Bottom())
	dx := x - bbRoot.Left()
	subScope.Transform(dx, 0, 1.0, 1.0)
	bbRoot = bbRoot.Translate(dx, 0.0)

	return bbRoot, nil
}
