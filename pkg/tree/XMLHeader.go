package tree

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type XMLHeader struct {
	tagData
}

func NewXMLHeader(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	// selfclosing=true, but in a special way, so actually irrelevant
	td, err := newTag("?xml", true, attr, ctx)
	if err != nil {
		return nil, err
	}

	return &XMLHeader{td}, nil
}

func (t *XMLHeader) Validate() error {
	if len(t.children) != 0 {
		errCtx := t.Context()
		return errCtx.NewError("Error: unexpected XML header children")
	}

	return nil
}

func (t *XMLHeader) Write(indent string, nl, tab string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("<?xml ")
	b.WriteString(t.writeAttributes())
	b.WriteString("?>")
	b.WriteString(nl)

	return b.String()
}
