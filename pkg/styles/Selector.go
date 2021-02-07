package styles

import (
  "errors"
  "strings"

	"github.com/computeportal/wtsuite/pkg/parsers"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
	"github.com/computeportal/wtsuite/pkg/tree"
)

type Selector interface {
  Extend(extra *tokens.String) ([]Selector, error) // eg. pseudo selector, or child
  Match(tag tree.Tag) []tree.Tag // returns empty list if no match
  Write() string
}

type SelectorData struct {
  elementName string // *, div, body, html (can't be empty)
  class string // can't  be combined with id, *.class is simplified as .class
  id string // can't be combined with class, *#id is simplified as .id
  filters []AttrFilter // also part of search queries

  pseudoClasses []PseudoClass
  pseudoElement string // there can only be one, and can't be extended if there is already a pseudoelement

  descendant *SelectorData
  sibling *SelectorData
  immediate bool // sibling or descendant
}

func (s *SelectorData) Copy() *SelectorData {
  return &SelectorData{
    s.elementName,
    s.class,
    s.id,
    s.filters,
    s.pseudoClasses,
    s.pseudoElement,
    s.descendant,
    s.sibling,
    s.immediate,
  }
}

// for extend
func (s *SelectorData) setDescendant(desc *SelectorData, imm bool) *SelectorData {
  if desc != nil {
    cpy := s.Copy()
    if cpy.descendant != nil {
      cpy.descendant = cpy.descendant.setDescendant(desc, imm)
    } else {
      cpy.descendant = desc
      cpy.immediate = imm
    }
    return cpy
  } else {
    return s
  }
}

func (s *SelectorData) setSibling(sib *SelectorData, imm bool) *SelectorData {
  if sib != nil {
    cpy := s.Copy()
    if s.sibling != nil {
      cpy.sibling = cpy.sibling.setSibling(sib, imm)
    } else {
      cpy.sibling = sib
      cpy.immediate = imm
    }
    return cpy
  } else {
    return s
  }
}

func (s *SelectorData) addFilters(filters []AttrFilter) *SelectorData {
  if filters != nil && len(filters) > 0 {
    cpy := s.Copy()
    if cpy.descendant != nil {
      cpy.descendant = cpy.descendant.addFilters(filters)
    } else if s.sibling != nil {
      cpy.sibling = cpy.sibling.addFilters(filters)
    } else {
      cpy.filters = append(cpy.filters, filters...)
    }
    return cpy
  } else {
    return s
  }
}

func (s *SelectorData) addPseudoClasses(pseudoClasses []PseudoClass) *SelectorData {
  if pseudoClasses != nil && len(pseudoClasses) > 0 {
    cpy := s.Copy()
    if cpy.descendant != nil {
      cpy.descendant = cpy.descendant.addPseudoClasses(pseudoClasses)
    } else if s.sibling != nil {
      cpy.sibling = cpy.sibling.addPseudoClasses(pseudoClasses)
    } else {
      cpy.pseudoClasses = append(cpy.pseudoClasses, pseudoClasses...)
    }
    return cpy
  } else {
    return s
  }
}

func (s *SelectorData) setPseudoElement(pseudoElement string) (*SelectorData, error) {
  if pseudoElement != "" {
    cpy := s.Copy()
    var err error
    if cpy.descendant != nil {
      cpy.descendant, err = cpy.descendant.setPseudoElement(pseudoElement)
    } else if s.sibling != nil {
      cpy.sibling, err = cpy.sibling.setPseudoElement(pseudoElement)
    } else {
      if cpy.pseudoElement == "" {
        cpy.pseudoElement = pseudoElement
      } else {
        err = errors.New("Error: pseudoElement already set")
      }
    }
    return cpy, err
  } else {
    return s, nil
  }
}

// split by comma
func parseList(query *tokens.String) ([][]raw.Token, error) {
  ctx := query.Context()
  if ctx.Content() != query.Value() {
    src := context.NewSource(query.Value())
    ctx = context.NewContext(src, ctx.Path())
  }

  p, err := parsers.NewCSSSelectorParser(query.Value(), ctx)
  if err != nil {
    return nil, err
  }

  ts, err := p.Tokenize()
  if err != nil {
    return nil, err
  }

  tss := raw.SplitBySymbol(ts, patterns.COMMA)

  return tss, nil
}

