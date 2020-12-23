package raw

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type LiteralColor struct {
	r, g, b, a int
	TokenData
}

func NewLiteralColor(s string, ctx context.Context) (*LiteralColor, error) {
	parseHex1 := func(ss string, prevErr error) (int, error) {
		v, err := strconv.ParseInt("0x"+ss+ss, 0, 64)
		if err != nil {
			return 0, err
		} else {
			return int(v), prevErr
		}
	}

	parseHex2 := func(ss string, prevErr error) (int, error) {
		v, err := strconv.ParseInt("0x"+ss, 0, 64)
		if err != nil {
			return 0, err
		} else {
			return int(v), prevErr
		}
	}

	r, g, b, a := 255, 255, 255, 255

	s = s[1:] // cut off the `#`

	n := len(s)
	var err error = nil
	switch n {
	case 3, 4:
		r, err = parseHex1(s[0:1], err)
		g, err = parseHex1(s[1:2], err)
		b, err = parseHex1(s[2:3], err)

		if n == 4 {
			a, err = parseHex1(s[3:4], err)
		}
	case 6, 8:
		r, err = parseHex2(s[0:2], err)
		g, err = parseHex2(s[2:4], err)
		b, err = parseHex2(s[4:6], err)

		if n == 8 {
			a, err = parseHex2(s[6:8], err)
		}
	default:
		err = errors.New("")
	}

	if err != nil {
		err = ctx.NewError("Syntax Error: invalid color")
	}

	return &LiteralColor{r, g, b, a, TokenData{ctx}}, err
}

func (t *LiteralColor) Values() (r, g, b, a int) {
	r, g, b, a = t.r, t.g, t.b, t.a
	return r, g, b, a
}

func (t *LiteralColor) Dump(indent string) string {
	s := fmt.Sprintf("%d, %d, %d, %d", t.r, t.g, t.b, t.a)
	return indent + "LiteralColor(" + s + ")\n"
}

func IsLiteralColor(t Token) bool {
	_, ok := t.(*LiteralColor)
	return ok
}

func AssertLiteralColor(t Token) (*LiteralColor, error) {
	if cl, ok := t.(*LiteralColor); ok {
		return cl, nil
	}

	errCtx := t.Context()
	return nil, errCtx.NewError("Error: expected literal color")
}
