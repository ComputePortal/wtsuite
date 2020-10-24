package values

import (
	"../../context"
)

type IDBRequestProperties struct {
	result Value // XXX: can result be nil?
	PropertiesData
}

func NewIDBRequestProperties(result Value, ctx context.Context) Properties {
	return &IDBRequestProperties{result, newPropertiesData(ctx)}
}

func (p *IDBRequestProperties) Copy(cache CopyCache) Properties {
	var resultCpy Value = nil
	if p.result != nil {
		resultCpy = p.result.Copy(cache)
	}

	innerCpy := p.PropertiesData.copy(cache)

	return &IDBRequestProperties{resultCpy, innerCpy}
}

func (p *IDBRequestProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*IDBRequestProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	if p.result == nil && other.result == nil {
		return &IDBRequestProperties{p.result, newPropertiesDataWithContent(mergedProps, p.Context())}
	} else if p.result == nil && other.result != nil {
		return nil
	} else if p.result != nil && other.result == nil {
		return nil
	}

	mergedResult := p.result.Merge(other.result)
	if mergedResult == nil {
		return nil
	}

	return &IDBRequestProperties{mergedResult, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func AssertIDBRequestProperties(p_ Properties) *IDBRequestProperties {
	if p, ok := p_.(*IDBRequestProperties); ok {
		return p
	} else {
		panic("not IDBRequestProperties")
	}
}

func (p *IDBRequestProperties) Result() Value {
	return p.result
}

func (p *IDBRequestProperties) LoopNestedPrototypes(fn func(Prototype)) {
	p.PropertiesData.LoopNestedPrototypes(fn)

	p.result.LoopNestedPrototypes(fn)
}
