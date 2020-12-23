package math

import (
	"math"
	"strconv"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Align struct {
	minHorSpacing      float64
	minVerSpacing      float64
	constantHorSpacing bool
	constantVerSpacing bool
	eqs                [][]Token // first index is the row, second index is the column
	TokenData
}

func NewAlign(minHorSpacing, minVerSpacing float64, constantHorSpacing, constantVerSpacing bool,
	eqs [][]Token, ctx context.Context) (*Align, error) {
	if len(eqs) == 0 {
		panic("can't have zero rows")
	}

	if len(eqs[0]) == 0 {
		panic("can't have zero cols")
	}

	return &Align{minHorSpacing, minVerSpacing, constantHorSpacing, constantVerSpacing,
		eqs, newTokenData(ctx)}, nil
}

func newAlign(minHorSpacing, minVerSpacing float64, constantHorSpacing, constantVerSpacing bool,
	eqs [][]Token, ctx context.Context) *Align {
	al, err := NewAlign(minHorSpacing, minVerSpacing, constantHorSpacing, constantVerSpacing, eqs, ctx)
	if err != nil {
		panic(err)
	}

	return al
}

func (t *Align) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)

	b.WriteString("Align\n")

	for i, row := range t.eqs {
		b.WriteString("|-")
		b.WriteString("row ")
		b.WriteString(strconv.Itoa(i))
		b.WriteString("\n")

		for _, eq := range row {
			b.WriteString(eq.Dump(indent + "|   "))
			b.WriteString("\n")
		}

		if i < len(t.eqs)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (t *Align) maxRowHeight(bbs [][]boundingbox.BB, iRow int) float64 {
	h := 0.0

	for _, bb := range bbs[iRow] {
		h = math.Max(h, bb.Height())
	}

	return h
}

func (t *Align) maxHeight(bbs [][]boundingbox.BB) float64 {
	h := 0.0

	for iRow, _ := range bbs {
		h = math.Max(h, t.maxRowHeight(bbs, iRow))
	}

	return h
}

func (t *Align) maxColWidth(bbs [][]boundingbox.BB, iCol int) float64 {
	w := 0.0

	for _, bbsRow := range bbs {
		w = math.Max(w, bbsRow[iCol].Width())
	}

	return w
}

func (t *Align) maxWidth(bbs [][]boundingbox.BB) float64 {
	w := 0.0

	for _, bbsRow := range bbs {
		for _, bb := range bbsRow {
			w = math.Max(w, bb.Width())
		}
	}

	return w
}

// this is the distance between the baselines
func (t *Align) rowHeights(bbs [][]boundingbox.BB) []float64 {
	heights := make([]float64, len(t.eqs))
	if len(t.eqs) == 0 {
		panic("rowHeights should never be 0 length")
	}

	if t.constantVerSpacing {
		d := t.maxHeight(bbs) + t.minVerSpacing

		for i, _ := range t.eqs {
			heights[i] = d
		}
	} else {
		for i, _ := range t.eqs {
			heights[i] = t.maxRowHeight(bbs, i) + t.minVerSpacing
		}
	}

	return heights
}

func (t *Align) colWidths(bbs [][]boundingbox.BB) []float64 {
	if len(t.eqs[0]) == 0 {
		panic("colWidths should never be 0 length")
	}

	widths := make([]float64, len(t.eqs[0]))

	if t.constantHorSpacing {
		d := t.maxWidth(bbs) + t.minHorSpacing

		for i, _ := range t.eqs[0] {
			widths[i] = d
		}
	} else {
		for i, _ := range t.eqs[0] {
			widths[i] = t.maxColWidth(bbs, i) + t.minHorSpacing
		}
	}

	return widths
}

func (t *Align) accumulate(arr []float64) []float64 {
	res := make([]float64, len(arr))

	if len(res) == 0 {
		panic("should never be 0 length")
	}

	res[0] = 0.0

	for i, a := range arr {

		if i < len(arr)-1 {
			res[i+1] = res[i] + a
		}
	}

	return res
}

func (t *Align) GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error) {
	subScopes := make([][]SubScope, 0)
	bbs := make([][]boundingbox.BB, 0)

	for _, row := range t.eqs {
		subScopeRow := []SubScope{}
		bbsRow := []boundingbox.BB{}
		for _, eq := range row {
			subScope := scope.NewSubScope()

			bb, err := eq.GenerateTags(subScope, 0.0, 0.0)
			if err != nil {
				return nil, err
			}

			subScopeRow = append(subScopeRow, subScope)
			bbsRow = append(bbsRow, bb)
		}

		subScopes = append(subScopes, subScopeRow)
		bbs = append(bbs, bbsRow)
	}

	hs := t.accumulate(t.rowHeights(bbs)) // first is zero, last is lost
	ws := t.accumulate(t.colWidths(bbs))

	var bbTotal boundingbox.BB = nil
	for i, row := range t.eqs {
		subScopeRow := subScopes[i]
		bbsRow := bbs[i]
		for j, _ := range row {
			dx := x + ws[j]
			dy := y + hs[i]

			subScopeRow[j].Transform(dx, dy, 1.0, 1.0)

			bb := bbsRow[j].Translate(dx, dy)
			if bbTotal == nil {
				bbTotal = bb
			} else {
				bbTotal = boundingbox.Merge(bbTotal, bb)
			}
		}
	}

	return bbTotal, nil
}
