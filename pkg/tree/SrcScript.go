package tree

import (
	"strings"

	"../tokens/context"
)

type SrcScript struct {
	src string
	LeafTag
}

func NewSrcScript(src string, ctx context.Context) (*SrcScript, error) {
	return &SrcScript{src, NewLeafTag(ctx)}, nil
}

func (t *SrcScript) Write(indent string, nl, tab string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("<script type=\"text/javascript\" src=\"")
	b.WriteString(t.src)
	b.WriteString("\"></script>")
	b.WriteString(nl)

	return b.String()
}
