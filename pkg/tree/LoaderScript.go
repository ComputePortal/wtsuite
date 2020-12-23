package tree

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type LoaderScript struct {
	content string
	LeafTag
}

func NewLoaderScript(content string, ctx context.Context) (*LoaderScript, error) {
	return &LoaderScript{content, NewLeafTag(ctx)}, nil
}

func (t *LoaderScript) Write(indent string, nl, tab string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("<script>")
	b.WriteString(nl)

	b.WriteString("function onload(){")
	b.WriteString(nl)

	b.WriteString(t.content)
	b.WriteString("}")
	b.WriteString(nl)

	b.WriteString("window.addEventListener(\"load\",onload,false);")
	b.WriteString(nl)

	b.WriteString(indent)
	b.WriteString("</script>")

	return b.String()
}
