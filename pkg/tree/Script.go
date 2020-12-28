package tree

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

// reuse tagData's write functions
type Script struct {
	attributes *tokens.StringDict
	content    string // can be empty if src attribute isnt set
	LeafTag
}

func NewScript(attr *tokens.StringDict, content string, ctx context.Context) (Tag, error) {
	return &Script{attr, content, NewLeafTag(ctx)}, nil
}

func (t *Script) CollectScripts(idMap IDMap, classMap ClassMap, bundle *scripts.InlineBundle) error {
	srcToken_, hasSrc := t.attributes.Get("src")

	if t.content != "" && hasSrc {
		errCtx := t.attributes.Context()
		return errCtx.NewError("Error: can't have both src and inline content")
	}

	if t.content == "" && !hasSrc {
		errCtx := t.attributes.Context()
		return errCtx.NewError("Error: can't have neither src and inline content")
	}

	ctx := t.Context()
	if t.content != "" {
		script, err := scripts.NewInlineScript(t.content)
		if err != nil {
			return err
		}

		bundle.Append(script)
	} else {
		srcToken, err := tokens.AssertString(srcToken_)
		if err != nil {
			return err
		}

		script, err := scripts.NewSrcScript(srcToken.Value())
		if err != nil {
			if err.Error() == "not found" {
				errCtx := ctx
				return errCtx.NewError("Error: '" + srcToken.Value() + "' not found")
			} else {
				return err
			}
		}

		bundle.Append(script)
	}

	return nil
}

func (t *Script) Write(indent string, nl, tab string) string {
	var b strings.Builder

  b.WriteString(indent)
  b.WriteString("<script>")
  b.WriteString(nl)
  b.WriteString(t.content)
  b.WriteString(nl)
  b.WriteString(indent)
  b.WriteString("</script>")
  b.WriteString(nl)

  return b.String()
}
