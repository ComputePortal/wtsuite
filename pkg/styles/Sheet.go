package styles

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/directives"
)

type Sheet interface {
  Append(r Rule)
  IsEmpty() bool
  Len() int
  //IsNotLazy() bool
  Write(compr bool, nl string, tab string) (string, error)
  ExpandNested() (Sheet, error) // expanding a second time does nothing
  //ExpandLazy(root *tree.Root) (Sheet, error)
}

type SheetData struct {
  rules []Rule
}

func NewSheet() Sheet {
  return &SheetData{make([]Rule, 0)}
}

func (s *SheetData) Append(r Rule) {
  s.rules = append(s.rules, r)
}

func (s *SheetData) IsEmpty() bool {
  return len(s.rules) == 0
}

func (s *SheetData) Len() int {
  return len(s.rules)
}

/*func (s *SheetData) IsNotLazy() bool {
  for _, r := range rules {
    if r.IsLazy() {
      return false
    }
  }

  return true
}*/

func (s *SheetData) Write(compr bool, nl string, tab string) (string, error) {
  var b strings.Builder

  if directives.MATH_FONT_URL != "" {
    b.WriteString(writeMathFontFace(directives.MATH_FONT_URL))
  }

  for _, r := range s.rules {
    inner, err := r.Write("", nl, tab)
    if err != nil {
      return "", err
    }
    b.WriteString(inner)
  }

  res := b.String()

  if compr {
    res = compress(res)
  }

  return res, nil
}

func (s *SheetData) ExpandNested() (Sheet, error) {
  rules := make([]Rule, 0)

  for _, r := range s.rules {
    expanded, err := r.ExpandNested()
    if err != nil {
      return nil, err
    }

    rules = append(rules, expanded...)
  }

  return &SheetData{rules}, nil
}
