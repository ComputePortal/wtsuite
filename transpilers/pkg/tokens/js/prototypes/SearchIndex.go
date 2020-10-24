package prototypes

import (
	"../values"

	"../../context"
)

var SearchIndex *BuiltinPrototype = allocBuiltinPrototype()

func generateSearchIndexPrototype() bool {
	*SearchIndex = BuiltinPrototype{
		"SearchIndex", nil,
		map[string]BuiltinFunction{
			"ignore": NewNormal(String, Boolean),
			"page": NewNormalFunction(Int,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewObject(map[string]values.Value{
						"url":     NewString(ctx),
						"title":   NewString(ctx),
						"content": NewArray([]values.Value{NewString(ctx)}, ctx),
					}, ctx), nil
				}),
			"match": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewSet([]values.Value{NewInt(ctx)}, ctx), nil
				}),
			"matchPrefix": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewSet([]values.Value{NewInt(ctx)}, ctx), nil
				}),
			"matchSuffix": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewSet([]values.Value{NewInt(ctx)}, ctx), nil
				}),
			"matchSubstring": NewNormalFunction(String,
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewSet([]values.Value{NewInt(ctx)}, ctx), nil
				}),
			// fuzzy is just simple levenshtein edit distance
			"fuzzy": NewNormalFunction(&And{String, Int},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewArray([]values.Value{NewSet([]values.Value{NewInt(ctx)}, ctx)}, ctx), nil
				}),
			"fuzzyPrefix": NewNormalFunction(&And{String, Int},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewArray([]values.Value{NewSet([]values.Value{NewInt(ctx)}, ctx)}, ctx), nil
				}),
			"fuzzySuffix": NewNormalFunction(&And{String, Int},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewArray([]values.Value{NewSet([]values.Value{NewInt(ctx)}, ctx)}, ctx), nil
				}),
			"fuzzySubstring": NewNormalFunction(&And{String, Int},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					return NewArray([]values.Value{NewSet([]values.Value{NewInt(ctx)}, ctx)}, ctx), nil
				}),
			"onready": NewSetterFunction(&Function{},
				func(stack values.Stack, this *values.Instance, args []values.Value,
					ctx context.Context) (values.Value, error) {
					arg := args[0]
					if err := arg.EvalMethod(stack.Parent(), []values.Value{}, ctx); err != nil {
						return nil, err
					}

					return nil, nil
				}),
		},
		NewConstructorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if err := CheckInputs(&And{String, &Opt{Object}}, args, ctx); err != nil {
				return nil, err
			}

			return NewInstance(SearchIndex, ctx), nil
		}),
	}

	return true
}

var _SearchIndexOk = generateSearchIndexPrototype()
