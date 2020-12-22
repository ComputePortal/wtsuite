package boundingbox

import (
	"fmt"
)

type BB interface {
	Width() float64
	Height() float64
	Left() float64
	Top() float64
	Right() float64
	Bottom() float64

	Translate(dx float64, dy float64) BB
	Scale(sx float64, sy float64) BB

	Dump() string
}

type BBData struct {
	xMin, yMin float64
	xMax, yMax float64
}

func NewBB(xMin, yMin, xMax, yMax float64) BB {
	return &BBData{xMin, yMin, xMax, yMax}
}

func (bb *BBData) Width() float64 {
	return bb.xMax - bb.xMin
}

func (bb *BBData) Height() float64 {
	return bb.yMax - bb.yMin
}

func (bb *BBData) Left() float64 {
	return bb.xMin
}

func (bb *BBData) Right() float64 {
	return bb.xMax
}

func (bb *BBData) Top() float64 {
	return bb.yMin
}

func (bb *BBData) Bottom() float64 {
	return bb.yMax
}

func Merge(bbs ...BB) BB {
	if len(bbs) == 0 {
		panic("expected at least one bb")
	}

	xMin := bbs[0].Left()
	xMax := bbs[0].Right()

	yMin := bbs[0].Top()
	yMax := bbs[0].Bottom()

	for i, bb := range bbs {
		if i == 0 {
			continue
		}

		if bb.Left() < xMin {
			xMin = bb.Left()
		}

		if bb.Top() < yMin {
			yMin = bb.Top()
		}

		if bb.Right() > xMax {
			xMax = bb.Right()
		}

		if bb.Bottom() > yMax {
			yMax = bb.Bottom()
		}
	}

	return NewBB(xMin, yMin, xMax, yMax)
}

func (bb *BBData) Translate(dx, dy float64) BB {
	return NewBB(bb.Left()+dx, bb.Top()+dy, bb.Right()+dx, bb.Bottom()+dy)
}

func (bb *BBData) Scale(sx, sy float64) BB {
	// relative to origin
	return NewBB(bb.Left()*sx, bb.Top()*sy, bb.Right()*sx, bb.Bottom()*sy)
}

func (bb *BBData) Dump() string {
	return fmt.Sprintf("(bb:%g,%g->%g,%g)", bb.Left(), bb.Top(), bb.Right(), bb.Bottom())
}
