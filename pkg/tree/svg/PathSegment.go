package svg

import (
	"reflect"

	"../../tokens/context"
)

type PathSegment interface {
	Length() float64
	// f between 0 and 1
	Position(f float64) (float64, float64)
	// f between 0 and 1
	Tangent(f float64) (float64, float64) // pointing from start to end, normalized

	Reverse() PathSegment
	Context() context.Context
}

type PathSegmentData struct {
	start, stop PathCommand
	x0, y0      float64
	x1, y1      float64
}

func (s *PathSegmentData) Context() context.Context {
	return context.MergeContexts(s.start.Context(), s.stop.Context())
}

func mirrorControlPoint(x0, y0, xIn, yIn float64) (float64, float64) {
	xOut := x0 + (x0 - xIn)
	yOut := y0 + (y0 - yIn)

	return xOut, yOut
}

// same symbol names as wikipedia page
func intersect(x1, y1, x2, y2, x3, y3, x4, y4 float64) (float64, float64) {
	den := (x1-x2)*(y3-y4) - (y1-y2)*(x3-x4)

	cof12 := (x1*y2 - y1*x2)
	cof34 := (x3*y4 - y3*x4)

	x := (cof12*(x3-x4) - cof34*(x1-x2)) / den
	y := (cof12*(y3-y4) - cof34*(y1-y2)) / den

	return x, y
}

func GenerateSegments(pcs []PathCommand, ctx context.Context) ([]PathSegment, error) {
	result := make([]PathSegment, 0)

	x, y := 0.0, 0.0

	if len(pcs) < 2 {
		errCtx := ctx
		return nil, errCtx.NewError("Error: path too short\n")
	}

	switch pc := pcs[0].(type) {
	case *MoveTo:
		x, y = pc.x, pc.y
	case *MoveBy:
		x, y = x+pc.dx, y+pc.dy
	default:
		errCtx := pc.Context()
		return nil, errCtx.NewError("Error: expected m or M as first command")
	}

	prev := pcs[0]

	var pathStartX, pathStartY float64
	for i := 1; i < len(pcs); i++ {
		switch pc := pcs[i].(type) {
		case *MoveTo:
			x, y = pc.x, pc.y
			pathStartX, pathStartY = x, y
		case *MoveBy:
			x, y = x+pc.dx, y+pc.dy
			pathStartX, pathStartY = x, y
		case *LineTo:
			result = append(result, NewLine(prev, pc, x, y, pc.x, pc.y))
			x, y = pc.x, pc.y
		case *LineBy:
			result = append(result, NewLine(prev, pc, x, y, x+pc.dx, y+pc.dy))
			x, y = x+pc.dx, y+pc.dy
		case *HorTo:
			result = append(result, NewLine(prev, pc, x, y, pc.x, y))
			x = pc.x
		case *VerTo:
			result = append(result, NewLine(prev, pc, x, y, x, pc.y))
			y = pc.y
		case *HorBy:
			result = append(result, NewLine(prev, pc, x, y, x+pc.dx, y))
			x = x + pc.dx
		case *VerBy:
			result = append(result, NewLine(prev, pc, x, y, x, y+pc.dy))
			y = y + pc.dy
		case *Close:
			x, y = pathStartX, pathStartY
		case *QuadraticTo:
			result = append(result, NewQuadratic(prev, pc, x, y, pc.x, pc.y, pc.x1, pc.y1))
			x, y = pc.x, pc.y
		case *QuadraticBy:
			result = append(result, NewQuadratic(prev, pc, x, y, x+pc.dx, y+pc.dy, x+pc.dx1, y+pc.dy1))
			x, y = x+pc.dx, y+pc.dy
		case *ExtraQuadraticTo:
			xa, ya := x, y
			if q, ok := result[len(result)-1].(*Quadratic); ok {
				xa, ya = mirrorControlPoint(x, y, q.xa, q.ya)
			}
			result = append(result, NewQuadratic(prev, pc, x, y, pc.x, pc.y, xa, ya))
			x, y = pc.x, pc.y
		case *ExtraQuadraticBy:
			xa, ya := x, y
			if q, ok := result[len(result)-1].(*Quadratic); ok {
				xa, ya = mirrorControlPoint(x, y, q.xa, q.ya)
			}
			result = append(result, NewQuadratic(prev, pc, x, y, x+pc.dx, y+pc.dy, xa, ya))
			x, y = x+pc.dx, y+pc.dy
		case *CubicTo:
			result = append(result, NewCubic(prev, pc, x, y, pc.x, pc.y, pc.x1, pc.y1, pc.x2, pc.y2))
			x, y = pc.x, pc.y
		case *CubicBy:
			result = append(result, NewCubic(prev, pc,
				x, y,
				x+pc.dx, y+pc.dy,
				x+pc.dx1, y+pc.dy1,
				x+pc.dx2, y+pc.dy2))
			x, y = x+pc.dx, y+pc.dy
		case *ExtraCubicTo:
			xa, ya := x, y
			switch curve := result[len(result)-1].(type) {
			case *Quadratic:
				xa, ya = mirrorControlPoint(x, y, curve.xa, curve.ya)
			case *Cubic:
				xa, ya = mirrorControlPoint(x, y, curve.xb, curve.yb)
			}
			result = append(result, NewCubic(prev, pc, x, y, pc.x, pc.y, xa, ya, pc.x2, pc.y2))
			x, y = pc.x, pc.y
		case *ExtraCubicBy:
			xa, ya := x, y
			switch curve := result[len(result)-1].(type) {
			case *Quadratic:
				xa, ya = mirrorControlPoint(x, y, curve.xa, curve.ya)
			case *Cubic:
				xa, ya = mirrorControlPoint(x, y, curve.xb, curve.yb)
			}
			result = append(result, NewCubic(prev, pc, x, y, x+pc.dx, y+pc.dy, xa, ya, x+pc.dx2, y+pc.dy2))
			x, y = x+pc.dx, y+pc.dy
		case *ArcTo:
			result = append(result, NewArc(prev, pc, x, y, pc.x, pc.y, pc.rx, pc.ry, pc.xAxisRot,
				pc.largeArc, pc.positiveSweep))
			x, y = pc.x, pc.y
		case *ArcBy:
			result = append(result, NewArc(prev, pc, x, y, x+pc.dx, y+pc.dy, pc.rx, pc.ry, pc.xAxisRot,
				pc.largeArc, pc.positiveSweep))
			x, y = x+pc.dx, y+pc.dy
		default:
			panic("unhandled")
		}
	}

	return result, nil
}

