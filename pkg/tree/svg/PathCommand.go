package svg

import (
	"fmt"

	"../../tree"

	"../../tokens/context"
)

var (
	COMPRESS = false
)

// other operations are done by inference
type PathCommand interface {
	Write() string
	Compress(tree.SVGViewBox)
	Context() context.Context
}

type PathCommandData struct {
	ctx context.Context
}

func (c *PathCommandData) Context() context.Context {
	return c.ctx
}

func writePair(x, y float64) string {
	space := " "
	if y < 0.0 {
		space = ""
	}

	return fmt.Sprintf("%g%s%g", x, space, y)
}

func writeSingle(x float64) string {
	return fmt.Sprintf("%g", x)
}

func writeAdditionalPair(x, y float64) string {
	space := " "
	if x < 0.0 {
		space = ""
	}

	return fmt.Sprintf("%s%s", space, writePair(x, y))
}

func writeAdditionalSingle(x float64) string {
	space := " "
	if x < 0.0 {
		space = ""
	}
	return fmt.Sprintf("%s%s", space, writeSingle(x))
}

type MoveTo struct {
	PathCommandData
	x, y float64
}

func NewMoveTo(x, y float64, ctx context.Context) *MoveTo {
	return &MoveTo{PathCommandData{ctx}, x, y}
}

func (c *MoveTo) Compress(vb tree.SVGViewBox) {
	c.x, c.y = vb.CompressXY(c.x, c.y)
}

func (c *MoveTo) Write() string {
	return fmt.Sprintf("M%s", writePair(c.x, c.y))
}

type MoveBy struct {
	PathCommandData
	dx, dy float64
}

func NewMoveBy(dx, dy float64, ctx context.Context) *MoveBy {
	return &MoveBy{PathCommandData{ctx}, dx, dy}
}

func (c *MoveBy) Compress(vb tree.SVGViewBox) {
	c.dx, c.dy = vb.CompressDXDY(c.dx, c.dy)
}

func (c *MoveBy) Write() string {
	return fmt.Sprintf("m%s", writePair(c.dx, c.dy))
}

type LineTo struct {
	PathCommandData
	x, y float64
}

func NewLineTo(x, y float64, ctx context.Context) *LineTo {
	return &LineTo{PathCommandData{ctx}, x, y}
}

func (c *LineTo) Compress(vb tree.SVGViewBox) {
	c.x, c.y = vb.CompressXY(c.x, c.y)
}

func (c *LineTo) Write() string {
	return fmt.Sprintf("L%s", writePair(c.x, c.y))
}

type LineBy struct {
	PathCommandData
	dx, dy float64
}

func NewLineBy(dx, dy float64, ctx context.Context) *LineBy {
	return &LineBy{PathCommandData{ctx}, dx, dy}
}

func (c *LineBy) Compress(vb tree.SVGViewBox) {
	c.dx, c.dy = vb.CompressDXDY(c.dx, c.dy)
}

func (c *LineBy) Write() string {
	return fmt.Sprintf("l%s", writePair(c.dx, c.dy))
}

type HorTo struct {
	PathCommandData
	x float64
}

func NewHorTo(x float64, ctx context.Context) *HorTo {
	return &HorTo{PathCommandData{ctx}, x}
}

func (c *HorTo) Compress(vb tree.SVGViewBox) {
	c.x = vb.CompressX(c.x)
}

func (c *HorTo) Write() string {
	return fmt.Sprintf("H%s", writeSingle(c.x))
}

type VerTo struct {
	PathCommandData
	y float64
}

func NewVerTo(y float64, ctx context.Context) *VerTo {
	return &VerTo{PathCommandData{ctx}, y}
}

func (c *VerTo) Compress(vb tree.SVGViewBox) {
	c.y = vb.CompressY(c.y)
}

func (c *VerTo) Write() string {
	return fmt.Sprintf("V%s", writeSingle(c.y))
}

type HorBy struct {
	PathCommandData
	dx float64
}

func NewHorBy(dx float64, ctx context.Context) *HorBy {
	return &HorBy{PathCommandData{ctx}, dx}
}

func (c *HorBy) Compress(vb tree.SVGViewBox) {
	c.dx = vb.CompressDX(c.dx)
}

func (c *HorBy) Write() string {
	return fmt.Sprintf("h%s", writeSingle(c.dx))
}

type VerBy struct {
	PathCommandData
	dy float64
}

func NewVerBy(dy float64, ctx context.Context) *VerBy {
	return &VerBy{PathCommandData{ctx}, dy}
}

