package svg

import (
	"math"
)

type Cubic struct {
	PathSegmentData
	xa, ya float64 // control point
	xb, yb float64 // control point
}

func NewCubic(start, stop PathCommand, x0, y0, x1, y1 float64, xa, ya, xb, yb float64) *Cubic {
	return &Cubic{PathSegmentData{start, stop, x0, y0, x1, y1}, xa, ya, xb, yb}
}

func (s *Cubic) Length() float64 {
	// sample using ten divisions (11 samples)
	l := 0.0

	n := 10
	df := 1.0 / float64(n)
	for i := 0; i < n; i++ {
		fStart := float64(i) * df
		fEnd := float64(i+1) * df

		xStart, yStart := s.Position(fStart)
		xEnd, yEnd := s.Position(fEnd)

		dx := xEnd - xStart
		dy := yEnd - yStart
		l += math.Sqrt(dx*dx + dy*dy)
	}

	return l
}

func (s *Cubic) Position(f float64) (float64, float64) {
	f_ := 1.0 - f

	xqa, yqa := quadraticPosition(s.x0, s.y0, s.xa, s.ya, s.xb, s.yb, f)

	xqb, yqb := quadraticPosition(s.xa, s.ya, s.xb, s.yb, s.x1, s.y1, f)

	x := f_*xqa + f*xqb
	y := f_*yqa + f*yqb

	return x, y
}

func (s *Cubic) Tangent(f float64) (float64, float64) {
	f_ := 1.0 - f

	tx := 3.0*f_*f_*(s.xa-s.x0) + 6.0*f_*f*(s.xb-s.xa) + 3.0*f*f*(s.x1-s.xb)
	ty := 3.0*f_*f_*(s.ya-s.y0) + 6.0*f_*f*(s.yb-s.ya) + 3.0*f*f*(s.y1-s.yb)

	norm := math.Sqrt(tx*tx + ty*ty)
	return tx / norm, ty / norm
}

func (s *Cubic) Reverse() PathSegment {
	return NewCubic(s.stop, s.start, s.x1, s.y1, s.x0, s.y0, s.xb, s.yb, s.xa, s.ya)
}
