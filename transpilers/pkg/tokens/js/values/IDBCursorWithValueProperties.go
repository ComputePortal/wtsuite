package values

import (
	"../../context"
)

type IDBCursorWithValueProperties struct {
	value Value
	PropertiesData
}

func NewIDBCursorWithValueProperties(v Value, ctx context.Context) Properties {
	return &IDBCursorWithValueProperties{v, newPropertiesData(ctx)}
}

func (p *IDBCursorWithValueProperties) Copy(cache CopyCache) Properties {
	var vCpy Value = nil
	if p.value != nil {
		vCpy = p.value.Copy(cache)
	}

	innerCpy := p.PropertiesData.copy(cache)

	return &IDBCursorWithValueProperties{vCpy, innerCpy}
}

func (p *IDBCursorWithValueProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*IDBCursorWithValueProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	if p.value == nil && other.value == nil {
		return &IDBCursorWithValueProperties{p.value, newPropertiesDataWithContent(mergedProps, p.Context())}
	} else if p.value == nil && other.value != nil {
		return nil
	} else if p.value != nil && other.value == nil {
		return nil
	}

	mergedValue := p.value.Merge(other.value)
	if mergedValue == nil {
		return nil
	}

	return &IDBCursorWithValueProperties{mergedValue, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func (p *IDBCursorWithValueProperties) Value() Value {
	return p.value
}

func (p *IDBCursorWithValueProperties) LoopNestedPrototypes(fn func(Prototype)) {
	p.PropertiesData.LoopNestedPrototypes(fn)

	p.value.LoopNestedPrototypes(fn)
}

func AssertIDBCursorWithValueProperties(p_ Properties) *IDBCursorWithValueProperties {
	if p, ok := p_.(*IDBCursorWithValueProperties); ok {
		return p
	} else {
		panic("not IDBCursorWithValueProperties")
	}
}
