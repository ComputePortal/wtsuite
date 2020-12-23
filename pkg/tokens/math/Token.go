package math

import (
	"github.com/computeportal/wtsuite/pkg/tokens/math/boundingbox"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Token interface {
	Dump(indent string) string // to inspect the syntax-tree
	Context() context.Context
	GenerateTags(scope Scope, x float64, y float64) (boundingbox.BB, error)
}

type TokenData struct {
	ctx context.Context
}

type DumpableData struct {
	name string
}

func newTokenData(ctx context.Context) TokenData {
	return TokenData{ctx}
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
