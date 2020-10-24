package prototypes

import (
	"../values"

	"../../context"
)

var EventTarget *BuiltinPrototype = allocBuiltinPrototype()

func AddEventListener(stack values.Stack, this *values.Instance,
	args []values.Value, ctx context.Context) (values.Value, error) {
	eventProto := Event
	if str, ok := args[0].LiteralStringValue(); ok {
		switch str {
		case "click", "dblclick", "mousedown", "mouseup", "mousemove", "mouseover", "mouseleave", "mouseenter":
			eventProto = MouseEvent
		case "keydown", "keyup":
			eventProto = KeyboardEvent
		case "wheel":
			eventProto = WheelEvent
		case "hashchange":
			eventProto = HashChangeEvent
		}
	}

	target := this
	if this.IsInstanceOf(Window) {
		target = NewInstance(HTMLElement, ctx)
	}

	event := NewAltEvent(eventProto, target, ctx)

  if err := args[1].EvalMethod(stack.Parent(), []values.Value{event},
    ctx); err != nil {
    return nil, err
  }

	return nil, nil
}

func generateEventTargetPrototype() bool {
	*EventTarget = BuiltinPrototype{
		"EventTarget", nil,
		map[string]BuiltinFunction{
			"dispatchEvent":    NewMethodLikeNormal(Event, Boolean),
			"addEventListener": NewNormalFunction(&And{String, &Function{}}, AddEventListener),
		},
		nil,
	}

	return true
}

var _EventTargetOk = generateEventTargetPrototype()
