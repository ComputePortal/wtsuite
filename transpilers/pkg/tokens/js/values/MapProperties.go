package values

import (
	"../../context"
)

type MapProperties struct {
	keys  []Value
	items []Value
	PropertiesData
}

func NewMapProperties(ctx context.Context) *MapProperties {
	return &MapProperties{[]Value{}, []Value{}, newPropertiesData(ctx)}
}

func (p *MapProperties) Copy(cache CopyCache) Properties {
	keys := copyValueList(p.keys, cache)
	items := copyValueList(p.items, cache)

	innerCpy := p.PropertiesData.copy(cache)

	return &MapProperties{keys, items, innerCpy}
}

func (p *MapProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*MapProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	// TODO: key-item pairs can be in differing order
	mergedKeys := mergeValueListsSameOrder(p.keys, other.keys)
	if mergedKeys == nil {
		return nil
	}

	mergedItems := mergeValueListsSameOrder(p.items, other.items)
	if mergedItems == nil {
		return nil
	}

	return &MapProperties{mergedKeys, mergedItems, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func (p *MapProperties) AppendKey(key Value) {
	p.keys = append(p.keys, key)
}

func (p *MapProperties) AppendItem(item Value) {
	p.items = append(p.items, item)
}

func (p *MapProperties) GetKey() Value {
	if len(p.keys) == 0 {
		return NewAllNull(p.Context())
	} else {
		return NewMulti(p.keys, p.Context())
	}
}

func (p *MapProperties) GetItem() Value {
	if len(p.items) == 0 {
		return NewAllNull(p.Context())
	} else {
		return NewMulti(p.items, p.Context())
	}
}

func (p *MapProperties) GetKeys() []Value {
	return p.keys
}

func (p *MapProperties) GetItems() []Value {
	return p.items
}

func AssertMapProperties(p_ Properties) *MapProperties {
	p, ok := p_.(*MapProperties)
	if ok {
		return p
	} else {
		panic("not MapProperties")
	}
}

func (p *MapProperties) LoopNestedPrototypes(fn func(Prototype)) {
	p.PropertiesData.LoopNestedPrototypes(fn)

	for _, key := range p.keys {
		key.LoopNestedPrototypes(fn)
	}

	for _, item := range p.items {
		item.LoopNestedPrototypes(fn)
	}
}
