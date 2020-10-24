package values

import (
	"../../context"
)

type BooleanProperties struct {
	isLiteral bool
	value     bool
	PropertiesData
}

func NewBooleanProperties(isLiteral bool, value bool, ctx context.Context) Properties {
	return &BooleanProperties{isLiteral, value, newPropertiesData(ctx)}
}

func (p *BooleanProperties) Copy(cache CopyCache) Properties {
	innerCpy := p.PropertiesData.copy(cache)

	return &BooleanProperties{p.isLiteral, p.value, innerCpy}
}

func (p *BooleanProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*BooleanProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	if p.isLiteral && other.isLiteral && p.value == other.value {
		return &BooleanProperties{true, p.value, newPropertiesDataWithContent(mergedProps, p.Context())}
	}

	return &BooleanProperties{false, false, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func (p *BooleanProperties) RemoveLiteralness() *BooleanProperties {
	return &BooleanProperties{false, false, newPropertiesDataWithContent(p.props, p.Context())}
}

func (p *BooleanProperties) LiteralValue() (bool, bool) {
	return p.value, p.isLiteral
}
