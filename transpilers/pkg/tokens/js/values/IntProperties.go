package values

import (
	"../../context"
)

type IntProperties struct {
	isLiteral bool
	value     int
	PropertiesData
}

func NewIntProperties(isLiteral bool, value int, ctx context.Context) Properties {
	return &IntProperties{isLiteral, value, newPropertiesData(ctx)}
}

func (p *IntProperties) Copy(cache CopyCache) Properties {
	innerCpy := p.PropertiesData.copy(cache)

	return &IntProperties{p.isLiteral, p.value, innerCpy}
}

func (p *IntProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*IntProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	if p.isLiteral && other.isLiteral && p.value == other.value {
		return &IntProperties{true, p.value, newPropertiesDataWithContent(mergedProps, p.Context())}
	} else {
		return &IntProperties{false, 0, newPropertiesDataWithContent(mergedProps, p.Context())}
	}
}

func (p *IntProperties) RemoveLiteralness() *IntProperties {
	return &IntProperties{false, 0, newPropertiesDataWithContent(p.props, p.Context())}
}

func (p *IntProperties) LiteralValue() (int, bool) {
	return p.value, p.isLiteral
}

func AssertIntProperties(p_ Properties) *IntProperties {
	p, ok := p_.(*IntProperties)
	if ok {
		return p
	} else {
		panic("not IntProperties")
	}
}
