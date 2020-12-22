package html

import (
	"fmt"
	"math"
	"reflect"

	"../context"
)

var PX_PER_REM = 0

type Float struct {
	value float64
	unit  string
	TokenData
}

func NewValueFloat(value float64, ctx context.Context) *Float {
	return NewValueUnitFloat(value, "", ctx)
}

func NewValueUnitFloat(value float64, unit string, ctx context.Context) *Float {
	return &Float{value, unit, TokenData{ctx}}
}

func (t *Float) Value() float64 {
	return t.value
}

func (t *Float) Unit() string {
	return t.unit
}

func (t *Float) Eval(scope Scope) (Token, error) {
	return t, nil
}

func formatPx(fl float64) string {
	if math.Mod(fl, 1.0) == 0.0 {
		return fmt.Sprintf("%gpx", fl)
	} else {
		return fmt.Sprintf("%.01fpx", fl)
	}
}

func (t *Float) Write() string {
	// experimental:
	if t.unit == "rem" && PX_PER_REM != 0 {
		return formatPx(t.value * float64(PX_PER_REM))
	} else if t.unit == "px" {
		return formatPx(t.value)
	} else {
		return fmt.Sprintf("%g%s", t.value, t.unit)
	}
}

func (t *Float) Dump(indent string) string {
	return indent + "Float(" + t.Write() + ")\n"
}

func IsFloat(t Token) bool {
	_, ok := t.(*Float)
	return ok
}

// unit can be wild card
func AssertFloat(t Token, unit string) (*Float, error) {
	if x, ok := t.(*Float); !ok {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected Float")
	} else {
		if unit == "*" {
			return x, nil
		} else if x.Unit() != unit {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: expected Float with unit " + unit)
		} else {
			return x, nil
		}
	}
}

func AssertFractionFloat(t Token) (*Float, error) {
	if f, ok := t.(*Float); ok {
		if f.Unit() == "%" {
			return NewValueFloat(f.Value()*0.01, f.Context()), nil
		} else if f.Unit() == "" {
			if f.Value() >= 0.0 && f.Value() <= 1.0 {
				return f, nil
			} else {
				errCtx := f.Context()
				return nil, errCtx.NewError(fmt.Sprintf("Error: float '%f' is out of range [0,1]", f.Value()))
			}
		} else {
			errCtx := t.Context()
			return nil, errCtx.NewError("Error: expected unitless or '%' Float, got '" + f.Unit() + "'")
		}
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected Float")
	}
}

// unitless!
func AssertIntOrFloat(t Token) (*Float, error) {
	if IsInt(t) {
		i, err := AssertInt(t)
		if err != nil {
			panic(err)
		}

		return NewValueFloat(float64(i.Value()), t.Context()), nil
	} else if !IsFloat(t) {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected number (int or float, got " + reflect.TypeOf(t).String() + ")")
	} else {
		return AssertFloat(t, "")
	}
}

// always unitless
func IsIntOrFloat(t Token) bool {
	if IsInt(t) {
		return true
	} else if fl, ok := t.(*Float); ok && fl.Unit() == "" {
		return true
	} else {
		return false
	}
}

// can contain units, ints converted to non-unit float
func AssertAnyIntOrFloat(t Token) (*Float, error) {
	if f, ok := t.(*Float); ok {
		return f, nil
	} else if i, ok := t.(*Int); ok {
		return NewValueFloat(float64(i.Value()), i.Context()), nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected Float or Int")
	}
}

func (a *Float) IsSame(other Token) bool {
	if b, ok := other.(*Float); ok {
		return a.value == b.value && a.unit == b.unit
	} else {
		return false
	}
}
