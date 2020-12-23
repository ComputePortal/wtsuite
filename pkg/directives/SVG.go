package directives

import (
	"fmt"
	"math"
	"net/url"
	"strings"

	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/svg"
)

func BuildSVG(scope Scope, node Node, tag *tokens.Tag) error {
	return buildTree(scope, node, SVG, tag, "")
}

func buildSVGPathInternal(node Node, attr *tokens.StringDict, ctx context.Context) error {
	d, err := tokens.DictString(attr, "d")
	if err != nil {
		return err
	}

	pcs, err := svg.ParsePathString(d.Value(), d.Context())
	if err != nil {
		return err
	}

	attr.Delete("d")

	childTag, err := svg.BuildPath(attr, pcs, ctx)
	if err != nil {
		return err
	}

	return node.AppendChild(childTag)
}

func buildSVGPath(scope Scope, node Node, tag *tokens.Tag) error {
	if err := tag.AssertEmpty(); err != nil {
		return err
	}

	subScope := NewSubScope(scope)

	attr, err := buildAttributes(subScope, tag, []string{"d"})
	if err != nil {
		return err
	}

	return buildSVGPathInternal(node, attr, tag.Context())
}

// return a filled head if sw == 0.0, otherwise stroked
func buildSVGArrowHead(l, x0, y0, tx, ty, sw float64, ctx context.Context) (tree.Tag, error) {
	nx := -ty
	ny := tx

	h := 0.66 * l

	angle := math.Atan(0.5 * h / l)

	tipOffset := 0.5 * sw / math.Sin(angle)

	backOffsetT := 0.5 * sw * math.Sin(angle)
	backOffsetN := sw/math.Cos(angle) - 0.5*sw*math.Cos(angle)

	/*pcs := []svg.PathCommand{
		svg.NewMoveTo(x0, y0, ctx),
		svg.NewLineBy(-l*tx-0.5*h*nx, -l*ty-0.5*h*ny, ctx),
		svg.NewLineBy(h*nx, h*ny, ctx),
		svg.NewClose(ctx),
	}*/

	offsetL := l - backOffsetT
	offsetH := 0.5*h - backOffsetN
	pcs := []svg.PathCommand{
		svg.NewMoveTo(x0-offsetL*tx-offsetH*nx, y0-offsetL*ty-offsetH*ny, ctx),
		svg.NewLineTo(x0-tipOffset*tx, y0-tipOffset*ty, ctx),
		svg.NewLineTo(x0-offsetL*tx+offsetH*nx, y0-offsetL*ty+offsetH*ny, ctx),
	}

	attr := tokens.NewEmptyStringDict(ctx)
	styleVal := tokens.NewEmptyStringDict(ctx)

	if sw == 0.0 {
		strokeVal := tokens.NewValueString("none", ctx)
		styleVal.Set("stroke", strokeVal)
	} else {
		fillVal := tokens.NewValueString("none", ctx)
		styleVal.Set("fill", fillVal)
	}

	attr.Set("style", styleVal)

	tag, err := svg.BuildPath(attr, pcs, ctx)
	if err != nil {
		panic(err)
	}

	return tag, nil
}

func searchSVGStrokeWidth(node Node, scope Scope, attr *tokens.StringDict,
	ctx context.Context) (*tokens.Float, error) {
	swToken_, err := SearchStyle(node, scope, attr, "stroke-width", ctx)
	if err != nil {
		return nil, err
	}

	if tokens.IsNull(swToken_) {
		return nil, ctx.NewError("Error: no stroke-width found, dont know how big to make the arrow")
	} else {
		return tokens.AssertIntOrFloat(swToken_)
	}
}

