package values

import (
	"../../context"
)

var VERBOSITY = 0

type Value interface {
	Context() context.Context

	TypeName() string
  Check(other Value, ctx context.Context) error

	EvalConstructor(args []Value, ctx context.Context) (Value, error)
	EvalFunction(args []Value, preferMethod bool, ctx context.Context) (Value, error) // method returns nil Value

	GetMember(key string, includePrivate bool, ctx context.Context) (Value, error)
  SetMember(key string, includePrivate bool, arg Value, ctx context.Context) error

	LiteralBooleanValue() (bool, bool)
  LiteralIntValue() (int, bool)
	LiteralStringValue() (string, bool)
}

type ValueData struct {
  // don't store Type here because we need access to specific parts depending on the Value
	ctx context.Context
}

func (v *ValueData) Context() context.Context {
	return v.ctx
}

func (v *ValueData) LiteralBooleanValue() (bool, bool) {
	return false, false
}

func (v *ValueData) LiteralIntValue() (int, bool) {
	return 0, false
}

func (v *ValueData) LiteralStringValue() (string, bool) {
	return "", false
}

// TODO: use parent prototypes too
func CommonValue(vs []Value, ctx context.Context) Value {
  if len(vs) == 0 {
    return NewAny(ctx)
  } else if len(vs) == 1 {
    return RemoveLiteralness(vs[0])
  } else {
    for i, v := range vs {
      v = RemoveLiteralness(v)
      found := true
      for j, other := range vs {
        other = RemoveLiteralness(other)
        if i != j {
          if err := v.Check(other, ctx); err != nil {
            found = false
            break
          }
        }
      }

      if found {
        return NewContextValue(v, ctx)
      }
    }

    return NewAny(ctx)
  }
}
