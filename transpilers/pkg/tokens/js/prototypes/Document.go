package prototypes

import (
	"../values"

	"../../context"
)

var Document *BuiltinPrototype = allocBuiltinPrototype()

func generateDocumentPrototype() bool {
	*Document = BuiltinPrototype{
		"Document", EventTarget,
		map[string]BuiltinFunction{
			"activeElement":  NewGetter(HTMLElement),
			"body":           NewGetter(HTMLElement),
			"cookie":         NewGetterSetter(String),
			"createTextNode": NewNormal(String, Text),
			"createElement": NewNormalFunction(&And{String, &Opt{Object}},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					if str, ok := args[0].LiteralStringValue(); ok {
						switch str {
						case "a":
							return NewInstance(HTMLLinkElement, ctx), nil
						case "canvas":
							return NewInstance(HTMLCanvasElement, ctx), nil
						case "div":
							return NewInstance(HTMLElement, ctx), nil
						case "img":
							return NewInstance(HTMLImageElement, ctx), nil
						case "input":
							return NewInstance(HTMLInputElement, ctx), nil
						}
					}

					return NewInstance(HTMLElement, ctx), nil
				}),
			"documentElement": NewGetter(HTMLElement),
      "fonts":           NewGetter(FontFaceSet),
			"getElementById": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					errCtx := args[0].Context()
					isMain := false
					if isMainBool, ok := this.Properties().GetProperty(".main"); ok {
						mainBoolVal, ok1 := isMainBool.LiteralBooleanValue()
						if ok1 && mainBoolVal {
							isMain = true
						}
					}

					if str, ok := args[0].LiteralStringValue(); ok && isMain && str != "" {
						vif := stack.GetViewInterface()

						return vif.GetElemTypeInstance(stack, str, errCtx)
					} else {
						return NewInstance(HTMLElement, ctx), nil
						//return nil, errCtx.NewError("Error: expected literal string")
					}
				}),
			"getVariable": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					errCtx := args[0].Context()
					if str, ok := args[0].LiteralStringValue(); ok {
						vif := stack.GetViewInterface()

						return vif.GetVarTypeInstance(stack, str, errCtx)
					} else {
						return nil, errCtx.NewError("Error: expected literal string")
					}
				}),
			"hidden": NewGetter(Boolean),
			"newElement": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {

					errCtx := args[0].Context()
					if str, ok := args[0].LiteralStringValue(); ok {
						vif := stack.GetViewInterface()

						return vif.GetDefTypeInstance(stack, str, errCtx)
					} else {
						return nil, errCtx.NewError("Error: expected literal string")
					}
				}),
			"querySelector": NewNormal(String, HTMLElement),
			"referrer":      NewGetterSetter(String),
			"title":         NewGetterSetter(String),
		},
		nil,
	}

	return true
}

var _DocumentOk = generateDocumentPrototype()
