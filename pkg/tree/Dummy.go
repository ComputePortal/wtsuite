package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Dummy struct {
	tagData
}

func newDummy(attr *tokens.StringDict, ctx context.Context) (*Dummy, error) {
	td, err := newTag("dummy", false, attr, ctx)
	if err != nil {
		return nil, err
	}

	return &Dummy{td}, nil
}

func NewDummy(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	return newDummy(attr, ctx)
}

func NewSVGDummy(attr *tokens.StringDict, ctx context.Context) (SVGTag, error) {
	return newDummy(attr, ctx)
}

func (t *Dummy) Validate() error {
	panic("shouldnt be called")
}

func (t *Dummy) Write(indent string, nl, tab string) string {
	panic("shoudlnt be called")
}

func (t *Dummy) Compress(vb SVGViewBox) {
	panic("hopefully folded before compression")
}

func (t *Dummy) Minify() bool {
	panic("hopefully folded before compression")
}
