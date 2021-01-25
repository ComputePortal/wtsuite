package html

import (
  "strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

// the number of children of a tag aren't available until after the tag has been fully processed, but the number of children can be useful inside the css (eg. to set the max height of a flexbox so that wrapping creates a column effect)
// hence all attributes should be evaluated one final time using scope==nil before collecting the css rules, writing the result

type Lazy struct {
  fn func(FinalTag) (Token, error) // fn shouldn't have side-effects
  TokenData
}

func NewLazy(fn func(FinalTag) (Token, error), ctx context.Context) *Lazy {
  return &Lazy{fn, TokenData{ctx}}
}

func (t *Lazy) Eval(scope Scope) (Token, error) {
  return t, nil
}

func (t *Lazy) EvalLazy(tag FinalTag) (Token, error) {
  return t.fn(tag)
}

func (t *Lazy) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("Lazy\n")

  return b.String()
}

func (t *Lazy) IsSame(other_ Token) bool {
  if other, ok := other_.(*Lazy); ok {
    return other == t 
  }

  return false
}

func IsLazy(t Token) bool {
  _, ok := t.(*Lazy)
  return ok
}
