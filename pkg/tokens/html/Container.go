package html

import (
	"../context"
)

// implemented specifically by directives.PStyle
// TODO: should this be implemented by RawDict, List etc? (right Get() for those containers is a little different)
type Indexable interface {
  Get(scope Scope, x Token, ctx context.Context) (Token, error)
}

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

func IsIndexable(t Token) bool {
	_, ok := t.(Indexable)
	return ok
}

func AssertIndexable(t Token) (Indexable, error) {
	if IsIndexable(t) {
		if res, ok := t.(Indexable); ok {
			return res, nil
		} else {
			panic("bad container")
		}
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected indexable")
	}
}
