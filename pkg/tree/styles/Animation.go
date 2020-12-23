package styles

import (
	"fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/raw"
)

var (
	_uanim = 0

	_animDigests = make(map[string]string)
)

type Keyframe struct {
	pos   float64
	state map[string]string
}

type Keyframes struct {
	name   string
	frames []Keyframe
}

func NewUniqueAnimationName(digest string) string {
	if id, ok := _animDigests[digest]; ok {
		return id
	} else {
		res := fmt.Sprintf("_a%d", _uanim)

		_uanim++

		_animDigests[digest] = res

		return res
	}
}

func NewKeyframes(name string, attr *tokens.StringDict) (*Keyframes, error) {
	frames := make([]Keyframe, 0)

	if err := attr.Loop(func(k *tokens.String, v tokens.Token, last bool) error {
		stateToken, err := tokens.AssertStringDict(v)
		if err != nil {
			return err
		}

		state, err := dictToStringMap(stateToken, attr.Context())
		if err != nil {
			return err
		}

		switch strings.TrimSpace(k.Value()) {
		case "from":
			frames = append(frames, Keyframe{0.0, state})
		case "to":
			frames = append(frames, Keyframe{1.0, state})
		default:
			posToken, err := raw.NewLiteralFloat(k.Value(), k.InnerContext())
			if err != nil {
				return err
			}

			switch posToken.Unit() {
			case "%":
				frames = append(frames, Keyframe{posToken.Value() / 100.0, state})
			case "":
				frames = append(frames, Keyframe{posToken.Value(), state})
			default:
				errCtx := posToken.Context()
				return errCtx.NewError("Error: expected unitless or %")
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Keyframes{name, frames}, nil
}

func (r *Keyframes) writeStart(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("@keyframes ")
	b.WriteString(r.name)
	b.WriteString("{")
	b.WriteString(NL)

	return b.String()
}

func (r *Keyframes) writeFrames(indent string) string {
	var b strings.Builder

	n := len(r.frames)
	for i, frame := range r.frames {
		b.WriteString(indent)

		if i < n-1 || frame.pos != 1.0 {
			b.WriteString(fmt.Sprintf("%0.0f", frame.pos))
			b.WriteString("%")
		} else {
			b.WriteString("100%")
		}

		b.WriteString("{")

		b.WriteString(stringMapToString(frame.state, "", ""))

		b.WriteString("}")
		b.WriteString(NL)
	}

	return b.String()
}

func (r *Keyframes) writeStop(indent string) string {
	return indent + "}" + NL
}

func (r *Keyframes) Write(indent string) string {
	var b strings.Builder

	b.WriteString(r.writeStart(indent))
	b.WriteString(r.writeFrames(indent + TAB))
	b.WriteString(r.writeStop(indent))

	return b.String()
}

func (r *Keyframes) ExpandNested(sel Selector) ([]Rule, error) {
	panic("this is the result of ExpandNested() (can't expand twice)")
}

func Animation(sel Selector, args []string, v tokens.Token, ctx context.Context) ([]Rule, error) {
	if len(args) != 0 {
		return nil, ctx.NewError("Error: expected 0 arguments")
	}

	if tokens.IsNull(v) {
		return []Rule{}, nil
	}

	attr, err := tokens.AssertStringDict(v)
	if err != nil {
		return nil, err
	}

	var keyframesToken *tokens.StringDict = nil // needed
	var durationToken *tokens.Float = nil       // needed
	var timingFnToken *tokens.String = nil      // optional
	var delayToken *tokens.Float = nil          // optional
	var nItersToken tokens.Primitive = nil      // optional
	var directionToken *tokens.String = nil     // optional
	var fillModeToken *tokens.String = nil      // optional
	var playStateToken *tokens.String = nil     // optional

	if err := attr.Loop(func(k *tokens.String, v tokens.Token, last bool) error {
		var err error
		switch k.Value() {
		case "keyframes":
			// TODO: maybe MixedDict is better?
			keyframesToken, err = tokens.AssertStringDict(v)
			if err != nil {
				return err
			}
		case "duration":
			durationToken, err = tokens.AssertFloat(v, "s")
			if err != nil {
				return err
			}
		case "timing-function":
			timingFnToken, err = tokens.AssertString(v)
			if err != nil {
				return err
			}
		case "delay":
			delayToken, err = tokens.AssertFloat(v, "s")
			if err != nil {
				return err
			}
		case "iteration-count":
			nItersToken, err = tokens.AssertPrimitive(v)
			if err != nil {
				return err
			}
		case "direction":
			directionToken, err = tokens.AssertString(v)
			if err != nil {
				return err
			}
		case "fill-mode":
			fillModeToken, err = tokens.AssertString(v)
			if err != nil {
				return err
			}
		case "play-state":
			playStateToken, err = tokens.AssertString(v)
			if err != nil {
				return err
			}
		case "from", "to":
			errCtx := k.Context()
			return errCtx.NewError("Error: from/to must be wrapped in \"keyframes\"")
		default:
			errCtx := k.Context()
			return errCtx.NewError("Error: key not recognized")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	if keyframesToken == nil {
		errCtx := attr.Context()
		return nil, errCtx.NewError("Error: keyframes not defined")
	}

	if durationToken == nil {
		errCtx := attr.Context()
		return nil, errCtx.NewError("Error: duration not defined")
	}

	name := NewUniqueAnimationName(keyframesToken.Dump(""))

	fields := make([]string, 0)

	fields = append(fields, durationToken.Write())
	if timingFnToken != nil {
		fields = append(fields, timingFnToken.Write())
	}
	if delayToken != nil {
		fields = append(fields, delayToken.Write())
	}
	if nItersToken != nil {
		fields = append(fields, nItersToken.Write())
	}
	if directionToken != nil {
		fields = append(fields, directionToken.Write())
	}
	if fillModeToken != nil {
		fields = append(fields, fillModeToken.Write())
	}
	if playStateToken != nil {
		fields = append(fields, playStateToken.Write())
	}

	fields = append(fields, name)

	animAttr := map[string]string{
		"animation": strings.Join(fields, " "),
	}

	keyframes, err := NewKeyframes(name, keyframesToken)
	if err != nil {
		return nil, err
	}

	return []Rule{keyframes, NewSelectorRule(sel, animAttr)}, nil
}

var _animationOk = registerAtFunction("animation", Animation)