func ParseSelectorList(query *tokens.String) ([]Selector, error) {
  tss, err := parseList(query)
  if err != nil {
    return nil, err
  }

  sels := make([]Selector, 0)
  for _, ts_ := range tss {
    if len(ts_) < 1 {
      errCtx := query.Context()
      return nil, errCtx.NewError("Error: bad css selector list")
    }

    sel, err := ParseSelector(ts_)
    if err != nil {
      return nil, err
    }

    sels = append(sels, sel)
  }

  return sels, nil
}

func assertNonClassOrIDWord(t raw.Token) (*raw.Word, error) {
  w, err := raw.AssertWord(t)
  if err != nil {
    return nil, err
  }

  if strings.Contains(w.Value(), ".") || strings.Contains(w.Value(), "#") {
    errCtx := w.Context()
    return nil, errCtx.NewError("Error: bad css selector word")
  }

  return w, nil
}

func parseNameToken(t raw.Token) (string, string, string, error) {
  elementName := ""
  class := ""
  id := ""

  ctx := t.Context()

  if raw.IsSymbol(t, "*") {
    elementName = "*"
  } else {
    w_, err := raw.AssertWord(t)
    if err != nil {
      return "", "", "", err
    }
    
    w := w_.Value()

    if strings.HasPrefix(w, ".") {
      if strings.Contains(w, "#") {
        return "", "", "", ctx.NewError("Error: can't mix id and class selector")
      }

      class = w[1:]
    } else if strings.Contains(w, ".") {
      parts := strings.Split(w, ".")
      if len(parts) != 2 {
        return "", "", "", ctx.NewError("Error: bad class selector")
      }

      elementName = parts[0]
      class = parts[1]

      if strings.Contains(elementName, "#") {
        return "", "", "", ctx.NewError("Error: can't mix id and class selector")
      } else if strings.Contains(class, "#") {
        return "", "", "", ctx.NewError("Error: can't mix id and class selector")
      }
    } else if strings.HasPrefix(w, "#") {
      if strings.Contains(w, ".") {
        return "", "", "", ctx.NewError("Error: can't mix id and class selector")
      }

      id = w[1:]
    } else if strings.Contains(w, "#") {
      parts := strings.Split(w, "#") 
      if len(parts) != 2 {
        return "", "", "", ctx.NewError("Error: bad class selector")
      }

      elementName = parts[0]
      id = parts[1]

      if strings.Contains(elementName, ".") {
        return "", "", "", ctx.NewError("Error: can't mix id and class selector")
      } else if strings.Contains(id, ".") {
        return "", "", "", ctx.NewError("Error: can't mix id and class selector")
      }
    } else {
      elementName = w
    }
  }

  return elementName, class, id, nil
}

