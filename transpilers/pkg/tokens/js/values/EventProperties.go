package values

import (
	"../../context"
)

type EventProperties struct {
	target Value
	PropertiesData
}

func NewEventProperties(target Value, ctx context.Context) Properties {
	return &EventProperties{target, newPropertiesData(ctx)}
}

func (p *EventProperties) Copy(cache CopyCache) Properties {
	innerCpy := p.PropertiesData.copy(cache)

	targetCpy := p.target.Copy(cache)

	return &EventProperties{targetCpy, innerCpy}
}

func (p *EventProperties) Merge(other_ Properties) Properties {
	other, ok := other_.(*EventProperties)
	if !ok {
		return nil
	}

	mergedProps := p.merge(other.props)
	if mergedProps == nil {
		return nil
	}

	mergedTarget := p.target.Merge(other.target)
	if mergedTarget == nil {
		return nil
	}

	return &EventProperties{mergedTarget, newPropertiesDataWithContent(mergedProps, p.Context())}
}

func (p *EventProperties) Target() Value {
	return p.target
}

func AssertEventProperties(props Properties) *EventProperties {
	if ep, ok := props.(*EventProperties); ok {
		return ep
	} else {
		panic("not event properties")
	}
}

func (p *EventProperties) LoopNestedPrototypes(fn func(Prototype)) {
	p.PropertiesData.LoopNestedPrototypes(fn)

	p.target.LoopNestedPrototypes(fn)
}
