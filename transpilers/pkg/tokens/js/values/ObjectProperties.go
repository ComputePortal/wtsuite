package values

import (
	"../../context"
)

// from literal object
type ObjectProperties struct {
	isLiteral bool
	items     []Value
	PropertiesData
}

func NewObjectProperties(props map[string]Value, ctx context.Context) *ObjectProperties {
	isLiteral := props != nil
	res := &ObjectProperties{isLiteral, make([]Value, 0), newPropertiesData(ctx)}
	if props != nil {
		res.props = props
	}

	return res
}

func (p *ObjectProperties) IsLocked(s Stack) bool {
	return false
}

func (p *ObjectProperties) Copy(cache CopyCache) Properties {
	innerCpy := p.PropertiesData.copy(cache)

	itemsCpy := copyValueList(p.items, cache)

	return &ObjectProperties{p.isLiteral, itemsCpy, innerCpy}
}

func (p *ObjectProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*ObjectProperties)
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
			return &ObjectProperties{true, mergedItems, newPropertiesDataWithContent(mergedProps, p.Context())}
		}
	}

	allItems := append(p.items, other.items...)
	return &ObjectProperties{false, allItems, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func (p *ObjectProperties) LiteralValues() (map[string]Value, bool) {
	return p.props, p.isLiteral
}

func (p *ObjectProperties) RemoveLiteralness() {
	p.isLiteral = false
}

func (p *ObjectProperties) GetItem() Value {
	vs := make([]Value, 0)

	for _, v := range p.props {
		vs = append(vs, v)
	}

	if !p.isLiteral {
		for _, v := range p.items {
			vs = append(vs, v)
		}
	}

	if len(vs) == 0 {
		return NewAllNull(p.Context())
	}

	res := NewMulti(vs, p.Context())

	/*if _, ok := res.LiteralStringValue(); ok && !p.isLiteral {
		for _, vDump := range vs {
			str, _ := vDump.LiteralStringValue()
			fmt.Println(vDump.TypeName(), str, len(p.props), len(p.items))
		}
		panic("literal, yet not literal")
	}*/

	return res
}

func (p *ObjectProperties) AppendItem(item Value) {
	p.isLiteral = false

	p.items = append(p.items, item)
}

func (p *ObjectProperties) GetProperty(key string) (Value, bool) {
	if p.isLiteral {
		return p.PropertiesData.GetProperty(key)
	} else {
		return p.GetItem(), true
	}
}

func (p *ObjectProperties) SetProperty(key string, v Value) error {
	if p.isLiteral {
		return p.PropertiesData.SetProperty(key, v)
	} else {
		p.items = append(p.items, v)
		return nil
	}
}

func AssertObjectProperties(p_ Properties) *ObjectProperties {
	if p, ok := p_.(*ObjectProperties); ok {
		return p
	} else {
		panic("not ObjectProperties")
	}
}

func (p *ObjectProperties) LoopNestedPrototypes(fn func(Prototype)) {
	p.PropertiesData.LoopNestedPrototypes(fn)

	for _, item := range p.items {
		item.LoopNestedPrototypes(fn)
	}
}