// use a dummy *SelectorData to return all the vlues
func parseFiltersPseudoAndDescendants(ts []raw.Token) (*SelectorData, error) {
  filters := make([]AttrFilter, 0)
  pseudoClasses := make([]PseudoClass, 0)
  pseudoElement := ""

  var err error
  var descendant *SelectorData = nil
  var sibling *SelectorData = nil
  immediate := false

  for i := 0; i < len(ts); i++ {
    t := ts[i]

    done := false

    switch {
    case raw.IsBracketsGroup(t):
      if len(pseudoClasses) != 0 {
        errCtx := t.Context()
        return nil, errCtx.NewError("Error: attribute selectors must come before pseudoclasses")
      }

      filter, err := ParseAttrFilter(t)
      if err != nil {
        return nil, err
      }

      filters = append(filters, filter)
    case raw.IsSymbol(t, ":"):
      if pseudoElement != "" {
        errCtx := t.Context()
        return nil, errCtx.NewError("Error: pseudoclasses must come before pseudo element")
      }

      var pseudoClass PseudoClass
      var err error

      if i < len(ts) - 2 && raw.IsAnyWord(ts[i+1]) && raw.IsParensGroup(ts[i+2]) {
        pseudoClass, err = ParsePseudoClass(ts[i+1:i+3])
        i += 2
      } else if i < len(ts) - 1 && raw.IsAnyWord(ts[i+1]) {
        pseudoClass, err = ParsePseudoClass(ts[i+1:i+2])
        i += 1
      } else {
        errCtx := t.Context()
        return nil, errCtx.NewError("Error: bad pseudo class")
      }

      if err != nil {
        return nil, err
      }

      pseudoClasses = append(pseudoClasses, pseudoClass)
    case raw.IsSymbol(t, "::"):
      if i != len(ts) - 2 {
        errCtx := t.Context()
        return nil, errCtx.NewError("Error: pseudo element must come last")
      }

      pseudoElementToken, err := assertNonClassOrIDWord(ts[i+1])
      if err != nil {
        return nil, err
      }

      pseudoElement = pseudoElementToken.Value()
      i += 1
    case raw.IsSymbol(t, ">"):
      if i == len(ts) - 1 {
        errCtx := t.Context()
        return nil, errCtx.NewError("Error: expected tokens after")
      }

      immediate = true

      descendant, err = ParseSelector(ts[i+1:])
      if err != nil {
        return nil, err
      }

      done = true
    case raw.IsSymbol(t, "~"):
      if i == len(ts) - 1 {
        errCtx := t.Context()
        return nil, errCtx.NewError("Error: expected tokens after")
      }

      sibling, err = ParseSelector(ts[i+1:])
      if err != nil {
        return nil, err
      }

      done = true
    case raw.IsSymbol(t, "+"):
      if i == len(ts) - 1 {
        errCtx := t.Context()
        return nil, errCtx.NewError("Error: expected tokens after")
      }

      immediate = true

      sibling, err = ParseSelector(ts[i+1:])
      if err != nil {
        return nil, err
      }

      done = true
    case raw.IsAnyWord(t) || raw.IsSymbol(t, "*"):
      descendant, err = ParseSelector(ts[i:])
      if err != nil {
        return nil, err
      }

      done = true
    default:
      errCtx := t.Context()
      return nil, errCtx.NewError("Error: bad css selector token")
    }

    if done {
      break
    }
  }

  return &SelectorData{"", "", "", filters, pseudoClasses, pseudoElement, descendant, sibling, immediate}, nil
}

func ParseSelector(ts []raw.Token) (*SelectorData, error) {
  elementName, class, id, err := parseNameToken(ts[0])
  if err != nil {
    return nil, err
  }

  if len(ts) == 1 {
    return &SelectorData{elementName, class, id, []AttrFilter{}, []PseudoClass{}, "", nil, nil, false}, nil
  } else {
    sel, err := parseFiltersPseudoAndDescendants(ts[1:])
    if err != nil {
      return nil, err
    }

    sel.elementName = elementName
    sel.class = class
    sel.id = id

    return sel, nil
  }
}

func (s *SelectorData) extend(ts []raw.Token) (Selector, error) {
  var res *SelectorData = nil

  switch {
  case raw.IsSymbol(ts[0], ">"):
    if len(ts) < 2 {
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: expected more tokens")
    }

    descendant, err := ParseSelector(ts[1:])
    if err != nil {
      return nil, err
    }

    res = s.setDescendant(descendant, true)
  case raw.IsSymbol(ts[0], "+"):
    if len(ts) < 2 {
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: expected more tokens")
    }

    sibling, err := ParseSelector(ts[1:])
    if err != nil {
      return nil, err
    }

    res = s.setSibling(sibling, true)
  case raw.IsSymbol(ts[0], "~"):
    if len(ts) < 2 {
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: expected more tokens")
    }

    sibling, err := ParseSelector(ts[1:])
    if err != nil {
      return nil, err
    }

    res = s.setSibling(sibling, false)
  case raw.IsSymbol(ts[0], "*") || raw.IsAnyWord(ts[0]):
    descendant, err := ParseSelector(ts)
    if err != nil {
      return nil, err
    }

    res = s.setDescendant(descendant, false)
  case raw.IsBracketsGroup(ts[0]) || raw.IsSymbol(ts[0], ":") || raw.IsSymbol(ts[0], "::"):
    sel, err := parseFiltersPseudoAndDescendants(ts)
    if err != nil {
      return nil, err
    }

    res = s.addFilters(sel.filters)
    res = res.addPseudoClasses(sel.pseudoClasses)
    res, err = res.setPseudoElement(sel.pseudoElement)
    if err != nil {
      errCtx := ts[0].Context()
      return nil, errCtx.NewError("Error: pseudoElement already set")
    }

    res = res.setDescendant(sel.descendant, sel.immediate)
    res = res.setSibling(sel.sibling, sel.immediate)
  default:
    errCtx := ts[0].Context()
    return nil, errCtx.NewError("Error: bad css selector extension")
  }

  return res, nil
}

