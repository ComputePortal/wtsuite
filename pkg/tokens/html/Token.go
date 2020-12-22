package html

import (
	"../context"
)

var (
	VERBOSITY = 0
)

type Scope interface {
	// caller can differ from scope
	Eval(key string, args []Token, ctx context.Context) (Token, error)
	Permissive() bool
}

type Token interface {
	Dump(indent string) string // to inspect the syntax-tree
	Eval(scope Scope) (Token, error)
	IsSame(other Token) bool
	Context() context.Context
}

type TokenData struct {
	ctx context.Context
}

type DumpableData struct {
	name string
}

func NewDumpableData(name string) DumpableData {
	return DumpableData{name}
}

func (t *DumpableData) Dump(indent string) string {
	return indent + t.name + "\n"
}

func (t *TokenData) Context() context.Context {
	return t.ctx
}

func MergeContexts(ts ...Token) context.Context {
	ctxs := make([]context.Context, len(ts))

	for i, t := range ts {
		ctxs[i] = t.Context()
	}

	return context.MergeContexts(ctxs...)
}
