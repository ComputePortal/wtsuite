package directives

import (
	"fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/math"
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/svg"
)

type MathNode struct {
	NodeData
}

func NewMathNode(parentTag tree.Tag, parentNode Node) *MathNode {
	ctx := parentTag.Context()
	gattr := tokens.NewEmptyStringDict(ctx)
	gtag, err := svg.BuildTag("g", gattr, ctx)
	if err != nil {
		panic(err)
	}

	parentTag.AppendChild(gtag)

	return &MathNode{newNodeData(gtag, parentNode)}
}

func (n *MathNode) Transform(x, y, sx, sy float64) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("translate(%g,%g)", x, y))
	b.WriteString(" ")
	b.WriteString(fmt.Sprintf("scale(%g,%g)", sx, sy))

	ctx := n.tag.Context()
	transformToken := tokens.NewValueString(b.String(), ctx)
	attr := n.tag.Attributes()

	attr.Set("transform", transformToken)
}

func (n *MathNode) buildScaledMathText(x float64, y float64, fontSize float64, centerX float64, scaleX float64, value string,
	ctx context.Context) error {
	attr := tokens.NewEmptyStringDict(ctx)

	xToken := tokens.NewValueFloat(x, ctx)
	attr.Set("x", xToken)

	yToken := tokens.NewValueFloat(y, ctx)
	attr.Set("y", yToken)

	fontSizeToken := tokens.NewValueFloat(fontSize, ctx)
	attr.Set("font-size", fontSizeToken)

	tag, err := svg.BuildTag("text", attr, ctx)
	if err != nil {
		return err
	}
	tag.AppendChild(tree.NewText(value, ctx))

	if scaleX == 1.0 {
		n.tag.AppendChild(tag)
	} else {
		gattr := tokens.NewEmptyStringDict(ctx)

		var b strings.Builder
		b.WriteString("translate(")
		b.WriteString(fmt.Sprintf("%g", centerX))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%g", y))
		b.WriteString(")")
		b.WriteString(" ")
		b.WriteString("scale(")
		b.WriteString(fmt.Sprintf("%g", scaleX))
		b.WriteString(",1.0)")
		b.WriteString(" ")
		b.WriteString("translate(")
		b.WriteString(fmt.Sprintf("%g", -centerX))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%g", -y))
		b.WriteString(")")

		transformToken := tokens.NewValueString(b.String(), ctx)

		gattr.Set("transform", transformToken)
		gtag, err := svg.BuildTag("g", gattr, ctx)
		if err != nil {
			return err
		}

		gtag.AppendChild(tag)
		n.tag.AppendChild(gtag)
	}

	return nil
}

func (n *MathNode) BuildMathText(x float64, y float64, fontSize float64, value string,
	ctx context.Context) error {
	return n.buildScaledMathText(x, y, fontSize, 0.0, 1.0, value, ctx)
}

func (n *MathNode) BuildMathPath(d string, ctx context.Context) error {
	attr := tokens.NewEmptyStringDict(ctx)

	dToken := tokens.NewValueString(d, ctx)
	attr.Set("d", dToken)

	return buildSVGPathInternal(n, attr, ctx)
}

func (n *MathNode) NewSubScope() math.SubScope {
	return NewMathNode(n.tag, n)
}

func (n *MathNode) getNode() Node {
  return n
}
