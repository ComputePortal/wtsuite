package prototypes

import (
	"../values"

	"../../context"
)

var Window *BuiltinPrototype = allocBuiltinPrototype()

// also usable for setInterval and requestIdleCallback
func WindowSetTimeout(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	// keep this CheckInputs, so it is available outside the prototype
	if err := CheckInputs(&And{&Function{}, &Opt{Number}}, args, ctx); err != nil {
		return nil, err
	}

	if err := args[0].EvalMethod(stack, []values.Value{}, ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func WindowRequestIdleCallback(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	// keep this CheckInputs, so it is available outside the prototype
	if err := CheckInputs(&And{&Function{}, &Opt{Object}}, args, ctx); err != nil {
		return nil, err
	}

	// TODO: should we check that the Object contains the timeout key?

	if err := args[0].EvalMethod(stack, []values.Value{}, ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

func WindowFetch(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	// keep this CheckInputs, so it is available outside the prototype
	if err := CheckInputs(String, args, ctx); err != nil {
		return nil, err
	}

	return NewResolvedPromise(NewInstance(Response, ctx), ctx)
}

func generateWindowPrototype() bool {
	*Window = BuiltinPrototype{
		"Window", EventTarget,
		map[string]BuiltinFunction{
			"atob":             NewNormal(String, String),
			"btoa":             NewNormal(String, String),
			"blur":             NewNormal(&None{}, nil),
			"close":            NewNormal(&None{}, nil),
			"devicePixelRatio": NewGetter(Number),
			"fetch":            NewNormalFunction(String, WindowFetch),
			"focus":            NewNormal(&None{}, nil),
			"getComputedStyle": NewNormal(HTMLElement, CSSStyleDeclaration),
			"indexedDB":        NewGetter(IDBFactory),
			"innerHeight":      NewGetter(Number),
			"innerWidth":       NewGetter(Number),
			"localStorage":     NewGetter(Storage),
			"location":         NewGetter(Location),
			"open":             NewMethodLikeNormal(&And{String, &Opt{String}}, Window),
			"requestAnimationFrame": NewNormalFunction(&Function{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {

					timeStamp := NewInstance(Number, ctx)

					if err := args[0].EvalMethod(stack, []values.Value{timeStamp},
						ctx); err != nil {
						return nil, err
					}

					return nil, nil
				}),
			"requestIdleCallback": NewNormalFunction(&And{&Function{}, &Opt{Object}}, WindowRequestIdleCallback),
			"scrollTo":            NewNormal(&And{Number, Number}, nil),
			"scrollX":             NewGetter(Number),
			"scrollY":             NewGetter(Number),
			"sessionStorage":      NewGetter(Storage),
			"setInterval":         NewNormalFunction(&And{&Function{}, &Opt{Number}}, WindowSetTimeout),
			"setTimeout":          NewNormalFunction(&And{&Function{}, &Opt{Number}}, WindowSetTimeout),
		},
		nil,
	}

	return true
}

var _WindowOk = generateWindowPrototype()
