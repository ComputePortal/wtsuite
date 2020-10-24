package prototypes

import (
	"../values"

	"../../context"
)

var HTMLCanvasElement *BuiltinPrototype = allocBuiltinPrototype()

func NewHTMLCanvasElement(ctx context.Context) *values.Instance {
	return NewInstance(HTMLCanvasElement, ctx)
}

func generateHTMLCanvasElementPrototype() bool {
	*HTMLCanvasElement = BuiltinPrototype{
		"HTMLCanvasElement", HTMLElement,
		map[string]BuiltinFunction{
			"getContext": NewNormalFunction(&And{String, &Opt{Object}},
				func(stack values.Stack, this *values.Instance,
					args []values.Value, ctx context.Context) (values.Value, error) {
					str, ok := args[0].LiteralStringValue()
					if !ok {
						return nil, ctx.NewError("Error: expected literal string value as argument")
					}

					switch str {
					case "2d":
						return NewInstance(CanvasRenderingContext2D, ctx), nil
					case "webgl":
						return NewInstance(WebGLRenderingContext, ctx), nil
					default:
						return nil, ctx.NewError("Error: expected '2d' or 'webgl', got '" + str + "'")
					}
				}),
			"height":    NewGetterSetter(Int),
			"toDataURL": NewNormal(&And{&Opt{String}, &Opt{Number}}, String), // defaults: type="image/png", quality=0.92
			"width":     NewGetterSetter(Int),
		},
		NewConstructorGeneratorFunction(nil,
			func(stack values.Stack, keys []string, args []values.Value,
				ctx context.Context) (values.Value, error) {
				if keys != nil || args != nil {
					return nil, ctx.NewError("Error: unexpected content types")
				}
				return NewInstance(HTMLCanvasElement, ctx), nil
			}),
	}

	return true
}

var _HTMLCanvasElementOk = generateHTMLCanvasElementPrototype()
