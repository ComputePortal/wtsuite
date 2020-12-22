package svg

import (
	"math"
)

type Line struct {
	PathSegmentData
}

func NewLine(start, stop PathCommand, x0, y0, x1, y1 float64) *Line {
	return &Line{PathSegmentData{start, stop, x0, y0, x1, y1}}
}

func (s *Line) Length() float64 {
	dx := s.x1 - s.x0
	dy := s.y1 - s.y0

	return math.Sqrt(dx*dx + dy*dy)
}

func (s *Line) Position(f float64) (float64, float64) {
	xf := s.x0*(1.0-f) + s.x1*f
	yf := s.y0*(1.0-f) + s.y1*f

	return xf, yf
}

func (s *Line) Tangent(f float64) (float64, float64) {
	tx := s.x1 - s.x0
	ty := s.y1 - s.y0

	l := math.Sqrt(tx*tx + ty*ty)
	tx = tx / l
	ty = ty / l

	return tx, ty
}

func (s *Line) Reverse() PathSegment {
	return NewLine(s.stop, s.start, s.x1, s.y1, s.x0, s.y0)
}
