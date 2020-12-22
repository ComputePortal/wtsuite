package svg

import (
	"math"
)

type Arc struct {
	PathSegmentData
	ArcCoreData
	xc, yc float64
}

func NewArc(start, stop PathCommand, x0, y0, x1, y1 float64, rx, ry, xAxisRot float64,
	largeArc, positiveSweep bool) *Arc {
	arc := &Arc{PathSegmentData{start, stop, x0, y0, x1, y1},
		ArcCoreData{rx, ry, xAxisRot, largeArc, positiveSweep},
		0.0, 0.0}

	arc.xc, arc.yc = arc.centre()

	return arc
}

func (s *Arc) rotateCoords(xIn, yIn float64) (float64, float64) {
	beta := math.Pi * s.xAxisRot / 180.0

	xOut := xIn*math.Cos(beta) - yIn*math.Sin(beta)
	yOut := xIn*math.Sin(beta) + yIn*math.Cos(beta)

	return xOut, yOut
}

func (s *Arc) unrotateCoords(xIn, yIn float64) (float64, float64) {
	beta := math.Pi * s.xAxisRot / 180.0

	xOut := xIn*math.Cos(beta) + yIn*math.Sin(beta)
	yOut := -xIn*math.Sin(beta) + yIn*math.Cos(beta)

	return xOut, yOut
}

// give error if ellipse is too small for chosen points
func (s *Arc) centre() (float64, float64) {
	// solve a quadratic equation for xc, substituting yc by A + B*xc

	// rotate input points by xAxisRot
	x0, y0 := s.rotateCoords(s.x0, s.y0)
	x1, y1 := s.rotateCoords(s.x1, s.y1)

	dx, dy := x1-x0, y1-y0

	rx2 := s.rx * s.rx
	ry2 := s.ry * s.ry
	e2 := ry2 / rx2

	// 2*yc*(y1-y0) == Ay + By*xc
	Ay := e2*(x1*x1-x0*x0) + (y1*y1 - y0*y0)
	By := -e2 * 2.0 * dx

	// solve A*xc^2 + B*xc + C = 0
	// first multiply fitting eq for x0 by (y1-y0)^2
	dy2 := dy * dy
	Ax := dy2/rx2 + 0.25*By*By/ry2
	Bx := -2.0*x0*dy2/rx2 - By*y0*dy/ry2 + 0.5*Ay*By/ry2
	Cx := dy2*x0*x0/rx2 + dy2*y0*y0/ry2 - Ay*y0*dy/ry2 + 0.25*Ay*Ay/ry2 - dy2

	Dx := Bx*Bx - 4.0*Ax*Cx
	if Dx < 0.0 {
		panic("ellipse too small, TODO: write proper escape")
	}

	xc0 := (-Bx - math.Sqrt(Dx)) / (2.0 * Ax)
	xc1 := (-Bx + math.Sqrt(Dx)) / (2.0 * Ax)

	yc0 := (Ay + By*xc0) / (2.0 * dy)
	yc1 := (Ay + By*xc1) / (2.0 * dy)

	if (s.positiveSweep == s.largeArc) == ((x0-xc0)*(y1-yc0)-(x1-xc0)*(y0-yc0) > 0) {
		return s.unrotateCoords(xc0, yc0)
	} else {
		return s.unrotateCoords(xc1, yc1)
	}
}

func (s *Arc) angle(x, y float64) float64 {
	x, y = s.rotateCoords(x-s.xc, y-s.yc)

	theta := math.Acos(x / s.rx)

	if math.IsNaN(theta) {
		panic("bad ellipse")
	}

	if y < 0.0 {
		theta = 2.0*math.Pi - theta
	}

	return theta
}

func (s *Arc) angles() (float64, float64) {
	return s.angle(s.x0, s.y0), s.angle(s.x1, s.y1)
}

func (s *Arc) lengthCondition(cond func(l float64, theta float64) bool) float64 {
	theta0, theta1 := s.angles()
	n := 50.0
	l := 0.0

	if s.positiveSweep {
		if theta1 < theta0 {
			theta1 = theta1 + 2.0*math.Pi
		}

		dtheta := (theta1 - theta0) / n
		for theta := theta0; theta <= theta1; theta = theta + dtheta {
			dx := s.rx * math.Sin(theta)
			dy := s.ry * math.Cos(theta)
			l += math.Sqrt(dx*dx+dy*dy) * dtheta

			if cond(l, theta) {
				return l
			}
		}

	} else {
		if theta1 > theta0 {
			theta1 = theta1 - 2.0*math.Pi
		}

		dtheta := (theta0 - theta1) / n
		for theta := theta0; theta >= theta1; theta = theta - dtheta {
			dx := s.rx * math.Sin(theta)
			dy := s.ry * math.Cos(theta)
			l += math.Sqrt(dx*dx+dy*dy) * dtheta
			if cond(l, theta) {
				return l
			}
		}
	}

	return l
}

func (s *Arc) Length() float64 {
	return s.lengthCondition(func(l float64, theta float64) bool { return false })
}

func (s *Arc) Position(f float64) (float64, float64) {
	L := s.Length()

	var x, y float64
	ok := false

	lf := L * f
	s.lengthCondition(func(l float64, theta float64) bool {
		if l > lf {
			x = s.rx * math.Cos(theta)
			y = s.ry * math.Sin(theta)

			ok = true
			return true
		} else {
			return false
		}
	})

	if !ok {
		panic("f never reached")
	}

	return x, y
}

func (s *Arc) Tangent(f float64) (float64, float64) {
	L := s.Length()

	var dx, dy float64
	ok := false

	lf := L * f
	s.lengthCondition(func(l float64, theta float64) bool {
		if l > lf {
			dx = -s.rx * math.Sin(theta)
			dy = s.ry * math.Cos(theta)

			ok = true
			return true
		} else {
			return false
		}
	})

	if !ok {
		panic("f never reached")
	}

	// normalize
	d := math.Sqrt(dx*dx + dy*dy)

	if !s.positiveSweep {
		d = d * -1.0
	}

	return dx / d, dy / dy
}

func (s *Arc) Reverse() PathSegment {
	return NewArc(s.stop, s.start, s.x1, s.y1, s.x0, s.y0, s.rx, s.ry, s.xAxisRot, s.largeArc, s.positiveSweep)
}
