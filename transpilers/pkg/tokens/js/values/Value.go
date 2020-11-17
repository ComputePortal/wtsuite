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
func CommonValue(vs []Value) Value {
  ctx := context.NewDummyContext()

  if len(vs) == 0 {
    panic("expected at least 1")
  } else if len(vs) == 1 {
    return vs[0]
  } else {
    for i, v := range vs {
      found := true
      for j, other := range vs {
        if i != j {
          if err := v.Check(other, ctx); err != nil {
            found = false
            break
          }
        }
      }

      if found {
        return v
      }
    }

    return NewAny(vs[0].Context())
  }
}
