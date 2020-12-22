package svg

import (
	"../../tree"

	"../../tokens/context"
	tokens "../../tokens/html"
)

type Group struct {
	tree.SVGTagData
}

func NewGroup(attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	svgTag, err := tree.NewSVGTagData("g", attr, ctx)
	return &Group{svgTag}, err
}

func (t *Group) Validate() error {
	// TODO: parse transforms
	return t.SVGTagData.Validate()
}

func (t *Group) Compress(vb tree.SVGViewBox) {
	// TODO: transform SVGViewBox
	t.SVGTagData.Compress(vb)
}

func (t *Group) Minify() bool {
	attr := t.Attributes()

	attr.Delete("aria-label")
	attr.Delete("inkscape:label")
	attr.Delete("inkscape:groupmode")

	t.SetID("") // probably not used

	// remove all font attributes if no text detected as direct child
	/*hasTextChild := false
	for _, child := range t.Children() {
		if _, ok := child.(*tree.Text); ok {
			hasTextChild = true
		}
	}

	MinifyStyles(attr, hasTextChild)*/

	// XXX XXX XXX
	// for inkscape minification we can remove the style tag altogether, but this is not generally applicable!
	attr.Delete("style")

	t.SVGTagData.Minify()

	return t.NumChildren() == 0
}
