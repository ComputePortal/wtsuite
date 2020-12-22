package html

type Primitive interface {
	Token
	Write() string
}

func IsPrimitive(t Token) bool {
	switch t.(type) {
	case *Int, *Float, *String, *Flag, *Color:
		return true
	default:
		return false
	}
}

func AssertPrimitive(t Token) (Primitive, error) {
	if IsPrimitive(t) {
		if res, ok := t.(Primitive); ok {
			return res, nil
		} else {
			panic("bad primitive")
		}
	}

	errCtx := t.Context()
	err := errCtx.NewError("Error: expected primitive")
	return nil, err
}
