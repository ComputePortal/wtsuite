package raw

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type NL struct {
  TokenData
}

func NewNL(ctx context.Context) *NL {
  return &NL{TokenData{ctx}}
}

func (t *NL) Dump(indent string) string {
  return indent + "NL\n"
}

func IsNL(t Token) bool {
  _, ok := t.(*NL)
  return ok
}
