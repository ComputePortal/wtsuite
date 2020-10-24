package prototypes

import (
	"fmt"

	"../values"

	"../../context"
)

type EventPrototype struct {
	BuiltinPrototype
}

var Event *EventPrototype = allocEventPrototype()

func allocEventPrototype() *EventPrototype {
	return &EventPrototype{BuiltinPrototype{
		"", nil,
		map[string]BuiltinFunction{},
		nil,
	}}
}

func (p *EventPrototype) Check(args []interface{}, pos int, ctx context.Context) (int, error) {
	return CheckPrototype(p, args, pos, ctx)
}

func (p *EventPrototype) HasAncestor(other_ values.Interface) bool {
	if other, ok := other_.(*EventPrototype); ok {
		if other == p {
			return true
		} else {
			parent := p.GetParent()
			if parent != nil {
				return parent.HasAncestor(other_)
			} else {
				return false
			}
		}
	} else {
		_, ok = other_.IsImplementedBy(p)
		return ok
	}
}

// TODO: do other Events also need this function?
func (p *EventPrototype) CastInstance(v *values.Instance, typeChildren []*values.NestedType, ctx context.Context) (values.Value, error) {
	newV_, ok := v.ChangeInstanceInterface(p, false, true)
	if !ok {
		return nil, ctx.NewError("Error: " + v.TypeName() + " doesn't inherit from " + p.Name())
	}

	newV, ok := newV_.(*values.Instance)
	if !ok {
		panic("unexpected")
	}

	if typeChildren == nil {
		return newV, nil
	} else {
		if len(typeChildren) != 1 {
			return nil, ctx.NewError(fmt.Sprintf("Error: Event expects 1 type child, got %d", len(typeChildren)))
		}

		typeChild := typeChildren[0]

		// now cast all the items
		props := values.AssertEventProperties(newV.Properties())

		target := props.Target()
		var err error
		target, err = target.Cast(typeChild, ctx)
		if err != nil {
			return nil, err
		}

		newV = NewEvent(target, ctx)
		return newV, nil
	}
}

func NewEvent(target values.Value, ctx context.Context) *values.Instance {
	return values.NewInstance(Event, values.NewEventProperties(target, ctx), ctx)
}

func NewAltEvent(proto values.Prototype, target values.Value,
	ctx context.Context) *values.Instance {
	return values.NewInstance(proto, values.NewEventProperties(target, ctx), ctx)
}

func generateEventPrototype() bool {
	*Event = EventPrototype{BuiltinPrototype{
		"Event", nil,
		map[string]BuiltinFunction{
			"preventDefault":           NewNormal(&None{}, nil),
			"stopPropagation":          NewNormal(&None{}, nil),
			"stopImmediatePropagation": NewNormal(&None{}, nil),
			"target": NewGetterFunction(func(stack values.Stack, this *values.Instance,
				args []values.Value, ctx context.Context) (values.Value, error) {
				props := values.AssertEventProperties(this.Properties())

				target := props.Target()
				if target == nil {
					return nil, ctx.NewError("Error: event.target unset")
				}
				return target, nil
			}),
		},
		NewConstructorFunction(func(stack values.Stack, args []values.Value,
			ctx context.Context) (values.Value, error) {
			if err := CheckInputs(&And{String, &Opt{Object}}, args, ctx); err != nil {
				return nil, err
			}

			return NewEvent(nil, ctx), nil
		}),
	}}

	return true
}

var _EventOk = generateEventPrototype()
