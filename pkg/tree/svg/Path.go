package svg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type Path struct {
	tree.SVGTagData
	d []PathCommand
}

func BuildPath(attr *tokens.StringDict, d []PathCommand, ctx context.Context) (tree.Tag, error) {
	svgTag, err := tree.NewSVGTagData("path", attr, ctx)
	return &Path{svgTag, d}, err
}

func (t *Path) Validate() error {
	return nil
}

func (t *Path) Write(indent string, nl, tab string) string {
	var d strings.Builder

	for _, pc := range t.d {
		d.WriteString(pc.Write())
	}

	valueToken := tokens.NewValueString(d.String(), t.d[0].Context())
	t.SetTmpAttribute("d", valueToken)

	result := t.SVGTagData.Write(indent, nl, tab)

	t.RemoveTmpAttribute("d")

	return result
}

func (t *Path) Compress(vb tree.SVGViewBox) {
	attr := t.Attributes()

	CompressStyles(attr, vb)

	for _, pc := range t.d {
		pc.Compress(vb)
	}
}

func (t *Path) Minify() bool {
	if len(t.d) == 0 {
		return true
	}

	attr := t.Attributes()
	attr.Delete("inkscape:connector-curvature")
	attr.Delete("sodipodi:nodetypes")

	MinifyStyles(attr, false)

	t.SetID("") // probably not used

	return false
}

// parsing of d
func assertNumber(fs []string, i int, prevError error, ctx context.Context) (float64, int, error) {
	if prevError != nil {
		return 0.0, i, prevError
	}

	if i > len(fs)-1 {
		err := ctx.NewError("Error: expected more fields")
		panic(err)
		return 0, i, err
	}

	s := fs[i]

	x, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, i, ctx.NewError(fmt.Sprintf("Error: '"+s+"' is not a valid number (position %d)", i))
	}

	if x > 0.0 && s[0] == '-' {
		panic("number string is negative, but parsed float isnt")
	}

	return x, i + 1, nil
}

func assertBool(fs []string, i int, prevError error, ctx context.Context) (bool, int, error) {
	if prevError != nil {
		return false, i, prevError
	}

	if i > len(fs)-1 {
		return false, i, ctx.NewError("Error: expected more fields")
	}

	switch fs[i] {
	case "0":
		return false, i + 1, nil
	case "1":
		return true, i + 1, nil
	default:
		errCtx := ctx
		return false, i, errCtx.NewError("Error: expected 0 or 1, got " + fs[i])
	}
}