func (s *SelectorData) Extend(extra *tokens.String) ([]Selector, error) {
  tss, err := parseList(extra)
  if err != nil {
    return nil, err
  }

  sels := make([]Selector, 0)
  for _, ts := range tss {
    if len(ts) == 0 {
      errCtx := extra.Context()
      return nil, errCtx.NewError("Error: expected some tokens")
    }

    sel, err := s.extend(ts)
    if err != nil {
      return nil, err
    }

    sels = append(sels, sel)
  }

  return sels, nil
}

// see if name or attr match, don't yet look at descendants or siblings
func (s *SelectorData) match(tag tree.Tag) bool {
  if s.elementName != "" && s.elementName != "*" {
    if s.elementName != tag.Name() {
      return false
    }
  }

  if s.class != "" {
    foundClass := false
    visTag, ok := tag.(tree.VisibleTag)
    if ok {
      classes := visTag.GetClasses()
      for _, class := range classes {
        if class == s.class {
          foundClass = true
          break
        }
      }
    }

    if !foundClass {
      return false
    }
  }

  if s.id != "" {
    if tag.GetID() != s.id {
      return false
    }
  }

  for _, f := range s.filters {
    if !f.Match(tag.Attributes()) {
      return false
    }
  }

  return true
}

func (s *SelectorData) Match(tag tree.Tag) []tree.Tag {
  if s.match(tag) {
    if s.sibling == nil && s.descendant == nil {
      return []tree.Tag{tag}
    } else if s.sibling != nil {
      siblings := tag.Siblings()
      if len(siblings) == 0 {
        return []tree.Tag{}
      } else if s.immediate {
        return s.sibling.Match(siblings[0])
      } else {
        res := []tree.Tag{}
        for _, sib := range siblings {
          sibRes := s.sibling.Match(sib)
          res = append(res, sibRes...)
        }
        return res
      }
    } else if s.descendant != nil {
      res := []tree.Tag{}
      for _, child := range tag.Children() {
        if s.immediate {
          if !s.descendant.match(child) {
            continue
          }
        }

        childRes := s.descendant.Match(child)
        res = append(res, childRes...)
      }

      return res
    } else {
      panic("algo error")
    }
  } else {
    // ignore this tag, and look at any of the children
    res := []tree.Tag{}
    for _, child := range tag.Children() {
      childRes := s.Match(child)
      res = append(res, childRes...)
    }

    return res
  }
}

func (s *SelectorData) Write() string {
  var b strings.Builder

  if s.elementName == "*" {
    if s.class == "" && s.id == "" {
      b.WriteString("*")
    }
  } else if s.elementName != "" {
    b.WriteString(s.elementName)
  }

  if s.class != "" {
    b.WriteString(".")
    b.WriteString(s.class)
  }

  if s.id != "" {
    b.WriteString("#")
    b.WriteString(s.id)
  }

  for _, f := range s.filters {
    b.WriteString("[")
    b.WriteString(f.Write())
    b.WriteString("]")
  }

  for _, p := range s.pseudoClasses {
    b.WriteString(":")
    b.WriteString(p.Write())
  }

  if s.pseudoElement != "" {
    b.WriteString("::")
    b.WriteString(s.pseudoElement)
  }

  if s.descendant != nil {
    if s.immediate {
      b.WriteString(" > ")
    } else {
      b.WriteString(" ")
    }

    b.WriteString(s.descendant.Write())
  } else if s.sibling != nil {
    if s.immediate {
      b.WriteString(" + ")
    } else {
      b.WriteString(" ~ ")
    }

    b.WriteString(s.sibling.Write())
  }

  return b.String()
}
