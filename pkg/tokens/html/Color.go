package html

import (
	"fmt"
	//"strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Color struct {
	r, g, b, a int
	TokenData
}

func NewValueColor(r, g, b, a int, ctx context.Context) *Color {
	return &Color{r, g, b, a, TokenData{ctx}}
}

func NewColor(r, g, b, a int, ctx context.Context) (*Color, error) {
	return NewValueColor(r, g, b, a, ctx), nil
}

func (t *Color) Values() (r, g, b, a int) {
	r, g, b, a = t.r, t.g, t.b, t.a
	return r, g, b, a
}

func (t *Color) FloatValues() (float64, float64, float64, float64) {
	return float64(t.r) / 255.0, float64(t.g) / 255.0, float64(t.b) / 255.0, float64(t.a) / 255.0
}

func (t *Color) IsPrimitive() bool {
	return true
}

func (t *Color) Eval(scope Scope) (Token, error) {
	return t, nil
}

func (t *Color) EvalLazy(tag FinalTag) (Token, error) {
	return t, nil
}

func (t *Color) Write() string {
	formatHex := func(i int) string {
		return fmt.Sprintf("%02x", i)
	}

	/*formatInt := func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	}*/

	if t.a == 255 {

		s := formatHex(t.r) + formatHex(t.g) + formatHex(t.b)
		if s[0] == s[1] && s[2] == s[3] && s[4] == s[5] {
			return "#" + s[0:1] + s[2:3] + s[4:5]
		} else {
			return "#" + s
		}
	} else {
		//return "rgba(" + formatInt(t.r) + "," + formatInt(t.g) + "," + formatInt(t.b) + "," + formatInt(t.a) + ")"
		return "#" + formatHex(t.r) + formatHex(t.g) + formatHex(t.b) + formatHex(t.a)
	}
}

func (t *Color) Dump(indent string) string {
	return indent + "Color(" + t.Write() + ")\n"
}

func IsColor(t Token) bool {
	_, ok := t.(*Color)
	return ok
}

func (a *Color) IsSame(other Token) bool {
	if b, ok := other.(*Color); ok {
		return (a.r == b.r) && (a.g == b.g) && (a.b == b.b) && (a.a == b.a)
	} else {
		return false
	}
}

func AssertColor(t Token) (*Color, error) {
	if c, ok := t.(*Color); ok {
		return c, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected Color")
	}
}
