package values

import (
	"../../context"
)

type NumberProperties struct {
	isLiteral bool
	value     float64
	PropertiesData
}

func NewNumberProperties(isLiteral bool, value float64, ctx context.Context) Properties {
	return &NumberProperties{isLiteral, value, newPropertiesData(ctx)}
}

func (p *NumberProperties) Copy(cache CopyCache) Properties {
	innerCpy := p.PropertiesData.copy(cache)

	return &NumberProperties{p.isLiteral, p.value, innerCpy}
}

func (p *NumberProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*NumberProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	if p.isLiteral && other.isLiteral && p.value == other.value {
		return &NumberProperties{true, p.value, newPropertiesDataWithContent(mergedProps, p.Context())}
	} else {
		return &NumberProperties{false, 0.0, newPropertiesDataWithContent(mergedProps, p.Context())}
	}
}

func (p *NumberProperties) RemoveLiteralness() *NumberProperties {
	return &NumberProperties{false, 0.0, newPropertiesDataWithContent(p.props, p.Context())}
}

func (p *NumberProperties) LiteralValue() (float64, bool) {
	return p.value, p.isLiteral
}

func AssertNumberProperties(p_ Properties) *NumberProperties {
	p, ok := p_.(*NumberProperties)
	if ok {
		return p
	} else {
		panic("not NumberProperties")
	}
}
