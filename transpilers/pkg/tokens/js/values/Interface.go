package values

import (
	"../../context"
)

type Interface interface {
  Name() string

  Context() context.Context

  Check(other Interface, ctx context.Context) error

  // get extended interfaces in case of js.Interface, get implements interfaces in case of js.Class
  GetInterfaces() ([]Interface, error)

  // get all prototypes that implement this interface (actual prototypes dont need include themselves) (used by InstanceOf.Write())
  GetPrototypes() ([]Prototype, error)

  // returns nil if it doesnt exist
  GetInstanceMember(key string, includePrivate bool, ctx context.Context) (Value, error)

  SetInstanceMember(key string, includePrivate bool, arg Value, ctx context.Context) error
}

// returns nil if not an Instance with an Interface
func GetInterface(v_ Value) Interface {
  v_ = UnpackContextValue(v_)

  switch v := v_.(type) {
  case *Instance:
    return v.GetInterface()
  default:
    return nil
  }
}
