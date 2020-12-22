package html

import ()

type Dict interface {
	Container
	IsEmpty() bool
}

func IsDict(t Token) bool {
	_, ok := t.(Dict)
	return ok
}

func AssertDict(t Token) (Dict, error) {
	if d, ok := t.(Dict); ok {
		return d, nil
	} else {
		errCtx := d.Context()
		return nil, errCtx.NewError("Error: expected a dict")
	}
}
