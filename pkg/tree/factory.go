package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

var gb TagBuilder = &GenBuilder{false}
var gbInline TagBuilder = &GenBuilder{true}

var ncb TagBuilder = &NoChildrenBuilder{} // children can still be appended without error though

var table = map[string]TagBuilder{
	"!doctype": &fnBuilder{NewDocType},
	"!DOCTYPE": &fnBuilder{NewDocType},
	"a":        gbInline,
  "abbr":     gbInline,
  "address":  gbInline,
  "area":     gb,
  "article":  gb,
  "aside":    gb,
  "audio":    gb,
	"b":        gbInline,
	"base":     &fnBuilder{NewBase},
  "bdi":      gbInline,
  "bdo":      gbInline,
  "blockquote": gb,
	"body":     &fnBuilder{NewBody},
	"br":       &fnBuilder{NewBr},
	"button":   gb,
	"canvas":   gb,
  "caption":  gb,
  "cite":     gbInline,
  "code":     gbInline,
  "col":      &fnBuilder{NewCol},
  "colgroup": gb,
  "data":     gb,
  "datalist": gb,
  "dd":       gbInline,
  "del":      gbInline,
  "details":  gb,
  "dfn":      gbInline,
  "dialog":   gb,
	"div":      &fnBuilder{NewDiv},
  "dl":       gb,
  "dt":       gbInline,
	"dummy":    &fnBuilder{NewDummy}, // empty tags are just collapsed
	"em":       gbInline,
  "embed":    gb,
  "fieldset": gb,
  "figcaption": gb,
  "figure":   gb,
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
  "hr":       &fnBuilder{NewHr},
	"html":     &fnBuilder{NewHTML},
	"i":        gbInline,
	"iframe":   gb,
	"img":      &fnBuilder{NewImg},
	"input":    &fnBuilder{NewInput},
  "ins":      gbInline,
  "kbd":      gbInline,
	"label":    gb,
  "legend":   gb,
	"li":       gb,
	"link":     &fnBuilder{NewLink},
	"main":     gb,
  "map":      gb,
  "mark":     gbInline,
	"meta":     &fnBuilder{NewMeta},
  "meter":    gb,
	"nav":      gb,
  "noscript": gb,
  "object":   gb,
	"ol":       gb,
  "optgroup": gb,
	"option":   gb,
  "output":   gb,
	"p":        gbInline,
  "param":    gb,
  "picture":  gb,
  "pre":      gb,
  "progress": gb,
  "q":        gbInline,
  "rp":       gb,
  "rt":       gb,
  "ruby":     gb,
  "s":        gbInline,
  "samp":     gbInline,
  //"script":   gb, // done via directive instead
  "section":  gb,
	"select":   gb,
  "small":    gbInline,
  "source":   gb,
	"span":     gbInline,
  "strong":   gbInline,
  "sub":      gbInline,
  "summary":  gb,
  "sup":      gbInline,
	"svg":      &fnBuilder{NewSVG},
	"table":    gb,
	"tbody":    gb,
	"td":       gbInline,
  "template": gb, // differs from directive!
	"textarea": gbInline,
	"tfoot":    gb,
	"th":       gbInline,
	"thead":    gb,
  "time":     gbInline,
	"title":    &fnBuilder{NewTitle},
	"tr":       gbInline,
  "track":    gb,
  "u":        gbInline,
	"ul":       gb,
	"var":      gb,
  "video":    gb,
  "wbr":      &fnBuilder{NewWbr},
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