func ShortenEnds(pcs []PathCommand, dStart, dEnd float64, ctx context.Context) ([]PathCommand, error) {
	result := make([]PathCommand, 0)

	ss, err := GenerateSegments(pcs, ctx)
	if err != nil {
		return nil, err
	}

	if dStart < 0.0 || dEnd < 0.0 {
		panic("implicit lengthening not yet supported")
	}

	if dStart > 0.0 {
		var x, y float64
		var x0, y0 float64

		switch curve := ss[0].(type) {
		case *Line:
			l := curve.Length()

			f := dStart / l
			if f > 1.0 {
				errCtx := context.MergeContexts(curve.Context(), curve.stop.Context())
				return nil, errCtx.NewError("Error: segment not long enough to shorten")
			}

			x0, y0 = curve.Position(0.0)
			x, y = curve.Position(f)

		case *Arc:
			l := curve.Length()

			f := dStart / l
			if f > 1.0 {
				errCtx := context.MergeContexts(curve.start.Context(), curve.stop.Context())
				return nil, errCtx.NewError("Error: segment not long enough to shorten")
			}

			x0, y0 = curve.Position(0.0)
			x, y = curve.Position(f)
		case *Quadratic:
			l := curve.Length()

			f := dStart / l

			if f > 1.0 {
				errCtx := context.MergeContexts(curve.start.Context(), curve.stop.Context())
				return nil, errCtx.NewError("Error: segment not long enough to shorten")
			}

			x0, y0 = curve.Position(0.0)
			x, y = curve.Position(f)
		default:
			errCtx := ss[0].Context()
			return nil, errCtx.NewError("Error: not supported for shortening (" + reflect.TypeOf(ss[0]).String() + ")")
		}

		switch move := pcs[0].(type) {
		case *MoveTo:
			result = append(result, NewMoveTo(x, y, move.Context()))
		case *MoveBy:
			result = append(result, NewMoveTo(x, y, move.Context()))
		default:
			errCtx := ss[0].Context()
			return nil, errCtx.NewError("Error: expected M or m")
		}

		replaceSecond := false
		switch second := pcs[1].(type) {
		case *HorBy:
			result = append(result, NewHorBy(second.dx-(x-x0), second.Context()))
			replaceSecond = true
		case *VerBy:
			result = append(result, NewVerBy(second.dy-(y-y0), second.Context()))
			replaceSecond = true
		case *LineBy:
			result = append(result, NewLineBy(second.dx-(x-x0),
				second.dy-(y-y0), second.Context()))
			replaceSecond = true
		case *ArcBy:
			result = append(result, NewArcBy(second.dx-(x-x0), second.dy-(y-y0),
				second.rx, second.ry, second.xAxisRot, second.largeArc, second.positiveSweep, second.Context()))
			replaceSecond = true
		}

		if replaceSecond {
			result = append(result, pcs[2:len(pcs)-1]...)
		} else {
			result = append(result, pcs[1:len(pcs)-1]...)
		}
	} else {
		result = append(result, pcs[0:len(pcs)-1]...)
	}

	if dEnd > 0.0 {
		var x, y float64
		var dx, dy float64
		var xa, ya float64
		var dxa, dya float64
		switch curve := ss[len(ss)-1].(type) {
		case *Line:
			l := curve.Length()

			f := (1.0 - dEnd/l)
			if f < 0.0 {
				errCtx := context.MergeContexts(curve.Context(), curve.stop.Context())
				return nil, errCtx.NewError("Error: segment not long enough to shorten")
			}

			x, y = curve.Position(f)
			dx, dy = x-curve.x0, y-curve.y0
		case *Arc:
			l := curve.Length()

			f := (1.0 - dEnd/l)
			if f < 0.0 {
				errCtx := context.MergeContexts(curve.start.Context(), curve.stop.Context())
				return nil, errCtx.NewError("Error: segment not long enough to shorten")
			}

			x, y = curve.Position(f)
			dx, dy = x-curve.x0, y-curve.y0
		case *Quadratic:
			l := curve.Length()

			f := (1.0 - dEnd/l)
			if f < 0.0 {
				errCtx := context.MergeContexts(curve.start.Context(), curve.stop.Context())
				return nil, errCtx.NewError("Error: segment not long enough to shorten")
			}

			x, y = curve.Position(f)

			tx, ty := curve.Tangent(f)
			t0x, t0y := curve.Tangent(0.0)
			// intersect these two tangents to determine the new control point
			xa, ya = intersect(x, y, x+tx, y+ty, curve.x0, curve.y0, curve.x0+t0x, curve.y0+t0y)

			dxa, dya = xa-curve.x0, ya-curve.y0
			dx, dy = x-curve.x0, y-curve.y0
		default:
			errCtx := ss[len(ss)-1].Context()
			return nil, errCtx.NewError("Error: not supported for shortening (" + reflect.TypeOf(ss[len(ss)-1]).String() + ")")
		}

		switch ec := pcs[len(pcs)-1].(type) {
		case *LineTo:
			result = append(result, NewLineTo(x, y, ec.Context()))
		case *LineBy:
			result = append(result, NewLineBy(dx, dy, ec.Context()))
		case *HorTo:
			result = append(result, NewHorTo(x, ec.Context()))
		case *HorBy:
			result = append(result, NewHorBy(dx, ec.Context()))
		case *VerTo:
			result = append(result, NewVerTo(y, ec.Context()))
		case *VerBy:
			result = append(result, NewVerBy(dy, ec.Context()))
		case *ArcTo:
			result = append(result, NewArcTo(x, y, ec.rx, ec.ry, ec.xAxisRot,
				ec.largeArc, ec.positiveSweep, ec.Context()))
		case *ArcBy:
			result = append(result, NewArcBy(dx, dy, ec.rx, ec.ry, ec.xAxisRot,
				ec.largeArc, ec.positiveSweep, ec.Context()))
		case *QuadraticTo:
			result = append(result, NewQuadraticTo(xa, ya, x, y, ec.Context()))
		case *QuadraticBy:
			result = append(result, NewQuadraticBy(dxa, dya, dx, dy, ec.Context()))
		default:
			errCtx := ss[len(ss)-1].Context()
			return nil, errCtx.NewError("Error: expected line or arc")
		}
	} else {
		result = append(result, pcs[len(pcs)-1])
	}

	return result, nil
}

func SegmentsLength(segments []PathSegment) float64 {
	l := 0.0
	for _, s := range segments {

		l += s.Length()
	}

	return l
}

func SegmentsPosition(segments []PathSegment, t float64) (float64, float64) {
	tPrev := 0.0
	for _, s := range segments {
		lSeg := s.Length()
		if t < tPrev+lSeg {
			return s.Position((t - tPrev) / lSeg)
		} else {
			tPrev += lSeg
		}
	}

	return segments[len(segments)-1].Position(1.0)
}