func buildSVGArrow(scope Scope, node Node, tag *tokens.Tag) error {
	if err := tag.AssertEmpty(); err != nil {
		return err
	}

	ctx := tag.Context()

	subScope := NewSubScope(scope)

	attr, err := buildAttributes(subScope, tag, []string{"d", "type"})
	if err != nil {
		return err
	}

	d, err := tokens.DictString(attr, "d")
	if err != nil {
		return err
	}

	pcs, err := svg.ParsePathString(d.Value(), d.Context())
	if err != nil {
		return err
	}

	attr.Delete("d")

	// TODO: get size from stroke-width
	sw, err := searchSVGStrokeWidth(node, scope, attr, tag.Context())
	if err != nil {
		return err
	}

	headSize := sw.Value() * 5

	typeToken, err := tokens.DictString(attr, "type")
	if err != nil {
		return err
	}

	attr.Delete("type")

	typeStr := typeToken.Value()

	hasStart := false
	hasStop := false
	switch typeStr {
	case "<":
		hasStart = true
	case ">":
		hasStop = true
	case "<>":
		hasStart = true
		hasStop = true
	default:
		errCtx := typeToken.Context()
		return errCtx.NewError("Error: invalid type '" + typeStr + "'")
	}

	tags := make([]tree.Tag, 0)
	segments, err := svg.GenerateSegments(pcs, d.Context())
	if err != nil {
		return err
	}

	pathOffset := sw.Value() / math.Sin(math.Atan(0.5*0.66))
	startLen := 0.0
	if hasStart {
		startLen = pathOffset

		tipX, tipY := segments[0].Position(0.0)
		tanX, tanY := segments[0].Tangent(0.0)

		startArrow, err := buildSVGArrowHead(headSize, tipX, tipY, -tanX, -tanY, sw.Value(), ctx)
		if err != nil {
			return err
		}

		tags = append(tags, startArrow)

		startArrowFilled, err := buildSVGArrowHead(headSize, tipX, tipY, -tanX, -tanY, 0.0, ctx)
		if err != nil {
			return err
		}

		tags = append(tags, startArrowFilled)
	}

	stopLen := 0.0
	if hasStop {
		stopLen = pathOffset

		iLast := len(segments) - 1
		tipX, tipY := segments[iLast].Position(1.0)
		tanX, tanY := segments[iLast].Tangent(1.0)

		stopArrow, err := buildSVGArrowHead(headSize, tipX, tipY, tanX, tanY, sw.Value(), ctx)
		if err != nil {
			return err
		}

		tags = append(tags, stopArrow)

		stopArrowFilled, err := buildSVGArrowHead(headSize, tipX, tipY, tanX, tanY, 0.0, ctx)
		if err != nil {
			return err
		}

		tags = append(tags, stopArrowFilled)
	}

	pcs, err = svg.ShortenEnds(pcs, startLen, stopLen, d.Context())
	if err != nil {
		return err
	}

	pathAttr_, err := attr.Copy(ctx)
	if err != nil {
		return err
	}
	pathAttr, ok := pathAttr_.(*tokens.StringDict)
	if !ok {
		panic("unexpected")
	}
	pathStyle := tokens.NewEmptyStringDict(ctx)
	pathStyle.Set("fill", tokens.NewValueString("none", ctx))
	pathAttr.Set("style", pathStyle)
	pathTag, err := svg.BuildPath(pathAttr, pcs, ctx)
	if err != nil {
		return err
	}

	// actually prepend the tag, so it is below the arrowheads
	tags = append([]tree.Tag{pathTag}, tags...)

	// wrap everything in a group
	groupTag, err := svg.BuildTag("g", attr, ctx)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		groupTag.AppendChild(tag)
	}

	return node.AppendChild(groupTag)
}

func svgToURI(tag tree.Tag, ctx context.Context) (tokens.Token, error) {
	tag.FoldDummy()
	if tag.Name() == "dummy" {
		if tag.NumChildren() != 1 {
			errCtx := tag.Context()
			return nil, errCtx.NewError(fmt.Sprintf("Error: expected 1 child, got %d\n", tag.NumChildren()))
		}
		tag = tag.Children()[0]
	}

	if tag.Name() != "svg" {
		return nil, ctx.NewError("Error: expected svg image, got " + tag.Name())
	}

	tagAttr := tag.Attributes()

	// also add the xmlns info
	xmlnsKeyToken := tokens.NewValueString("xmlns", ctx)
	xmlnsValToken := tokens.NewValueString("http://www.w3.org/2000/svg", ctx)

	xmlnsXlinkKeyToken := tokens.NewValueString("xmlns:xlink", ctx)
	xmlnsXlinkValToken := tokens.NewValueString("http://www.w3.org/1999/xlink", ctx)

	tagAttr.Set(xmlnsKeyToken, xmlnsValToken)
	tagAttr.Set(xmlnsXlinkKeyToken, xmlnsXlinkValToken)

	var b strings.Builder

	b.WriteString("url('data:image/svg+xml;utf8,")
	b.WriteString(url.PathEscape(tag.Write("", "", "")))
	b.WriteString("')")

	return tokens.NewValueString(b.String(), ctx), nil
}

func evalSVGURI(scope Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

	args, err := functions.EvalArgs(scope, args)
	if err != nil {
		return nil, err
	}

	templateToken, err := tokens.AssertString(args[0])
	if err != nil {
		return nil, err
	}

	if !scope.HasTemplate(templateToken.Value()) {
		errCtx := templateToken.Context()
		return nil, errCtx.NewError("Error: template '" + templateToken.Value() + "' not found")
	}

	// get the def
	attrToken, err := tokens.AssertStringDict(args[1])
	if err != nil {
		return nil, err
	}

	imageTag := tokens.NewTag(templateToken.Value(), attrToken.ToRaw(), []*tokens.Tag{}, ctx)

	uriNode := NewURINode()
	if err := BuildTemplate(scope, uriNode, imageTag); err != nil {
		return nil, err
	}

	if uriNode.tag == nil {
		errCtx := ctx
		return nil, errCtx.NewError("Error: no tag appended")
	}

	return svgToURI(uriNode.tag, ctx)
}

var _svgOk = registerDirective("svg", BuildSVG)
