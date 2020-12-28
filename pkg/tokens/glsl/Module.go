package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Module struct {
  dependencies []*LiteralString // relative paths!

  Block
}

func NewModule(ctx context.Context) *Module {
  return &Module{
    make([]*LiteralString, 0),
    newBlock(ctx),
  }
}