// recursive
func parsePathCommand(key string, fs []string, i int, result []PathCommand, ctx context.Context) ([]PathCommand, error) {
	if key == "" {
		key = fs[i]
		i++
	}

	var err error = nil
	var x, y, x1, y1, x2, y2, dx, dy, dx1, dy1, dx2, dy2 float64
	switch key {
	case "M":
		x, i, err = assertNumber(fs, i, err, ctx)
		y, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewMoveTo(x, y, ctx))
		switch {
		case i > len(fs)-1:
			return result, nil
		case !patterns.SVGPATH_LETTER_REGEXP.MatchString(fs[i]):
			return parsePathCommand("L", fs, i, result, ctx)
		}
	case "m":
		dx, i, err = assertNumber(fs, i, err, ctx)
		dy, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewMoveBy(dx, dy, ctx))
		switch {
		case i > len(fs)-1:
			return result, nil
		case !patterns.SVGPATH_LETTER_REGEXP.MatchString(fs[i]):
			return parsePathCommand("l", fs, i, result, ctx)
		}
	case "L":
		x, i, err = assertNumber(fs, i, err, ctx)
		y, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewLineTo(x, y, ctx))
	case "l":
		dx, i, err = assertNumber(fs, i, err, ctx)
		dy, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewLineBy(dx, dy, ctx))
	case "H":
		x, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewHorTo(x, ctx))
	case "h":
		dx, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewHorBy(dx, ctx))
	case "V":
		y, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewVerTo(y, ctx))
	case "v":
		dy, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewVerBy(dy, ctx))
	case "Q":
		x1, i, err = assertNumber(fs, i, err, ctx)
		y1, i, err = assertNumber(fs, i, err, ctx)
		x, i, err = assertNumber(fs, i, err, ctx)
		y, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewQuadraticTo(x1, y1, x, y, ctx))
	case "q":
		dx1, i, err = assertNumber(fs, i, err, ctx)
		dy1, i, err = assertNumber(fs, i, err, ctx)
		dx, i, err = assertNumber(fs, i, err, ctx)
		dy, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewQuadraticBy(dx1, dy1, dx, dy, ctx))
	case "T":
		x, i, err = assertNumber(fs, i, err, ctx)
		y, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewExtraQuadraticTo(x, y, ctx))
	case "t":
		dx, i, err = assertNumber(fs, i, err, ctx)
		dy, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewExtraQuadraticBy(dx, dy, ctx))
	case "C":
		x1, i, err = assertNumber(fs, i, err, ctx)
		y1, i, err = assertNumber(fs, i, err, ctx)
		x2, i, err = assertNumber(fs, i, err, ctx)
		y2, i, err = assertNumber(fs, i, err, ctx)
		x, i, err = assertNumber(fs, i, err, ctx)
		y, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewCubicTo(x1, y1, x2, y2, x, y, ctx))
	case "c":
		dx1, i, err = assertNumber(fs, i, err, ctx)
		dy1, i, err = assertNumber(fs, i, err, ctx)
		dx2, i, err = assertNumber(fs, i, err, ctx)
		dy2, i, err = assertNumber(fs, i, err, ctx)
		dx, i, err = assertNumber(fs, i, err, ctx)
		dy, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewCubicBy(dx1, dy1, dx2, dy2, dx, dy, ctx))
	case "S":
		x2, i, err = assertNumber(fs, i, err, ctx)
		y2, i, err = assertNumber(fs, i, err, ctx)
		x, i, err = assertNumber(fs, i, err, ctx)
		y, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewExtraCubicTo(x2, y2, x, y, ctx))
	case "s":
		dx2, i, err = assertNumber(fs, i, err, ctx)
		dy2, i, err = assertNumber(fs, i, err, ctx)
		dx, i, err = assertNumber(fs, i, err, ctx)
		dy, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewExtraCubicBy(dx2, dy2, dx, dy, ctx))
	case "A":
		var rx, ry, xAxisRot float64
		var largeArc, positiveSweep bool
		rx, i, err = assertNumber(fs, i, err, ctx)
		ry, i, err = assertNumber(fs, i, err, ctx)
		xAxisRot, i, err = assertNumber(fs, i, err, ctx)
		largeArc, i, err = assertBool(fs, i, err, ctx)
		positiveSweep, i, err = assertBool(fs, i, err, ctx)
		x, i, err = assertNumber(fs, i, err, ctx)
		y, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewArcTo(x, y, rx, ry, xAxisRot, largeArc, positiveSweep, ctx))
	case "a":
		var rx, ry, xAxisRot float64
		var largeArc, positiveSweep bool
		rx, i, err = assertNumber(fs, i, err, ctx)
		ry, i, err = assertNumber(fs, i, err, ctx)
		xAxisRot, i, err = assertNumber(fs, i, err, ctx)
		largeArc, i, err = assertBool(fs, i, err, ctx)
		positiveSweep, i, err = assertBool(fs, i, err, ctx)
		dx, i, err = assertNumber(fs, i, err, ctx)
		dy, i, err = assertNumber(fs, i, err, ctx)
		result = append(result, NewArcBy(dx, dy, rx, ry, xAxisRot, largeArc, positiveSweep, ctx))
	case "z", "Z":
		result = append(result, NewClose(ctx))
	default:
		errCtx := ctx
		return nil, errCtx.NewError("Error: unrecognized path command '" + key + "'")
	}

	if err != nil {
		return nil, err
	}

	switch {
	case i > len(fs)-1:
		return result, nil
	case patterns.SVGPATH_LETTER_REGEXP.MatchString(fs[i]):
		key = fs[i]
		i++
		return parsePathCommand(key, fs, i, result, ctx)
	default:
		return parsePathCommand(key, fs, i, result, ctx)
	}
}

func ParsePathString(d string, ctx context.Context) ([]PathCommand, error) {
	// replace all letters by spacing around them
	d = patterns.SVGPATH_LETTER_REGEXP.ReplaceAllString(d, " $0 ")

	// add a space to left of all minus signs
	d = patterns.SVGPATH_MINUS_REGEXP.ReplaceAllString(d, "$1 $2")

	// replace all commas by spaces
	d = patterns.COMMA_REGEXP.ReplaceAllString(d, " ")

	fs := strings.Fields(d)

  fsFinal := make([]string, 0)
  // add a zero to all fields starting with a dot, and split if there is more than one dot
  for _, f_ := range fs {
    f := f_
    if strings.HasPrefix(patterns.PERIOD, f) {
      f = "0" + f
    } 

    parts := strings.Split(f, patterns.PERIOD)

    if len(parts) == 1 {
      fsFinal = append(fsFinal, parts[0])
    } else if len(parts) == 2 {
      fsFinal = append(fsFinal, f)
    } else {
      fsFinal = append(fsFinal, parts[0] + patterns.PERIOD + parts[1])

      for _, rem := range parts[2:] {
        fsFinal = append(fsFinal, "0" + patterns.PERIOD + rem)
      }
    }
  }

	if len(fsFinal) == 0 {
		return nil, ctx.NewError("Error: bad path")
	}

	result := make([]PathCommand, 0)


	return parsePathCommand("", fsFinal, 0, result, ctx)
}
