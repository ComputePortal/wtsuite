package prototypes

import (
	"../values"

	"../../context"
)

var HTMLCollection *HTMLCollectionPrototype = allocHTMLCollectionPrototype()

type HTMLCollectionPrototype struct {
	BuiltinPrototype
}

func allocHTMLCollectionPrototype() *HTMLCollectionPrototype {
	return &HTMLCollectionPrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

// exactly the same as BuiltinPrototype.Check, but with *HTMLCollectionPrototype receiver
func (p *HTMLCollectionPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *HTMLCollectionPrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*HTMLCollectionPrototype); ok {
		if other == p {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (p *HTMLCollectionPrototype) GetIndex(stack values.Stack, this *values.Instance,
	index values.Value, ctx context.Context) (values.Value, error) {
	return NewInstance(HTMLElement, ctx), nil
}

func (p *HTMLCollectionPrototype) SetIndex(stack values.Stack, this *values.Instance,
	index values.Value, arg values.Value, ctx context.Context) error {
	return CheckInputs(&And{Int, HTMLElement}, []values.Value{index, arg}, ctx)
}

func generateHTMLCollectionPrototype() bool {
	*HTMLCollection = HTMLCollectionPrototype{BuiltinPrototype{
		"HTMLCollection", nil,
		map[string]BuiltinFunction{
			"item":      NewNormal(Int, HTMLElement),
			"length":    NewGetter(Int),
			"namedItem": NewNormal(String, HTMLElement),
		},
		NewNoContentGenerator(HTMLCollection),
	},
	}

	return true
}

var _HTMLCollectionOk = generateHTMLCollectionPrototype()
