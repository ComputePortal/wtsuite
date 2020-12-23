package html

import (
	"fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// for text tag 'attr' and 'children' are nil, and 'text' is set
// for non-text tag 'text' isnt set
type Tag struct {
	name     string
	attr     *RawDict
	children []*Tag
	text     string
	textCtx  context.Context
  isDirective bool
	TokenData
}

func NewTag(name string, attr *RawDict, children []*Tag, ctx context.Context) *Tag {
	return &Tag{name, attr, children, "", ctx, false, TokenData{ctx}}
}

func NewTextTag(text string, ctx context.Context) *Tag {
	return &Tag{"", nil, nil, text, ctx, false, TokenData{ctx}}
}

func NewScriptTag(name string, attr *RawDict, text string, textCtx, ctx context.Context) *Tag {
	return &Tag{name, attr, nil, text, textCtx, false, TokenData{ctx}}
}

func NewDirectiveTag(name string, attr *RawDict, children []*Tag, ctx context.Context) *Tag {
	return &Tag{name, attr, children, "", ctx, true, TokenData{ctx}}
}

func (t *Tag) Dump(indent string) string {
	var b strings.Builder

	if t.IsText() {
		b.WriteString(indent)
		b.WriteString("TEXT:")
		b.WriteString(t.text)
		b.WriteString("\n")
	} else if t.IsScript() {
		b.WriteString(indent)
		b.WriteString("<")
		b.WriteString(t.name)
		b.WriteString("\n")
		b.WriteString(indent)
		b.WriteString(t.text)
		b.WriteString("\n")
	} else {
		b.WriteString(indent)
		b.WriteString("<")
		b.WriteString(t.name)
		b.WriteString("\n")
		b.WriteString(indent)
		b.WriteString("  ##START OF TAG ATTR\n")
		b.WriteString(t.attr.Dump(indent + "  "))
		b.WriteString(indent)
		b.WriteString("  ##END OF TAG ATTR\n")
		for _, tag := range t.children {
			b.WriteString(tag.Dump(indent + "  "))
		}
	}

	return b.String()
}

func (t *Tag) Eval(scope Scope) (Token, error) {
	panic("Tags should never be evaluated")
}

func (t *Tag) IsSame(other Token) bool {
	panic("Tags should never be compared")
}

func (t *Tag) IsText() bool {
	return t.attr == nil && t.children == nil
}

func (t *Tag) IsScript() bool {
	// XXX: this is a bad place to check all of this
	return t.attr != nil && (t.children == nil || len(t.children) == 0) && t.Name() == "script"
}

func (t *Tag) Name() string {
	return t.name
}

func (t *Tag) Attributes(posNames []string) (*StringDict, error) {
	result := NewEmptyStringDict(t.attr.Context())

  hasKWArgs := false

	// duplicates dont give an error anymore: so beware
	for _, keyVal := range t.attr.items {
		key := keyVal.key
		val := keyVal.value

		if IsInt(key) {
      if hasKWArgs {
        errCtx := val.Context()
        return nil, errCtx.NewError("Error: positional attributes must come first")
      }

			iKey, err := AssertInt(key)
			if err != nil {
				panic(err)
			}

			i := iKey.Value()
			if i > len(posNames)-1 {
				errCtx := iKey.Context()
				return nil, errCtx.NewError(fmt.Sprintf("Error: positional arg %d (%s) out of range (expected %d positional args)", i+1, val.Dump(""), len(posNames)))
			}

			posName := posNames[i]
			strKey := NewValueString(posName, key.Context())
			result.Set(strKey, val)
		} else {
      hasKWArgs = true
			result.Set(key, val)
		}
	}

	return result, nil
}

func (t *Tag) RawAttributes() *RawDict {
	return t.attr
}

func (t *Tag) Children() []*Tag {
	return t.children
}

func (t *Tag) AppendChild(c *Tag) error {
	if t.IsText() || t.IsScript() {
		errCtx := c.Context()
		return errCtx.NewError("Error: cannot be nested in text/script")
	}

	t.children = append(t.children, c)

	return nil
}

func (t *Tag) Text() string {
	return t.text
}

func (t *Tag) TextContext() context.Context {
	return t.textCtx
}

func (t *Tag) IsEmpty() bool {
	return (t.children == nil || len(t.children) == 0) && t.text == ""
}

func (t *Tag) AssertEmpty() error {
	if !t.IsEmpty() {
		errCtx := t.Context()
		return errCtx.NewError("Error: unexpected content")
	}

	return nil
}

func (t *Tag) AssertNoAttributes() error {
	if t.attr != nil && t.attr.Len() != 0 {
		errCtx := t.attr.Context()
		return errCtx.NewError("Error: unexpected attributes")
	} else {
		return nil
	}
}

func (t *Tag) IsDirective() bool {
  return t.isDirective
}
