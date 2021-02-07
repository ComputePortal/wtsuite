package styles

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

type PseudoClass interface {
  Write() string
}

type GenericPseudoClass struct {
  name string
  parensContent string // uses simple context print, so far from ideal
}

func NewGenericPseudoClass(name string) *GenericPseudoClass {
  return &GenericPseudoClass{name, ""}
}

func NewGenericPseudoClassWithArgs(name string, args string) *GenericPseudoClass {
  return &GenericPseudoClass{name, args}
}

func (p *GenericPseudoClass) Write() string {
  var b strings.Builder

  b.WriteString(p.name)
  if p.parensContent != "" {
    b.WriteString("(")
    b.WriteString(p.parensContent)
    b.WriteString(")")
  }

  return b.String()
}


func ParsePseudoClass(ts []raw.Token) (PseudoClass, error) {
  if len(ts) == 0 {
    panic("empty list")
  }

  nameToken, err := assertNonClassOrIDWord(ts[0])
  if err != nil {
    return nil, err
  }


  if len(ts) == 2 {
    parens, err := raw.AssertParensGroup(ts[1])
    if err != nil {
      return nil, err
    }

    switch nameToken.Value() {
    default:
      pCtx := parens.Context()
      parensContent := pCtx.Content()
      parensContent = parensContent[1:len(parensContent)-1]

      return NewGenericPseudoClassWithArgs(nameToken.Value(), parensContent), nil
    }
  } else if len(ts) > 2 {
    errCtx := ts[2].Context()
    return nil, errCtx.NewError("Error: unexpected token")
  } else {
    return NewGenericPseudoClass(nameToken.Value()), nil
  }
}
