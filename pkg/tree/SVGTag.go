package tree

import (
	"reflect"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type SVGTag interface {
	Tag

	// TODO: take g transforms into account when going down the tree
	Compress(vb SVGViewBox)

	// returns true if can be removed
	Minify() bool // for svg documents created by drawing programs
}

type SVGTagData struct {
	VisibleTagData
}

func NewSVGTagData(name string, attr *tokens.StringDict, ctx context.Context) (SVGTagData, error) {
	visTag, err := NewVisibleTag(name, false, attr, ctx)
	if err != nil {
		return SVGTagData{}, err
	}

	return SVGTagData{visTag}, err
}

func (t *SVGTagData) Validate() error {
	for _, child := range t.Children() {
		if _, ok := child.(SVGTag); !ok {
			errCtx := child.Context()
			return errCtx.NewError("Error: not an svg tag " + child.Name() + reflect.TypeOf(child).String())
		}
	}

	return t.ValidateChildren()
}

func (t *SVGTagData) Compress(vb SVGViewBox) {
	for _, child_ := range t.Children() {
		if child, ok := child_.(SVGTag); ok {
			child.Compress(vb)
		} else {
			panic("should've been checked earlier")
		}
	}
}

func (t *SVGTagData) Minify() bool {
	result := make([]Tag, 0)
	for _, child_ := range t.Children() {
		if child, ok := child_.(SVGTag); ok {
			if !child.Minify() {
				result = append(result, child)
			}
		} else {
			panic("should've been checked earlier")
		}
	}

	t.DeleteAllChildren()
	for _, child := range result {
		t.AppendChild(child)
	}

	return false
}
