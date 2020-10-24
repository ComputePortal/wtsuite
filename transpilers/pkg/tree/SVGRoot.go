package tree

import (
	"fmt"

	"../tokens/context"
)

type SVGRoot struct {
	tagData
}

func NewSVGRoot(ctx context.Context) *SVGRoot {
	return &SVGRoot{tagData{"", "", false, nil, nil, make([]Tag, 0), ctx}}
}

func (t *SVGRoot) getXMLHeaderAndSVG() (*XMLHeader, *SVG, error) {
	var xmlHeader *XMLHeader = nil
	var svg *SVG = nil

	for _, child := range t.children {
		switch tt := child.(type) {
		case *XMLHeader:
			if xmlHeader != nil {
				errCtx := context.MergeContexts(child.Context(), xmlHeader.Context())
				return nil, nil, errCtx.NewError("SVG Error: ?xml header defined twice")
			} else if svg != nil {
				errCtx := context.MergeContexts(child.Context(), svg.Context())
				return nil, nil, errCtx.NewError("SVG Error: svg defined before ?xml header")
			}

			xmlHeader = tt
		case *SVG:
			if svg != nil {
				errCtx := context.MergeContexts(child.Context(), svg.Context())
				return nil, nil, errCtx.NewError("SVG Error: svg defined twice")
			}

			svg = tt
		default:
			errCtx := child.Context()
			return nil, nil, errCtx.NewError("SVG Error: expected only ?xml header or svg")
		}
	}

	if xmlHeader == nil {
		err := t.ctx.NewError(fmt.Sprintf("SVG Error: no ?xml header defined (nChildren: %d)",
			len(t.children)))
		return nil, nil, err
	}

	if svg == nil {
		return nil, nil, t.ctx.NewError("SVG Error: no svg defined")
	}

	return xmlHeader, svg, nil
}

func (t *SVGRoot) Validate() error {
	xmlHeader, svg, err := t.getXMLHeaderAndSVG()
	if err != nil {
		return err
	}

	if err := xmlHeader.Validate(); err != nil {
		return err
	}

	if err := svg.Validate(); err != nil {
		return err
	}

	return err
}

func (t *SVGRoot) Minify() {
	_, svg, err := t.getXMLHeaderAndSVG()
	if err != nil {
		panic("should've been caught before")
	}

	if svg.Minify() {
		panic("should be false")
	}
}

func (t *SVGRoot) Write(indent string, nl string, tab string) string {
	return t.writeChildren(indent, nl, tab)
}
