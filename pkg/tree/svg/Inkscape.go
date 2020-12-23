package svg

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type SodiPodiNamedView struct {
	tree.SVGTagData
}

func NewSodiPodiNamedView(attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	svgTag, err := tree.NewSVGTagData("sodipodi:namedview", attr, ctx)
	if err != nil {
		return nil, err
	}

	return &SodiPodiNamedView{svgTag}, nil
}

func (t *SodiPodiNamedView) Minify() bool {
	return true
}

type Defs struct {
	tree.SVGTagData
}

func NewDefs(attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	svgTag, err := tree.NewSVGTagData("defs", attr, ctx)
	if err != nil {
		return nil, err
	}

	return &Defs{svgTag}, nil
}

func (t *Defs) Minify() bool {
	return t.NumChildren() == 0
}

type Metadata struct {
	tree.SVGTagData
}

func NewMetadata(attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	svgTag, err := tree.NewSVGTagData("metadata", attr, ctx)
	if err != nil {
		return nil, err
	}

	return &Metadata{svgTag}, nil
}

func (t *Metadata) Minify() bool {
	return true
}

type InkscapeGrid struct {
	tree.SVGTagData
}

func NewInkscapeGrid(attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	svgTag, err := tree.NewSVGTagData("inkscape:grid", attr, ctx)
	if err != nil {
		return nil, err
	}

	return &InkscapeGrid{svgTag}, nil
}
