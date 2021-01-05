package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Structable interface {
  Name() string
  CheckConstruction(args []Value, ctx context.Context) error
  GetMember(key string, ctx context.Context) (Value, error)
  SetMember(key string, arg Value, ctx context.Context) error
}
