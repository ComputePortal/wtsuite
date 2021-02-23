package svg

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type TagBuilder interface {
	Build(key string, attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error)
}

type fnBuilder struct {
	fn func(*tokens.StringDict, context.Context) (tree.SVGTag, error)
}

type GenBuilder struct {
}

type NoChildrenBuilder struct {
}

type TextChildBuilder struct {
}

var gb TagBuilder = &GenBuilder{}

var ncb TagBuilder = &NoChildrenBuilder{} // children can still be appended without error though

var table = map[string]TagBuilder{
	"cc:Work":        gb,
	"circle":         ncb,
	"dc:format":      gb,
	"dc:title":       gb,
	"dc:type":        gb,
	"defs":           &fnBuilder{NewDefs},
	"dummy":          &fnBuilder{tree.NewSVGDummy},
	"ellipse":        ncb,
	"feGaussianBlur": ncb,
	"filter":         gb,
	"g":              &fnBuilder{NewGroup},
	"image":          ncb,
	"inkscape:grid":  &fnBuilder{NewInkscapeGrid},
	"metadata":       &fnBuilder{NewMetadata},
	//"path":         ncb, // done via directive, and then generic
	"pattern":            gb,
	"rdf:RDF":            gb,
	"rect":               &fnBuilder{NewRect},
	"sodipodi:namedview": &fnBuilder{NewSodiPodiNamedView},
	"svg":                gb,
	"text":               gb, // not the same as Text!
  "title":              gb,
}

func IsTag(key string) bool {
	_, ok := table[key]
	return ok
}

func BuildTag(key string, attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	b, ok := table[key]

	if !ok {
		panic("not found '" + key + "'") // default to something generic?
	}

	return b.Build(key, attr, ctx)
}

func (b *fnBuilder) Build(key string, attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	return b.fn(attr, ctx)
}

func (b *GenBuilder) Build(key string, attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	return NewGeneric(key, attr, ctx)
}

func (b *NoChildrenBuilder) Build(key string, attr *tokens.StringDict, ctx context.Context) (tree.SVGTag, error) {
	// TODO: this is a bad builder because children are appended later anyway

	return NewGeneric(key, attr, ctx)
}