func (c *VerBy) Compress(vb tree.SVGViewBox) {
	c.dy = vb.CompressDY(c.dy)
}

func (c *VerBy) Write() string {
	return fmt.Sprintf("v%s", writeSingle(c.dy))
}

type Close struct {
	PathCommandData
}

func NewClose(ctx context.Context) *Close {
	return &Close{PathCommandData{ctx}}
}

func (c *Close) Compress(vb tree.SVGViewBox) {
	return
}

func (c *Close) Write() string {
	return "z"
}

type QuadraticTo struct {
	PathCommandData
	x1, y1 float64 // control point
	x, y   float64
}

func NewQuadraticTo(x1, y1, x, y float64, ctx context.Context) *QuadraticTo {
	return &QuadraticTo{PathCommandData{ctx}, x1, y1, x, y}
}

func (c *QuadraticTo) Compress(vb tree.SVGViewBox) {
	c.x1, c.y1 = vb.CompressXY(c.x1, c.y1)
	c.x, c.y = vb.CompressXY(c.x, c.y)
}

func (c *QuadraticTo) Write() string {
	return fmt.Sprintf("Q%s%s", writePair(c.x1, c.y1), writeAdditionalPair(c.x, c.y))
}

type QuadraticBy struct {
	PathCommandData
	dx1, dy1 float64
	dx, dy   float64
}

func NewQuadraticBy(dx1, dy1, dx, dy float64, ctx context.Context) *QuadraticBy {
	return &QuadraticBy{PathCommandData{ctx}, dx1, dy1, dx, dy}
}

func (c *QuadraticBy) Compress(vb tree.SVGViewBox) {
	c.dx1, c.dy1 = vb.CompressDXDY(c.dx1, c.dy1)
	c.dx, c.dy = vb.CompressDXDY(c.dx, c.dy)
}

func (c *QuadraticBy) Write() string {
	return fmt.Sprintf("q%s%s", writePair(c.dx1, c.dy1), writeAdditionalPair(c.dx, c.dy))
}

type ExtraQuadraticTo struct {
	PathCommandData
	x, y float64
}

func NewExtraQuadraticTo(x, y float64, ctx context.Context) *ExtraQuadraticTo {
	return &ExtraQuadraticTo{PathCommandData{ctx}, x, y}
}

func (c *ExtraQuadraticTo) Compress(vb tree.SVGViewBox) {
	c.x, c.y = vb.CompressXY(c.x, c.y)
}

func (c *ExtraQuadraticTo) Write() string {
	return fmt.Sprintf("T%s", writePair(c.x, c.y))
}

type ExtraQuadraticBy struct {
	PathCommandData
	dx, dy float64
}

func NewExtraQuadraticBy(dx, dy float64, ctx context.Context) *ExtraQuadraticBy {
	return &ExtraQuadraticBy{PathCommandData{ctx}, dx, dy}
}

func (c *ExtraQuadraticBy) Compress(vb tree.SVGViewBox) {
	c.dx, c.dy = vb.CompressDXDY(c.dx, c.dy)
}

func (c *ExtraQuadraticBy) Write() string {
	return fmt.Sprintf("t%s", writePair(c.dx, c.dy))
}

type CubicTo struct {
	PathCommandData
	x1, y1 float64
	x2, y2 float64
	x, y   float64
}

func NewCubicTo(x1, y1, x2, y2, x, y float64, ctx context.Context) *CubicTo {
	return &CubicTo{PathCommandData{ctx}, x1, y1, x2, y2, x, y}
}

func (c *CubicTo) Compress(vb tree.SVGViewBox) {
	c.x1, c.y1 = vb.CompressXY(c.x1, c.y1)
	c.x2, c.y2 = vb.CompressXY(c.x2, c.y2)
	c.x, c.y = vb.CompressXY(c.x, c.y)
}

func (c *CubicTo) Write() string {
	return fmt.Sprintf("C%s%s%s",
		writePair(c.x1, c.y1),
		writeAdditionalPair(c.x2, c.y2),
		writeAdditionalPair(c.x, c.y))
}

type CubicBy struct {
	PathCommandData
	dx1, dy1 float64
	dx2, dy2 float64
	dx, dy   float64
}

func NewCubicBy(dx1, dy1, dx2, dy2, dx, dy float64, ctx context.Context) *CubicBy {
	return &CubicBy{PathCommandData{ctx}, dx1, dy1, dx2, dy2, dx, dy}
}

