package html

import (
	"fmt"
	"strings"
)

func writeJSON(argToken Token) (string, error) {
	var b strings.Builder

	switch tt := argToken.(type) {
	case *StringDict:
		b.WriteString("{")
		if err := tt.Loop(func(key *String, val Token, last bool) error {
			b.WriteString("\"")
			b.WriteString(key.Value())
			b.WriteString("\":")
			sub, err := writeJSON(val)
			if err != nil {
				return err
			}
			b.WriteString(sub)
			if !last {
				b.WriteString(",")
			}
			return nil
		}); err != nil {
			return b.String(), err
		}
		b.WriteString("}")
	case *IntDict:
		b.WriteString("{")
		if err := tt.Loop(func(key *Int, val Token, last bool) error {
			b.WriteString("\"")
			b.WriteString(fmt.Sprintf("%d", key.Value()))
			b.WriteString("\":")
			sub, err := writeJSON(val)
			if err != nil {
				return err
			}
			b.WriteString(sub)
			if !last {
				b.WriteString(",")
			}
			return nil
		}); err != nil {
			return b.String(), err
		}
		b.WriteString("}")
	case *List:
		b.WriteString("[")
		if err := tt.Loop(func(i int, item Token, last bool) error {
			sub, err := writeJSON(item)
			if err != nil {
				return err
			}
			b.WriteString(sub)
			if !last {
				b.WriteString(",")
			}
			return nil
		}); err != nil {
			return b.String(), err
		}
		b.WriteString("]")
	case *Int:
		b.WriteString(tt.Write())
	case *Float:
		if tt.Unit() != "" {
			b.WriteString("\"")
		}
		b.WriteString(tt.Write())
		if tt.Unit() != "" {
			b.WriteString("\"")
		}
	case *Bool:
		if tt.Value() {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case *Null:
		b.WriteString("null")
	case Primitive:
		b.WriteString("\"")
		b.WriteString(tt.Write())
		b.WriteString("\"")
	default:
		errCtx := tt.Context()
		return b.String(), errCtx.NewError("Error: expected string, dict or list, or element or constructor")
	}

	return b.String(), nil
}
