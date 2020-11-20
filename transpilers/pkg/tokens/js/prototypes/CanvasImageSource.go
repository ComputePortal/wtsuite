package prototypes

import (
  "../values"

  "../../context"
)

type CanvasImageSource struct {
  AbstractBuiltinInterface
}

var canvasImageSourceInterface values.Interface = &CanvasImageSource{newAbstractBuiltinInterface("CanvasImageSource")}

func NewCanvasImageSourceInterface() values.Interface {
  return canvasImageSourceInterface
}

func NewCanvasImageSource(ctx context.Context) values.Value {
  return values.NewInstance(NewCanvasImageSourceInterface(), ctx)
}

// interfaces usually check in a different way, but this isn't really an interface
func (p *CanvasImageSource) Check(other_ values.Interface, ctx context.Context) error {
  if _, ok := other_.(*CanvasImageSource); ok {
    return nil
  } else if proto, ok := other_.(values.Prototype); ok {
    checkProtos, _ := p.GetPrototypes()
    for _, checkProto := range checkProtos {
      if err := checkProto.Check(proto, ctx); err == nil {
        return nil
      }
    }
  } 

  return ctx.NewError("Error: not an CanvasImageSource")
}

func (p *CanvasImageSource) GetPrototypes() ([]values.Prototype, error) {
  return []values.Prototype{
    //NewCSSImageValuePrototype(),
    NewHTMLCanvasElementPrototype(),
    NewHTMLImageElementPrototype(),
    //NewHTMLVideoElementPrototype(),
    //NewImageBitmapPrototype(),
    //NewOffscreenCanvasPrototype(),
    //NewSVGImageElementPrototype(),
  }, nil
}

func (p *CanvasImageSource) IsUniversal() bool {
  return false
}