func (c *CubicBy) Compress(vb tree.SVGViewBox) {
	c.dx1, c.dy1 = vb.CompressDXDY(c.dx1, c.dy1)
	c.dx2, c.dy2 = vb.CompressDXDY(c.dx2, c.dy2)
	c.dx, c.dy = vb.CompressDXDY(c.dx, c.dy)
}

func (c *CubicBy) Write() string {
	return fmt.Sprintf("c%s%s%s",
		writePair(c.dx1, c.dy1),
		writeAdditionalPair(c.dx2, c.dy2),
		writeAdditionalPair(c.dx, c.dy))
}

type ExtraCubicTo struct {
	PathCommandData
	x2, y2 float64
	x, y   float64
}

func NewExtraCubicTo(x2, y2, x, y float64, ctx context.Context) *ExtraCubicTo {
	return &ExtraCubicTo{PathCommandData{ctx}, x2, y2, x, y}
}

func (c *ExtraCubicTo) Compress(vb tree.SVGViewBox) {
	c.x2, c.y2 = vb.CompressXY(c.x2, c.y2)
	c.x, c.y = vb.CompressXY(c.x, c.y)
}

func (c *ExtraCubicTo) Write() string {
	return fmt.Sprintf("S%s%s", writePair(c.x2, c.y2), writeAdditionalPair(c.x, c.y))
}

type ExtraCubicBy struct {
	PathCommandData
	dx2, dy2 float64
	dx, dy   float64
}

func NewExtraCubicBy(dx2, dy2, dx, dy float64, ctx context.Context) *ExtraCubicBy {
	return &ExtraCubicBy{PathCommandData{ctx}, dx2, dy2, dx, dy}
}

func (c *ExtraCubicBy) Compress(vb tree.SVGViewBox) {
	c.dx2, c.dy2 = vb.CompressDXDY(c.dx2, c.dy2)
	c.dx, c.dy = vb.CompressDXDY(c.dx, c.dy)
}

func (c *ExtraCubicBy) Write() string {
	return fmt.Sprintf("s%s%s", writePair(c.dx2, c.dy2), writeAdditionalPair(c.dx, c.dy))
}

type ArcCoreData struct {
	rx, ry                  float64
	xAxisRot                float64
	largeArc, positiveSweep bool
}

func (c *ArcCoreData) Compress(vb tree.SVGViewBox) {
	c.rx, c.ry = vb.CompressDXDY(c.rx, c.ry)
	c.xAxisRot = vb.CompressScalar(c.xAxisRot)
}

func (c *ArcCoreData) writeCore() string {
	las := "0"
	if c.largeArc {
		las = "1"
	}

	pos := "0"
	if c.positiveSweep {
		pos = "1"
	}

	return fmt.Sprintf("%s%s %s %s",
		writePair(c.rx, c.ry),
		writeAdditionalSingle(c.xAxisRot),
		las,
		pos)
}

type ArcTo struct {
	ArcCoreData
	PathCommandData
	x, y float64
}

func NewArcTo(x, y float64, rx, ry, xAxisRot float64, largeArc, positiveSweep bool, ctx context.Context) *ArcTo {
	return &ArcTo{ArcCoreData{rx, ry, xAxisRot, largeArc, positiveSweep}, PathCommandData{ctx}, x, y}
}

func (c *ArcTo) Compress(vb tree.SVGViewBox) {
	c.ArcCoreData.Compress(vb)
	c.x, c.y = vb.CompressXY(c.x, c.y)
}

func (c *ArcTo) Write() string {
	return fmt.Sprintf("A%s%s", c.ArcCoreData.writeCore(), writeAdditionalPair(c.x, c.y))
}

type ArcBy struct {
	ArcCoreData
	PathCommandData
	dx, dy float64
}

func NewArcBy(dx, dy float64, rx, ry, xAxisRot float64, largeArc, positiveSweep bool, ctx context.Context) *ArcBy {
	return &ArcBy{ArcCoreData{rx, ry, xAxisRot, largeArc, positiveSweep}, PathCommandData{ctx}, dx, dy}
}

func (c *ArcBy) Compress(vb tree.SVGViewBox) {
	c.ArcCoreData.Compress(vb)
	c.dx, c.dy = vb.CompressDXDY(c.dx, c.dy)
}

func (c *ArcBy) Write() string {
	return fmt.Sprintf("a%s%s", c.ArcCoreData.writeCore(), writeAdditionalPair(c.dx, c.dy))
}
