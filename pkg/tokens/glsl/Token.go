package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Token interface {
  Dump(intent string) string
  Context() context.Context
}

type TokenData struct {
  ctx context.Context
}

func newTokenData(ctx context.Context) TokenData {
  return TokenData{ctx}
}

func (t *TokenData) Context() context.Context {
  return t.ctx
}
