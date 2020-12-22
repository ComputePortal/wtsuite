package svg

import (
	"../../tree"

	"../../tokens/context"
	tokens "../../tokens/html"
)

type Rect struct {
	tree.SVGTagData
}

func NewRect(attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	svgTag, err := tree.NewSVGTagData("rect", attr, ctx)
	return &Rect{svgTag}, err
}

func (t *Rect) Validate() error {
	attr := t.Attributes()

	if err := AssertFloatAttribute(attr, "x"); err != nil {
		return err
	}

	if err := AssertFloatAttribute(attr, "y"); err != nil {
		return err
	}

	if err := AssertFloatAttribute(attr, "width"); err != nil {
		return err
	}

	if err := AssertFloatAttribute(attr, "height"); err != nil {
		return err
	}

	if t.NumChildren() != 0 {
		errCtx := t.Context()
		return errCtx.NewError("Error: unexpected children")
	}

	return t.SVGTagData.Validate()
}

func (t *Rect) Compress(vb tree.SVGViewBox) {
	attr := t.Attributes()

	CompressFloatAttribute(attr, "x", vb.CompressX)
	CompressFloatAttribute(attr, "y", vb.CompressY)
	CompressFloatAttribute(attr, "width", vb.CompressX)
	CompressFloatAttribute(attr, "height", vb.CompressY)

	CompressStyles(attr, vb)

	t.SVGTagData.Compress(vb)
}

func (t *Rect) Minify() bool {
	attr := t.Attributes()
	t.SetID("") // probably not used

	MinifyStyles(attr, false)

	return false
}
