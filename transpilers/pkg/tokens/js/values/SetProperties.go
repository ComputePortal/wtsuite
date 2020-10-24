package values

import (
	"../../context"
)

type SetProperties struct {
	items []Value
	PropertiesData
}

func NewSetProperties(ctx context.Context) *SetProperties {
	return &SetProperties{make([]Value, 0), newPropertiesData(ctx)}
}

func (p *SetProperties) Copy(cache CopyCache) Properties {
	items := copyValueList(p.items, cache)

	innerCpy := p.PropertiesData.copy(cache)

	return &SetProperties{items, innerCpy}
}

func (p *SetProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*SetProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	mergedItems := mergeValueLists(p.items, other.items)
	if mergedItems != nil {
		return &SetProperties{mergedItems, newPropertiesDataWithContent(mergedProps, p.Context())}
	}

	allItems := append(p.items, other.items...)
	return &SetProperties{allItems, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func (p *SetProperties) AppendItem(item Value) {
	p.items = append(p.items, item)
}

func (p *SetProperties) GetItems() []Value {
	return p.items
}

func AssertSetProperties(p_ Properties) *SetProperties {
	p, ok := p_.(*SetProperties)
	if ok {
		return p
	} else {
		panic("not SetProperties")
	}
}

func (p *SetProperties) LoopNestedPrototypes(fn func(Prototype)) {
	for _, item := range p.items {
		item.LoopNestedPrototypes(fn)
	}
}
