package prototypes

import (
  "../values"

  "../../context"
)

func checkParent(p values.Prototype, other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(values.Prototype); ok {
    otherParent, err := other.GetParent()
    if err == nil {
      if otherParent != nil {
        if err := p.Check(otherParent, ctx); err == nil {
          return nil
        }
      }
    }

  } 

  return ctx.NewError("Error: expected " + p.Name() + ", got " + other_.Name())
}
