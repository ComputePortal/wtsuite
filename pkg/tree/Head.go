package tree

import (
	"reflect"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Head struct {
	tagData
}

func NewHead(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	td, err := newTag("head", false, attr, ctx)
	if err != nil {
		return nil, err
	}
	return &Head{td}, nil
}

func (t *Head) Validate() error {
	for _, child := range t.children {
		switch ct := child.(type) {
		case *Meta, *Title, *Link, *Script, *Style, *LoaderScript, *SrcScript:
			// ok
		default:
			errCtx := ct.Context()
			return errCtx.NewError("HTML Error: expected meta, title, link, script or style (got: " + strings.ToLower(reflect.TypeOf(ct).String()[6:]) + ")")
		}
	}
	return t.ValidateChildren()
}
