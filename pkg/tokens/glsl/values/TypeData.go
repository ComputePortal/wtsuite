package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type TypeData struct {
  name string
  ctx context.Context
}

func newTypeData(name string, ctx context.Context) TypeData {
  return TypeData{name, ctx}
}

func (v *TypeData) Context() context.Context {
  return v.ctx
}

func (v *TypeData) TypeName() string {
  return "type(" + v.name + ")"
}

func (v *TypeData) GetMember(key string, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get member of type")
}

func (v *TypeData) SetMember(key string, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set member of type")
}

func (v *TypeData) GetIndex(idx *LiteralInt, ctx context.Context) (Value, error) {
  return nil, ctx.NewError("Error: can't get index of type")
}

func (v *TypeData) SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error {
  return ctx.NewError("Error: can't set index of type")
}

func (v *TypeData) LiteralIntValue() (int, bool) {
  return 0, false
}

func (v *TypeData) Length() int {
  return 0
}
