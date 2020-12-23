package values

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Prototype interface {
	Interface

  IsAbstract() bool

  IsFinal() bool

  // returns nil if it doesn't have a parent
  GetParent() (Prototype, error)

  // return nil if it doesnt exist
  GetClassMember(key string, includePrivate bool, ctx context.Context) (Value, error)

  // return nil if constructor doesn't exist
	GetClassValue() (*Class, error)
}

// returns nil if not an Instance with an Interface
func GetPrototype(v_ Value) Prototype {
  v_ = UnpackContextValue(v_)

  switch v := v_.(type) {
  case *Instance:
    interf := v.GetInterface()
    if proto, ok := interf.(Prototype); ok {
      return proto
    } else {
      return nil
    }
  case *LiteralIntInstance:
    interf := v.GetInterface()
    if proto, ok := interf.(Prototype); ok {
      return proto
    } else {
      return nil
    }
  case *LiteralBooleanInstance:
    interf := v.GetInterface()
    if proto, ok := interf.(Prototype); ok {
      return proto
    } else {
      return nil
    }
  case *LiteralStringInstance:
    interf := v.GetInterface()
    if proto, ok := interf.(Prototype); ok {
      return proto
    } else {
      return nil
    }
  default:
    return nil
  }
}

func PrototypeIsAncestorOf(parent Prototype, child Prototype) bool {
  for child != nil {
    if parent == child {
      return true
    }

    var err error
    child, err = child.GetParent()
    if err != nil {
      return false
    }
  }

  return false
}
