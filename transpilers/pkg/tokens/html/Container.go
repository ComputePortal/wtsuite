package html

import (
	"../context"
)

type Container interface {
	Token
	Len() int
	Copy(ctx context.Context) (Token, error)
	LoopValues(func(Token) error) error // indices in list or keys dict are ignored
}

func IsContainer(t Token) bool {
	_, ok := t.(Container)
	return ok
}

func AssertContainer(t Token) (Container, error) {
	if IsContainer(t) {
		if res, ok := t.(Container); ok {
			return res, nil
		} else {
			panic("bad container")
		}
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected container (dict or list)")
	}
}
