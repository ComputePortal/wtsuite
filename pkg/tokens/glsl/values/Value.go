package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Value interface {
  Context() context.Context

  TypeName() string
  Check(other Value, ctx context.Context) error // types can check values, values can check values

  Instantiate(ctx context.Context) (Value, error) // only works for types, err if not a type
  EvalFunction(args []Value, ctx context.Context) (Value, error) // method returns nil Value

  GetMember(key string, ctx context.Context) (Value, error)
  SetMember(key string, arg Value, ctx context.Context) error

  GetIndex(idx *LiteralInt, ctx context.Context) (Value, error)
  SetIndex(idx *LiteralInt, arg Value, ctx context.Context) error

  LiteralIntValue() (int, bool)
  Length() int // mostly 1, more than 1 for vectors and arrays
}
