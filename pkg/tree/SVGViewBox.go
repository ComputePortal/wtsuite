package tree

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// wrt the viewBox
var REL_PRECISION = 7
var ABS_PRECISION = 5

var GEN_PRECISION = 3 // == -math.Log10(0.001)

type SVGViewBox interface {
	CompressXY(x float64, y float64) (float64, float64)
	CompressX(x float64) float64
	CompressY(y float64) float64
	CompressDXDY(dx float64, dy float64) (float64, float64)
	CompressDX(dx float64) float64
	CompressDY(dy float64) float64
	CompressScalar(s float64) float64

	CompressSelf() string
}

type SVGViewBoxData struct {
	xMin, yMin     float64
	xMax, yMax     float64
	width, height  float64
	precX, precY   float64
	precDX, precDY float64
	precGen        float64
}

func NewViewBoxFromString(s string, ctx context.Context) (SVGViewBox, error) {
	fs := strings.Fields(s)
	if len(fs) != 4 {
		return nil, ctx.NewError("Error: expected 4 viewBox numbers")
	}

	xMin, err := strconv.ParseFloat(fs[0], 64)
	if err != nil {
		return nil, ctx.NewError("Error: failed to parse first number")
	}

	yMin, err := strconv.ParseFloat(fs[1], 64)
	if err != nil {
		return nil, ctx.NewError("Error: failed to parse second number")
	}

	width, err := strconv.ParseFloat(fs[2], 64)
	if err != nil {
		return nil, ctx.NewError("Error: failed to parse third number")
	}

	height, err := strconv.ParseFloat(fs[3], 64)
	if err != nil {
		return nil, ctx.NewError("Error: failed to parse fourth number")
	}

	if width <= 0.0 {
		return nil, ctx.NewError("Error: non-positive width")
	}

	if height <= 0.0 {
		return nil, ctx.NewError("Error: non-positive width")
	}

	precX := math.Pow10(int(math.Round(math.Log10(width))) - ABS_PRECISION)
	precY := math.Pow10(int(math.Round(math.Log10(height))) - ABS_PRECISION)
	precDX := math.Pow10(int(math.Round(math.Log10(width))) - REL_PRECISION)
	precDY := math.Pow10(int(math.Round(math.Log10(height))) - REL_PRECISION)
	precGen := math.Pow10(int(math.Round(-float64(GEN_PRECISION))))

	return &SVGViewBoxData{xMin, yMin, xMin + width, yMin + height, width, height, precX, precY, precDX, precDY, precGen}, nil
}

// eg. prec = 0.00001
func RoundScalar(s float64, prec float64) float64 {
	result := math.Round(s/prec) * prec

	result, _ = strconv.ParseFloat(fmt.Sprintf("%0.8f", result), 64)

	return result
}

func (vb *SVGViewBoxData) CompressScalar(s float64) float64 {
	return RoundScalar(s, vb.precGen)
}

func (vb *SVGViewBoxData) CompressX(x float64) float64 {
	return RoundScalar(x, vb.precX)
}

func (vb *SVGViewBoxData) CompressY(y float64) float64 {
	return RoundScalar(y, vb.precY)
}

func (vb *SVGViewBoxData) CompressXY(x float64, y float64) (float64, float64) {
	return vb.CompressX(x), vb.CompressY(y)
}

func (vb *SVGViewBoxData) CompressDX(dx float64) float64 {
	return RoundScalar(dx, vb.precDX)
}

func (vb *SVGViewBoxData) CompressDY(dy float64) float64 {
	return RoundScalar(dy, vb.precDY)
}

func (vb *SVGViewBoxData) CompressDXDY(dx float64, dy float64) (float64, float64) {
	return vb.CompressDX(dx), vb.CompressDY(dy)
}

func (vb *SVGViewBoxData) CompressSelf() string {
	xMin := vb.CompressX(vb.xMin)
	yMin := vb.CompressY(vb.yMin)

	width := vb.CompressX(vb.width)
	height := vb.CompressY(vb.height)

	return fmt.Sprintf("%g %g %g %g", xMin, yMin, width, height)
}
