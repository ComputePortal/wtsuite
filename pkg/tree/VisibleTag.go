package tree

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type VisibleTag interface {
	Name() string
	GetID() string
	SetID(string)
	GetClasses() []string
	SetClasses([]string)


	Context() context.Context

	CollectIDs(IDMap) error
}

type VisibleTagData struct {
	classes []string
	tagData
}

func NewVisibleTag(name string, selfClosing bool, attr *tokens.StringDict, ctx context.Context) (VisibleTagData, error) {
  classes := []string{}
  if classToken_, ok := attr.Get("class"); ok {
    if !tokens.IsNull(classToken_) {
      classToken, err := tokens.AssertString(classToken_)
      if err != nil {
        return VisibleTagData{}, err
      }

      classes = strings.Fields(classToken.Value())
    }

    attr.Delete("class")
  }

	td, err := newTag(name, selfClosing, attr, ctx)
	if err != nil {
		return VisibleTagData{}, err
	}

	return VisibleTagData{classes, td}, nil
}

func (t *VisibleTagData) GetClasses() []string {
	return t.classes[:]
}

func (t *VisibleTagData) SetClasses(c []string) {
	t.classes = c
}

func (t *VisibleTagData) CollectIDs(idMap IDMap) error {
	if t.id != "" {
		if idMap.Has(t.id) {
			other := idMap.Get(t.id)
			errCtx := t.Context()
			err := errCtx.NewError("Error: id " + t.id + " already defined")
			context.PrependContextString(err, "Info: defined here", other.Context())
			return err
		} else {
			idMap.Set(t.id, t)
		}
	}

	for _, child := range t.children {
		switch c := child.(type) {
		case VisibleTag:
			if err := c.CollectIDs(idMap); err != nil {
				return err
			}
		default:
			errCtx := child.Context()
			return errCtx.NewError("Error: not a visible tag (collecting ids in VisibleTagData)")
		}

	}

	return nil
}

func (t *VisibleTagData) rewriteAttributes() string {
	var b strings.Builder

	if t.id != "" {
		b.WriteString(" id=")
		b.WriteString("\"")
		b.WriteString(t.id)
		b.WriteString("\"")
	}

	if len(t.classes) > 0 {
		b.WriteString(" class=\"")
		b.WriteString(strings.Join(t.classes, " "))
		b.WriteString("\"")
	}

	b.WriteString(t.tagData.writeAttributes())

	return b.String()
}

func (t *VisibleTagData) SetTmpAttribute(key string, valueToken tokens.Token) {
	keyToken := tokens.NewValueString(key, t.Context())
	t.attributes.Set(keyToken, valueToken)
}

func (t *VisibleTagData) RemoveTmpAttribute(key string) {
	t.attributes.Delete(key)
}

func (t *VisibleTagData) write(wrapAutoHref bool, indent string, nl, tab string) string {
	hasID := (t.id != "")
	hasClasses := len(t.classes) > 0

	if hasID  {
		valueToken := tokens.NewValueString(t.id, t.Context())
		t.SetTmpAttribute("id", valueToken)
	}

	if hasClasses {
		valueToken := tokens.NewValueString(strings.Join(t.classes, " "), t.Context())
		t.SetTmpAttribute("class", valueToken)
	}

	result := t.tagData.write(wrapAutoHref, indent, nl, tab)

	if hasID {
		t.RemoveTmpAttribute("id")
	}

	if hasClasses {
		t.RemoveTmpAttribute("class")
	}

	return result
}

func (t *VisibleTagData) Write(indent string, nl, tab string) string {
	return t.write(true, indent, nl, tab)
}

func (t *VisibleTagData) ToJSType() string {
	// original type is eg. Input, but due to assertions it ends up as VisibleTagData
	if t.name == "input" || t.name == "textarea" {
		return "HTMLInputElement"
	} else if t.name == "img" {
		return "HTMLImageElement"
	} else if t.name == "select" {
		return "HTMLSelectElement"
	} else if t.name == "canvas" {
		return "HTMLCanvasElement"
	} else if t.name == "iframe" {
		return "HTMLIFrameElement"
	} else {
		return "HTMLElement"
	}
}

func (t *VisibleTagData) InnerHTML() string {
	hasText := false
	for _, child := range t.children {
		if _, ok := child.(*Text); ok {
			hasText = true
			break
		} else if child.Name() == "b" || child.Name() == "i" {
			hasText = true
			break
		}
	}

	if hasText {
		return t.writeChildren("", "", "")
	} else {
		return ""
	}
}

func (t *VisibleTagData) WriteWrappedAutoHref(indent string, nl, tab string) string {
	var b strings.Builder

	hrefToken, hasHref := t.attributes.Get("href")

	autoLink := AUTO_LINK && hasHref && !tokens.IsNull(hrefToken)

	// also apply id to a tag, instead of img tag
	id := t.id
	if autoLink {
		href, err := tokens.DictString(t.attributes, "href")
		if err != nil {
			panic("should've been caught before")
		}

		t.attributes.Delete("href")

		b.WriteString("<a href=\"")
		b.WriteString(href.Value())
		b.WriteString("\"")

		if id != "" {
			b.WriteString(" id=\"")
			b.WriteString(id)
			b.WriteString("\"")
		}
		b.WriteString("><div style=\"display: flex;\">")

		t.id = ""
	}

	b.WriteString(t.write(false, "", nl, tab))

	if autoLink {
		t.id = id
		b.WriteString("</div></a>")
	}

	return b.String()
}
