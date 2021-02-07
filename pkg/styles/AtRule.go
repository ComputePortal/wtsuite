package styles

import (
  "strings"
)

type AtRule struct {
  sel *AtSelector // linked rules
  rules []Rule // regular rules
}

func NewAtRule(sel *AtSelector, rules []Rule) *AtRule {
  return &AtRule{sel, rules}
}

func (r *AtRule) SetParent(parent *AtSelector) {
  r.sel.SetParent(parent)
}

// not called the first time. The subsequent times nothing should change
func (r *AtRule) ExpandNested() ([]Rule, error) {
  return []Rule{r}, nil
}

func (r *AtRule) Write(indent string, nl string, tab string) (string, error) {
  var b strings.Builder

  // collect all the at selectors
  sels := collectAtSelectors(r.sel)

  for i, sel := range sels {
    b.WriteString(indent)
    for j := 0; j < i; j++ {
      b.WriteString(tab)
    }
    b.WriteString(sel.Write())
    b.WriteString(" {")
    b.WriteString(nl)
  }

  innerIndent := indent
  for j := 0; j < len(sels); j++ {
    innerIndent += tab
  }

  for _, rule := range r.rules  {
    inner, err := rule.Write(innerIndent, nl, tab)
    if err != nil {
      return "", err
    }
    b.WriteString(inner)
    b.WriteString(nl)
  }

  for i, _ := range sels {
    b.WriteString(indent)

    for j := len(sels)-i-1; j >= 0; j-- {
      b.WriteString(tab)
    }

    b.WriteString("}")
    b.WriteString(nl)
  }

  return b.String(), nil
}
