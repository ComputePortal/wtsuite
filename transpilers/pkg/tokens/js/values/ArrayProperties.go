package values

import (
	"../../context"
)

type ArrayProperties struct {
	isLiteral bool // inplace change to Array (eg. in-place sort)
	items     []Value
	PropertiesData
}

func NewArrayProperties(isLiteral bool, items []Value, ctx context.Context) Properties {
	if !isLiteral {
		items = UnpackMulti(items)
	}

	return &ArrayProperties{isLiteral, items, newPropertiesData(ctx)}
}

func (p *ArrayProperties) Copy(cache CopyCache) Properties {
	innerCpy := p.PropertiesData.copy(cache)

	itemsCpy := copyValueList(p.items, cache)

	return &ArrayProperties{p.isLiteral, itemsCpy, innerCpy}
}

func (p *ArrayProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*ArrayProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	if p.isLiteral && other.isLiteral {
		mergedItems := mergeValueLists(p.items, other.items)
		if mergedItems != nil {
			return &ArrayProperties{true, mergedItems, newPropertiesDataWithContent(mergedProps, p.Context())}
		}
	}

	allItems := append(p.items, other.items...)
	return &ArrayProperties{false, allItems, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func (p *ArrayProperties) Length() (int, bool) {
	if p.isLiteral {
		return len(p.items), true
	} else {
		return 0, false
	}
}

func (p *ArrayProperties) GetItem() Value {
	if len(p.items) == 0 {
		//if !p.isLiteral {
		//hereCtx := p.Context()
		//panic(hereCtx.NewError("empty array that is not a literal"))
		// eg. Array() construction
		return NewAllNull(p.Context())
		//} else {
		//panic("use LiteralValues() instead")
		//}
	} else {
		return NewMulti(p.items, p.Context())
	}
}

func (p *ArrayProperties) LiteralValues() ([]Value, bool) {
	if p.isLiteral {
		return p.items, true
	} else {
		return nil, false
	}
}

func (p *ArrayProperties) ChangeOrder() {
	p.isLiteral = false
}

func (p *ArrayProperties) AppendItem(item Value) {
	p.isLiteral = false

	p.items = append(p.items, item)
}

func AssertArrayProperties(p_ Properties) *ArrayProperties {
	p, ok := p_.(*ArrayProperties)
	if ok {
		return p
	} else {
		panic("not ArrayProperties")
	}
}

func (p *ArrayProperties) RemoveLiteralness() *ArrayProperties {
	return &ArrayProperties{false, p.items, newPropertiesDataWithContent(p.props, p.Context())}
}

func (p *ArrayProperties) LoopNestedPrototypes(fn func(Prototype)) {
	p.PropertiesData.LoopNestedPrototypes(fn)

	for _, item := range p.items {
		item.LoopNestedPrototypes(fn)
	}
}
