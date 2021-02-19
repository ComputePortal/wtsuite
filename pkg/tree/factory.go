package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

var gb TagBuilder = &GenBuilder{false}
var gbInline TagBuilder = &GenBuilder{true}

var ncb TagBuilder = &NoChildrenBuilder{} // children can still be appended without error though

var table = map[string]TagBuilder{
	"!DOCTYPE": &fnBuilder{NewDocType},
	"a":        gbInline,
  "address":  gbInline,
  "article":  gb,
	"b":        gbInline,
	"i":        gbInline,
  "blockquote": gb,
	"br":       &fnBuilder{NewBr},
	"button":   gb,
	"body":     &fnBuilder{NewBody},
	"canvas":   gb,
  "code":     gbInline,
  "dd":       gb,
  "dl":       gb,
	"dummy":    &fnBuilder{NewDummy}, // empty tags are just collapsed
	"div":      &fnBuilder{NewDiv},
	"em":       gbInline,
	"footer":   gb,
	"form":     gb,
	"h1":       gbInline,
	"h2":       gbInline,
	"h3":       gbInline,
	"h4":       gbInline,
	"h5":       gbInline,
	"h6":       gbInline,
	"head":     &fnBuilder{NewHead},
	"header":   gb,
	"html":     &fnBuilder{NewHTML},
	"iframe":   gb,
	"img":      &fnBuilder{NewImg},
	"input":    &fnBuilder{NewInput},
	"label":    gb,
	"li":       gb,
	"link":     &fnBuilder{NewLink},
	"main":     gb,
	"meta":     &fnBuilder{NewMeta},
	"nav":      gb,
	"ol":       gb,
	"option":   gb,
	"p":        gbInline,
  "section":  gb,
	"select":   gb,
	"span":     gbInline,
	"svg":      &fnBuilder{NewSVG},
	"table":    gb,
	"tbody":    gb,
	"td":       gbInline,
	"textarea": gbInline,
	"tfoot":    gb,
	"th":       gbInline,
	"thead":    gb,
	"title":    &fnBuilder{NewTitle},
	"tr":       gbInline,
	"var":      gb,
	"ul":       gb,
	"?xml":     &fnBuilder{NewXMLHeader},
}

type TagBuilder interface {
	Build(key string, attr *tokens.StringDict, ctx context.Context) (Tag, error)
}

type fnBuilder struct {
	fn func(*tokens.StringDict, context.Context) (Tag, error)
}

func (b *fnBuilder) Build(key string, attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	return b.fn(attr, ctx)
}

type GenBuilder struct {
	inline bool
}

type NoChildrenBuilder struct {
}

// generic
func (b *GenBuilder) Build(key string, attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	return NewGeneric(key, attr, b.inline, ctx)
}

func (b *NoChildrenBuilder) Build(key string, attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	return NewGeneric(key, attr, false, ctx)
}

func IsTag(key string) bool {
	_, ok := table[key]

	return ok
}

func BuildTag(key string, attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	b, ok := table[key]

	if !ok {
		return nil, ctx.NewError("Error: tag " + key + " not found")
	}

	return b.Build(key, attr, ctx)
}
