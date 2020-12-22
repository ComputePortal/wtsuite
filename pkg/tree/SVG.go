package tree

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type SVG struct {
	VisibleTagData
}

func NewSVG(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag("svg", false, attr, ctx)
	return &SVG{visTag}, err
}

func (t *SVG) Validate() error {
	if viewBoxToken_, ok := t.attributes.Get("viewBox"); ok {
		viewBoxToken, err := tokens.AssertString(viewBoxToken_)
		if err != nil {
			return err
		}

		if _, err := NewViewBoxFromString(viewBoxToken.Value(), viewBoxToken.Context()); err != nil {
			return err
		}
	} else {
		errCtx := t.attributes.Context()
		return errCtx.NewError("Error: no viewBox attribute")
	}

	for _, child_ := range t.Children() {
		if _, ok := child_.(SVGTag); !ok {
			errCtx := child_.Context()
			return errCtx.NewError("Error: not an svg tag '" + child_.Name() + "'")
		}
	}

	return t.ValidateChildren()
}

func (t *SVG) Compress(vb SVGViewBox) {
	// compress viewBox itself
	viewBoxStr, err := tokens.DictString(t.attributes, "viewBox")
	if err != nil {
		panic("should've been detected before")
	}

	vb_, err := NewViewBoxFromString(viewBoxStr.Value(), viewBoxStr.Context())
	if err != nil {
		panic("should've been detected before")
	}

	compressedString := vb_.CompressSelf()
	viewBoxToken := tokens.NewValueString(compressedString, viewBoxStr.Context())
	t.attributes.Set("viewBox", viewBoxToken)

	for _, child_ := range t.Children() {
		if child, ok := child_.(SVGTag); ok {
			child.Compress(vb)
		} else {
			panic("should've been caught before")
		}
	}
}

func (t *SVG) Minify() bool {
	attr := t.Attributes()
	attr.Delete("xmlns:dc")
	attr.Delete("xmlns:cc")
	attr.Delete("xmlns:rdf")
	attr.Delete("xmlns:svg")
	attr.Delete("xmlns:sodipodi")
	attr.Delete("xmlns:inkscape")
	t.id = "" // probably not used
	attr.Delete("inkscape:version")
	attr.Delete("sodipodi:docname")
	attr.Delete("width") // probaly needs be defined outside anyway
	attr.Delete("height")

	result := make([]Tag, 0)
	for _, child_ := range t.children {
		if child, ok := child_.(SVGTag); ok {
			if !child.Minify() {
				result = append(result, child)
			}
		} else {
			panic("should;ve been caught before")
		}
	}

	t.children = result

	// never remove svg tag
	return false
}

func (t *SVG) Write(indent string, nl, tab string) string {
	if COMPRESS_NUMBERS {
		viewBoxStr, err := tokens.DictString(t.attributes, "viewBox")
		if err != nil {
			panic("should've been detected before")
		}

		vb, err := NewViewBoxFromString(viewBoxStr.Value(), viewBoxStr.Context())
		if err != nil {
			panic("should've been detected before")
		}

		t.Compress(vb)
	}

	return t.VisibleTagData.WriteWrappedAutoHref(indent, nl, tab)
}
