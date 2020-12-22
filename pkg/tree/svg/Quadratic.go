package svg

import (
	"math"
)

type Quadratic struct {
	PathSegmentData
	xa, ya float64 // control point
}

func NewQuadratic(start, stop PathCommand, x0, y0, x1, y1 float64, xa, ya float64) *Quadratic {
	return &Quadratic{PathSegmentData{start, stop, x0, y0, x1, y1}, xa, ya}
}

// dirty approximation (probably good enough for all visual purposes)
func (s *Quadratic) Length() float64 {
	dx01 := s.x1 - s.x0
	dy01 := s.y1 - s.y0
	l01 := math.Sqrt(dx01*dx01 + dy01*dy01)

	dx0a := s.xa - s.x0
	dy0a := s.ya - s.y0
	l0a := math.Sqrt(dx0a*dx0a + dy0a*dy0a)

	dxa1 := s.x1 - s.xa
	dya1 := s.y1 - s.ya
	l1a := math.Sqrt(dxa1*dxa1 + dya1*dya1)

	return 2.0/3.0*l01 + 1.0/3.0*(l0a+l1a)
}

func quadraticPosition(x0, y0, x1, y1, x2, y2, f float64) (float64, float64) {
	f_ := 1.0 - f

	x := f_*(f_*x0+f*x1) + f*(f_*x1+f*x2)
	y := f_*(f_*y0+f*y1) + f*(f_*y1+f*y2)

	return x, y
}

func (s *Quadratic) Position(f float64) (float64, float64) {
	return quadraticPosition(s.x0, s.y0, s.xa, s.ya, s.x1, s.y1, f)
}

func (s *Quadratic) Tangent(f float64) (float64, float64) {
	tx := 2.0*(1.0-f)*(s.xa-s.x0) + 2.0*f*(s.x1-s.xa)
	ty := 2.0*(1.0-f)*(s.ya-s.y0) + 2.0*f*(s.y1-s.ya)

	norm := math.Sqrt(tx*tx + ty*ty)
	return tx / norm, ty / norm
}

func (s *Quadratic) Reverse() PathSegment {
	return NewQuadratic(s.stop, s.start, s.x1, s.y1, s.x0, s.y0, s.xa, s.ya)
}
