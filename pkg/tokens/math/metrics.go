package math

import (
	"fmt"

	"./boundingbox"
	"./serif"
)

const (
	lineThickness        = 0.066
	plusMinusWidth       = 0.564
	plusMinusFracYOffset = 0.286 // from baseline upwards
	plusMinusXOffset     = 0.030
	extraAddSubSpacing   = 0.300

	genericBinSymbolWidth   = 0.564
	genericBinSymbolSpacing = 0.300
	genericMinHorSpacing    = 2.0
	genericMinVerSpacing    = 0.2
)

func unicodeAdvanceWidth(unicode int) float64 {
	if w, ok := serif.AdvanceWidths[unicode]; ok {
		return float64(w) / float64(serif.UnitsPerEm)
	} else {
		panic(fmt.Sprintf("character 0x%x not yet available in advance width map", unicode))
	}
}

func unicodeBB(unicode int) boundingbox.BB {
	if bb, ok := serif.Bounds[unicode]; ok {
		return boundingbox.NewBB(
			float64(bb.Left())/float64(serif.UnitsPerEm),
			float64(bb.Top())/float64(serif.UnitsPerEm),
			float64(bb.Right())/float64(serif.UnitsPerEm),
			float64(bb.Bottom())/float64(serif.UnitsPerEm),
		)
	} else {
		panic(fmt.Sprintf("character 0x%x not yet available in bb map", unicode))
	}
}

func minMax(isMin bool, xs []float64) float64 {
	b := false

	xMinMax := 0.0

	for _, x := range xs {
		if !b {
			xMinMax = x
			b = true
		} else if isMin && (x < xMinMax) {
			xMinMax = x
		} else if (!isMin) && (x > xMinMax) {
			xMinMax = x
		}
	}

	if !b {
		panic("empty list")
	}

	return xMinMax
}

func max(xs ...float64) float64 {
	return minMax(false, xs)
}

func min(xs ...float64) float64 {
	return minMax(true, xs)
}
