package values

import (
	"fmt"

	"../../context"
)

type StringProperties struct {
	isLiteral bool
	value     string
	PropertiesData
}

func NewStringProperties(isLiteral bool, value string, ctx context.Context) Properties {
	return &StringProperties{isLiteral, value, newPropertiesData(ctx)}
}

func (p *StringProperties) Copy(cache CopyCache) Properties {
	innerCpy := p.PropertiesData.copy(cache)

	return &StringProperties{p.isLiteral, p.value, innerCpy}
}

func (p *StringProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*StringProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	if p.isLiteral && other.isLiteral && p.value == other.value && p.value != "" {
		if p.value == "" {
			fmt.Println("empty string literal")
		}
		return &StringProperties{true, p.value, newPropertiesDataWithContent(mergedProps, p.Context())}
	} else {
		return &StringProperties{false, "", newPropertiesDataWithContent(mergedProps, p.Context())}
	}
}

func (p *StringProperties) RemoveLiteralness() *StringProperties {
	return &StringProperties{false, "", newPropertiesDataWithContent(p.props, p.Context())}
}

func (p *StringProperties) LiteralValue() (string, bool) {
	return p.value, p.isLiteral
}

func AssertStringProperties(p_ Properties) *StringProperties {
	if p, ok := p_.(*StringProperties); ok {
		return p
	} else {
		panic("not StringProperties")
	}
}
