package styles

import (
  "strings"

	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

type AttrFilter interface {
  Match(attr *tokens.StringDict) bool
  Write() string
}

type AttrFilterData struct {
  name string
  value string
  ci bool
}

func (f *AttrFilterData) getValue(attr *tokens.StringDict) (string, bool) {
  v_, ok := attr.Get(f.name)
  if ok && tokens.IsPrimitive(v_) {
    v, err := tokens.AssertPrimitive(v_)
    if err != nil {
      panic(err)
    }

    return v.Write(), true
  }

  return "", false
}

func (f *AttrFilterData) getComp(v string) (string, string) {
  if f.ci {
    return strings.ToLower(v), strings.ToLower(f.value)
  } else {
    return v, f.value
  }
}

func (f *AttrFilterData) write(op string) string {
  var b strings.Builder

  b.WriteString(f.name)
  b.WriteString(op)
  b.WriteString("\"")
  b.WriteString(f.value)
  b.WriteString("\"")

  if f.ci {
    b.WriteString(" i")
  }

  return b.String()
}

type PlainAttrFilter struct {
  AttrFilterData
}

func NewPlainAttrFilter(name string) *PlainAttrFilter {
  // ci doesnt matter here
  return &PlainAttrFilter{AttrFilterData{name, "", false}}
}

func (f *PlainAttrFilter) Match(attr *tokens.StringDict) bool {
  _, ok := f.getValue(attr)
  return ok
}

func (f *PlainAttrFilter) Write() string {
  return f.name
}

type ExactAttrFilter struct {
  AttrFilterData
}

// value should be in quotes in the final code
func NewExactAttrFilter(name string, value string, ci bool) *ExactAttrFilter {
  return &ExactAttrFilter{AttrFilterData{name, value, ci}}
}

func (f *ExactAttrFilter) Match(attr *tokens.StringDict) bool {
  v, ok := f.getValue(attr)
  if ok {
    v, comp := f.getComp(v)

    return v == comp
  } 

  return false
}

func (f *ExactAttrFilter) Write() string {
  return f.write("=")
}

type FieldAttrFilter struct {
  AttrFilterData
}

func NewFieldAttrFilter(name string, value string, ci bool) *FieldAttrFilter {
  return &FieldAttrFilter{AttrFilterData{name, value, ci}}
}

func (f *FieldAttrFilter) Match(attr *tokens.StringDict) bool {
  v, ok := f.getValue(attr)
  if ok {
    fields := strings.Fields(v)
    for _, field := range fields {
      if field == f.value {
        return true
      }
    }
  }

  return false
}

func (f *FieldAttrFilter) Write() string {
  return f.write("~=")
}

type SubcodeAttrFilter struct {
  AttrFilterData
}

func NewSubcodeAttrFilter(name string, value string, ci bool) *SubcodeAttrFilter {
  return &SubcodeAttrFilter{AttrFilterData{name, value, ci}}
}

func (f *SubcodeAttrFilter) Match(attr *tokens.StringDict) bool {
  v, ok := f.getValue(attr)
  if ok {
    v, comp := f.getComp(v)

    if v == comp || strings.HasPrefix(v, comp + "-") {
      return true
    }
  }

  return false
}

func (f *SubcodeAttrFilter) Write() string {
  return f.write("|=")
}

type PrefixAttrFilter struct {
  AttrFilterData
}

func NewPrefixAttrFilter(name string, value string, ci bool) *PrefixAttrFilter {
  return &PrefixAttrFilter{AttrFilterData{name, value, ci}}
}

func (f *PrefixAttrFilter) Match(attr *tokens.StringDict) bool {
  v, ok := f.getValue(attr)
  if ok {
    v, comp := f.getComp(v)

    if strings.HasPrefix(v, comp) {
      return true
    }
  }

  return false
}

func (f *PrefixAttrFilter) Write() string {
  return f.write("^=")
}

type SuffixAttrFilter struct {
  AttrFilterData
}

func NewSuffixAttrFilter(name string, value string, ci bool) *SuffixAttrFilter {
  return &SuffixAttrFilter{AttrFilterData{name, value, ci}}
}

func (f *SuffixAttrFilter) Match(attr *tokens.StringDict) bool {
  v, ok := f.getValue(attr)
  if ok {
    v, comp := f.getComp(v)
    
    if strings.HasSuffix(v, comp) {
      return true
    }
  }

  return false
}

func (f *SuffixAttrFilter) Write() string {
  return f.write("$=")
}

type ContainsAttrFilter struct {
  AttrFilterData
}

func NewContainsAttrFilter(name string, value string, ci bool) *ContainsAttrFilter {
  return &ContainsAttrFilter{AttrFilterData{name, value, ci}}
}

func (f *ContainsAttrFilter) Match(attr *tokens.StringDict) bool {
  v, ok := f.getValue(attr)
  if ok {
    v, comp := f.getComp(v) 

    if strings.Contains(v, comp) {
      return true
    }
  }

  return false
}

func (f *ContainsAttrFilter) Write() string {
  return f.write("*=")
}

func ParseAttrFilter(t_ raw.Token) (AttrFilter, error) {
  t, err := raw.AssertBracketsGroup(t_)
  if err != nil {
    return nil, err
  }

  if len(t.Fields) != 1 {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: unexpected separator")
  }

  inner := t.Fields[0]

  if len(inner) != 1 && len(inner) != 3 && len(inner) != 4 {
    errCtx := t.Context()
    return nil, errCtx.NewError("Error: bad attribute selector")
  }

  w, err := assertNonClassOrIDWord(inner[0])
  if err != nil {
    return nil, err
  }

  ci := false
  if len(inner) == 4 {
    ciToken, err := raw.AssertWord(inner[3])
    if err != nil {
      return nil, err
    }

    if ciToken.Value() != "i" && ciToken.Value() != "I" {
      errCtx := ciToken.Context()
      return nil, errCtx.NewError("Error: expected i or I")
    }

    ci = true
  }

  if len(inner) > 2 {
    valueToken, err := raw.AssertLiteralString(inner[2])
    if err != nil {
      return nil, err
    }
    
    value := valueToken.Value()

    switch {
    case raw.IsSymbol(inner[1], "="):
      return NewExactAttrFilter(w.Value(), value, ci), nil
    case raw.IsSymbol(inner[1], "~="):
      return NewFieldAttrFilter(w.Value(), value, ci), nil
    case raw.IsSymbol(inner[1], "|="):
      return NewSubcodeAttrFilter(w.Value(), value, ci), nil
    case raw.IsSymbol(inner[1], "^="):
      return NewPrefixAttrFilter(w.Value(), value, ci), nil
    case raw.IsSymbol(inner[1], "$="):
      return NewSuffixAttrFilter(w.Value(), value, ci), nil
    case raw.IsSymbol(inner[1], "*="):
      return NewContainsAttrFilter(w.Value(), value, ci), nil
    default:
      errCtx := inner[1].Context()
      return nil, errCtx.NewError("Error: unrecognized attribute selector op (" + inner[1].Dump("") + ")")
    }
  } else {
    return NewPlainAttrFilter(w.Value()), nil
  }
}
