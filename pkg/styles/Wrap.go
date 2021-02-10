package styles

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

// @wrap and @wrap-siblings are special at-rules that insert <div id=... class=.....> tags into the tree
// they do this before the other rules are evaluated
// TODO: as js query to combine with cross-lang template import in js, must be runtime because we don't know where template might be injected, all wrap statements must be executed?
// only @wrap can have runtime js equivalent, @wrap-siblings not
type WrapRule struct {
  sels []Selector
  elementName string // defaults to div
  class string
  id string
  siblings bool
  ctx context.Context
}

// no bubbling, so sel must be nil
func newWrapRule(sel Selector, key *tokens.String, attr *tokens.StringDict, siblings bool) ([]Rule, error) {
  if sel != nil {
    errCtx := key.Context()
    return nil, errCtx.NewError("Error: @wrap must be top-level (can't bubble)")
  }

  elementName := "div"
  class := ""
  id := ""

  if err := attr.Loop(func(key_ *tokens.String, val_ tokens.Token, last bool) error {
    switch key_.Value() {
    case "type":
      elementNameToken, err := tokens.AssertString(val_)
      if err != nil {
        return err
      }

      elementName = elementNameToken.Value()
      if !tree.IsTag(elementName) || strings.HasPrefix(elementName, "!") || elementName == "" {
        errCtx := elementNameToken.Context()
        return errCtx.NewError("Error: \"" + elementName + "\" is not a valid tag type")
      }
    case "class":
      classToken, err := tokens.AssertString(val_)
      if err != nil {
        return err
      }

      class = strings.TrimSpace(classToken.Value())
      if class == "" {
        errCtx := classToken.Context()
        return errCtx.NewError("Error: class can't be empty")
      }
    case "id":
      idToken, err := tokens.AssertString(val_)
      if err != nil {
        return err
      }

      id = strings.TrimSpace(idToken.Value())
      if id == "" {
        errCtx := idToken.Context()
        return errCtx.NewError("Error: id can't be empty")
      }
    default:
      errCtx := key_.Context()
      return errCtx.NewError("Error: \"" + key_.Value() + "\" is invalid key in @wrap rule")
    }

    return nil
  }); err != nil {
    return nil, err
  }

  // finally parse the selector, which should be the remainder of the key
  sels, err := ParseSelectorList(key)
  if err != nil {
    return nil, err
  }

  return []Rule{&WrapRule{sels, elementName, class, id, siblings, attr.Context()}}, nil
}

func NewWrapRule(sel Selector, key *tokens.String, attr *tokens.StringDict) ([]Rule, error) {

  return newWrapRule(sel, key.TrimLeft("@wrap"), attr, false)
}

func NewWrapSiblingsRule(sel Selector, key *tokens.String, attr *tokens.StringDict) ([]Rule, error) {
  return newWrapRule(sel, key.TrimLeft("@wrap-siblings"), attr, true)
}

func (r *WrapRule) ExpandNested() ([]Rule, error) {
  return []Rule{r}, nil
}

func (r *WrapRule) Write(indent string, nl string, tab string) (string, error) {
  return "", nil
}

// returns -1 if not found
func (r *WrapRule) findTag(lst []tree.Tag, t tree.Tag) int {
  for i, test := range lst {
    if test == t {
      return i
    }
  }

  return -1
}

func (r *WrapRule) extractTag(lst []tree.Tag, i int) (tree.Tag, []tree.Tag) {
  if i == 0 {
    return lst[0], lst[1:]
  } else if i == len(lst) - 1 {
    return lst[i], lst[0:i]
  } else {
    return lst[i], append(lst[0:i], lst[i+1:]...)
  }
}

// result also includes t
func (r *WrapRule) collectImmediateSiblings(t tree.Tag, matched []tree.Tag) ([]tree.Tag, []tree.Tag) {
  sibs := t.Siblings()

  i0 := r.findTag(sibs, t)
  if i0 == -1 {
    // t should be part of sibs
    panic("algo error")
  }

  iStart := i0
  for iStart > 0 {
    if i := r.findTag(matched, sibs[iStart-1]); i != -1 {
      _, matched = r.extractTag(matched, i)
      iStart -= 1
    } else {
      break
    }
  }

  iStop := i0 + 1
  for iStop < len(sibs) {
    if i := r.findTag(matched, sibs[iStop]); i != -1 {
      _, matched = r.extractTag(matched, i)
      iStop += 1
    } else {
      break
    }
  }

  return sibs[iStart:iStop], matched
}

func (t *WrapRule) wrap(tags []tree.Tag) error {
  if len(tags) == 0 {
    return nil
  }
  // 1. cfreate new tag

  attr := tokens.NewEmptyStringDict(t.ctx)
  newTagData, err := tree.NewVisibleTag(t.elementName, false, attr, t.ctx)
  if err != nil {
    panic(err)
  }

  newTag := &newTagData

  if t.class != "" {
    newTag.SetClasses([]string{t.class})
  }

  if t.id != "" {
    newTag.SetID(t.id)
  }

  parent := tags[0].Parent()
  sibs := parent.Children()

  // delete all current tags
  iPivot := -1
  for i, sib := range sibs {
    if sib == tags[0] {
      iPivot = i
      break
    }
  }

  if iPivot == -1 {
    panic("algo error")
  }

  for _, tag := range tags {
    if tag.Name() == "html" || tag.Name() == "body" {
      errCtx := t.ctx
      return errCtx.NewError("Error: can't wrap body or html")
    }

    parent.DeleteChild(iPivot)

    newTag.AppendChild(tag)
  }

  parent.InsertChild(iPivot, newTag)
  return nil
}

func (r *WrapRule) ApplyWrap(tag tree.Tag) error {
  matched := []tree.Tag{}
  for _, sel := range r.sels {
    matched = append(matched, sel.Match(tag)...)
  }

  // XXX: list based for now, there might be a performance benefit with a map for very large lists 
  for len(matched) > 0 {
    toBeWrapped := matched[0:1]
    matched = matched[1:]

    if r.siblings {
      toBeWrapped, matched = r.collectImmediateSiblings(toBeWrapped[0], matched)
    }

    if err := r.wrap(toBeWrapped); err != nil {
      return err
    }
  }

  return nil
}

var _wrapOk = registerAtRuleGen("wrap", NewWrapRule)
var _wrapSiblingsOk = registerAtRuleGen("wrap-siblings", NewWrapSiblingsRule)
