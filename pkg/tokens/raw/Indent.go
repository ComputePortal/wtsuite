package raw

import (
  "strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Indent struct {
  n int
  TokenData
}

func NewIndent(n int, ctx context.Context) *Indent {
  return &Indent{n, TokenData{ctx}}
}

func (t *Indent) Dump(indent string) string {
  return indent + "Indent(" + strconv.Itoa(t.n) + ")\n"
}

func (t *Indent) N() int {
  return t.n
}

func IsIndent(t Token) bool {
  _, ok := t.(*Indent)
  return ok
}

func IsWhitespace(t Token) bool {
  return IsIndent(t) || IsNL(t)
}

func AssertIndent(t_ Token) (*Indent, error) {
  t, ok := t_.(*Indent)
  if ok {
    return t, nil
  } else {
    errCtx := t_.Context()
    return nil, errCtx.NewError("Error: expected an indent")
  }
}
