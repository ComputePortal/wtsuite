package raw

import (
	"fmt"
	"strconv"

	"../context"
	"../patterns"
)

type LiteralFloat struct {
	value float64
	unit  string
	TokenData
}

func parseValueUnit(s string, ctx context.Context) (float64, string, error) {
	n := len(s)

	vs := s // unparsed value string
	unit := ""
	if !patterns.EndsWithDigit(s) {
		var ok bool
		if unit, ok = patterns.ExtractUnit(s); !ok {
			errCtx := ctx
			return 0.0, "", errCtx.NewError("Syntax Error: invalid unit (" + s + ")")
		}

		vs = s[0 : n-len(unit)]
	}

	if v, err := strconv.ParseFloat(vs, 64); err != nil {
		// TODO: give valid units as info
		errCtx := ctx.NewContext(0, len(s))
		return v, unit, errCtx.NewError("Syntax Error: invalid float")
	} else {
		return v, unit, nil
	}
}

func NewValueLiteralFloat(value float64, ctx context.Context) (*LiteralFloat, error) {
	return NewValueUnitLiteralFloat(value, "", ctx)
}

func NewValueUnitLiteralFloat(value float64, unit string, ctx context.Context) (*LiteralFloat, error) {
	return &LiteralFloat{value, unit, TokenData{ctx}}, nil
}

func NewLiteralFloat(s string, ctx context.Context) (*LiteralFloat, error) {
	value, unit, err := parseValueUnit(s, ctx)
	if err != nil {
		return nil, err
	}

	return NewValueUnitLiteralFloat(value, unit, ctx)
}

func (t *LiteralFloat) Value() float64 {
	return t.value
}

func (t *LiteralFloat) Unit() string {
	return t.unit
}

func (t *LiteralFloat) Dump(indent string) string {
	s := fmt.Sprintf("%g%s", t.value, t.unit)

	return indent + "LiteralFloat(" + s + ")\n"
}

func IsLiteralFloat(t Token) bool {
	_, ok := t.(*LiteralFloat)
	return ok
}

// unit can be '*' wild card
func AssertLiteralFloat(t Token, unit string) (*LiteralFloat, error) {
	if f, ok := t.(*LiteralFloat); !ok {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected literal float")
	} else {
		if unit == "*" {
			return f, nil
		} else if f.Unit() != unit {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: expected literal float with unit " + unit)
		} else {
			return f, nil
		}
	}
}
